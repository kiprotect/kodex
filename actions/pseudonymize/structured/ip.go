// KIProtect (Community Edition - CE) - Privacy & Security Engineering Platform
// Copyright (C) 2020  KIProtect GmbH (HRB 208395B) - Germany
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

package structured

import (
	"fmt"
	"net"
)

type IPAddr struct {
	mask    uint
	hasMask bool
	v6      bool
	CompositeListType
}

func MakeIPAddr(params interface{}) (CompositeType, error) {
	return &IPAddr{}, nil
}

func (ip *IPAddr) Copy() CompositeType {
	listCopy := ip.CompositeListType.Copy()
	listCopyType, _ := listCopy.(*CompositeListType)
	return &IPAddr{
		mask:              ip.mask,
		hasMask:           ip.hasMask,
		CompositeListType: *listCopyType,
	}
}

func (ip *IPAddr) Unmarshal(format string, data interface{}) error {
	ipStr, ok := data.(string)

	if !ok {
		byteArray, ok := data.([]byte)
		if ok {
			ipStr = string(byteArray)
		} else {
			return fmt.Errorf("expected a string or byte array as input")
		}
	}

	ipAddr, ipNet, err := net.ParseCIDR(ipStr)

	if err != nil {
		ipAddr = net.ParseIP(ipStr)
		if ipAddr == nil {
			return fmt.Errorf("No valid CIDR or IP address found! %s", ipStr)
		}
		ip.hasMask = false
		if ipAddr.To4() != nil {
			ip.mask = 32
		} else {
			ip.v6 = true
			ip.mask = 128
		}
	} else {
		ip.hasMask = true
		mask, _ := ipNet.Mask.Size()
		ip.mask = uint(mask)
	}

	if ipAddr.To4() != nil {
		ipAddr = ipAddr.To4()
	} else if ipAddr.To16() != nil {
		ip.v6 = true
		ipAddr = ipAddr.To16()
	}

	nonZeroErr := fmt.Errorf("nonzero bits found beyond the indicated netmask")

	m := ip.mask / 8
	k := ip.mask % 8
	if k != 0 {
		if ipAddr[m]^(ipAddr[m]&(0xFF>>(8-k))) != 0 {
			return nonZeroErr
		}
		m += 1
	}
	for n := m; n < uint(len(ipAddr)); n++ {
		if ipAddr[n] != 0 {
			return nonZeroErr
		}
	}

	elements := make([]Type, 0)
	ipElement := MakeIPAddress(ipAddr[:m], ip.mask)
	elements = append(elements, ipElement)
	ip.SetSubtypes(elements)
	return nil
}

func (ip *IPAddr) Marshal(format string) (interface{}, error) {

	ipAddress, ok := ip.subtypes[0].(IPElementIf)
	if !ok {
		return nil, fmt.Errorf("No valid IP address found! %s", ipAddress)
	}
	value := ipAddress.Value()
	var fullValue []byte
	if len(value) == 16 {
		fullValue = make([]byte, 16)
	} else {
		fullValue = make([]byte, 4)
	}
	copy(fullValue, value)
	if ip.hasMask {
		ipNet := net.IPNet{
			IP:   net.IP(fullValue),
			Mask: net.CIDRMask(int(ip.mask), len(fullValue)*8),
		}
		return ipNet.String(), nil
	}
	return net.IP(fullValue).String(), nil
}
