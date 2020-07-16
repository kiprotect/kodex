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

package kiprotect

import (
	"encoding/json"
	"fmt"
	"github.com/kiprotect/go-helpers/forms"
	"github.com/kiprotect/go-helpers/maps"
)

type ActionConfig interface {
	Model
	Action() (Action, error)
	Project() Project
	SetData(interface{}) error
	Data() interface{}
	ConfigData() map[string]interface{}
	SetConfigData(map[string]interface{}) error
	Name() string
	Description() string
	ActionType() string
	SetName(string) error
	SetDescription(string) error
	SetActionType(string) error
}

type IsActionConfig struct{}

func (i IsActionConfig) Validate(value interface{}, values map[string]interface{}) (interface{}, error) {
	config, ok := maps.ToStringMap(value)
	if !ok {
		return nil, fmt.Errorf("not a string map")
	}

	if configParams, err := ActionConfigForm.Validate(config); err != nil {
		return nil, err
	} else {
		return configParams, nil
	}
}

var ActionConfigForm = forms.Form{
	ErrorMsg: "invalid data encountered in the action config",
	Fields: []forms.Field{
		{
			Name: "name",
			Validators: append([]forms.Validator{
				forms.IsRequired{}}, NameValidators...),
		},
		{
			Name: "type",
			Validators: []forms.Validator{
				forms.IsRequired{},
				forms.IsString{},
			},
		},
		{
			Name: "config",
			Validators: []forms.Validator{
				forms.IsStringMap{},
			},
		},
		{
			Name: "data",
			Validators: []forms.Validator{
				forms.IsOptional{},
				forms.IsStringMap{},
			},
		},
		{
			Name: "description",
			Validators: append([]forms.Validator{
				forms.IsOptional{Default: ""}}, DescriptionValidators...),
		},
	},
}

/* Base Functionality */

type BaseActionConfig struct {
	Self     ActionConfig
	Project_ Project
}

func (b *BaseActionConfig) Type() string {
	return "action"
}

func (b *BaseActionConfig) Project() Project {
	return b.Project_
}

func (b *BaseActionConfig) checkActionConfig(actionType string, config map[string]interface{}) error {
	definitions := b.Self.Project().Controller().Definitions()
	if definition, ok := definitions.ActionDefinitions[actionType]; !ok {
		return fmt.Errorf("invalid action type: %s", actionType)
	} else if _, err := definition.Maker(b.Self.Name(), b.Self.Description(), b.Self.ID(), config); err != nil {
		return err
	}
	return nil
}

func (b *BaseActionConfig) actionKeys() []string {
	actionKeys := make([]string, 0)
	definitions := b.Self.Project().Controller().Definitions()
	for key, _ := range definitions.ActionDefinitions {
		actionKeys = append(actionKeys, key)
	}
	return actionKeys
}

func (b *BaseActionConfig) Update(values map[string]interface{}) error {

	if params, err := ActionConfigForm.ValidateUpdate(values); err != nil {
		return err
	} else {
		actionType, ok := params["type"].(string)
		if !ok {
			actionType = b.Self.ActionType()
		}
		actionConfig, ok := params["config"].(map[string]interface{})
		if !ok {
			actionConfig = b.Self.ConfigData()
		}
		// we validate the new config/type combination
		if err := b.checkActionConfig(actionType, actionConfig); err != nil {
			return err
		}
		return b.update(params)
	}

}

func (b *BaseActionConfig) Create(values map[string]interface{}) error {

	if params, err := ActionConfigForm.Validate(values); err != nil {
		return err
	} else {
		actionType := params["type"].(string)
		actionConfig := params["config"].(map[string]interface{})
		if err := b.checkActionConfig(actionType, actionConfig); err != nil {
			return err
		}
		return b.update(params)
	}

}

func (b *BaseActionConfig) Action() (Action, error) {
	definitions := b.Self.Project().Controller().Definitions()
	return MakeAction(b.Self.Name(), b.Self.Description(), b.Self.ActionType(), b.Self.ID(), b.Self.ConfigData(), definitions)
}

func (b *BaseActionConfig) update(params map[string]interface{}) error {

	for key, value := range params {
		var err error
		switch key {
		case "name":
			err = b.Self.SetName(value.(string))
		case "description":
			err = b.Self.SetDescription(value.(string))
		case "type":
			err = b.Self.SetActionType(value.(string))
		case "config":
			err = b.Self.SetConfigData(value.(map[string]interface{}))
		case "data":
			err = b.Self.SetData(value)
		}
		if err != nil {
			return err
		}
	}

	return nil

}

func (b *BaseActionConfig) MarshalJSON() ([]byte, error) {

	data := map[string]interface{}{
		"name":        b.Self.Name(),
		"description": b.Self.Description(),
		"type":        b.Self.ActionType(),
		"data":        b.Self.Data(),
		"project":     b.Self.Project(),
		"config":      b.Self.ConfigData(),
	}

	for k, v := range JSONData(b.Self) {
		data[k] = v
	}

	return json.Marshal(data)
}
