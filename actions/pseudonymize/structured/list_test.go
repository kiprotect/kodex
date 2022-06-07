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
	"encoding/binary"
	"fmt"
	"math"
	// we only use this for testing and are aware math/rand is not a secure PRNG
	"math/rand"
	"testing"
	"time"
)

type RangedBitField struct {
	max    int64
	min    int64
	value  int64
	by     uint
	length uint
}

func (r *RangedBitField) Length() uint {
	return r.length
}

func (r *RangedBitField) Copy() Type {
	return &RangedBitField{
		max:    r.max,
		min:    r.min,
		value:  r.value,
		by:     r.by,
		length: r.length,
	}
}

func (r *RangedBitField) Equals(t Type) bool {
	rb, ok := t.(*RangedBitField)
	if !ok {
		return false
	}
	if r.max != rb.max || r.min != rb.min || r.value != rb.value {
		return false
	}
	return true
}

func (r *RangedBitField) Encode() (*BitArray, error) {
	ev := uint64(r.value - r.min)
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, ev)
	ba, err := MakeBitArrayFromBytes(b[:r.by], r.length)
	return ba, err
}

func (r *RangedBitField) Decode(b *BitArray) error {
	if b.Length() != r.Length() {
		return fmt.Errorf("invalid length")
	}
	bytes := b.Bytes()
	ba := make([]byte, 8)
	copy(ba[:r.by], bytes)
	ev := binary.LittleEndian.Uint64(ba)
	value := int64(ev) + r.min
	r.value = value
	return nil
}

func (r *RangedBitField) IsValid() bool {
	if r.value >= r.min && r.value <= r.max {
		return true
	}
	return false
}

func MakeCompositeTestType(ll int) (CompositeType, error) {
	var err error
	subtypes := make([]Type, ll)
	for i := 0; i < ll; i++ {
		min := int64(1 << uint(i*2))
		max := int64(1 << uint(i*2+1))
		value := min + rand.Int63n(max-min)
		if subtypes[i], err = MakeRangedBitField(min, max, value); err != nil {
			return nil, err
		}
	}
	ct := MakeCompositeListType()
	ctl := ct.(*CompositeListType)
	ctl.SetSubtypes(subtypes)
	return ct, nil
}

func MakeRangedBitField(min, max, value int64) (*RangedBitField, error) {
	length := uint(math.Ceil(math.Log2(float64(max - min + 1))))
	by := length / 8
	if length%8 != 0 {
		by += 1
	}
	if value < min || value > max {
		return nil, fmt.Errorf("out of bounds")
	}
	return &RangedBitField{
		max:    max,
		min:    min,
		value:  value,
		length: length,
		by:     by,
	}, nil
}

func TestCompositeType(t *testing.T) {

	rand.Seed(time.Now().UTC().UnixNano())

	_, err := MakeRangedBitField(0, 2000, 1400)
	if err != nil {
		t.Fatal(err)
	}

	ll := rand.Intn(20)

	ct, err := MakeCompositeTestType(ll)
	if err != nil {
		t.Fatal(err)
	}

	ba, err := ct.Encode()
	if err != nil {
		t.Fatal(err)
	}

	ct2, err := MakeCompositeTestType(ll)

	if err != nil {
		t.Fatal(err)
	}

	err = ct2.Decode(ba)
	if err != nil {
		t.Fatal(err)
	}

	if !ct.Equals(ct2) {
		t.Fatalf("values should be equal")
	}
}
