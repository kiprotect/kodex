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

package merengue

import (
	"testing"
)

func BenchmarkPseudonymize(b *testing.B) {
	ba := makeTestValue()
	keyA := []byte("key A")

	b.SetBytes(int64(len(ba)))

	for n := 0; n < b.N; n++ {
		Pseudonymize(ba, uint(len(ba)*8-4), keyA, Sha256)
	}
}

func BenchmarkDepseudonymize(b *testing.B) {
	ba := makeTestValue()
	keyA := []byte("key A")

	b.SetBytes(int64(len(ba)))

	res := Pseudonymize(ba, uint(len(ba)*8-4), keyA, Sha256)
	for n := 0; n < b.N; n++ {
		Depseudonymize(res, uint(len(ba)*8-4), keyA, Sha256)
	}
}

func BenchmarkBidirectionalPseudonymize(b *testing.B) {
	ba := makeTestValue()
	keyA := []byte("key A")
	keyB := []byte("key B")

	b.SetBytes(int64(len(ba)))

	for n := 0; n < b.N; n++ {
		PseudonymizeBidirectional(ba, uint(len(ba)*8-4), keyA, keyB, Sha256)
	}
}

func BenchmarkBidirectionalDepseudonymize(b *testing.B) {
	ba := makeTestValue()
	keyA := []byte("key A")
	keyB := []byte("key B")

	b.SetBytes(int64(len(ba)))

	res := PseudonymizeBidirectional(ba, uint(len(ba)*8-4), keyA, keyB, Sha256)
	for n := 0; n < b.N; n++ {
		DepseudonymizeBidirectional(res, uint(len(ba)*8-4), keyA, keyB, Sha256)
	}
}

func makeTestValue() []byte {
	ba := []byte("foo bar")
	ba[len(ba)-1] &= 0xFF >> 4
	return ba
}
