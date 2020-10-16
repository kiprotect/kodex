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
	"fmt"
	"testing"
)

func TestIntegerMarshalling(t *testing.T) {
	i, err := MakeInteger(map[string]interface{}{"min": int64(0), "max": int64(100000)})
	if err != nil {
		t.Fatal(err)
	}
	err = i.Unmarshal("", 34324)
	if err != nil {
		t.Fatal(err)
	}
	ii, err := i.Marshal("")
	if err != nil {
		t.Fatal(err)
	}
	if ii != int64(34324) {
		t.Fatalf("Not equal")
	}
}

func TestValidRangePseudonymization(t *testing.T) {
	max := int64(4000)
	min := int64(1000)
	i, err := MakeInteger(map[string]interface{}{"min": min, "max": max})
	if err != nil {
		t.Fatal(err)
	}

	key := []byte("foobar")

	m := make(map[int64]bool)

	for n := min; n <= max; n++ {

		err = i.Unmarshal("", n)
		if err != nil {
			t.Fatal(err)
		}

		res, err := PSH(i, key)
		if err != nil {
			t.Fatal(err)
		}

		iRes, ok := res.(*Integer)

		if !ok {
			t.Fatal(fmt.Errorf("should be an Integer"))
		}

		ni, err := iRes.Marshal("")

		nii, ok := ni.(int64)

		// no destination value should be produced more than once
		if _, ok := m[nii]; ok {
			t.Fatalf("Mapped more than once: %d", nii)
		}

		m[nii] = true

		if !ok {
			t.Fatal(fmt.Errorf("should be an int64"))
		}

		if err != nil {
			t.Fatal(err)
		}

		if nii < min || nii > max {
			t.Fatal(fmt.Errorf("result out of bounds"))
		}

		ci, err := DPSH(res, key)

		if err != nil {
			t.Fatal(err)
		}

		if !ci.Equals(i) {
			t.Fatal(fmt.Errorf("should be equal"))
		}
	}
	// each possible value should be produced exactly once
	for n := min; n <= max; n++ {
		if _, ok := m[n]; !ok {
			t.Fatalf("Not mapped: %d", n)
		}
	}

}

func TestIntegerPseudonymization(t *testing.T) {

	i, err := MakeInteger(map[string]interface{}{"min": int64(400), "max": int64(10000)})
	if err != nil {
		t.Fatal(err)
	}

	n := int64(1223)

	err = i.Unmarshal("", n)
	if err != nil {
		t.Fatal(err)
	}

	key := []byte("foobar")
	res, err := PSH(i, key)
	if err != nil {
		t.Fatal(err)
	}

	resi, ok := res.(*Integer)
	if !ok {
		t.Fatal(fmt.Errorf("should be an Integer"))
	}

	ri, err := resi.Marshal("")

	if ri == n {
		t.Fatal(fmt.Errorf("should (normally) be different"))
	}

	ci, err := DPSH(res, key)

	if err != nil {
		t.Fatal(err)
	}

	cii, ok := ci.(*Integer)

	if !ok {
		t.Fatal(fmt.Errorf("should be an Integer"))
	}

	ni, err := cii.Marshal("")

	if err != nil {
		t.Fatal(err)
	}

	if ni != n {
		t.Fatal(err)
	}

	if !ci.Equals(i) {
		t.Fatalf("Should be equal")
	}

	ci, err = DPSH(res, []byte("wrong-key"))

	if err != nil {
		t.Fatal(err)
	}

	if ci.Equals(i) {
		t.Fatalf("Should be different")
	}
}
