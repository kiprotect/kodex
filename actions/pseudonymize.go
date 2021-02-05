// Kodex (Community Edition - CE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2021  KIProtect GmbH (HRB 208395B) - Germany
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
	pseudonymizer pseudonymize.Pseudonymizer
	key           string
	method        string
	config        map[string]interface{}
}

func (p *PseudonymizeTransformation) Undoable(item *kodex.Item) bool {
	return true
}

func (p *PseudonymizeTransformation) process(item *kodex.Item, writer kodex.ChannelWriter, f func(interface{}) (interface{}, error)) (*kodex.Item, error) {
	value, ok := item.Get(p.key)
	if !ok {
		return nil, fmt.Errorf("key %s missing", p.key)
	}
	newValue, err := f(value)
	if err != nil {
		return nil, err
	}
	item.Set(p.key, newValue)
	return item, err
}

func (p *PseudonymizeTransformation) Undo(item *kodex.Item, writer kodex.ChannelWriter) (*kodex.Item, error) {
	return p.process(item, writer, p.pseudonymizer.Depseudonymize)
}

func (p *PseudonymizeTransformation) Do(item *kodex.Item, writer kodex.ChannelWriter) (*kodex.Item, error) {
	return p.process(item, writer, p.pseudonymizer.Pseudonymize)
}

func (p *PseudonymizeTransformation) GenerateParams(key, salt []byte) error {
	return p.pseudonymizer.GenerateParams(key, salt)
}

func (p *PseudonymizeTransformation) SetParams(params interface{}) error {
	return p.pseudonymizer.SetParams(params)
}

func (p *PseudonymizeTransformation) Params() interface{} {
	return p.pseudonymizer.Params()
}

func MakePseudonymizeAction(name, description string, id []byte, config map[string]interface{}) (kodex.Action, error) {

	params, err := PseudonymizeConfigForm.Validate(config)

	if err != nil {
		return nil, err
	}

	method := params["method"].(string)

	psMaker, ok := pseudonymize.Pseudonymizers[method]

	if !ok {
		return nil, fmt.Errorf("Unknown pseudonymizer method %s", method)
	}

	ps, err := psMaker(config)

	if err != nil {
		return nil, err
	}

	return &PseudonymizeTransformation{
		pseudonymizer: ps,
		method:        method,
		key:           params["key"].(string),
		BaseAction:    kodex.MakeBaseAction(name, description, "pseudonymize", id, config),
	}, nil

}

var PseudonymizeConfigForm = forms.Form{
	ErrorMsg: "invalid data encountered in the 'pseudonymize' form",
	Fields: []forms.Field{
		{
			Name: "method",
			Validators: []forms.Validator{
				forms.IsRequired{},
				forms.IsString{},
				forms.IsIn{
					Choices: []interface{}{"merengue", "structured"},
				},
			},
		},
		{
			Name: "key",
			Validators: []forms.Validator{
				forms.IsRequired{},
				forms.IsString{},
			},
		},
	},
}
