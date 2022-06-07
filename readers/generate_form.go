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

package readers

import (
	"encoding/json"
	"github.com/kiprotect/go-helpers/forms"
	"reflect"
)

func keys(m interface{}) []interface{} {
	mt := reflect.ValueOf(m)
	keys := make([]interface{}, 0)
	for _, k := range mt.MapKeys() {
		keys = append(keys, k.Interface())
	}
	return keys
}

type IsGenerator struct{}

func (g IsGenerator) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"types": []string{},
	})
}

var LiteralForm = forms.Form{
	Fields: []forms.Field{
		forms.Field{
			Name: "value",
			Validators: []forms.Validator{
				forms.IsRequired{},
			},
		},
	},
}

var TimestampForm = forms.Form{
	Fields: []forms.Field{
		forms.Field{
			Name: "format",
			Validators: []forms.Validator{
				forms.IsString{},
				forms.IsIn{Choices: []interface{}{"rfc3339", "unix"}},
			},
		},
	},
}

var GeneratorForm = forms.Form{
	ErrorMsg: "invalid data encountered in the field config form",
	Fields: []forms.Field{
		forms.Field{
			Name: "type",
			Validators: []forms.Validator{
				forms.IsString{},
				forms.IsIn{Choices: keys(generators)},
			},
		},
		forms.Field{
			Name: "config",
			Validators: []forms.Validator{
				forms.IsOptional{Default: map[string]interface{}{}},
				forms.IsStringMap{},
			},
		},
	},
}

func (g IsGenerator) Validate(input interface{}, values map[string]interface{}) (interface{}, error) {
	params := input.(map[string]interface{})
	gt := params["type"].(string)
	gc := params["config"].(map[string]interface{})
	if generator, err := generators[gt](gc); err != nil {
		return nil, err
	} else {
		return generator, nil
	}
}

var GenerateForm = forms.Form{
	ErrorMsg: "invalid data encountered in the generate reader form",
	Fields: []forms.Field{
		{
			Name: "fields",
			Validators: []forms.Validator{
				forms.IsStringMap{
					Form: &forms.Form{
						Fields: []forms.Field{
							forms.Field{
								Name: "*",
								Validators: []forms.Validator{
									forms.IsStringMap{
										Form: &GeneratorForm,
									},
									IsGenerator{},
								},
							},
						},
					},
				},
			},
		},
		{
			Name: "frequency",
			Validators: []forms.Validator{
				forms.IsOptional{Default: 1000},
				forms.IsFloat{HasMin: true, HasMax: true, Min: 1e-3, Max: 1e6},
			},
		},
	},
}
