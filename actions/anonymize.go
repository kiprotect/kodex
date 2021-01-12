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
	"github.com/kiprotect/kodex/actions/anonymize"
)

type AnonymizeAction struct {
	kodex.BaseAction
	anonymizer anonymize.Anonymizer
	key        string
	method     string
	config     map[string]interface{}
}

func (p *AnonymizeAction) Undoable(interface{}) bool {
	return true
}

func (p *AnonymizeAction) process(item *kodex.Item, f func(interface{}) (interface{}, error)) (*kodex.Item, error) {
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

func (p *AnonymizeAction) Do(item *kodex.Item, writer kodex.ChannelWriter) (*kodex.Item, error) {
	return p.anonymizer.Anonymize(item, writer)
}

func (p *AnonymizeAction) GenerateParams(key, salt []byte) error {
	return nil
}

func (p *AnonymizeAction) SetParams(params interface{}) error {
	return nil
}

func (p *AnonymizeAction) Params() interface{} {
	return nil
}

func (p *AnonymizeAction) Setup() error {
	return p.anonymizer.Setup()
}

func (p *AnonymizeAction) Teardown() error {
	return p.anonymizer.Teardown()
}

func (p *AnonymizeAction) Finalize(writer kodex.ChannelWriter) ([]*kodex.Item, error) {
	return p.anonymizer.Finalize(writer)
}

func (p *AnonymizeAction) Reset() error {
	return p.anonymizer.Reset()
}

func MakeAnonymizeAction(name, description string, id []byte, config map[string]interface{}) (kodex.Action, error) {

	params, err := AnonymizeConfigForm.Validate(config)

	if err != nil {
		return nil, err
	}

	method := params["method"].(string)

	anonymizerMaker, ok := anonymize.Anonymizers[method]

	if !ok {
		return nil, fmt.Errorf("Unknown anonymizer method %s", method)
	}

	anonymizer, err := anonymizerMaker(name, id, config)

	if err != nil {
		return nil, err
	}

	if baseAction, err := kodex.MakeBaseAction(name, description, "anonymize", id, config); err != nil {
		return nil, err
	} else {
		return &AnonymizeAction{
			anonymizer: anonymizer,
			method:     method,
			BaseAction: baseAction,
		}, nil
	}

}

var AnonymizeConfigForm = forms.Form{
	ErrorMsg: "invalid data encountered in the 'anonymize' form",
	Fields: []forms.Field{
		{
			Name: "method",
			Validators: []forms.Validator{
				forms.IsRequired{},
				forms.IsString{},
				forms.IsIn{
					Choices: []interface{}{"aggregate"},
				},
			},
		},
	},
}
