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
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/kiprotect/go-helpers/forms"
)

type ActionDefinition struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Internal    bool        `json:"internal"`
	Maker       ActionMaker `json:"-"`
	Data        interface{} `json:"data"`
	Form        forms.Form  `json:"form"`
}

type ActionMaker func(name, description string, id []byte, config map[string]interface{}) (Action, error)
type ActionDefinitions map[string]ActionDefinition

type Action interface {
	ConfigGroup() map[string]interface{}
	ConfigHash() ([]byte, error)
	ParameterGroup(item *Item) (*ParameterGroup, error)
	Params() interface{}
	HasParams() bool
	SetParams(interface{}) error
	GenerateParams(key, salt []byte) error
	ID() []byte
	Name() string
	Description() string
	Type() string
	Config() map[string]interface{}
	Setup() error
	Teardown() error
}

type Schedule struct {
}

type DoableAction interface {
	Do(*Item, ChannelWriter) (*Item, error)
}

type ConfigurableAction interface {
	DoWithConfig(*Item, ChannelWriter, Config) (*Item, error)
}

type ScheduledAction interface {
	// Get notified about a callback
	Callback(interface{}) error
	// Schedule a callback
	Schedule(Schedule) error
}

type StatefulAction interface {
	// Finalizes the action
	Finalize() error
}

type UndoableAction interface {
	Undoable(*Item) bool
	Undo(*Item, ChannelWriter) (*Item, error)
}

/* Base Functionality */

type BaseAction struct {
	Name_        string
	Description_ string
	Type_        string
	ID_          []byte
	Config_      map[string]interface{}
	configHash   []byte
}

type ActionSpecification struct {
	Name, Description, Type string
	ID                      []byte
	Config                  map[string]interface{}
}

var ActionSpecificationForm = forms.Form{
	Fields: []forms.Field{
		forms.Field{
			Name: "name",
			Validators: []forms.Validator{
				forms.IsOptional{Default: ""},
				forms.IsString{},
			},
		},
		forms.Field{
			Name: "description",
			Validators: []forms.Validator{
				forms.IsOptional{Default: ""},
				forms.IsString{},
			},
		},
		forms.Field{
			Name: "id",
			Validators: []forms.Validator{
				forms.IsOptional{
					DefaultGenerator: func() interface{} { return RandomID() },
				},
				forms.IsBytes{},
			},
		},
		forms.Field{
			Name: "type",
			Validators: []forms.Validator{
				forms.IsRequired{},
				forms.IsString{},
			},
		},
		forms.Field{
			Name: "config",
			Validators: []forms.Validator{
				forms.IsRequired{},
				forms.IsStringMap{},
			},
		},
	},
}

func MakeAction(name, description, actionType string, id []byte, config map[string]interface{}, definitions Definitions) (Action, error) {
	if definition, ok := definitions.ActionDefinitions[actionType]; !ok {
		return nil, fmt.Errorf("unknown action type: %s", actionType)
	} else {
		return definition.Maker(name, description, id, config)
	}
}

func MakeActions(specs []ActionSpecification, definitions Definitions) ([]Action, error) {
	actions := make([]Action, len(specs))
	for i, specification := range specs {
		actionDefinition, ok := definitions.ActionDefinitions[specification.Type]
		if !ok {
			return nil, fmt.Errorf("unknown action type: %s", specification.Type)
		}
		action, err := actionDefinition.Maker(specification.Name, specification.Description, specification.ID, specification.Config)
		if err != nil {
			return nil, err
		}
		actions[i] = action
	}
	return actions, nil
}

func MakeBaseAction(name, description, actionType string, id []byte, config map[string]interface{}) (BaseAction, error) {
	return BaseAction{
		Description_: description,
		Name_:        name,
		Type_:        actionType,
		ID_:          id,
		Config_:      config,
	}, nil
}

func (b *BaseAction) HasParams() bool {
	return true
}

func (b *BaseAction) ConfigHash() ([]byte, error) {
	if b.configHash != nil {
		return b.configHash, nil
	}
	configHash, err := StructuredHash(map[string]interface{}{
		"config": b.ConfigGroup(),
		"type":   b.Type(),
	})
	if err != nil {
		return nil, err
	}
	b.configHash = configHash
	return configHash, nil
}

func (b *BaseAction) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"name":        b.Name_,
		"description": b.Description_,
		"type":        b.Type_,
		"id":          hex.EncodeToString(b.ID_),
		"config":      b.Config_,
	})
}

func (b *BaseAction) Setup() error {
	return nil
}

func (b *BaseAction) Teardown() error {
	return nil
}

func (b *BaseAction) ConfigGroup() map[string]interface{} {
	return b.Config_
}

func (b *BaseAction) Type() string {
	return b.Type_
}

func (b *BaseAction) Name() string {
	return b.Name_
}

func (b *BaseAction) ID() []byte {
	return b.ID_
}

func (b *BaseAction) Config() map[string]interface{} {
	return b.Config_
}

func (b *BaseAction) Description() string {
	return b.Description_
}

// Returns the parameter group for a specific item
func (b *BaseAction) ParameterGroup(item *Item) (*ParameterGroup, error) {
	return &ParameterGroup{
		hash: []byte("default"),
		data: map[string]interface{}{},
	}, nil
}
