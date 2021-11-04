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

package kodex

import (
	"encoding/json"
	"github.com/kiprotect/go-helpers/forms"
)

type StreamStatus string

const (
	ActiveStream   StreamStatus = "active"
	DisabledStream StreamStatus = "disabled"
	TestingStream  StreamStatus = "testing"
)

type StreamStats struct {
	ItemFrequency float64
	IdleFraction  float64
	Executors     int64
}

type Stream interface {
	Processable // Processable includes Model
	PriorityModel
	Configs() ([]Config, error)
	Config([]byte) (Config, error)
	MakeConfig(id []byte) Config

	AddSource(Source, SourceStatus) error
	RemoveSource(Source) error
	Sources() (map[string]SourceMap, error)

	Status() StreamStatus
	SetStatus(StreamStatus) error
	Name() string
	SetName(string) error
	Description() string
	SetDescription(string) error

	SetData(interface{}) error
	Data() interface{}

	Project() Project
}

/* Base Functionality */

// BaseStream contains useful common functionality that should be shared by
// all implementations of the interface, such as validation.
type BaseStream struct {
	Self     Stream
	Project_ Project
}

func (b *BaseStream) Type() string {
	return "stream"
}

func (b *BaseStream) Update(values map[string]interface{}) error {

	if params, err := StreamForm.ValidateUpdate(values); err != nil {
		return err
	} else {
		return b.update(params)
	}

}

func (b *BaseStream) Create(values map[string]interface{}) error {

	if params, err := StreamForm.Validate(values); err != nil {
		return err
	} else {
		return b.update(params)
	}

}

func (b *BaseStream) MarshalJSON() ([]byte, error) {

	data := map[string]interface{}{
		"name":        b.Self.Name(),
		"status":      b.Self.Status(),
		"description": b.Self.Description(),
		"project":     b.Self.Project(),
		"data":        b.Self.Data(),
	}

	for k, v := range JSONData(b.Self) {
		data[k] = v
	}

	return json.Marshal(data)
}

func (b *BaseStream) update(params map[string]interface{}) error {

	for key, value := range params {
		var err error
		switch key {
		case "status":
			err = b.Self.SetStatus(StreamStatus(value.(string)))
		case "name":
			err = b.Self.SetName(value.(string))
		case "description":
			err = b.Self.SetDescription(value.(string))
		case "data":
			err = b.Self.SetData(value)
		}
		if err != nil {
			return err
		}
	}

	return nil

}

var NameValidators = []forms.Validator{
	forms.IsString{MinLength: 2, MaxLength: 100},
}

var DescriptionValidators = []forms.Validator{
	forms.IsString{MaxLength: 256},
}

var IsValidStreamStatus = forms.IsIn{
	Choices: []interface{}{
		string(ActiveStream),
		string(DisabledStream),
		string(TestingStream)},
}

var StreamForm = forms.Form{
	ErrorMsg: "invalid data encountered in the stream form",
	Fields: []forms.Field{
		{
			Name: "name",
			Validators: append([]forms.Validator{
				forms.IsRequired{}}, NameValidators...),
		},
		{
			Name: "status",
			Validators: []forms.Validator{
				forms.IsOptional{Default: string(ActiveStream)},
				IsValidStreamStatus,
			},
		},
		{
			Name: "description",
			Validators: append([]forms.Validator{
				forms.IsOptional{Default: ""}}, DescriptionValidators...),
		},
		{
			Name:       "data",
			Validators: []forms.Validator{forms.IsOptional{}, forms.IsStringMap{}},
		},
	},
}

func (b *BaseStream) Project() Project {
	return b.Project_
}
