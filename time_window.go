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

package kodex

import (
	"time"
)

type TimeWindow struct {
	From int64
	To   int64
	Type string
}

type TimeWindowFunc func(int64) TimeWindow

func Second(value int64) TimeWindow {
	t := time.Unix(value/1e9, value%1e9).UTC()
	from := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), 0, t.Location())
	to := from.Add(time.Second * 1)
	return TimeWindow{
		From: from.UnixNano(),
		To:   to.UnixNano(),
		Type: "second",
	}
}

func Minute(value int64) TimeWindow {
	t := time.Unix(value/1e9, value%1e9).UTC()
	from := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, t.Location())
	to := from.Add(time.Minute * 1)
	return TimeWindow{
		From: from.UnixNano(),
		To:   to.UnixNano(),
		Type: "minute",
	}
}

func Hour(value int64) TimeWindow {
	t := time.Unix(value/1e9, value%1e9).UTC()
	from := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, t.Location())
	to := from.Add(time.Hour * 1)
	return TimeWindow{
		From: from.UnixNano(),
		To:   to.UnixNano(),
		Type: "hour",
	}
}

func Day(value int64) TimeWindow {
	t := time.Unix(value/1e9, value%1e9).UTC()
	from := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	to := from.AddDate(0, 0, 1)
	return TimeWindow{
		From: from.UnixNano(),
		To:   to.UnixNano(),
		Type: "day",
	}
}

func Week(value int64) TimeWindow {
	t := time.Unix(value/1e9, value%1e9).UTC()
	from := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	wd := (int(t.Weekday()) - 1) % 7 // weekday starting from Monday
	if wd < 0 {
		wd += 7
	}
	from = from.AddDate(0, 0, -wd)
	to := from.AddDate(0, 0, 7)
	return TimeWindow{
		From: from.UnixNano(),
		To:   to.UnixNano(),
		Type: "week",
	}
}

func Month(value int64) TimeWindow {
	t := time.Unix(value/1e9, value%1e9).UTC()
	from := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
	to := from.AddDate(0, 1, 0)
	return TimeWindow{
		From: from.UnixNano(),
		To:   to.UnixNano(),
		Type: "month",
	}
}

func MakeTimeWindow(t int64, twType string) TimeWindow {
	switch twType {
	case "second":
		return Second(t)
	case "minute":
		return Minute(t)
	case "hour":
		return Hour(t)
	case "day":
		return Day(t)
	case "week":
		return Week(t)
	case "month":
		return Month(t)
	}
	return TimeWindow{}
}

func (t *TimeWindow) Copy() TimeWindow {
	return TimeWindow{
		Type: t.Type,
		From: t.From,
		To:   t.To,
	}
}

func (t *TimeWindow) IncreaseBy(n int64) {
	from := time.Unix(t.From/1e9, t.From%1e9).UTC()
	to := time.Unix(t.To/1e9, t.To%1e9).UTC()
	switch t.Type {
	case "second":
		from = from.Add(time.Second * time.Duration(n))
		to = to.Add(time.Second * time.Duration(n))
	case "minute":
		from = from.Add(time.Minute * time.Duration(n))
		to = to.Add(time.Minute * time.Duration(n))
	case "hour":
		from = from.Add(time.Hour * time.Duration(n))
		to = to.Add(time.Hour * time.Duration(n))
	case "day":
		from = from.AddDate(0, 0, int(n))
		to = to.AddDate(0, 0, int(n))
	case "week":
		from = from.AddDate(0, 0, 7*int(n))
		to = to.AddDate(0, 0, 7*int(n))
	case "month":
		from = from.AddDate(0, int(n), 0)
		to = to.AddDate(0, int(n), 0)
	}
	t.From = from.UnixNano()
	t.To = to.UnixNano()
}
