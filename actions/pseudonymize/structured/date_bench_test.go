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

package structured

import (
	"testing"
)

func makeTestDate(b *testing.B) CompositeType {
	format := "%(2000-2100)Y-%m-%d"
	date, err := MakeDate(nil)
	if err != nil {
		b.Fatal(err)
	}
	if err := date.Unmarshal(format, "2008-01-10"); err != nil {
		b.Fatal(err)
	}
	return date
}

func BenchmarkDatePseudonymization(b *testing.B) {
	pseudonymizationBenchmark(b, makeTestDate(b), PSH)
}

func BenchmarkDateDepseudonymization(b *testing.B) {
	depseudonymizationBenchmark(b, makeTestDate(b), PSH, DPSH)
}

func BenchmarkBidirectionalDatePseudonymization(b *testing.B) {
	pseudonymizationBenchmark(b, makeTestDate(b), PS)
}

func BenchmarkBidirectionalDateDepseudonymization(b *testing.B) {
	depseudonymizationBenchmark(b, makeTestDate(b), PS, DPS)
}
