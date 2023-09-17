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

package actions

import (
	"fmt"
	"github.com/kiprotect/go-helpers/forms"
	"github.com/kiprotect/kodex"
	"github.com/kiprotect/kodex/actions/pseudonymize"
)

type PseudonymizeTransformation struct {
	kodex.BaseAction
	Pseudonymizer pseudonymize.Pseudonymizer
	Key           string
	Method        string
}

func (p *PseudonymizeTransformation) Undoable(item *kodex.Item) bool {
	return true
}

func (p *PseudonymizeTransformation) process(item *kodex.Item, writer kodex.ChannelWriter, f func(interface{}) (interface{}, error)) (*kodex.Item, error) {
	value, ok := item.Get(p.Key)
	if !ok {
		return nil, fmt.Errorf("key %s missing", p.Key)
	}
	newValue, err := f(value)
	if err != nil {
		return nil, err
	}
	item.Set(p.Key, newValue)
	return item, err
}

func (p *PseudonymizeTransformation) Undo(item *kodex.Item, writer kodex.ChannelWriter) (*kodex.Item, error) {
	return p.process(item, writer, p.Pseudonymizer.Depseudonymize)
}

func (p *PseudonymizeTransformation) Do(item *kodex.Item, writer kodex.ChannelWriter) (*kodex.Item, error) {
	return p.process(item, writer, p.Pseudonymizer.Pseudonymize)
}

func (p *PseudonymizeTransformation) GenerateParams(key, salt []byte) error {
	return p.Pseudonymizer.GenerateParams(key, salt)
}

func (p *PseudonymizeTransformation) SetParams(params interface{}) error {
	return p.Pseudonymizer.SetParams(params)
}

func (p *PseudonymizeTransformation) Params() interface{} {
	return p.Pseudonymizer.Params()
}

func MakePseudonymizeAction(spec kodex.ActionSpecification) (kodex.Action, error) {

	if spec.Config == nil || len(spec.Config) == 0 {
		spec.Config = map[string]any{
			"method": "merengue",
			"config": map[string]any{},
		}
	}

	params, err := PseudonymizeConfigForm.Validate(spec.Config)

	if err != nil {
		return nil, err
	}

	method := params["method"].(string)

	psMaker, ok := pseudonymize.Pseudonymizers[method]

	if !ok {
		return nil, fmt.Errorf("Unknown pseudonymizer method %s", method)
	}

	var mapConfig map[string]any

	if _, ok := spec.Config["config"]; !ok {
		// we fall back to the parent config...
		mapConfig = spec.Config
	} else {
		// we convert the config to a map
		mapConfig, _ = params["config"].(map[string]any)
	}

	ps, err := psMaker(mapConfig)

	if err != nil {
		return nil, err
	}

	return &PseudonymizeTransformation{
		Pseudonymizer: ps,
		Method:        method,
		Key:           params["key"].(string),
		BaseAction:    kodex.MakeBaseAction(spec, "pseudonymize"),
	}, nil

}

var PseudonymizeConfigForm = forms.Form{
	ErrorMsg: "invalid data encountered in the 'pseudonymize' form",
	Fields: []forms.Field{
		{
			Name:        "method",
			Description: "The pseudonymization method to use. Structured pseudonymization will preserve the data format and (partial) structure of the input data when pseudonymizing. Merengue pseudonymization will produce unstructured pseudonyms instead.",
			Validators: []forms.Validator{
				forms.IsRequired{},
				forms.IsIn{
					Choices: []interface{}{"merengue", "structured"},
				},
			},
		},
		{
			Name:        "key",
			Description: "The key of the attribute to pseudonymize ('_' by default).",
			Validators: []forms.Validator{
				forms.IsOptional{Default: "_"},
				forms.IsString{},
			},
		},
		{
			Name: "config",
			Validators: []forms.Validator{
				forms.IsOptional{},
				forms.Switch{
					Key: "method",
					Cases: map[string][]forms.Validator{
						"merengue": []forms.Validator{
							forms.IsStringMap{
								Form: &pseudonymize.MerengueConfigForm,
							},
						},
						"structured": []forms.Validator{
							forms.IsStringMap{
								Form: &pseudonymize.StructuredPseudonymizerForm,
							},
						},
					},
				},
			},
		},
	},
}
