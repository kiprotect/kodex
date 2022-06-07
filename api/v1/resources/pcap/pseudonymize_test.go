// Kodex (Community Edition - CE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2022  KIProtect GmbH (HRB 208395B) - Germany
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package pcap

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/kiprotect/kodex/api/decorators"
	"github.com/kiprotect/kodex/api/v1/resources/pcap/pcapgo"
	"github.com/kiprotect/kodex/helpers"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func getRouter() (*gin.Engine, error) {
	paths, fS, err := helpers.SettingsPaths()

	if err != nil {
		return nil, err
	}

	settings, err := helpers.Settings(paths, fS)

	if err != nil {
		return nil, err
	}

	router := gin.New()
	router.Use(decorators.WithSettings(settings))

	router.POST("/protect", Protect)
	router.POST("/unprotect", Unprotect)

	return router, nil
}

type TestStruct struct {
	Input        []byte
	URLParams    map[string]string
	ResponseCode int
}

func equalPCAPs(source []byte, destination []byte) (bool, error) {

	readerIn, err := pcapgo.NewReader(bytes.NewReader(source))

	if err != nil {
		return false, err
	}

	readerOut, err := pcapgo.NewReader(bytes.NewReader(destination))

	if err != nil {
		return false, err
	}

	for {
		dataIn, _, errIn := readerIn.ReadPacketData()
		if errIn != nil && errIn != io.EOF {
			return false, errIn
		}
		dataOut, _, errOut := readerOut.ReadPacketData()
		if errOut != nil && errOut != io.EOF {
			return false, errOut
		}

		if errIn == io.EOF {
			if errOut == io.EOF {
				return true, nil
			}
			return false, nil
		}

		if dataIn == nil && dataOut != nil || dataIn != nil && dataOut == nil {
			return false, nil
		}

		packetIn := gopacket.NewPacket(dataIn, layers.LayerTypeEthernet, gopacket.DecodeOptions{Lazy: true, NoCopy: true})
		packetOut := gopacket.NewPacket(dataOut, layers.LayerTypeEthernet, gopacket.DecodeOptions{Lazy: true, NoCopy: true})

		ipv4In := packetIn.Layer(layers.LayerTypeIPv4)
		ipv6In := packetIn.Layer(layers.LayerTypeIPv6)
		arpIn := packetIn.Layer(layers.LayerTypeARP)
		ipv4Out := packetOut.Layer(layers.LayerTypeIPv4)
		ipv6Out := packetOut.Layer(layers.LayerTypeIPv6)
		arpOut := packetOut.Layer(layers.LayerTypeARP)

		if ipv4In != nil {
			if ipv4Out == nil {
				return false, nil
			}
			if !bytes.Equal(ipv4In.(*layers.IPv4).LayerContents(), ipv4Out.(*layers.IPv4).LayerContents()) {
				return false, nil
			}
		} else if ipv6In != nil {
			if ipv6Out == nil {
				return false, nil
			}
			if !bytes.Equal(ipv6In.(*layers.IPv6).LayerContents(), ipv6Out.(*layers.IPv6).LayerContents()) {
				return false, nil
			}
		} else if arpIn != nil {
			if arpOut == nil {
				return false, nil
			}
			if !bytes.Equal(arpIn.(*layers.ARP).SourceProtAddress, arpOut.(*layers.ARP).SourceProtAddress) {
				return false, nil
			}
		}
	}
	return true, nil
}

func TestPseudonymize(t *testing.T) {

	router, err := getRouter()

	if err != nil {
		t.Fatal(err)
	}

	tests := []TestStruct{
		TestStruct{
			Input:        smallPCAP,
			ResponseCode: 200,
			URLParams: map[string]string{
				"format": "pcap",
			},
		},
		TestStruct{
			Input:        largePCAP,
			ResponseCode: 200,
			URLParams: map[string]string{
				"format": "pcap",
			},
		},
	}

	key1 := "foozball"
	key2 := "fuuzball"

	for _, test := range tests {
		reader := bytes.NewReader(test.Input)
		req, _ := http.NewRequest("POST", "/protect", reader)
		q := req.URL.Query()

		for key, value := range test.URLParams {
			q.Add(key, value)
		}

		q.Add("key", key1)

		req.URL.RawQuery = q.Encode()
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if resp.Code != test.ResponseCode {
			t.Error(fmt.Errorf("Response code should be %d but is %d", test.ResponseCode, resp.Code))
			continue
		}

		if test.ResponseCode != 200 {
			continue
		}

		q.Del("key")
		q.Add("key", key2)
		req.URL.RawQuery = q.Encode()
		reader.Seek(0, 0)

		resp2 := httptest.NewRecorder()
		router.ServeHTTP(resp2, req)

		if n, err := equalPCAPs(resp2.Body.Bytes(), resp.Body.Bytes()); err != nil {
			t.Fatal(err)
		} else if n {
			t.Error("different keys should not be equal")
			continue
		}

		// we make sure the destination is not identical to the source
		if n, err := equalPCAPs(test.Input, resp.Body.Bytes()); err != nil {
			t.Fatal(err)
		} else if n {
			t.Error("should not be equal")
			continue
		}

		reader = bytes.NewReader(resp.Body.Bytes())

		req, _ = http.NewRequest("POST", "/unprotect", reader)
		q = req.URL.Query()

		for key, value := range test.URLParams {
			q.Add(key, value)
		}

		q.Add("key", key1)

		req.URL.RawQuery = q.Encode()
		resp = httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if resp.Code != 200 {
			t.Error(fmt.Errorf("Depseudonymization response code should be 200 but is %d", resp.Code))
			continue
		}

		if n, _ := equalPCAPs(resp.Body.Bytes(), test.Input); err != nil {
			t.Fatal(err)
		} else if !n {
			t.Error("should be equal")
			continue
		}

	}

}
