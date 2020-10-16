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
				forms.IsList{},
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
