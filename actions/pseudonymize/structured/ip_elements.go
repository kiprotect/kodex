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
	"bytes"
	"fmt"
)

type IPElement struct {
	value  []byte
	length uint
}

type IPElementIf interface {
	Value() []byte
	Length() uint
}

func (h *IPElement) Length() uint {
	return h.length
}

func (h *IPElement) Value() []byte {
	return h.value
}

func (h *IPElement) Equals(t Type) bool {
	ht, ok := t.(IPElementIf)
	if !ok {
		return false
	}
	if h.Length() != ht.Length() {
		return false
	}
	return bytes.Equal(h.Value(), ht.Value())
}

func (h *IPElement) Decode(b *BitArray) error {
	if b.Length() != h.Length() {
		return fmt.Errorf("Invalid length of byte array (%d vs %d)", b.Length(), h.Length())
	}

	h.value = b.Bytes()
	return nil
}

func (h *IPElement) Encode() (*BitArray, error) {
	return MakeBitArrayFromBytes(h.value, h.Length())
}

type IPAddress struct {
	IPElement
}

func (h *IPAddress) IsValid() bool {
	return true
}

func (h *IPAddress) Copy() Type {
	hc := *h
	return &hc
}

func MakeIPAddress(value []byte, mask uint) Type {
	return &IPAddress{
		IPElement: IPElement{
			value:  value,
			length: mask,
		},
	}
}
