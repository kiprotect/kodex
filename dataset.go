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

package kodex

import (
	"encoding/json"
	"fmt"
	"github.com/kiprotect/go-helpers/forms"
	"github.com/kiprotect/go-helpers/maps"
)

type Dataset interface {
	Model
	Project() Project
	Items() []map[string]any
	Data() any
	Name() string
	Description() string
	SetData(any) error
	SetName(string) error
	SetDescription(string) error
	SetItems([]map[string]any) error
}

type IsDataset struct{}

func (i IsDataset) Validate(value interface{}, values map[string]interface{}) (interface{}, error) {
	config, ok := maps.ToStringMap(value)
	if !ok {
		return nil, fmt.Errorf("not a string map")
	}

	if configParams, err := DatasetForm.Validate(config); err != nil {
		return nil, err
	} else {
		return configParams, nil
	}
}

var DatasetForm = forms.Form{
	ErrorMsg: "invalid data encountered in the action config",
	Fields: []forms.Field{
		{
			Name: "name",
			Validators: append([]forms.Validator{
				forms.IsRequired{}}, NameValidators...),
		},
		{
			Name: "items",
			Validators: []forms.Validator{
				forms.IsOptional{Default: []map[string]any{}},
				forms.IsList{
					Validators: []forms.Validator{
						forms.IsStringMap{},
					},
				},
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

type BaseDataset struct {
	Self     Dataset
	Project_ Project
}

func (b *BaseDataset) Type() string {
	return "dataset"
}

func (b *BaseDataset) Project() Project {
	return b.Project_
}

func (b *BaseDataset) Update(values map[string]interface{}) error {

	if params, err := DatasetForm.ValidateUpdate(values); err != nil {
		return err
	} else {
		return b.update(params)
	}

}

func (b *BaseDataset) Create(values map[string]interface{}) error {
	if params, err := DatasetForm.Validate(values); err != nil {
		return err
	} else {
		return b.update(params)
	}
}

func (b *BaseDataset) update(params map[string]interface{}) error {

	for key, value := range params {
		var err error
		switch key {
		case "name":
			err = b.Self.SetName(value.(string))
		case "description":
			err = b.Self.SetDescription(value.(string))
		case "items":
			err = b.Self.SetItems(value.([]map[string]any))
		case "data":
			err = b.Self.SetData(value)
		}
		if err != nil {
			return err
		}
	}

	return nil

}

func (b *BaseDataset) MarshalJSON() ([]byte, error) {

	data := map[string]interface{}{
		"name":        b.Self.Name(),
		"description": b.Self.Description(),
		"data":        b.Self.Data(),
		"items":       b.Self.Items(),
	}

	for k, v := range JSONData(b.Self) {
		data[k] = v
	}

	return json.Marshal(data)
}
