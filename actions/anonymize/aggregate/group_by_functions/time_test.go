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

package aggregate

import (
	"testing"
	"time"
)

type TimeTest struct {
	T       string
	ExpFrom string
	ExpTo   string
	Window  string
}

var timeTests = []TimeTest{
	{
		T:       "2017-01-04T23:21:22Z",
		ExpFrom: "2017-01-04T23:21:00Z",
		ExpTo:   "2017-01-04T23:22:00Z",
		Window:  "minute",
	},
	{
		T:       "2017-01-04T23:21:22+01:00",
		ExpFrom: "2017-01-04T22:21:00Z",
		ExpTo:   "2017-01-04T22:22:00Z",
		Window:  "minute",
	},
	{
		T:       "2017-01-04T23:21:22Z",
		ExpFrom: "2017-01-04T23:00:00Z",
		ExpTo:   "2017-01-05T00:00:00Z",
		Window:  "hour",
	},
	{
		T:       "2017-01-04T23:21:22-01:00",
		ExpFrom: "2017-01-05T00:00:00Z",
		ExpTo:   "2017-01-05T01:00:00Z",
		Window:  "hour",
	},
	{
		T:       "2017-01-04T23:21:22Z",
		ExpFrom: "2017-01-04T00:00:00Z",
		ExpTo:   "2017-01-05T00:00:00Z",
		Window:  "day",
	},
	{
		T:       "2017-01-04T23:21:22Z",
		ExpFrom: "2017-01-02T00:00:00Z",
		ExpTo:   "2017-01-09T00:00:00Z",
		Window:  "week",
	},
	{
		T:       "2017-01-02T23:21:22Z",
		ExpFrom: "2017-01-02T00:00:00Z",
		ExpTo:   "2017-01-09T00:00:00Z",
		Window:  "week",
	},
	{
		T:       "2017-01-08T23:21:22Z",
		ExpFrom: "2017-01-02T00:00:00Z",
		ExpTo:   "2017-01-09T00:00:00Z",
		Window:  "week",
	},
	{
		T:       "2017-01-08T23:21:22Z",
		ExpFrom: "2017-01-01T00:00:00Z",
		ExpTo:   "2017-02-01T00:00:00Z",
		Window:  "month",
	},
	{
		T:       "2017-02-08T23:21:22Z",
		ExpFrom: "2017-02-01T00:00:00Z",
		ExpTo:   "2017-03-01T00:00:00Z",
		Window:  "month",
	},
	{
		T:       "2017-01-08T23:21:22Z",
		ExpFrom: "2017-01-01T00:00:00Z",
		ExpTo:   "2018-01-01T00:00:00Z",
		Window:  "year",
	},
}

func TestMinute(t *testing.T) {
	for _, test := range timeTests {
		ot, _ := time.Parse(time.RFC3339, test.T)
		expFrom, _ := time.Parse(time.RFC3339, test.ExpFrom)
		expTo, _ := time.Parse(time.RFC3339, test.ExpTo)
		window, _ := TimeWindowFunctions[test.Window]
		mws := window(ot.UnixNano())
		if len(mws) != 1 {
			t.Errorf("Expected one window")
			continue
		}
		mw := mws[0]
		from := time.Unix(mw.FromTime/1e9, mw.FromTime%1e9).UTC()
		to := time.Unix(mw.ToTime/1e9, mw.ToTime%1e9).UTC()
		if from != expFrom {
			t.Errorf("Expected %s, got %s", expFrom, from)
		}
		if to != expTo {
			t.Errorf("Expected %s, got %s", expTo, to)
		}
	}
}

func TestUnixTime(t *testing.T) {
	v := int64(145434543)
	res, err := UnixParser(v)
	if err != nil {
		t.Error(err)
	}
	if res != v*1e9 {
		t.Errorf("Expected %d, got %d", v*1e9, res)
	}
	res, err = UnixParser(uint(v))
	if err != nil {
		t.Error(err)
	}
	if res != v*1e9 {
		t.Errorf("Expected %d, got %d", v*1e9, res)
	}
	res, err = UnixParser(uint64(v))
	if err != nil {
		t.Error(err)
	}
	if res != v*1e9 {
		t.Errorf("Expected %d, got %d", v*1e9, res)
	}
	res, err = UnixParser(int(v))
	if err != nil {
		t.Error(err)
	}
	if res != v*1e9 {
		t.Errorf("Expected %d, got %d", v*1e9, res)
	}
	f := UnixFormatter(res)
	fv, ok := f.(int64)
	if !ok {
		t.Errorf("Expected an int64")
	}
	if fv != v {
		t.Errorf("Expected %d, got %d", v, fv)
	}
}

func TestUnixNanoTime(t *testing.T) {
	v := int64(14543454323423432)
	res, err := UnixNanoParser(v)
	if err != nil {
		t.Error(err)
	}
	if res != v {
		t.Errorf("Expected %d, got %d", v, res)
	}
	res, err = UnixNanoParser(uint(v))
	if err != nil {
		t.Error(err)
	}
	if res != v {
		t.Errorf("Expected %d, got %d", v, res)
	}
	res, err = UnixNanoParser(uint64(v))
	if err != nil {
		t.Error(err)
	}
	if res != v {
		t.Errorf("Expected %d, got %d", v, res)
	}
	res, err = UnixNanoParser(int(v))
	if err != nil {
		t.Error(err)
	}
	if res != v {
		t.Errorf("Expected %d, got %d", v, res)
	}
	f := UnixNanoFormatter(res)
	fv, ok := f.(int64)
	if !ok {
		t.Errorf("Expected an int64")
	}
	if fv != v {
		t.Errorf("Expected %d, got %d", v, fv)
	}
}

func TestRfc3339Time(t *testing.T) {
	v := "2017-07-10T10:33:16Z"
	vu := int64(1499682796000000000)
	res, err := RFC3339Parser(v)
	if err != nil {
		t.Error(err)
	}
	if res != vu {
		t.Errorf("Expected %d, got %d", vu, res)
	}

	fRes := RFC3339Formatter(res)
	fResStr, ok := fRes.(string)
	if !ok {
		t.Errorf("Expected a string destination")
	}
	if fResStr != v {
		t.Errorf("Expected %s, got %s", v, fResStr)
	}
}

func TestRfc3339TimeFloat(t *testing.T) {
	v := "2017-07-10T10:33:16.123456789Z"
	vu := int64(1499682796123456789)
	res, err := RFC3339Parser(v)
	if err != nil {
		t.Error(err)
	}
	if res != vu {
		t.Errorf("Expected %d, got %d", vu, res)
	}

	fRes := RFC3339Formatter(res)
	fResStr, ok := fRes.(string)
	if !ok {
		t.Errorf("Expected a string destination")
	}
	if fResStr != v {
		t.Errorf("Expected %s, got %s", v, fResStr)
	}
}
