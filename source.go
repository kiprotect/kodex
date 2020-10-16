// Kodex (Community Edition - CE) - Privacy & Security Engineering Platform
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

package kodex

import (
	"encoding/json"
	"fmt"
	"github.com/kiprotect/go-helpers/forms"
)

type SourceStatus string

const (
	// If a source is active, we try to read from it
	ActiveSource SourceStatus = "active"
	// If it is disabled, we ignore it
	DisabledSource SourceStatus = "disabled"
)

type Source interface {
	Processable // Processable includes Model
	Reader() (Reader, error)
	ConfigData() map[string]interface{}
	SetConfigData(map[string]interface{}) error
	SourceType() string
	SetSourceType(string) error
	Name() string
	SetName(string) error
	Description() string
	SetDescription(string) error

	// Return all streams with a given source status for this source
	Streams(SourceStatus) ([]Stream, error)

	SetData(interface{}) error
	Data() interface{}

	Project() Project

	Service() Service
	SetService(Service) error
}

/* Base Functionality */

type BaseSource struct {
	Self     Source
	Project_ Project
	reader   Reader
}

func (b *BaseSource) Type() string {
	return "source"
}

func (b *BaseSource) Stats() (map[string]float64, error) {
	return nil, nil
}

func (b *BaseSource) Stat(name string) (float64, error) {
	return 0.0, nil
}

func (b *BaseSource) SetStat(name string, value float64) error {
	return nil
}

func (b *BaseSource) Project() Project {
	return b.Project_
}

func (b *BaseSource) checkSourceConfig(sourceType string, config map[string]interface{}) error {
	definitions := b.Self.Project().Controller().Definitions()
	if definition, ok := definitions.ReaderDefinitions[sourceType]; !ok {
		return fmt.Errorf("invalid source type: %s", sourceType)
	} else if _, err := definition.Maker(config); err != nil {
		return err
	}
	return nil
}

func (b *BaseSource) Update(values map[string]interface{}) error {

	if params, err := SourceForm.ValidateUpdate(values); err != nil {
		return err
	} else {
		sourceType, ok := params["type"].(string)
		if !ok {
			sourceType = b.Self.SourceType()
		}
		sourceConfig, ok := params["config"].(map[string]interface{})
		if !ok {
			sourceConfig = b.Self.ConfigData()
		}
		// we validate the new config/type combination
		if err := b.checkSourceConfig(sourceType, sourceConfig); err != nil {
			return err
		}
		return b.update(params)
	}

}

func (b *BaseSource) Create(values map[string]interface{}) error {

	if params, err := SourceForm.Validate(values); err != nil {
		return err
	} else {
		sourceType := params["type"].(string)
		sourceConfig := params["config"].(map[string]interface{})
		if err := b.checkSourceConfig(sourceType, sourceConfig); err != nil {
			return err
		}
		return b.update(params)
	}

}

func (b *BaseSource) MarshalJSON() ([]byte, error) {

	data := map[string]interface{}{
		"name":        b.Self.Name(),
		"description": b.Self.Description(),
		"type":        b.Self.SourceType(),
		"config":      b.Self.ConfigData(),
		"project":     b.Self.Project(),
		"data":        b.Self.Data(),
	}

	for k, v := range JSONData(b.Self) {
		data[k] = v
	}

	return json.Marshal(data)
}

func (b *BaseSource) update(params map[string]interface{}) error {

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
			err = b.Self.SetSourceType(value.(string))
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *BaseSource) Reader() (Reader, error) {
	if b.reader != nil {
		return b.reader, nil
	}
	reader, err := b.createReader()
	if err != nil {
		return nil, err
	}
	b.reader = reader
	return b.reader, nil
}

func (b *BaseSource) createReader() (Reader, error) {
	definitions := b.Self.Project().Controller().Definitions()
	sourceType := b.Self.SourceType()
	config := b.Self.ConfigData()
	definition, ok := definitions.ReaderDefinitions[sourceType]

	if !ok {
		return nil, fmt.Errorf("unknown reader type %s", sourceType)
	}

	return definition.Maker(config)

}

var SourceForm = forms.Form{
	ErrorMsg: "invalid data encountered in the source form",
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
