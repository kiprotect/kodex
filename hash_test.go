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

package kiprotect

import (
	"bytes"
	"testing"
)

func TestBasicHash(t *testing.T) {
	s := map[string]interface{}{
		"foo":   "bar",
		"zoo":   "db",
		"value": "another",
		"fooz":  []byte{10, 32, 111, 54, 63},
	}
	h1, err := StructuredHash(s)
	if err != nil {
		t.Fatal(err)
	}
	s["test"] = "another"
	h2, err := StructuredHash(s)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Equal(h1, h2) {
		t.Errorf("Hashes should be different")
	}
	s["substruct"] = map[string]interface{}{
		"foo":  "barbara",
		"test": "another one",
	}
	h3, err := StructuredHash(s)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Equal(h2, h3) {
		t.Errorf("Hashes should be different")
	}
	l := []string{"foo", "bar", "baz"}
	hl1, err := StructuredHash(l)
	if err != nil {
		t.Fatal(err)
	}
	l = []string{"bar", "foo", "baz"}
	hl2, err := StructuredHash(l)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Equal(hl1, hl2) {
		t.Errorf("Hashes should be different")
	}
}
