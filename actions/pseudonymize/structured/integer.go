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
	"github.com/kiprotect/go-helpers/maps"
	"math"
)

type Integer struct {
	max int64
	min int64
	CompositeListType
}

func MakeInteger(params interface{}) (CompositeType, error) {
	mv := int64(^uint64(0)>>1) - 1
	i := &Integer{
		min: 0,
		max: mv,
	}
	if params != nil {
		paramsMap, ok := maps.ToStringMap(params)
		if !ok {
			return nil, fmt.Errorf("Expected a map as parameters")
		}
		if paramsMap["min"] != nil {
			minValue, ok := paramsMap["min"].(int64)
			if !ok {
				return nil, fmt.Errorf("min should be of type int64")
			} else {
				i.min = minValue
			}
			maxValue, ok := paramsMap["max"].(int64)
			if !ok {
				return nil, fmt.Errorf("max should be of type int64")
			} else {
				i.max = int64(maxValue)
			}
		}
		if i.min < -mv || i.max > mv {
			return nil, fmt.Errorf("max/min out of bounds")
		}
	}
	return i, nil
}

func (d *Integer) Copy() CompositeType {
	listCopy := d.CompositeListType.Copy()
	listCopyType, _ := listCopy.(*CompositeListType)
	return &Integer{
		min:               d.min,
		max:               d.max,
		CompositeListType: *listCopyType,
	}
}

func (r *Integer) Marshal(format string) (interface{}, error) {
	if format != "" {
		return nil, fmt.Errorf("unsupported format: '%s'", format)
	}
	s, err := r.Get(0)
	if err != nil {
		return nil, err
	}
	i, ok := s.(*IntegerField)
	if !ok {
		return nil, fmt.Errorf("not an integer field")
	}
	return i.Value, nil
}

func (r *Integer) Unmarshal(format string, value interface{}) error {
	if format != "" {
		return fmt.Errorf("unsupported format: '%s'", format)
	}
	var err error
	subtypes := make([]Type, 1)
	var i64 int64
	intValue, ok := value.(int)
	if !ok {
		i64, ok = value.(int64)
		if !ok {
			// Golang will serialize JSON numbers into float64 values
			floatValue, ok := value.(float64)
			if !ok {
				return fmt.Errorf("expected an int64 value")
			}
			i64 = int64(floatValue)
		}
	} else {
		i64 = int64(intValue)
	}
	if subtypes[0], err = MakeIntegerField(r.min, r.max, i64); err != nil {
		return err
	}
	r.SetSubtypes(subtypes)
	return nil
}

type IntegerField struct {
	Max    int64
	Min    int64
	Value  int64
	by     uint
	length uint
}

func (r *IntegerField) Length() uint {
	return r.length
}

func (r *IntegerField) Copy() Type {
	return &IntegerField{
		Max:    r.Max,
		Min:    r.Min,
		Value:  r.Value,
		by:     r.by,
		length: r.length,
	}
}

func (r *IntegerField) Equals(t Type) bool {
	rb, ok := t.(*IntegerField)
	if !ok {
		return false
	}
	if r.Max != rb.Max || r.Min != rb.Min || r.Value != rb.Value {
		return false
	}
	return true
}

func Revert(inp []byte, l uint) []byte {
	bo := make([]byte, len(inp))
	var off, ind, indr, offr, g uint32
	ul := uint32(l)
	for g = 0; g < ul; g++ {
		ind = g / 8
		indr = (ul - g - 1) / 8
		offr = (ul - g - 1) % 8
		off = g % 8
		bo[ind] |= ((inp[indr] >> offr) & 1) << off
	}
	return bo
}

func (r *IntegerField) Encode() (*BitArray, error) {
	ev := uint64(r.Value - r.Min)
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, ev)
	b = Revert(b, r.length)
	ba, err := MakeBitArrayFromBytes(b[:r.by], r.length)
	return ba, err
}

func (r *IntegerField) Decode(b *BitArray) error {
	if b.Length() != r.Length() {
		return fmt.Errorf("invalid length")
	}
	bytes := b.Bytes()
	bytes = Revert(bytes, r.length)
	ba := make([]byte, 8)
	copy(ba[:r.by], bytes)
	ev := binary.LittleEndian.Uint64(ba)
	value := int64(ev) + r.Min
	r.Value = value
	return nil
}

func (r *IntegerField) IsValid() bool {
	if r.Value >= r.Min && r.Value <= r.Max {
		return true
	}
	return false
}

func MakeIntegerField(min, max, value int64) (*IntegerField, error) {
	length := uint(math.Ceil(math.Log2(float64(max - min + 1))))
	by := length / 8
	if length%8 != 0 {
		by += 1
	}
	if value < min || value > max {
		return nil, fmt.Errorf("integer %d is out of bounds (min: %d, max: %d)", value, min, max)
	}
	return &IntegerField{
		Max:    max,
		Min:    min,
		Value:  value,
		length: length,
		by:     by,
	}, nil
}
