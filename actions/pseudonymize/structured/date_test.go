// KIProtect (Community Edition - CE) - Privacy & Security Engineering Platform
// Copyright (C) 2020  KIProtect GmbH (HRB 208395B) - Germany
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
	"testing"
)

func TestBasics(t *testing.T) {
	hour := MakeHour(12)
	if !hour.IsValid() {
		t.Errorf("Hour should be valid")
	}

	invalidHour := MakeHour(54)

	if invalidHour.IsValid() {
		t.Errorf("Invalid hour should be invalid")
	}

	var dm Type = invalidHour
	hr, ok := dm.(*Hour)

	if dm.IsValid() {
		t.Errorf("Hour should be invalid")
	}

	if !ok {
		t.Errorf("Not an hour")
	}
	if hr.IsValid() {
		t.Errorf("Hour should still be invalid")
	}

}

func TestHourEncoding(t *testing.T) {
	hour := MakeHour(12)

	encoded, err := hour.Encode()
	if err != nil {
		t.Fatal(err)
	}

	hourNew := MakeHour(1)
	err = hourNew.Decode(encoded)
	if err != nil {
		t.Fatal(err)
	}

	if !hourNew.Equals(hour) {
		t.Errorf("Should be equal")
	}

}

func TestDateEncoding(t *testing.T) {

	date, err := MakeDate(nil)
	if err != nil {
		t.Fatal(err)
	}
	if err = date.Unmarshal("%Y-%m-%d", "2008-10-12"); err != nil {
		t.Fatal(err)
	}

	ba, err := date.Encode()

	if err != nil {
		t.Fatal(err)
	}

	// we must copy the elements
	date2 := date.Copy()

	err = date2.Decode(ba)

	if err != nil {
		t.Fatal(err)
	}

	if !date2.Equals(date) {
		t.Errorf("Dates should be equal!")
	}

}

func TestParsing(t *testing.T) {
	format := "%Y-%m-%d %%"
	fields, err := parseFormatString(format)
	if err != nil {
		t.Fatal(err)
	}
	_, err = parseDateString(fields, "2008-01-10 %")
	if err != nil {
		t.Error(err)
	}
}

func TestRangeParsing(t *testing.T) {
	format := "%(2000-2100)Y-%m-%d %%"
	fields, err := parseFormatString(format)
	if err != nil {
		t.Fatal(err)
	}
	field := fields[0]
	if !field.hasBounds || field.min != 2000 || field.max != 2100 {
		t.Errorf("Field bounds not correct")
	}
	_, err = parseDateString(fields, "2008-01-10 %")
	if err != nil {
		t.Error(err)
	}

	_, err = parseDateString(fields, "2108-01-10 %")
	if err == nil {
		t.Errorf("Should throw an error")
	}

}

func TestValididityCheck(t *testing.T) {
	date, err := MakeDate(nil)
	if err != nil {
		t.Fatal(err)
	}

	if err = date.Unmarshal("%Y-%m-%d", "2008-13-12"); err != nil {
		t.Fatal(err)
	}
	valid := date.IsValid()
	if valid[0] == false || valid[1] == true || valid[2] == true {
		t.Errorf("Only month and day should be invalid")
	}

	if date, err = MakeDate(nil); err != nil {
		t.Fatal(err)
	}
	if err = date.Unmarshal("%Y-%m-%d", "2007-02-29"); err != nil {
		t.Fatal(err)
	}
	valid = date.IsValid()
	if valid[0] == false || valid[1] == false || valid[2] == true {
		t.Errorf("Only day should be invalid")
	}

	if date, err = MakeDate(nil); err != nil {
		t.Fatal(err)
	}
	if err = date.Unmarshal("%Y-%m-%d", "2000-02-29"); err != nil {
		t.Fatal(err)
	}
	valid = date.IsValid()
	if valid[0] == false || valid[1] == false || valid[2] == false {
		t.Errorf("Everything should be valid")
	}
}

func TestMarshalAndUnmarshal(t *testing.T) {
	date, err := MakeDate(nil)
	if err != nil {
		t.Fatal(err)
	}
	origValue := "2008-10-12"
	if err = date.Unmarshal("%Y-%m-%d", origValue); err != nil {
		t.Fatal(err)
	}
	value, err := date.Marshal("%Y-%m-%d")
	if err != nil {
		t.Fatal(err)
	}
	str := value.(string)
	if str != origValue {
		t.Errorf("Should match original date")
	}
	value, err = date.Marshal("%d-%Y-%m foobar")
	if err != nil {
		t.Fatal(err)
	}
	str = value.(string)
	if str != "12-2008-10 foobar" {
		t.Errorf("Should match original date")
	}

	date2, err := MakeDate(nil)
	if err != nil {
		t.Fatal(err)
	}

	if err = date2.Unmarshal("%d-%Y-%m foobar", "12-2008-10 foobar"); err != nil {
		t.Fatal(err)
	}

	if !date2.Equals(date) {
		t.Errorf("Dates should be equal")
	}

	date3, err := MakeDate(nil)
	if err != nil {
		t.Fatal(err)
	}

	origValue = "2008-10-12 +0400"
	err = date3.Unmarshal("%Y-%m-%d %z", origValue)
	if err != nil {
		t.Fatal(err)
	}
	value, err = date3.Marshal("%Y-%m-%d %z")
	if err != nil {
		t.Fatal(err)
	}
	str = value.(string)
	if str != origValue {
		t.Errorf("Strings should be identical")
	}

}

func TestParseFormatString(t *testing.T) {
	format := "%Y-%m-%d %%"
	fields, err := parseFormatString(format)
	if err != nil {
		t.Fatal(err)
	}
	if len(fields) != 7 {
		t.Fatalf("Invalid number of fields")
	}
	if fields[0].field != 'Y' || fields[0].value != "" ||
		fields[1].field != 0 || fields[1].value != "-" ||
		fields[2].field != 'm' || fields[2].value != "" ||
		fields[3].field != 0 || fields[3].value != "-" ||
		fields[4].field != 'd' || fields[4].value != "" ||
		fields[5].field != 0 || fields[5].value != " " ||
		fields[6].field != 0 || fields[6].value != "%" {
		t.Fatalf("Invalid result")
	}
}

func TestInvalidFormatString(t *testing.T) {
	format := "%Y-%m-%d %Y %%"
	_, err := parseFormatString(format)
	if err == nil {
		t.Fatal(err)
	}
}
