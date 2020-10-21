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

package groupByFunctions

import (
	"fmt"
	"github.com/kiprotect/go-helpers/errors"
	"github.com/kiprotect/kodex"
	"time"
)

type GroupByTimeWindow struct {
}

type TimeWindow struct {
	FromTime int64
	ToTime   int64
	Type     string
}

type TimeWindowFunction func(int64) []*TimeWindow

func getItemTime(item *kodex.Item, field string, parser TimeParser) (int64, error) {
	value, ok := item.Get(field)
	kodex.Log.Info(field)
	if !ok {
		return 0, errors.MakeExternalError("time window value not defined",
			"VALUE-NOT-DEFINED",
			field,
			nil)
	}

	t, err := parser(value)

	if err != nil {
		return 0, errors.MakeExternalError("time window invalid",
			"TIME-WINDOW-INVALID",
			map[string]interface{}{
				"field": field,
				"value": value},
			err)
	}

	return t, nil
}

func MakeTimeWindowFunction(config map[string]interface{}) (GroupByFunction, error) {

	format := config["format"].(string)
	field := config["field"].(string)
	window := config["window"].(string)

	timeWindowFunction := TimeWindowFunctions[window]
	parser := TimeParsers[format]
	formatter := TimeFormatters["rfc3339"]

	return func(item *kodex.Item) ([]*GroupByValue, error) {
		if t, err := getItemTime(item, field, parser); err != nil {
			return nil, err
		} else {
			timeWindows := timeWindowFunction(t)
			groups := make([]*GroupByValue, len(timeWindows))
			for i, timeWindow := range timeWindows {
				groups[i] = &GroupByValue{
					Values: map[string]interface{}{
						"from": formatter(timeWindow.FromTime),
						"to":   formatter(timeWindow.ToTime),
					},
					Expiration: timeWindow.ToTime,
				}
			}
			return groups, nil
		}
	}, nil
}

func minute(value int64) []*TimeWindow {
	t := time.Unix(value/1e9, value%1e9).UTC()
	from := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, t.Location())
	to := from.Add(time.Minute * 1)
	return []*TimeWindow{&TimeWindow{
		FromTime: from.UnixNano(),
		ToTime:   to.UnixNano(),
	}}
}

func hour(value int64) []*TimeWindow {
	t := time.Unix(value/1e9, value%1e9).UTC()
	from := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, t.Location())
	to := from.Add(time.Hour * 1)
	return []*TimeWindow{&TimeWindow{
		FromTime: from.UnixNano(),
		ToTime:   to.UnixNano(),
	}}
}

func day(value int64) []*TimeWindow {
	t := time.Unix(value/1e9, value%1e9).UTC()
	from := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	to := from.AddDate(0, 0, 1)
	return []*TimeWindow{&TimeWindow{
		FromTime: from.UnixNano(),
		ToTime:   to.UnixNano(),
	}}
}

func week(value int64) []*TimeWindow {
	t := time.Unix(value/1e9, value%1e9).UTC()
	from := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	wd := (int(t.Weekday()) - 1) % 7 // weekday starting from Monday
	if wd < 0 {
		wd += 7
	}
	from = from.AddDate(0, 0, -wd)
	to := from.AddDate(0, 0, 7)
	return []*TimeWindow{&TimeWindow{
		FromTime: from.UnixNano(),
		ToTime:   to.UnixNano(),
	}}
}

func month(value int64) []*TimeWindow {
	t := time.Unix(value/1e9, value%1e9).UTC()
	from := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
	to := from.AddDate(0, 1, 0)
	return []*TimeWindow{&TimeWindow{
		FromTime: from.UnixNano(),
		ToTime:   to.UnixNano(),
	}}
}

func year(value int64) []*TimeWindow {
	t := time.Unix(value/1e9, value%1e9).UTC()
	from := time.Date(t.Year(), 1, 1, 0, 0, 0, 0, t.Location())
	to := from.AddDate(1, 0, 0)
	return []*TimeWindow{&TimeWindow{
		FromTime: from.UnixNano(),
		ToTime:   to.UnixNano(),
	}}
}

