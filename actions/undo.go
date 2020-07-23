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

package actions

import (
	"encoding/hex"
	"fmt"
	"github.com/kiprotect/go-helpers/forms"
	"github.com/kiprotect/kiprotect"
)

var UndoActionConfigForm = forms.Form{
	ErrorMsg: "invalid data encountered in the undo form",
	Fields: []forms.Field{
		forms.Field{
			Name: "actions",
			Validators: []forms.Validator{
				forms.IsOptional{
					Default: []kiprotect.ActionSpecification{},
				},
				forms.IsList{
					Validators: []forms.Validator{
						forms.IsStringMap{},
						kiprotect.IsActionSpecification{},
					},
				},
				kiprotect.IsActionSpecifications{},
			},
		},
	},
}

type UndoAction struct {
	kiprotect.BaseAction
	actionSpecs []kiprotect.ActionSpecification
	key, salt   []byte
}

func MakeUndoAction(name, description string, id []byte, config map[string]interface{}) (kiprotect.Action, error) {
	params, err := UndoActionConfigForm.Validate(config)
	if err != nil {
		return nil, err
	}
	if baseAction, err := kiprotect.MakeBaseAction(name, description, "undo", id, config); err != nil {
		return nil, err
	} else {
		return &UndoAction{
			BaseAction:  baseAction,
			actionSpecs: params["actions"].([]kiprotect.ActionSpecification),
		}, nil
	}
}

func (a *UndoAction) HasParams() bool {
	return false
}

func (a *UndoAction) Params() interface{} {
	return nil
}

func (a *UndoAction) GenerateParams(key, salt []byte) error {
	a.key, a.salt = key, salt
	return nil
}

func (a *UndoAction) SetParams(params interface{}) error {
	return nil
}

func (a *UndoAction) DoWithConfig(item *kiprotect.Item, writer kiprotect.ChannelWriter, config kiprotect.Config) (*kiprotect.Item, error) {
	var processor *kiprotect.Processor
	if a.key != nil {
		definitions := config.Stream().Project().Controller().Definitions()
		actions, err := kiprotect.MakeActions(a.actionSpecs, definitions)
		if err != nil {
			return nil, err
		}
		parameterSet, err := kiprotect.MakeParameterSet(actions, nil)
		if err != nil {
			return nil, err
		}
		processor, err = kiprotect.MakeProcessor(parameterSet, writer, config)
		processor.SetKey(a.key)
		processor.SetSalt(a.salt)
	} else {
		kipId, ok := item.Get("_kip")
		if !ok {
			return nil, fmt.Errorf("no parameter ID found")
		}
		kipIdStr, ok := kipId.(string)
		if !ok {
			return nil, fmt.Errorf("parameter ID is not a string")
		}
		kipIdBytes, err := hex.DecodeString(kipIdStr)
		if err != nil {
			return nil, fmt.Errorf("not a hex string")
		}
		parameterStore := config.Stream().Project().Controller().ParameterStore()
		parameterSet, err := parameterStore.ParameterSet(kipIdBytes)
		if err != nil || parameterSet == nil {
			return nil, fmt.Errorf("parameter set not found %s", kipIdStr)
		}
		processor, err = kiprotect.MakeProcessor(parameterSet, writer, config)
		if err != nil {
			return nil, err
		}
	}
	if err := processor.Setup(); err != nil {
		return nil, err
	}
	if newItems, err := processor.Undo([]*kiprotect.Item{item}, nil); err != nil {
		return nil, err
	} else {
		if len(newItems) == 1 {
			return newItems[0], nil
		} else {
			return nil, fmt.Errorf("expected a single item")
		}
	}
	return nil, nil
}

func (a *UndoAction) Setup() error {
	return nil
}

func (a *UndoAction) Teardown() error {
	return nil
}
