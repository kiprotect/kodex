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
	"time"
)

type Date struct {
	CompositeListType
}

func MakeDate(interface{}) (CompositeType, error) {
	return &Date{}, nil
}

func (d *Date) Copy() CompositeType {
	listCopy := d.CompositeListType.Copy()
	listCopyType, _ := listCopy.(*CompositeListType)
	return &Date{
		CompositeListType: *listCopyType,
	}
}

func (d *Date) Unmarshal(format string, data interface{}) error {
	formatFields, err := parseFormatString(format)
	if err != nil {
		return err
	}
	str, ok := data.(string)
	if !ok {
		byteArray, ok := data.([]byte)
		if ok {
			str = string(byteArray)
		} else {
			return fmt.Errorf("expected a string or byte array as input")
		}
	}
	elements, err := parseDateString(formatFields, str)
	if err != nil {
		return err
	}
	d.SetSubtypes(elements)
	return nil
}

func (d *Date) Marshal(format string) (interface{}, error) {
	formatFields, err := parseFormatString(format)
	if err != nil {
		return nil, err
	}
	return formatDate(formatFields, d.subtypes)
}

func (d *Date) IsValid() []bool {
	valid := d.CompositeListType.IsValid()
	var year, month, day int64
	var yearInd, monthInd, dayInd int
	yearInd, monthInd, dayInd = -1, -1, -1
	for i, field := range d.subtypes {
		dateField := field.(DateElementIf)
		switch dateField.Field() {
		case 'Y':
			year = dateField.Value()
			yearInd = i
			break
		case 'm':
			month = dateField.Value()
			monthInd = i
			break
		case 'd':
			day = dateField.Value()
			dayInd = i
			break
		}
	}
	if yearInd == -1 || monthInd == -1 || dayInd == -1 {
		return valid
	}
	// we check if the date is valid (e.g. if it is not February 31st)
	dateStr := fmt.Sprintf("%04d-%02d-%02d", year, month, day)
	if _, err := time.Parse("2006-01-02", dateStr); err != nil {
		valid[dayInd] = false
	}
	return valid
}
