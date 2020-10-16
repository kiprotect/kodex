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

package merengue

import (
	"bytes"
	"testing"
)

func TestPseudonymize(t *testing.T) {
	ba := []byte("foo bar")
	keyA := []byte("key A")
	keyB := []byte("key B")

	ba[len(ba)-1] &= 0xFF >> 4

	res := Pseudonymize(ba, uint(len(ba)*8-4), keyA, Sha256)
	dres := Depseudonymize(res, uint(len(ba)*8-4), keyA, Sha256)
	if !bytes.Equal(ba, dres) {
		t.Errorf("should be identical")
	}

	res2 := PseudonymizeBidirectional(ba, uint(len(ba)*8)-4, keyA, keyB, Sha256)
	dres2 := DepseudonymizeBidirectional(res2, uint(len(ba)*8-4), keyA, keyB, Sha256)
	if !bytes.Equal(ba, dres2) {
		t.Errorf("bidirectional result should be identical")
	}

}
