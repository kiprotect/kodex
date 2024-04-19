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

package metering

import (
	"github.com/kiprotect/go-helpers/forms"
	"github.com/kiprotect/kodex"
)

type InMemoryMeter struct {
	BaseModelMeter
}

var InMemoryMeterForm = forms.Form{
	ErrorMsg: "invalid data encountered in the Redis config form",
	Fields:   []forms.Field{},
}

func MakeInMemoryMeter(config map[string]interface{}) (*InMemoryMeter, error) {

	_, err := InMemoryMeterForm.Validate(config)
	if err != nil {
		return nil, err
	}

	meter := &InMemoryMeter{}
	meter.BaseModelMeter.Self = meter
	return meter, nil

}

// Add the given value to the metric
func (i *InMemoryMeter) Add(id string, name string, data map[string]string, tw kodex.TimeWindow, value int64) error {
	return nil
}

// Return the metric and its assigned quota
func (i *InMemoryMeter) Get(id string, name string, data map[string]string, tw kodex.TimeWindow) (*kodex.Metric, error) {
	return nil, nil
}

// Return metrics for a given ID and time interval
func (i *InMemoryMeter) Range(id string, from, to int64, name, twType string) ([]*kodex.Metric, error) {
	return nil, nil
}

func (i *InMemoryMeter) N(id string, to int64, n int64, name, twType string) ([]*kodex.Metric, error) {
	return nil, nil
}

func (i *InMemoryMeter) AddToModel(model kodex.Model, name string, tw kodex.TimeWindow, value int64) error {
	return nil
}

func (i *InMemoryMeter) ModelID(model kodex.Model) string {
	return ""
}
