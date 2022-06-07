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
	"time"
)

type GeneralizeAction struct {
	kodex.BaseAction
	config *GeneralizeConfig
}

var GeneralizeForm = forms.Form{
	ErrorMsg: "invalid data encountered in the generalize form",
	Fields: []forms.Field{
		{
			Name: "key",
			Validators: []forms.Validator{
				forms.IsOptional{Default: "_"},
				forms.IsString{},
			},
		},
		{
			Name: "type",
			Validators: []forms.Validator{
				forms.IsIn{Choices: []interface{}{"datetime"}},
			},
		},
		{
			Name: "input-format",
			Validators: []forms.Validator{
				forms.IsString{},
			},
		},
		{
			Name: "output-format",
			Validators: []forms.Validator{
				forms.IsString{},
			},
		},
	},
}

type GeneralizeConfig struct {
	Key          string `json:"key"`
	Type         string `json:"type"`
	InputFormat  string `json:"input-format"`
	OutputFormat string `json:"output-format"`
}

func MakeGeneralizeAction(spec kodex.ActionSpecification) (kodex.Action, error) {

	generalizeConfig := &GeneralizeConfig{}

	if params, err := GeneralizeForm.Validate(spec.Config); err != nil {
		return nil, err
	} else if err := GeneralizeForm.Coerce(generalizeConfig, params); err != nil {
		return nil, err
	} else {
		return &GeneralizeAction{
			BaseAction: kodex.MakeBaseAction(spec, "generalize"),
			config:     generalizeConfig,
		}, nil
	}
}

func (a *GeneralizeAction) Params() interface{} {
	return nil
}

func (a *GeneralizeAction) GenerateParams(key, salt []byte) error {
	return nil
}

func (a *GeneralizeAction) SetParams(params interface{}) error {
	return nil
}

func (a *GeneralizeAction) Do(item *kodex.Item, writer kodex.ChannelWriter) (*kodex.Item, error) {

	v, ok := item.Get(a.config.Key)

	if !ok {
		// key is missing
		return nil, fmt.Errorf("key missing")
	}

	if a.config.Type == "datetime" {
		s, ok := v.(string)

		if !ok {
			return nil, fmt.Errorf("expected a string value")
		}

		inputTime, err := time.Parse(a.config.InputFormat, s)

		if err != nil {
			return nil, err
		}

		item.Set(a.config.Key, inputTime.Format(a.config.OutputFormat))

	}

	return item, nil

}
