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
	"testing"
)

func makeTestIP(b *testing.B) CompositeType {
	ip, err := MakeIPAddr(nil)
	if err != nil {
		b.Fatal(err)
	}
	if err = ip.Unmarshal("", "123.121.21.1/25"); err != nil {
		b.Fatal(err)
	}
	return ip
}

func BenchmarkIPPseudonymization(b *testing.B) {
	pseudonymizationBenchmark(b, makeTestIP(b), PSH)
}

func BenchmarkIPDepseudonymization(b *testing.B) {
	depseudonymizationBenchmark(b, makeTestIP(b), PSH, DPSH)
}

func BenchmarkBidirectionalIPPseudonymization(b *testing.B) {
	pseudonymizationBenchmark(b, makeTestIP(b), PS)
}

func BenchmarkBidirectionalIPDepseudonymization(b *testing.B) {
	depseudonymizationBenchmark(b, makeTestIP(b), PS, DPS)
}
