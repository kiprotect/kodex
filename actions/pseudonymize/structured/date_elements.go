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
)

type DateElement struct {
	value          int64
	length         uint
	min, max       int64
	absMin, absMax int64
	field          byte
	order          uint
}

type DateElementIf interface {
	Value() int64
	Length() uint
	Order() uint
	SetRange(int64, int64) error
	Field() byte
}

func (h *DateElement) SetRange(min, max int64) error {
	if min < h.absMin || max > h.absMax {
		return fmt.Errorf("out of bounds (min: %d, max: %d)", h.absMin, h.absMax)
	}
	h.min = min
	h.max = max
	h.length = uint(math.Ceil(math.Log2(float64(max - min + 1))))
	return nil
}

func (h *DateElement) Order() uint {
	return h.order
}

func (h *DateElement) Length() uint {
	return h.length
}

func (h *DateElement) Field() byte {
	return h.field
}

func (h *DateElement) Value() int64 {
	return h.value
}

func (h *DateElement) Equals(t Type) bool {
	ht, ok := t.(DateElementIf)
	if !ok {
		return false
	}
	if h.Value() != ht.Value() || h.Length() != ht.Length() {
		return false
	}
	return true
}

func (h *DateElement) IsValid() bool {
	if h.min <= h.value && h.value <= h.max {
		return true
	}
	return false
}

func (h *DateElement) Decode(b *BitArray) error {

	if b.Length() != h.length {
		return fmt.Errorf("Invalid length of byte array (%d vs %d)", b.Length(), h.length)
	}

	be := b.Bytes()
	be = append(be, make([]byte, 8-len(be))...)
	vv := int64(binary.LittleEndian.Uint64(be))
	h.value = vv + h.min
	return nil
}

func (h *DateElement) Encode() (*BitArray, error) {
	b := make([]byte, 8)
	vv := uint64(h.value - h.min)
	binary.LittleEndian.PutUint64(b, vv)
	bl := h.length / 8
	if h.length%8 != 0 {
		bl += 1
	}
	return MakeBitArrayFromBytes(b[:bl], h.length)
}

/*
Year
*/

type Year struct {
	DateElement
}

func MakeYear(value int64) Type {
	return &Year{
		DateElement: DateElement{
			value:  value,
			length: 9,
			absMin: 0,
			absMax: 100000,
			min:    1800,
			max:    2200,
			order:  0,
			field:  'Y',
		},
	}
}

func (y *Year) Copy() Type {
	yc := *y
	return &yc
}

/*
Month
*/

type Month struct {
	DateElement
}

func MakeMonth(value int64) Type {
	return &Month{
		DateElement: DateElement{
			value:  value,
			length: 4,
			min:    1,
			max:    12,
			absMin: 1,
			absMax: 12,
			order:  1,
			field:  'm',
		},
	}
}

func (m *Month) Copy() Type {
	mc := *m
	return &mc
}

/*
Day
*/

type Day struct {
	DateElement
}

func MakeDay(value int64) Type {
	return &Day{
		DateElement: DateElement{
			value:  value,
			length: 5,
			order:  2,
			min:    1,
			max:    31,
			absMin: 1,
			absMax: 31,
			field:  'd',
		},
	}
}

func (d *Day) Copy() Type {
	dc := *d
	return &dc
}

/*
Time Zone
*/

var timeZones = map[int]int{
	-1200: 0,
	-1100: 1,
	-1000: 2,
	-930:  3,
	-900:  4,
	-800:  5,
	-700:  6,
	-600:  7,
	-500:  8,
	-400:  9,
	-330:  10,
	-300:  11,
	-200:  12,
	-100:  13,
	0:     14,
	100:   15,
	200:   16,
	300:   17,
	330:   18,
	400:   19,
	430:   20,
	500:   21,
	530:   22,
	545:   23,
	600:   24,
	630:   25,
	700:   26,
	800:   27,
	845:   28,
	900:   29,
	930:   30,
	1000:  31,
	1030:  32,
	1100:  33,
	1200:  34,
	1245:  35,
	1300:  36,
	1400:  37,
}

type TimeZone struct {
	DateElement
}

func MakeTimeZone(value int64) Type {
	ind, ok := timeZones[int(value)]
	if !ok {
		ind = -1
	}
	ind += 1
	return &TimeZone{
		DateElement: DateElement{
			value:  int64(ind),
			length: 6,
			order:  3,
			min:    1,
			max:    38,
			absMin: 1,
			absMax: 38,
			field:  'z',
		},
	}
}

func (t *TimeZone) Copy() Type {
	tc := *t
	return &tc
}

func (h *TimeZone) Value() int64 {
	if h.value == 0 {
		//this is an invalid time zone
		return -1300
	}
	for key, value := range timeZones {
		if int64(value) == h.value-1 {
			return int64(key)
		}
	}
	return -1300
}

/*
Hour
*/

type Hour struct {
	DateElement
}

func MakeHour(value int64) Type {
	return &Hour{
		DateElement: DateElement{
			value:  value,
			length: 5,
			min:    0,
			max:    23,
			absMin: 0,
			absMax: 23,
			order:  4,
			field:  'H',
		},
	}
}

func (h *Hour) Copy() Type {
	hc := *h
	return &hc
}

/*
Minute
*/

type Minute struct {
	DateElement
}

func MakeMinute(value int64) Type {
	return &Minute{
		DateElement: DateElement{
			value:  value,
			length: 6,
			min:    0,
			max:    59,
			absMin: 0,
			absMax: 59,
			order:  5,
			field:  'M',
		},
	}
}

func (h *Minute) Copy() Type {
	hc := *h
	return &hc
}

/*
Second
*/

type Second struct {
	DateElement
}

func MakeSecond(value int64) Type {
	return &Second{
		DateElement: DateElement{
			value:  value,
			length: 6,
			order:  6,
			min:    0,
			max:    59,
			absMin: 0,
			absMax: 59,
			field:  'S',
		},
	}
}

func (s *Second) Copy() Type {
	sc := *s
	return &sc
}

/*
Nanosecond
*/

type NanoSecond struct {
	DateElement
}

func MakeNanoSecond(value int64) Type {
	return &NanoSecond{
		DateElement: DateElement{
			value:  value,
			length: 30,
			order:  7,
			min:    0,
			max:    1e9 - 1,
			absMin: 0,
			absMax: 1e9 - 1,
			field:  'n',
		},
	}
}

func (n *NanoSecond) Copy() Type {
	nc := *n
	return &nc
}
