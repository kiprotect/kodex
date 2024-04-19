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
	"crypto/rand"
	"fmt"
	// we only use this for testing and are aware math/rand is not a secure PRNG
	mrand "math/rand"
	"testing"
)

func TestBitArrayBasics(t *testing.T) {
	for j := 0; j < 10; j++ {
		bytes := make([]byte, 3)
		_, err := rand.Read(bytes)
		if err != nil {
			t.Fatal(err)
		}
		ba2, err := MakeBitArrayFromBytes(bytes, 19)
		if err != nil {
			t.Fatal(err)
		}
		for i := uint(0); i < 20; i++ {
			ba1 := MakeBitArray(40)
			balo, err := ba1.Extract(0, i)
			if err != nil {
				t.Fatal(err)
			}
			baro, err := ba1.Extract(i+ba2.Length(), ba1.Length()-i-ba2.Length())
			if err != nil {
				t.Fatal(err)
			}
			err = ba1.Update(ba2, i)
			if err != nil {
				t.Fatal(err)
			}
			bae, err := ba1.Extract(i, ba2.length)
			if err != nil {
				t.Fatal(err)
			}
			if !bae.Equals(ba2) {
				t.Errorf("bit arrays should be equal")
			}
			bal, err := ba1.Extract(0, i)
			if err != nil {
				t.Fatal(err)
			}
			bar, err := ba1.Extract(i+ba2.Length(), ba1.Length()-i-ba2.Length())
			if err != nil {
				t.Fatal(err)
			}
			if !bal.Equals(balo) {
				t.Errorf("left part should be equal")
			}
			if !bar.Equals(baro) {
				t.Errorf("right part should be equal")
			}
		}
	}
}

func TestBitArrayRandomExtractAndUpdate(t *testing.T) {
	ba1 := MakeBitArray(100)
	for j := 0; j < 1000; j++ {
		l := uint(mrand.Intn(60) + 1)
		bl := l / 8
		if l%8 != 0 {
			bl += 1
		}
		bytes := make([]byte, bl)
		_, err := rand.Read(bytes)
		if err != nil {
			t.Fatal(err)
		}
		ba2, err := MakeBitArrayFromBytes(bytes, l)
		if err != nil {
			t.Fatal(err)
		}
		off := uint(mrand.Intn(30))
		bao := ba1.Copy()
		err = ba1.Update(ba2, off)
		if err != nil {
			t.Fatal(err)
		}
		bae, err := ba1.Extract(off, ba2.length)
		if err != nil {
			t.Fatal(err)
		}
		if !bae.Equals(ba2) {
			fmt.Println(off, ba1.AsString())
			fmt.Println(ba2.AsString())
			fmt.Println(bae.AsString())
			t.Errorf("bit arrays should be equal")
		}
		bal, err := ba1.Extract(0, off)
		if err != nil {
			t.Fatal(err)
		}
		balo, err := bao.Extract(0, off)
		if err != nil {
			t.Fatal(err)
		}
		if !balo.Equals(bal) {
			t.Fatal("left parts should be equal")
		}

		bar, err := ba1.Extract(off+ba2.Length(), ba1.Length()-off-ba2.Length())
		if err != nil {
			t.Fatal(err)
		}
		baro, err := bao.Extract(off+ba2.Length(), ba1.Length()-off-ba2.Length())
		if err != nil {
			t.Fatal(err)
		}
		if !baro.Equals(bar) {
			t.Fatal("right parts should be equal")
		}

	}
}
