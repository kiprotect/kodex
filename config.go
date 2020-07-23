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
	"github.com/kiprotect/go-helpers/errors"
	"github.com/kiprotect/go-helpers/forms"
)

type ConfigStatus string

const (
	ActiveConfig   ConfigStatus = "active"
	DisabledConfig ConfigStatus = "disabled"
	TestingConfig  ConfigStatus = "testing"
)

type ErrorPolicy string

const (
	AbortOnError ErrorPolicy = "abort"
	ReportErrors ErrorPolicy = "report"
	IgnoreErrors ErrorPolicy = "ignore"
)

type Config interface {
	Model

	Status() ConfigStatus
	Version() string
	Description() string
	Source() string
	Name() string
	Stream() Stream
	Data() interface{}

	ChannelWriter() (ChannelWriter, error)

	SetData(interface{}) error
	SetStatus(ConfigStatus) error
	SetVersion(string) error
	SetDescription(string) error
	SetSource(string) error
	SetName(string) error

	AddDestination(Destination, string, DestinationStatus) error
	RemoveDestination(Destination) error
	Destinations() (map[string][]DestinationMap, error)

	ActionConfigs() ([]ActionConfig, error)
	AddActionConfig(ActionConfig, int) error
	RemoveActionConfig(ActionConfig) error

	Processor() (*Processor, error)
}

/* Base Functionality */

type BaseChannelWriter struct {
	config Config
}

func (b *BaseChannelWriter) Write(channel string, items []*Item) error {
	destinations, err := b.getDestinations(channel)
	if err != nil {
		return err
	}
	var lastErr error
	for _, destination := range destinations {
		writer, err := destination.Writer()
		if err != nil {
			lastErr = err
			continue
		}
		if err := writer.Setup(b.config); err != nil {
			lastErr = err
			continue
		}
		if err := writer.Write(MakeBasicPayload(items, map[string]interface{}{}, false)); err != nil {
			lastErr = err
		}
	}
	return lastErr
}

func (b *BaseChannelWriter) getDestinations(channel string) ([]Destination, error) {
	destinations := make([]Destination, 0)
	if configDestinations, err := b.config.Destinations(); err != nil {
		return nil, err
	} else {
		for name, destinationMaps := range configDestinations {
			if name == channel {
				for _, destinationMap := range destinationMaps {
					destinations = append(destinations, destinationMap.Destination())
				}
			}
		}
	}
	return destinations, nil
}

func (b *BaseChannelWriter) Message(
	item *Item,
	data map[string]interface{},
	mt MessageType) error {

	destinations, err := b.config.Destinations()

	if err != nil {
		return err
	}

	for _, destinationMaps := range destinations {

		for _, destinationMap := range destinationMaps {
			if destinationMap.Status() != MessageDestination {
				continue
			}

			var itemData map[string]interface{}
			if item != nil {
				itemData = item.All()
			}

			messageItem := map[string]interface{}{
				"type": mt,
				"item": itemData,
				"data": data,
			}

			writer, err := destinationMap.Destination().Writer()

			if err != nil {
				return err
			}

			if err := writer.Setup(b.config); err != nil {
				return err
			}

			if err := writer.Write(MakeBasicPayload([]*Item{MakeItem(messageItem)}, map[string]interface{}{}, false)); err != nil {
				return err
			}
		}

	}
	return nil
}

func (b *BaseChannelWriter) Error(item *Item, itemError error) error {
	return b.sendErrorWarning(item, itemError, ErrorDestination)
}

func (b *BaseChannelWriter) Warning(item *Item, itemError error) error {
	return b.sendErrorWarning(item, itemError, WarningDestination)
}

func (b *BaseChannelWriter) sendErrorWarning(
	item *Item,
	itemError error,
	status DestinationStatus) error {

	destinations, err := b.config.Destinations()

	if err != nil {
		return err
	}

	for _, destinationMaps := range destinations {

		for _, destinationMap := range destinationMaps {
			if destinationMap.Status() != status {
				continue
			}

			var itemData map[string]interface{}
			if item != nil {
				itemData = item.All()
			}

			errorItem := map[string]interface{}{
				"item":         itemData,
				string(status): itemError,
			}

			writer, err := destinationMap.Destination().Writer()
			if err != nil {
				return err
			}

			if err := writer.Setup(b.config); err != nil {
				return err
			}

			if err := writer.Write(MakeBasicPayload([]*Item{MakeItem(errorItem)}, map[string]interface{}{}, false)); err != nil {
				return err
			}
		}

	}
	return nil
}

