// Kodex (Community Edition - CE) - Privacy & Security Engineering Platform
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
	"testing"
)

func TestSimplePseudonymization(t *testing.T) {

	ct, err := MakeCompositeTestType(12)
	if err != nil {
		t.Fatal(err)
	}

	key := []byte("foobar")
	res, err := PSH(ct, key)
	if err != nil {
		t.Fatal(err)
	}

	cn, err := DPSH(res, key)

	if err != nil {
		t.Fatal(err)
	}

	if !cn.Equals(ct) {
		t.Fatalf("Should be equal")
	}

	cn, err = DPSH(res, []byte("wrong-key"))

	if err != nil {
		t.Fatal(err)
	}

	if cn.Equals(ct) {
		t.Fatalf("Should be different")
	}
}

func TestSimpleBidirectionalPseudonymization(t *testing.T) {

	ct, err := MakeCompositeTestType(6)
	if err != nil {
		t.Fatal(err)
	}

	key := []byte("foobar")
	res, err := PS(ct, key)
	if err != nil {
		t.Fatal(err)
	}

	cn, err := DPS(res, key)

	if err != nil {
		t.Fatal(err)
	}

	if !cn.Equals(ct) {
		t.Fatalf("Should be equal")
	}

	cn, err = DPS(res, []byte("wrong-key"))

	if err != nil {
		t.Fatal(err)
	}

	if cn.Equals(ct) {
		t.Fatalf("Should be different")
	}
}
