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
)

type DestinationStatus string

const (
	ActiveDestination   DestinationStatus = "active"
	OnDemandDestination DestinationStatus = "on-demand"
	DisabledDestination DestinationStatus = "disabled"
	TestingDestination  DestinationStatus = "testing"
	ErrorDestination    DestinationStatus = "error"
	WarningDestination  DestinationStatus = "warning"
	MessageDestination  DestinationStatus = "message"
)

type Destination interface {
	Processable // Processable includes Model
	Writer() (Writer, error)
	ConfigData() map[string]interface{}
	SetConfigData(map[string]interface{}) error
	Name() string
	DestinationType() string
	SetDestinationType(string) error
	SetName(string) error
	Description() string
	SetDescription(string) error
	SetData(interface{}) error
	Data() interface{}
	Project() Project
}

/* Base Functionality */

type BaseDestination struct {
	Self     Destination
	Project_ Project
	writer   Writer
}

func (b *BaseDestination) Type() string {
	return "destination"
}

func (b *BaseDestination) checkDestinationConfig(destinationType string, config map[string]interface{}) error {
	definitions := b.Self.Project().Controller().Definitions()
	if definition, ok := definitions.WriterDefinitions[destinationType]; !ok {
		return fmt.Errorf("invalid destination type: %s", destinationType)
	} else if _, err := definition.Maker(config); err != nil {
		return err
	}
	return nil
}

func (b *BaseDestination) Update(values map[string]interface{}) error {

	if params, err := DestinationForm.ValidateUpdate(values); err != nil {
		return err
	} else {
		destinationType, ok := params["type"].(string)
		if !ok {
			destinationType = b.Self.DestinationType()
		}
		destinationConfig, ok := params["config"].(map[string]interface{})
		if !ok {
			destinationConfig = b.Self.ConfigData()
		}
		// we validate the new config/type combination
		if err := b.checkDestinationConfig(destinationType, destinationConfig); err != nil {
			return err
		}
		return b.update(params)
	}

}

func (b *BaseDestination) Create(values map[string]interface{}) error {

	if params, err := DestinationForm.Validate(values); err != nil {
		return err
	} else {
		destinationType := params["type"].(string)
		destinationConfig := params["config"].(map[string]interface{})
		if err := b.checkDestinationConfig(destinationType, destinationConfig); err != nil {
			return err
		}
		return b.update(params)
	}

}

func (b *BaseDestination) MarshalJSON() ([]byte, error) {

	data := map[string]interface{}{
		"name":        b.Self.Name(),
		"description": b.Self.Description(),
		"type":        b.Self.DestinationType(),
		"config":      b.Self.ConfigData(),
		"project":     b.Self.Project(),
		"data":        b.Self.Data(),
	}

	for k, v := range JSONData(b.Self) {
		data[k] = v
	}

	return json.Marshal(data)
}

func (b *BaseDestination) update(params map[string]interface{}) error {

	for key, value := range params {
		var err error
		switch key {
		case "name":
			err = b.Self.SetName(value.(string))
		case "description":
			err = b.Self.SetDescription(value.(string))
		case "config":
			err = b.Self.SetConfigData(value.(map[string]interface{}))
		case "data":
			err = b.Self.SetData(value)
		case "type":
			err = b.Self.SetDestinationType(value.(string))
		}
		if err != nil {
			return err
		}
	}

	return nil

}

func (b *BaseDestination) createWriter() (Writer, error) {
	definitions := b.Self.Project().Controller().Definitions()
	destinationType := b.Self.DestinationType()
	config := b.Self.ConfigData()
	definition, ok := definitions.WriterDefinitions[destinationType]
	if !ok {
		return nil, fmt.Errorf("unknown writer type: %s", destinationType)
	}

	return definition.Maker(config)

}

func (b *BaseDestination) Writer() (Writer, error) {
	if b.writer != nil {
		return b.writer, nil
	}
	writer, err := b.createWriter()
	if err != nil {
		return nil, err
	}
	b.writer = writer
	return b.writer, nil
}

var DestinationForm = forms.Form{
	ErrorMsg: "invalid data encountered in the destination form",
	Fields: []forms.Field{
		{
			Name: "name",
			Validators: append([]forms.Validator{
				forms.IsRequired{}}, NameValidators...),
		},
		{
			Name: "description",
			Validators: append([]forms.Validator{
				forms.IsOptional{Default: ""}}, DescriptionValidators...),
		},
		{
			Name:       "config",
			Validators: []forms.Validator{forms.IsStringMap{}},
		},
		{
			Name:       "data",
			Validators: []forms.Validator{forms.IsOptional{}, forms.IsStringMap{}},
		},
		{
			Name:       "type",
			Validators: []forms.Validator{forms.IsString{}},
		},
	},
}

func (b *BaseDestination) Project() Project {
	return b.Project_
}