type BaseConfig struct {
	Self    Config
	Stream_ Stream
}

func (b *BaseConfig) Type() string {
	return "config"
}

func (b *BaseConfig) Stream() Stream {
	return b.Stream_
}

func (b *BaseConfig) Update(values map[string]interface{}) error {

	if params, err := ConfigForm.ValidateUpdate(values); err != nil {
		return err
	} else {
		return b.update(params)
	}
}

func (b *BaseConfig) Create(values map[string]interface{}) error {

	if params, err := ConfigForm.Validate(values); err != nil {
		return err
	} else {
		return b.update(params)
	}
}

func (b *BaseConfig) MarshalJSON() ([]byte, error) {

	data := map[string]interface{}{
		"name":        b.Self.Name(),
		"description": b.Self.Description(),
		"source":      b.Self.Source(),
		"status":      b.Self.Status(),
		"version":     b.Self.Version(),
		"stream":      b.Self.Stream(),
		"data":        b.Self.Data(),
	}

	for k, v := range JSONData(b.Self) {
		data[k] = v
	}

	return json.Marshal(data)
}

func (b *BaseConfig) ChannelWriter() (ChannelWriter, error) {
	return &BaseChannelWriter{
		config: b.Self,
	}, nil
}

func (b *BaseConfig) Processor() (*Processor, error) {

	actionConfigs, err := b.Self.ActionConfigs()
	if err != nil {
		return nil, errors.MakeExternalError(
			"error getting action configs", "CREATE-PROCESSOR", nil, err)
	}

	actions := make([]Action, len(actionConfigs))

	for i, actionConfig := range actionConfigs {
		if action, err := actionConfig.Action(); err != nil {
			return nil, err
		} else {
			actions[i] = action
		}
	}

	parameterSet, err := MakeParameterSet(actions, b.Stream().Project().Controller().ParameterStore())
	if err != nil {
		return nil, err
	}
	channelWriter, err := b.Self.ChannelWriter()

	if err != nil {
		return nil, err
	}

	settings := b.Self.Stream().Project().Controller().Settings()
	processor, err := MakeProcessor(parameterSet, channelWriter, b.Self)

	if err != nil {
		return nil, err
	}

	if key, err := settings.Get("key"); err == nil {
		keyString, ok := key.(string)
		if !ok {
			return nil, fmt.Errorf("key is not a string")
		}
		processor.SetKey([]byte(keyString))
	}
	if salt, err := settings.Get("salt"); err == nil {
		saltStr, ok := salt.(string)
		if !ok {
			return nil, fmt.Errorf("salt is not a string")
		}
		processor.SetSalt([]byte(saltStr))
	}

	return processor, nil

}

func (b *BaseConfig) update(params map[string]interface{}) error {

	for key, value := range params {
		var err error
		switch key {
		case "name":
			err = b.Self.SetName(value.(string))
		case "description":
			err = b.Self.SetDescription(value.(string))
		case "status":
			err = b.Self.SetStatus(ConfigStatus(value.(string)))
		case "version":
			err = b.Self.SetVersion(value.(string))
		case "source":
			err = b.Self.SetSource(value.(string))
		case "data":
			err = b.Self.SetData(value)
		}
		if err != nil {
			return err
		}
	}

	return nil

}

var IsValidConfigStatus = forms.IsIn{
	Choices: []interface{}{
		string(ActiveConfig),
		string(DisabledConfig),
		string(TestingConfig)},
}

var ConfigForm = forms.Form{
	ErrorMsg: "invalid data encountered in the stream config form",
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
			Name: "status",
			Validators: []forms.Validator{
				forms.IsOptional{Default: string(DisabledConfig)},
				IsValidConfigStatus,
			},
		},
		{
			Name: "source",
			Validators: []forms.Validator{
				forms.IsOptional{}, forms.IsString{MinLength: 1, MaxLength: 40},
			},
		},
		{
			Name: "version",
			Validators: []forms.Validator{
				forms.IsOptional{}, forms.IsString{MinLength: 1, MaxLength: 40},
			},
		},
		{
			Name: "data",
			Validators: []forms.Validator{
				forms.IsOptional{}, forms.IsStringMap{},
			},
		},
	},
}
