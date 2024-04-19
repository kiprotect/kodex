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

package actions

import (
	"fmt"
	"github.com/kiprotect/go-helpers/forms"
	"github.com/kiprotect/kodex"
)

type FormAction struct {
	kodex.BaseAction
	actions []kodex.Action
	context *forms.FormDescriptionContext
	form    *forms.Form
}

var Validators = map[string]forms.ValidatorDefinition{}

type IsAction struct {
	Action kodex.Action   `json:"-"`
	Type   string         `json:"type"`
	Config map[string]any `json:"config"`
}

func (f *FormAction) Context() *forms.FormDescriptionContext {
	return f.context
}

func (i IsAction) Validate(input interface{}, values map[string]interface{}) (interface{}, error) {
	item := kodex.MakeItem(map[string]interface{}{"_": input})
	if newItem, err := i.Action.(kodex.DoableAction).Do(item, nil); err != nil {
		return nil, err
	} else if newItem != nil {
		v, _ := newItem.Get("_")
		return v, nil
	} else {
		return nil, nil
	}
}

var FormForm = forms.FormForm

var IsActionForm = forms.Form{
	Fields: []forms.Field{
		{
			Name: "type",
			Validators: []forms.Validator{
				forms.IsOptional{Default: "drop"},
				forms.IsString{},
			},
		},
		{
			Name: "config",
			Validators: []forms.Validator{
				forms.IsOptional{},
				forms.IsStringMap{},
			},
		},
	},
}

func MakeFormAction(spec kodex.ActionSpecification) (kodex.Action, error) {

	combinedValidators := map[string]forms.ValidatorDefinition{}

	for k, v := range Validators {
		combinedValidators[k] = v
	}

	for k, v := range forms.Validators {
		combinedValidators[k] = v
	}

	actions := make([]kodex.Action, 0)

	makeIsAction := func(config map[string]interface{}, context *forms.FormDescriptionContext) (forms.Validator, error) {
		isAction := &IsAction{}
		if params, err := IsActionForm.Validate(config); err != nil {
			return nil, fmt.Errorf("error validating action form: %v", err)
		} else if err := IsActionForm.Coerce(isAction, params); err != nil {
			return nil, fmt.Errorf("error coercing action form: %v", err)
		}
		// to do: better action name (?)
		if action, err := kodex.MakeAction(spec.Name, spec.Description, isAction.Type, spec.ID, isAction.Config, spec.Definitions); err != nil {
			return nil, fmt.Errorf("error making action: %v", err)
		} else if _, ok := action.(kodex.DoableAction); !ok {
			return nil, fmt.Errorf("undoable action")
		} else {
			isAction.Action = action

			// we append the action to the list of actions
			actions = append(actions, action)

		}

		return isAction, nil
	}

	combinedValidators["IsAction"] = forms.ValidatorDefinition{makeIsAction, IsActionForm}

	context := &forms.FormDescriptionContext{
		Validators: combinedValidators,
	}
	if form, err := forms.FromConfig(spec.Config, context); err != nil {
		return nil, err
	} else {
		return &FormAction{
			BaseAction: kodex.MakeBaseAction(spec, "form"),
			actions:    actions,
			context:    context,
			form:       form,
		}, nil
	}
}

func (a *FormAction) Form() *forms.Form {
	return a.form
}

func (a *FormAction) Params() interface{} {

	actionParams := make([]interface{}, 0, len(a.actions))

	for _, action := range a.actions {
		actionParams = append(actionParams, action.Params())
	}

	return actionParams
}

func (a *FormAction) GenerateParams(key, salt []byte) error {

	for _, action := range a.actions {
		if err := action.GenerateParams(key, salt); err != nil {
			return err
		}
	}

	return nil
}

func (a *FormAction) SetParams(params interface{}) error {

	paramsList, ok := params.([]interface{})

	if !ok {
		return fmt.Errorf("expected a list of parameters")
	}

	if len(paramsList) != len(a.actions) {
		return fmt.Errorf("action parameter list mismatch")
	}

	for i, action := range a.actions {

		if err := action.SetParams(paramsList[i]); err != nil {
			return err
		}
	}

	return nil
}

func (a *FormAction) Do(item *kodex.Item, writer kodex.ChannelWriter) (*kodex.Item, error) {
	if params, err := a.form.Validate(item.All()); err != nil {
		return nil, err
	} else {
		return kodex.MakeItem(params), nil
	}
}