func dayByHour(value int64) []*TimeWindow {
	windows := make([]*TimeWindow, 0)
	t := time.Unix(value/1e9, value%1e9).UTC()
	// we start 23 hours before the last full hour corresponding to t
	from := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, t.Location()).Add(-time.Hour * 23.0)
	// we end 24 hours later
	to := from.Add(time.Hour * 24.0)
	for i := time.Duration(0); i < 24; i++ {
		// we shift the time window by i hours and add it as well
		window := TimeWindow{
			FromTime: from.Add(time.Hour * i).UnixNano(),
			ToTime:   to.Add(time.Hour * i).UnixNano(),
		}
		windows = append(windows, &window)
	}
	return windows
}

func weekByDay(value int64) []*TimeWindow {
	windows := make([]*TimeWindow, 0)
	t := time.Unix(value/1e9, value%1e9).UTC()
	from := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).AddDate(0, 0, -6)
	to := from.AddDate(0, 0, 7)
	for i := 0; i < 7; i++ {
		window := TimeWindow{
			FromTime: from.AddDate(0, 0, i).UnixNano(),
			ToTime:   to.AddDate(0, 0, i).UnixNano(),
		}
		windows = append(windows, &window)
	}
	return windows
}

func monthByDay(value int64) []*TimeWindow {
	windows := make([]*TimeWindow, 0)
	t := time.Unix(value/1e9, value%1e9).UTC()
	from := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).AddDate(0, 0, -29)
	to := from.AddDate(0, 0, 30)
	for i := 0; i < 30; i++ {
		window := TimeWindow{
			FromTime: from.AddDate(0, 0, i).UnixNano(),
			ToTime:   to.AddDate(0, 0, i).UnixNano(),
		}
		windows = append(windows, &window)
	}
	return windows
}

var TimeWindowFunctions = map[string]TimeWindowFunction{
	"minute":       minute,
	"hour":         hour,
	"day":          day,
	"day-by-hour":  dayByHour,
	"week":         week,
	"week-by-day":  weekByDay,
	"month":        month,
	"month-by-day": monthByDay,
	"year":         year,
}

type TimeParser func(interface{}) (int64, error)

type TimeFormatter func(int64) interface{}

func RFC3339Formatter(input int64) interface{} {
	t := time.Unix(input/1e9, input%1e9).UTC()
	if input%1e9 != 0 {
		return t.Format(time.RFC3339Nano)
	} else {
		return t.Format(time.RFC3339)
	}
}

func RFC3339Parser(input interface{}) (int64, error) {
	strInput, ok := input.(string)
	if !ok {
		return 0, fmt.Errorf("expected a string")
	}
	t, err := time.Parse(time.RFC3339, strInput)
	if err != nil {
		return 0, err
	}
	return t.UnixNano(), nil
}

func getTime(input interface{}) (int64, error) {
	switch v := input.(type) {
	case int:
		return int64(v), nil
	case int64:
		return v, nil
	case uint64:
		return int64(v), nil
	case uint:
		return int64(v), nil
	case float64:
		return int64(v), nil
	}
	return 0, errors.MakeExternalError(
		"invalid type for time format",
		"INVALID-TYPE", nil, nil)

}

func UnixParser(input interface{}) (int64, error) {
	unixTime, err := getTime(input)
	if err != nil {
		return 0, err
	}
	return unixTime * 1e9, nil
}

func UnixFormatter(input int64) interface{} {
	return input / 1e9
}

func UnixNanoParser(input interface{}) (int64, error) {
	return getTime(input)
}

func UnixNanoFormatter(input int64) interface{} {
	return input
}

func UnixMilliParser(input interface{}) (int64, error) {
	t, err := getTime(input)
	if err != nil {
		return 0, err
	}
	return t * 1e6, nil
}

func UnixMilliFormatter(input int64) interface{} {
	return input / 1e6
}

var TimeParsers = map[string]TimeParser{
	"rfc3339":    RFC3339Parser,
	"unix":       UnixParser,
	"unix-nano":  UnixNanoParser,
	"unix-milli": UnixMilliParser,
}

var TimeFormatters = map[string]TimeFormatter{
	"rfc3339":    RFC3339Formatter,
	"unix":       UnixFormatter,
	"unix-nano":  UnixNanoFormatter,
	"unix-milli": UnixMilliFormatter,
}
