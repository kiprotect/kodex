// Kodex (Enterprise Edition - EE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2021  KIProtect GmbH (HRB 208395B) - All Rights Reserved

package pcap

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/kiprotect/go-helpers/forms"
	"github.com/kiprotect/kodex/api"
	"github.com/kiprotect/kodex/api/v1/resources/pcap/pcapgo"
	"github.com/kiprotect/kodex/actions/pseudonymize"
	"io"
	"io/ioutil"
	"net"
	"net/url"
	"regexp"
)

var PseudonymizeForm = forms.Form{
	ErrorMsg: "invalid data encountered in the PCAP pseudonymization form",
	Fields: []forms.Field{
		{
			Name: "key",
			Validators: []forms.Validator{
				forms.IsRequired{},
				forms.IsString{},
				forms.MatchesRegex{
					Regex: regexp.MustCompile(`^.{8,32}$`),
				},
			},
		},
		{
			Name: "preserve-subnets",
			Validators: []forms.Validator{
				forms.IsOptional{Default: true},
				forms.IsBoolean{},
			},
		},
		{
			Name: "format",
			Validators: []forms.Validator{
				forms.IsOptional{},
				forms.IsString{},
				forms.IsIn{
					Choices: []interface{}{"pcap"},
				},
			},
		},
	},
	Transforms: []forms.Transform{},
}

func process(c *gin.Context, params map[string]interface{}, f func(interface{}) (interface{}, error)) {

	if c.Request.Body == nil {
		c.JSON(400, map[string]string{"message": "no body data given"})
		return
	}

	// to do: make this more efficient
	pseudonymizeIP := func(in net.IP) (net.IP, error) {
		if res, err := f(in.String()); err != nil {
			return nil, err
		} else {
			strRes, ok := res.(string)
			if !ok {
				return nil, fmt.Errorf("Not a string")
			}
			out := net.ParseIP(strRes)
			if out == nil {
				return nil, fmt.Errorf("Not an IP")
			}
			return out, nil
		}
	}

	returnWithError := func(err error) {
		c.JSON(400, map[string]string{"message": "an error occurred", "error": err.Error()})
	}

	b, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		returnWithError(err)
		return
	}
	reader := bytes.NewReader(b)
	r, err := pcapgo.NewReader(reader)

	if err != nil {
		returnWithError(err)
		return
	}
	packetWriter := pcapgo.NewWriter(c.Writer)
	packetWriter.WriteFileHeader(r.Snaplen(), r.LinkType())
	for {
		data, ci, err := r.ReadPacketData()
		if err == io.EOF || data == nil {
			break
		}
		if err != nil {
			returnWithError(err)
			return
		}

		packet := gopacket.NewPacket(data, layers.LayerTypeEthernet, gopacket.DecodeOptions{Lazy: true, NoCopy: true})
		ethernet := packet.Layer(layers.LayerTypeEthernet).(*layers.Ethernet)
		ipv4 := packet.Layer(layers.LayerTypeIPv4)
		ipv6 := packet.Layer(layers.LayerTypeIPv6)
		arp := packet.Layer(layers.LayerTypeARP)

		buf := gopacket.NewSerializeBuffer()
		opts := gopacket.SerializeOptions{}

		if arp != nil {
			arpPacket := arp.(*layers.ARP)
			ipAddressSrc := net.IP(arpPacket.SourceProtAddress)
			resSrc, err := pseudonymizeIP(ipAddressSrc)
			if err != nil {
				returnWithError(err)
				return
			}
			if resSrc.To4() != nil {
				arpPacket.SourceProtAddress = resSrc.To4()
			} else {
				arpPacket.SourceProtAddress = resSrc
			}
			ipAddressDst := net.IP(arpPacket.DstProtAddress)
			resDst, err := pseudonymizeIP(ipAddressDst)
			if err != nil {
				returnWithError(err)
				return
			}
			if resDst.To4() != nil {
				arpPacket.DstProtAddress = resDst.To4()
			} else {
				arpPacket.DstProtAddress = resDst
			}
			if err := gopacket.SerializeLayers(buf, opts, ethernet, arpPacket, gopacket.Payload(arpPacket.Payload)); err != nil {
				returnWithError(err)
				return
			}
			bytes := buf.Bytes()
			nci := gopacket.CaptureInfo{
				Timestamp:      ci.Timestamp,
				CaptureLength:  len(bytes),
				Length:         len(bytes),
				InterfaceIndex: ci.InterfaceIndex,
				AncillaryData:  ci.AncillaryData,
			}
			if err := packetWriter.WritePacket(nci, bytes); err != nil {
				returnWithError(err)
				return
			}
		} else if ipv4 != nil {
			ipv4Packet := ipv4.(*layers.IPv4)
			resSrc, err := pseudonymizeIP(ipv4Packet.SrcIP)
			if err != nil {
				returnWithError(err)
				return
			}
			resDst, err := pseudonymizeIP(ipv4Packet.DstIP)
			if err != nil {
				returnWithError(err)
				return
			}
			ipv4Packet.SrcIP = resSrc.To4()
			ipv4Packet.DstIP = resDst.To4()
			if err := gopacket.SerializeLayers(buf, opts, ethernet, ipv4Packet, gopacket.Payload(ipv4Packet.Payload)); err != nil {
				returnWithError(err)
				return
			}
			bytes := buf.Bytes()
			nci := gopacket.CaptureInfo{
				Timestamp:      ci.Timestamp,
				CaptureLength:  len(bytes),
				Length:         len(bytes),
				InterfaceIndex: ci.InterfaceIndex,
				AncillaryData:  ci.AncillaryData,
			}
			if err := packetWriter.WritePacket(nci, bytes); err != nil {
				returnWithError(err)
				return
			}

		} else if ipv6 != nil {
			ipv6Packet := ipv6.(*layers.IPv6)
			resSrc, err := pseudonymizeIP(ipv6Packet.SrcIP)
			if err != nil {
				returnWithError(err)
				return
			}
			resDst, err := pseudonymizeIP(ipv6Packet.DstIP)
			if err != nil {
				returnWithError(err)
				return
			}
			ipv6Packet.SrcIP = resSrc
			ipv6Packet.DstIP = resDst
			if err := gopacket.SerializeLayers(buf, opts, ethernet, ipv6Packet, gopacket.Payload(ipv6Packet.Payload)); err != nil {
				returnWithError(err)
				return
			}
			bytes := buf.Bytes()
			nci := gopacket.CaptureInfo{
				Timestamp:      ci.Timestamp,
				CaptureLength:  len(bytes),
				Length:         len(bytes),
				InterfaceIndex: ci.InterfaceIndex,
				AncillaryData:  ci.AncillaryData,
			}
			if err := packetWriter.WritePacket(nci, bytes); err != nil {
				returnWithError(err)
				return
			}
		} else {
			if err := packetWriter.WritePacket(ci, data); err != nil {
				returnWithError(err)
				return
			}
		}
	}
	c.Writer.Flush()
}

