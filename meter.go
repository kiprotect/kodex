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

package kodex

type Metric struct {
	Name       string
	TimeWindow TimeWindow
	Value      int64
	Data       map[string]string
}

type Meter interface {
	// Add the given value to the metric
	Add(id string, name string, data map[string]string, tw TimeWindow, value int64) error
	// Return the metric and its assigned quota
	Get(id string, name string, data map[string]string, tw TimeWindow) (*Metric, error)
	// Return metrics for a given ID and time interval
	Range(id string, from, to int64, name, twType string) ([]*Metric, error)
	N(id string, to int64, n int64, name, twType string) ([]*Metric, error)
}

type ModelMeter interface {
	Meter
	// Add a given metric to a model (both time-based metrics and totals)
	AddToModel(model Model, name string, tw TimeWindow, value int64) error
	ModelID(model Model) string
}
