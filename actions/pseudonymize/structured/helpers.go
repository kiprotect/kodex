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

type BitArray struct {
	bytes  []byte
	length uint
	by     uint
}

func MakeBitArray(length uint) *BitArray {
	by := length / 8
	if length%8 != 0 {
		by += 1
	}
	return &BitArray{
		length: length,
		by:     by,
		bytes:  make([]byte, by),
	}
}
func (b *BitArray) Append(ba *BitArray) error {
	newLength := b.length + ba.length
	newBy := newLength / 8
	if newLength%8 != 0 {
		newBy += 1
	}
	if newBy > b.by {
		b.bytes = append(b.bytes, make([]byte, newBy-b.by)...)
	}
	b.length = newLength
	b.by = newBy
	return b.Update(ba, b.length-ba.length)
}

func (b *BitArray) AsString() string {
	s := ""
	bo := b.length % 8
	for i, bb := range b.bytes {
		if (uint(i) < b.by-1) || bo == 0 {
			for j := uint(0); j < 8; j++ {
				if ((bb >> j) & 1) == 1 {
					s = s + "1"
				} else {
					s = s + "0"
				}
			}
			if uint(i) != b.by-1 {
				s = s + "-"
			}
		} else {
			for j := uint(0); j < bo; j++ {
				if ((bb >> j) & 1) == 1 {
					s = s + "1"
				} else {
					s = s + "0"
				}
			}
		}
	}
	return s
}

func (b *BitArray) Equals(ba *BitArray) bool {
	if b.length != ba.length {
		return false
	}
	for i, bi := range b.bytes {
		if ba.bytes[i] != bi {
			return false
		}
	}
	return true
}

func MakeBitArrayFromBytes(bytes []byte, length uint) (*BitArray, error) {
	by := length / 8
	bo := length % 8
	if bo != 0 {
		by += 1
	}
	if by != uint(len(bytes)) {
		return nil, fmt.Errorf("Invalid length of byte array!")
	}
	if bo != 0 {
		bytes[by-1] &= 0xFF >> (8 - bo)
	}
	return &BitArray{
		bytes:  bytes,
		length: length,
		by:     by,
	}, nil
}

func (b *BitArray) Copy() *BitArray {
	nb := make([]byte, b.by)
	copy(nb, b.bytes)
	return &BitArray{
		bytes:  nb,
		length: b.length,
		by:     b.by,
	}
}

func (b *BitArray) Update(ba *BitArray, offset uint) error {
	/*
		* Get the start and end index for the update
		* For each segment, construct the update value
		* For the start and end segment, apply an additional mask

		Each segment will be
	*/

	if offset+ba.length > b.length {
		return fmt.Errorf("Out of bounds")
	}

	start := offset / 8
	stop := (offset + ba.length) / 8
	ob := (offset + ba.length) % 8
	if ob != 0 {
		stop += 1
	}
	n := stop - start
	bo := offset % 8
	for i := uint(0); i < n; i++ {
		ind := start + i
		if i == 0 {
			/*
				This is the start segment. The content will be the lower 8-bo bits
				of the first byte of ba.
			*/
			if ind == stop-1 && ob != 0 {
				b.bytes[ind] &= (0xFF >> (8 - bo)) | (0xFF << ob)
				b.bytes[ind] |= (ba.bytes[i] & (0xFF >> (8 - ob)) << bo)

			} else {
				b.bytes[ind] &= 0xFF >> (8 - bo)
				b.bytes[ind] |= ba.bytes[i] << bo
			}
		} else {
			/*
				This is an intermediate segment. The content will be the 8-bo
				upper bits of the i-1-th byte and the bo lower bits of the
				i-th byte of ba, shifted by 8-bo bits.
			*/
			if i < ba.by {
				if ind == stop-1 && ob != 0 {
					b.bytes[ind] &= (0xFF << ob)
					b.bytes[ind] |= (((ba.bytes[i] << bo) & 0xFF) | (ba.bytes[i-1] >> (8 - bo)))
				} else {
					b.bytes[ind] = ((ba.bytes[i] << bo) & 0xFF) | (ba.bytes[i-1] >> (8 - bo))
				}
			} else {
				if ind == stop-1 && ob != 0 {
					b.bytes[ind] &= 0xFF << ob
					b.bytes[ind] |= (ba.bytes[i-1] >> (8 - bo))
				} else {
					b.bytes[ind] &= 0xFF << bo
					b.bytes[ind] |= (ba.bytes[i-1] >> (8 - bo))
				}
			}
		}
	}
	return nil
}

func (b *BitArray) Extract(offset, length uint) (*BitArray, error) {

	if offset+length > b.length {
		return nil, fmt.Errorf("Out of bounds")
	}

	start := offset / 8
	stop := (offset + length) / 8
	ob := (offset + length) % 8
	if ob != 0 {
		stop += 1
	}

	n := stop - start
	bo := offset % 8

	ba := MakeBitArray(length)

	for i := uint(0); i < n; i++ {
		ind := start + i
		if i < ba.by {
			if ind+1 < b.by {
				ba.bytes[i] = (b.bytes[ind] >> bo) | (b.bytes[ind+1] << (8 - bo))
			} else {
				ba.bytes[i] = (b.bytes[ind] >> bo)
			}
		} else {
			ba.bytes[i-1] |= (b.bytes[ind] << (8 - bo))
		}
	}
	if length%8 != 0 {
		ba.bytes[ba.by-1] &= 0xFF >> (8 - (length % 8))
	}
	return ba, nil
}

func (b *BitArray) Bytes() []byte {
	return b.bytes
}

func (b *BitArray) Length() uint {
	return b.length
}