func Protect(c *gin.Context) {

	params, pseudonymizer := parse(c)

	if params == nil {
		return
	}

	process(c, params, pseudonymizer.Pseudonymize)
}

func Unprotect(c *gin.Context) {

	params, pseudonymizer := parse(c)

	if params == nil {
		return
	}

	process(c, params, pseudonymizer.Depseudonymize)

}

func queryToConfig(query url.Values) map[string]interface{} {
	config := make(map[string]interface{})
	for key, value := range query {
		if len(value) > 1 {
			config[key] = value
		} else if len(value) == 0 || len(value[0]) == 0 {
			config[key] = nil
		} else {
			config[key] = value[0]
		}
	}
	return config
}

func parse(c *gin.Context) (map[string]interface{}, pseudonymize.Pseudonymizer) {

	query := c.Request.URL.Query()

	params, err := PseudonymizeForm.Validate(queryToConfig(query))

	if err != nil {
		api.HandleError(c, 400, err)
		return nil, nil
	}

	psMaker, ok := pseudonymize.Pseudonymizers["structured"]

	if !ok {
		c.JSON(400, map[string]string{"message": "invalid pseudonymization method"})
		return nil, nil
	}

	pseudonymizer, err := psMaker(map[string]interface{}{
		"type":             "ip",
		"prefixpreserving": params["preserve-subnets"],
	})

	if err != nil {
		api.HandleError(c, 400, err)
		return nil, nil
	}

	key := []byte(params["key"].(string))

	if err := pseudonymizer.GenerateParams(key, nil); err != nil {
		api.HandleError(c, 500, err)
		return nil, nil
	}

	return params, pseudonymizer

}
