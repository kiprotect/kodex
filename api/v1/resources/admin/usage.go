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

package admin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kiprotect/go-helpers/forms"
	"github.com/kiprotect/kodex"
	"github.com/kiprotect/kodex/api"
	"github.com/kiprotect/kodex/api/helpers"
	"regexp"
	"sort"
	"strings"
	"time"
)

func UsageValidator(values map[string]interface{}, addError forms.ErrorAdder) error {
	if values["from"] != nil && values["to"] == nil || values["to"] != nil && values["from"] == nil {
		return fmt.Errorf("both from and to must be specified")
	}
	if values["from"] != nil && values["n"] != nil {
		return fmt.Errorf("cannot specify both n and from/to")
	}
	if values["n"] == nil && values["from"] == nil {
		return fmt.Errorf("you need to specify either n or from/to")
	}
	if values["from"] != nil {
		fromT := values["from"].(time.Time)
		toT := values["to"].(time.Time)
		if fromT.UnixNano() > toT.UnixNano() {
			return fmt.Errorf("from date must be before to date")
		}
	}
	return nil
}

var UsageForm = forms.Form{
	ErrorMsg: "invalid data encountered in the usage parameter form",
	Fields: []forms.Field{
		{
			Name: "type",
			Validators: []forms.Validator{
				forms.IsRequired{},
				forms.IsIn{Choices: []interface{}{"minute", "hour", "day", "week", "month"}},
			},
		},
		{
			Name: "name",
			Validators: []forms.Validator{
				forms.IsOptional{Default: ""},
				forms.MatchesRegex{Regexp: regexp.MustCompile(`^[\w\d\-]{0,50}$`)},
			},
		},
		{
			Name: "from",
			Validators: []forms.Validator{
				forms.IsOptional{},
				forms.IsTime{Format: "rfc3339", ToUTC: true},
			},
		},
		{
			Name: "to",
			Validators: []forms.Validator{
				forms.IsOptional{},
				forms.IsTime{Format: "rfc3339", ToUTC: true},
			},
		},
		{
			Name: "n",
			Validators: []forms.Validator{
				forms.IsOptional{},
				forms.IsInteger{HasMin: true, Min: 1, HasMax: true, Max: 500, Convert: true},
			},
		},
	},
	Transforms: []forms.Transform{},
	Validator:  UsageValidator,
}

type ApiValue struct {
	Name  string            `json:"name"`
	From  time.Time         `json:"from"`
	To    time.Time         `json:"to"`
	Data  map[string]string `json:"data"`
	Value int64             `json:"value"`
}

type Values struct {
	values []*ApiValue
}

func (f Values) Len() int {
	return len(f.values)
}

func (f Values) Less(i, j int) bool {
	r := (f.values[i].From).Sub(f.values[j].From)
	if r < 0 {
		return true
	}
	// if the from times match we compare the names
	if r == 0 {
		if strings.Compare(f.values[i].Name, f.values[j].Name) < 0 {
			return true
		}
	}
	return false
}

func (f Values) Swap(i, j int) {
	f.values[i], f.values[j] = f.values[j], f.values[i]

}

func UsageEndpoint(meter kodex.Meter, meterId string) func(*gin.Context) {

	return func(c *gin.Context) {

		config := helpers.QueryToConfig(c.Request.URL.Query())
		params, err := UsageForm.Validate(config)

		if err != nil {
			api.HandleError(c, 400, err)
			return
		}
		idObj, ok := c.Get(meterId)
		if !ok {
			api.HandleError(c, 500, fmt.Errorf("not meter ID defined"))
			return
		}
		id, ok := idObj.(string)
		if !ok {
			api.HandleError(c, 500, fmt.Errorf("invalid meter ID"))
			return
		}

		toTime := time.Now().UTC().UnixNano()

		var metrics []*kodex.Metric

		if params["n"] != nil {
			metrics, err = meter.N(id, toTime, params["n"].(int64), params["name"].(string), params["type"].(string))
		} else {
			fromT := params["from"].(time.Time)
			toT := params["to"].(time.Time)
			metrics, err = meter.Range(id, fromT.UnixNano(), toT.UnixNano(), params["name"].(string), params["type"].(string))
		}

		if err != nil {
			api.HandleError(c, 500, err)
			return
		}

		values := make([]*ApiValue, 0)

		for _, metric := range metrics {
			if metric.Name[0] == '_' {
				// we skip internal metrics (which start with a '_')
				continue
			}

			values = append(values, &ApiValue{
				From:  time.Unix(metric.TimeWindow.From/1e9, metric.TimeWindow.From%1e9).UTC(),
				To:    time.Unix(metric.TimeWindow.To/1e9, metric.TimeWindow.From%1e9).UTC(),
				Name:  metric.Name,
				Value: metric.Value,
				Data:  metric.Data,
			})
		}

		sortableValues := Values{values: values}

		sort.Sort(sortableValues)

		c.JSON(200, gin.H{"data": values, "params": params})
	}

}
