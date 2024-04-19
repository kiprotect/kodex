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

func makeTestInteger(b *testing.B) CompositeType {
	integer, err := MakeInteger(map[string]interface{}{"min": int64(0), "max": int64(10000)})
	if err != nil {
		b.Fatal(err)
	}
	if err := integer.Unmarshal("", 4000); err != nil {
		b.Fatal(err)
	}
	return integer
}

func BenchmarkIntegerPseudonymization(b *testing.B) {
	pseudonymizationBenchmark(b, makeTestInteger(b), PSH)
}

func BenchmarkIntegerDepseudonymization(b *testing.B) {
	depseudonymizationBenchmark(b, makeTestInteger(b), PSH, DPSH)
}

func BenchmarkBidirectionalIntegerPseudonymization(b *testing.B) {
	pseudonymizationBenchmark(b, makeTestInteger(b), PS)
}

func BenchmarkBidirectionalIntegerDepseudonymization(b *testing.B) {
	depseudonymizationBenchmark(b, makeTestInteger(b), PS, DPS)
}
