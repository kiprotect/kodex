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

package anonymize

import (
	"fmt"
	"github.com/kiprotect/go-helpers/errors"
	"github.com/kiprotect/go-helpers/forms"
	"github.com/kiprotect/kodex/actions/anonymize/aggregate"
	"github.com/kiprotect/kodex/actions/anonymize/aggregate/functions"
	"github.com/kiprotect/kodex/actions/anonymize/aggregate/group_by_functions"
)

type Function struct {
	Function aggregate.Function
	Name     string
	Config   map[string]interface{}
}

type IsFunction struct{}

func (i IsFunction) Validate(input interface{}, values map[string]interface{}) (interface{}, error) {
	name, ok := input.(string)
	if !ok {
		return nil, fmt.Errorf("expected a string")
	}
	f, ok := functions.Functions[name]
	if !ok {
		return nil, errors.MakeExternalError("unknown function", "AGGREGATE", name, nil)
	}
	// function config has already been validated
	config := values["config"].(map[string]interface{})
	function, err := f(config)
	if err != nil {
		return nil, errors.MakeExternalError("cannot initialize function", "AGGREGATE", name, err)
	}
	return Function{
		Name:     name,
		Function: function,
		Config:   config,
	}, nil
}

func timeWindowValues() []interface{} {
	values := make([]interface{}, 0)
	for key, _ := range groupByFunctions.TimeWindowFunctions {
		values = append(values, key)
	}
	return values
}

var TimeWindowForm = forms.Form{
	Fields: []forms.Field{
		{
			Name: "field",
			Validators: []forms.Validator{
				forms.IsString{},
			},
		},
		{
			Name: "window",
			Validators: []forms.Validator{
				forms.IsString{},
				forms.IsIn{Choices: timeWindowValues()},
			},
		},
	},
}

var GroupByForm = forms.Form{
	ErrorMsg: "invalid data encountered in the group-by-form",
	Fields: []forms.Field{
		{
			Name: "function",
			Validators: []forms.Validator{
				forms.IsIn{Choices: []interface{}{"time-window"}},
			},
		},
		{
			Name: "config",
			Validators: []forms.Validator{
				forms.IsStringMap{},
				forms.Switch{
					Key: "function",
					Cases: map[string][]forms.Validator{
						"time-window": {
							forms.IsStringMap{
								Form: &TimeWindowForm,
							},
						},
					},
				},
			},
		},
	},
}

var AggregateForm = forms.Form{
	ErrorMsg: "invalid data encountered in the aggregation config",
	Fields: []forms.Field{
		{
			Name: "config",
			Validators: []forms.Validator{
				forms.IsOptional{Default: map[string]interface{}{}},
				forms.IsStringMap{},
			},
		},
		{
			Name: "function",
			Validators: []forms.Validator{
				forms.IsRequired{},
				forms.IsString{},
				IsFunction{},
			},
		},
		{
			Name: "group-by",
			Validators: []forms.Validator{
				forms.IsRequired{},
				forms.IsList{
					Validators: []forms.Validator{
						forms.IsStringMap{
							Form: &GroupByForm,
						},
					},
				},
			},
		},
		{
			Name: "result-name",
			Validators: []forms.Validator{
				forms.IsOptional{},
				forms.IsString{},
			},
		},
		{
			Name: "finalize-after",
			Validators: []forms.Validator{
				forms.IsOptional{Default: 300},
				forms.IsInteger{Min: -1, HasMin: true},
			},
		},
		{
			Name: "channels",
			Validators: []forms.Validator{
				forms.IsOptional{Default: []string{}},
				forms.IsStringList{},
			},
		},
	},
}
