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
	"fmt"
	"regexp"
	"sort"
	"strconv"
)

var rangeRegex = regexp.MustCompile(`^\((\d*)-(\d*)\)`)

var formatRegexes = map[byte]*regexp.Regexp{
	'Y': regexp.MustCompile(`^\d{4}`),
	'm': regexp.MustCompile(`^\d{2}`),
	'd': regexp.MustCompile(`^\d{2}`),
	'M': regexp.MustCompile(`^\d{2}`),
	'S': regexp.MustCompile(`^\d{2}`),
	'H': regexp.MustCompile(`^\d{2}`),
	'n': regexp.MustCompile(`^\d{9}`),
	'z': regexp.MustCompile(`^(\+|-)\d{4}`),
}

var formatDefinitions = map[byte]func(int64) Type{
	'Y': MakeYear,
	'm': MakeMonth,
	'd': MakeDay,
	'M': MakeMinute,
	'S': MakeSecond,
	'H': MakeHour,
	'n': MakeNanoSecond,
	'z': MakeTimeZone,
}

var printFormats = map[byte]string{
	'Y': "%04d",
	'm': "%02d",
	'd': "%02d",
	'M': "%02d",
	'H': "%02d",
	'S': "%02d",
	'z': "%+05d",
	'n': "%09d",
}

type FormatField struct {
	value string
	field byte
	// optional minimum/maximum values
	min, max int64
	// true if bounds are given
	hasBounds bool
}

// Used for sorting date elements by their order
type Elements []Type

func (e Elements) Len() int {
	return len(e)
}

func (e Elements) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

func (e Elements) Less(i, j int) bool {
	del := e[i].(DateElementIf)
	der := e[j].(DateElementIf)
	return del.Order() < der.Order()
}

//helper function to parse date strings using a given format
func parseDateString(format []FormatField, dateString string) ([]Type, error) {
	/*
		Parses a date string following the given format
	*/
	p := 0
	elements := make([]Type, 0)
	for _, field := range format {
		if p >= len(dateString) {
			return nil, fmt.Errorf("unexpected ending")
		}
		if field.field == 0 {
			if p+len(field.value) > len(dateString) || len(dateString)-p < len(field.value) || dateString[p:p+len(field.value)] != field.value {
				return nil, fmt.Errorf("expected a '%s' at position %d", field.value, p)
			} else {
				p += len(field.value)
			}
		} else {
			// this is a field
			element, length, err := regexFieldParser(field, dateString[p:])
			if err != nil {
				return nil, err
			}
			p += length
			elements = append(elements, element)
		}
	}
	if p != len(dateString) {
		return nil, fmt.Errorf("extraneous input when parsing date: '%s'", dateString[p:])
	}
	sort.Sort(Elements(elements))
	return elements, nil
}

// helper function to format a date according to a format string
func formatDate(format []FormatField, elements []Type) (string, error) {
	str := ""
	for _, field := range format {
		if field.field == 0 {
			str += field.value
		} else {
			// we select the appropriate field from the subtypes and convert
			// it to a string representation
			var value int64
			found := false
			for _, element := range elements {
				de, ok := element.(DateElementIf)
				if !ok {
					continue
				}
				if de.Field() == field.field {
					value = de.Value()
					found = true
					break
				}
			}
			if !found {
				return "", fmt.Errorf("Field %%%s not found", string(field.field))
			}
			printFormat := printFormats[field.field]
			str += fmt.Sprintf(printFormat, value)
		}
	}
	return str, nil
}

func regexFieldParser(field FormatField, dateString string) (Type, int, error) {
	re := formatRegexes[field.field]
	if re == nil {
		return nil, 0, fmt.Errorf("invalid field name: %%%s", string(field.field))
	}
	str := re.FindString(dateString)
	if str == "" {
		return nil, 0, fmt.Errorf("does not match format for %%%s", string(field.field))
	}
	num, err := strconv.Atoi(str)
	if err != nil {
		return nil, 0, err
	}
	if field.hasBounds && (int64(num) < field.min || int64(num) > field.max) {
		return nil, 0, fmt.Errorf("field %%%s is out of bounds (must be between %d-%d but is %d)", string(field.field), field.min, field.max, num)
	}
	maker := formatDefinitions[field.field]
	element := maker(int64(num))
	if field.hasBounds {
		dateElement, ok := element.(DateElementIf)
		if !ok {
			return nil, 0, fmt.Errorf("error when converting to date element")
		}
		if err = dateElement.SetRange(field.min, field.max); err != nil {
			return nil, 0, fmt.Errorf("error when setting bounds for %%%s: %s", string(field.field), err.Error())
		}
	}
	return element, len(str), nil
}

func parseFormatString(format string) ([]FormatField, error) {
	fieldMap := make(map[byte]bool)
	formatFields := make([]FormatField, 0)
	i := 0
	lastI := i
	for {
		if i >= len(format) {
			break
		}
		if format[i] == '%' && i < len(format)-1 {
			var hasBounds bool
			var min, max int
			var err error
			if lastI != i {
				// we append the intermediate string to the list
				formatFields = append(formatFields, FormatField{field: 0, value: format[lastI:i]})
			}
			if format[i+1] == '%' {
				formatFields = append(formatFields, FormatField{value: "%", field: 0})
			} else {
				if format[i+1] == '(' {
					//there is an additional range argument here
					matches := rangeRegex.FindStringSubmatch(format[i+1:])
					if len(matches) == 0 {
						return nil, fmt.Errorf("Expected a parameter of the form (<from>-<to>), e.g. '(2000-2100)'")
					}
					hasBounds = true
					minStr, maxStr := matches[1], matches[2]
					min, err = strconv.Atoi(minStr)
					if err != nil {
						return nil, err
					}
					max, err = strconv.Atoi(maxStr)
					if err != nil {
						return nil, err
					}
					if min >= max {
						return nil, fmt.Errorf("The min value (%d) must be smaller than the max value (%d)", min, max)
					}
					i += len(matches[0])
				}
				if i >= len(format)-1 {
					return nil, fmt.Errorf("Expected a format specifier")
				}
				f := format[i+1]

				if formatRegexes[f] == nil {
					return nil, fmt.Errorf("Invalid field: %%%s", string(f))
				}
				if fieldMap[f] == true {
					return nil, fmt.Errorf("Field %%%s was repeated", string(f))
				}
				fieldMap[f] = true
				formatFields = append(formatFields, FormatField{value: "", field: f, hasBounds: hasBounds, min: int64(min), max: int64(max)})
			}
			i += 2
			lastI = i
			continue
		}
		i += 1
	}
	if lastI != i {
		formatFields = append(formatFields, FormatField{field: 0, value: format[lastI:i]})
	}
	return formatFields, nil
}
