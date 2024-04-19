// Kodex (Community Edition - CE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2024  KIProtect GmbH (HRB 208395B) - Germany
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

func pseudonymizationBenchmark(b *testing.B, c CompositeType, ps func(CompositeType, []byte) (CompositeType, error)) {

	key := []byte("foobar")

	ba, err := c.Encode()
	if err != nil {
		b.Fatal(err)
	}

	b.SetBytes(int64(ba.Length() / 8))

	for n := 0; n < b.N; n++ {
		_, err := ps(c, key)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func depseudonymizationBenchmark(b *testing.B, c CompositeType, ps func(CompositeType, []byte) (CompositeType, error), dps func(CompositeType, []byte) (CompositeType, error)) {

	key := []byte("foobar")

	ba, err := c.Encode()
	if err != nil {
		b.Fatal(err)
	}

	b.SetBytes(int64(ba.Length() / 8))

	res, err := ps(c, key)

	if err != nil {
		b.Fatal(err)
	}

	for n := 0; n < b.N; n++ {
		_, err := dps(res, key)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func makeStructuredTestType(b *testing.B) CompositeType {
	ct, err := MakeCompositeTestType(4)
	if err != nil {
		b.Fatal(err)
	}
	return ct
}

func BenchmarkStructuredPseudonymization(b *testing.B) {
	pseudonymizationBenchmark(b, makeStructuredTestType(b), PSH)
}

func BenchmarkStructuredDepseudonymization(b *testing.B) {
	depseudonymizationBenchmark(b, makeStructuredTestType(b), PSH, DPSH)
}

func BenchmarkBidirectionalStructuredPseudonymization(b *testing.B) {
	pseudonymizationBenchmark(b, makeStructuredTestType(b), PS)
}

func BenchmarkBidirectionalStructuredDepseudonymization(b *testing.B) {
	depseudonymizationBenchmark(b, makeStructuredTestType(b), PS, DPS)
}
