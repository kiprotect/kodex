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

package structured

import (
	"fmt"
)

type CompositeListType struct {
	subtypes []Type
}

type CompositeListTypeIf interface {
	Subtypes() []Type
}

func (d *CompositeListType) IsValid() []bool {
	valid := make([]bool, len(d.subtypes))
	for i, subtype := range d.subtypes {
		valid[i] = subtype.IsValid()
	}
	return valid
}

func (d *CompositeListType) Equals(dr CompositeType) bool {
	drl, ok := dr.(CompositeListTypeIf)
	if !ok {
		return false
	}
	dSubtypes := drl.Subtypes()
	if len(d.subtypes) != len(dSubtypes) {
		return false
	}
	for i, subtype := range dSubtypes {
		if !subtype.Equals(d.subtypes[i]) {
			return false
		}
	}
	return true
}

func (d *CompositeListType) Copy() CompositeType {
	subtypes := make([]Type, len(d.subtypes))
	for i, subtype := range d.subtypes {
		subtypes[i] = subtype.Copy()
	}
	return &CompositeListType{
		subtypes: subtypes,
	}
}

func (d *CompositeListType) EncodeSubtype(i int) (*BitArray, error) {
	if i >= len(d.subtypes) {
		return nil, fmt.Errorf("out of bounds")
	}
	return d.subtypes[i].Encode()
}

func (d *CompositeListType) DecodeSubtype(i int, b *BitArray) error {
	if i >= len(d.subtypes) {
		return fmt.Errorf("out of bounds")
	}
	return d.subtypes[i].Decode(b)
}

func (d *CompositeListType) Length(i int) (uint, error) {
	if i >= len(d.subtypes) {
		return 0, fmt.Errorf("Index out of range")
	}
	subtype := d.subtypes[i]
	return subtype.Length(), nil
}

func (d *CompositeListType) Offset(i int) (uint, error) {
	var off uint = 0
	if i >= len(d.subtypes) {
		return 0, fmt.Errorf("Index out of range")
	}
	for j, subtype := range d.subtypes {
		if j == i {
			return off, nil
		}
		off += subtype.Length()
	}
	panic("This should never happen...")
	return 0, fmt.Errorf("This should never happen")
}

func (d *CompositeListType) Decode(b *BitArray) error {
	var pos uint = 0
	for _, subtype := range d.subtypes {
		ba, err := b.Extract(pos, subtype.Length())
		if err != nil {
			return err
		}
		if err = subtype.Decode(ba); err != nil {
			return err
		}
		pos += subtype.Length()
	}
	return nil
}

func (d *CompositeListType) Encode() (*BitArray, error) {
	ba := MakeBitArray(0)
	for _, subtype := range d.subtypes {
		bai, err := subtype.Encode()
		if err != nil {
			return nil, err
		}
		if err = ba.Append(bai); err != nil {
			return nil, err
		}
	}
	return ba, nil
}

func (d *CompositeListType) Marshal(format string) (interface{}, error) {
	return nil, nil
}

func (d *CompositeListType) Unmarshal(format string, value interface{}) error {
	return nil
}

func MakeCompositeListType() CompositeType {
	return &CompositeListType{}
}

func (d *CompositeListType) SetSubtypes(types []Type) {
	d.subtypes = types
}

func (d *CompositeListType) Subtypes() []Type {
	return d.subtypes
}

func (d *CompositeListType) Get(i int) (Type, error) {
	if i >= len(d.subtypes) {
		return nil, fmt.Errorf("out of bounds")
	}
	return d.subtypes[i], nil
}
