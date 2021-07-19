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
)

type FormAction struct {
	kodex.BaseAction
	form *forms.Form
}

var Validators = map[string]forms.ValidatorMaker{}

type ParametrizedValidator interface {
	GenerateParams(key, salt []byte) error
	SetParams(params interface{}) error
	Params() interface{}
}

type IsAction struct {
	Action kodex.Action
	Type   string                 `json:"type"`
	Config map[string]interface{} `json:"config"`
}

func (i IsAction) Validate(input interface{}, values map[string]interface{}) (interface{}, error) {
	item := kodex.MakeItem(map[string]interface{}{"_": input})
	if newItem, err := i.Action.(kodex.DoableAction).Do(item, nil); err != nil {
		return nil, err
	} else {
		v, _ := newItem.Get("_")
		return v, nil
	}
}

var IsActionForm = forms.Form{
	Fields: []forms.Field{
		{
			Name: "type",
			Validators: []forms.Validator{
				forms.IsString{},
			},
		},
		{
			Name: "config",
			Validators: []forms.Validator{
				forms.IsStringMap{},
			},
		},
	},
}

func MakeFormAction(spec kodex.ActionSpecification) (kodex.Action, error) {
	combinedValidators := map[string]forms.ValidatorMaker{}

	for k, v := range Validators {
		combinedValidators[k] = v
	}

	for k, v := range forms.Validators {
		combinedValidators[k] = v
	}

	makeIsAction := func(config map[string]interface{}, context *forms.FormDescriptionContext) (forms.Validator, error) {
		isAction := &IsAction{}
		if params, err := IsActionForm.Validate(config); err != nil {
			return nil, err
		} else if err := IsActionForm.Coerce(isAction, params); err != nil {
			return nil, err
		}
		if action, err := kodex.MakeAction(spec.Name, spec.Description, isAction.Type, spec.ID, isAction.Config, spec.Definitions); err != nil {
			return nil, err
		} else if _, ok := action.(kodex.DoableAction); !ok {
			return nil, fmt.Errorf("undoable action")
		} else {
			isAction.Action = action
		}
		return isAction, nil
	}

	combinedValidators["IsAction"] = makeIsAction

	context := &forms.FormDescriptionContext{
		Validators: combinedValidators,
	}
	if form, err := forms.FromConfig(spec.Config, context); err != nil {
		return nil, err
	} else {
		return &FormAction{
			BaseAction: kodex.MakeBaseAction(spec, "form"),
			form:       form,
		}, nil
	}
}

func (a *FormAction) Params() interface{} {
	return nil
}

func (a *FormAction) GenerateParams(key, salt []byte) error {
	return nil
}

func (a *FormAction) SetParams(params interface{}) error {
	return nil
}

func (a *FormAction) Do(item *kodex.Item, writer kodex.ChannelWriter) (*kodex.Item, error) {
	if params, err := a.form.Validate(item.All()); err != nil {
		return nil, err
	} else {
		return kodex.MakeItem(params), nil
	}
}
