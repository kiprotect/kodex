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
	"github.com/kiprotect/go-helpers/forms"
)

type FieldSchema struct {
}

type DataSchema struct {
	Fields map[string]FieldSchema
}

type Schema interface {
	Model

	SetData(interface{}) error
	Data() interface{}

	Schema() *DataSchema
	SetSchema(*DataSchema) error
}

/* Base Functionality */

type BaseSchema struct {
	Self     Schema
	Project_ Project
}

func (b *BaseSchema) Type() string {
	return "schema"
}

func (b *BaseSchema) Update(values map[string]interface{}) error {

	if params, err := SchemaForm.ValidateUpdate(values); err != nil {
		return err
	} else {
		return b.update(params)
	}

}

func (b *BaseSchema) Create(values map[string]interface{}) error {

	if params, err := SchemaForm.Validate(values); err != nil {
		return err
	} else {
		return b.update(params)
	}

}

func (b *BaseSchema) MarshalJSON() ([]byte, error) {

	data := map[string]interface{}{
		"schema": b.Self.Schema(),
		"data":   b.Self.Data(),
	}

	for k, v := range JSONData(b.Self) {
		data[k] = v
	}

	return json.Marshal(data)
}

func (b *BaseSchema) update(params map[string]interface{}) error {

	for key, value := range params {
		var err error
		switch key {
		case "schema":
			err = b.Self.SetSchema(value.(*DataSchema))
		case "data":
			err = b.Self.SetData(value)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

type IsDataSchema struct{}

func (i IsDataSchema) Validate(value interface{}, values map[string]interface{}) (interface{}, error) {
	return value, nil
}

var SchemaForm = forms.Form{
	ErrorMsg: "invalid data encountered in the schema form",
	Fields: []forms.Field{
		{
			Name: "schema",
			Validators: []forms.Validator{
				forms.IsStringMap{},
				IsDataSchema{},
			},
		},
		{
			Name:       "data",
			Validators: []forms.Validator{forms.IsOptional{}, forms.IsStringMap{}},
		},
	},
}
