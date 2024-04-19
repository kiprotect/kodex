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

package merengue

import (
	"crypto/sha256"
	"encoding/binary"
)

var NULL = byte(0)

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

func Sha256(input []byte) []byte {
	h := sha256.Sum256(input)
	var out []byte = h[:]
	return out
}

func PseudonymizeBidirectional(inp []byte, l uint, keyA, keyB []byte, hash func([]byte) []byte) []byte {
	pa := Pseudonymize(inp, l, keyA, hash)
	return Pseudonymize(Revert(pa, l), l, keyB, hash)
}

func DepseudonymizeBidirectional(inp []byte, l uint, keyA, keyB []byte, hash func([]byte) []byte) []byte {
	pa := Depseudonymize(inp, l, keyB, hash)
	return Depseudonymize(Revert(pa, l), l, keyA, hash)
}

func Pseudonymize(inp []byte, l uint, key []byte, hash func([]byte) []byte) []byte {
	buf := make([]byte, 0, len(inp)+len(key)+4)
	ob := make([]byte, 0, len(inp))
	bs := make([]byte, 4)
	var off, ind, g uint32
	var b byte
	for g = 0; g < uint32(l); g++ {
		off = g % 8
		ind = g / 8
		buf = buf[:0]
		binary.LittleEndian.PutUint32(bs, g)
		buf = append(append(append(buf, ob...), key...), bs...)
		h := hash(buf)
		if off == 0 {
			ob = append(ob, NULL)
		}
		b = ((h[ind%uint32(len(h))] >> off) & 1) ^ ((inp[ind] >> off) & 1)
		ob[ind] |= b << off
	}
	return ob
}

func Depseudonymize(inp []byte, l uint, key []byte, hash func([]byte) []byte) []byte {
	fb := make([]byte, 0, len(inp))
	ob := make([]byte, 0, len(inp))
	buf := make([]byte, 0, len(inp)+len(key)+4)
	bs := make([]byte, 4)
	var ind, off, g uint32
	var b, bi byte
	for g = 0; g < uint32(l); g++ {
		off = g % 8
		ind = g / 8
		buf = buf[:0]
		binary.LittleEndian.PutUint32(bs, g)
		buf = append(append(append(buf, fb...), key...), bs...)
		h := hash(buf)
		if off == 0 {
			ob = append(ob, NULL)
			fb = append(fb, NULL)
		}
		bi = (inp[ind] >> off) & 1
		b = ((h[ind%uint32(len(h))] >> off) & 1) ^ bi
		ob[ind] |= b << off
		fb[ind] |= ((inp[ind] >> off) & 1) << off
	}
	return ob
}
