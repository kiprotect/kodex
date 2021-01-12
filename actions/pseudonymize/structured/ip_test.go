// Kodex (Community Edition - CE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2021  KIProtect GmbH (HRB 208395B) - Germany
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
	"testing"
)

func TestIPEncoding(t *testing.T) {
	ip := MakeIPAddress(net.ParseIP("127.0.0.1").To4(), 32)

	encoded, err := ip.Encode()
	if err != nil {
		t.Fatal(err)
	}

	ipNew := MakeIPAddress(net.ParseIP("127.0.1.1").To4(), 32)
	err = ipNew.Decode(encoded)
	if err != nil {
		t.Fatal(err)
	}

	if !ipNew.Equals(ip) {
		t.Errorf("Should be equal")
	}
}

func TestIPValididityCheck(t *testing.T) {
	ip, err := MakeIPAddr(nil)
	if err != nil {
		t.Fatal(err)
	}

	if err = ip.Unmarshal("", "123.121.21.1/24"); err == nil {
		t.Fatal(fmt.Errorf("Non-zero bits after netmask not detected"))
	}

	if err = ip.Unmarshal("", "123.121.21.7/27"); err != nil {
		t.Fatal(fmt.Errorf("Should not throw an error"))
	}

	if err = ip.Unmarshal("", "123.121.21.0/24"); err != nil {
		t.Fatal(err)
	}

	valid := ip.IsValid()
	if valid[0] == false {
		t.Fatal("IP should be valid")
	}

	marshalled, err := ip.Marshal("")
	if err != nil {
		t.Fatal(err)
	}

	if marshalled != "123.121.21.0/24" {
		fmt.Println(marshalled)
		t.Fatal("Issue with marshalling and unmarshalling the data")
	}

	err = ip.Unmarshal("", "123.121.21.12/128")
	if err == nil {
		t.Fatal("Bad subnet, should raise an error on Unmarshal.")
	}

	err = ip.Unmarshal("", "2001:db8:abcd:0012::0/64")

	valid = ip.IsValid()
	if valid[0] == false {
		t.Fatal("IPv6 should be valid")
	}
}
