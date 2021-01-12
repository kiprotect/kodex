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

type Project interface {
	Model
	Name() string
	SetName(string) error
	Description() string
	SetDescription(string) error
	Data() interface{}
	SetData(interface{}) error

	MakeActionConfig() ActionConfig
	MakeDestination() Destination
	MakeSource() Source
	MakeStream() Stream

	Controller() Controller
}

/* Base Functionality */

type BaseProject struct {
	Self        Project
	Controller_ Controller
}

func (b *BaseProject) Type() string {
	return "project"
}

func (b *BaseProject) Update(values map[string]interface{}) error {

	if params, err := ProjectForm.ValidateUpdate(values); err != nil {
		return err
	} else {
		return b.update(params)
	}
}

func (b *BaseProject) Create(values map[string]interface{}) error {

	if params, err := ProjectForm.Validate(values); err != nil {
		return err
	} else {
		return b.update(params)
	}
}

func (b *BaseProject) MarshalJSON() ([]byte, error) {

	data := map[string]interface{}{
		"name":        b.Self.Name(),
		"description": b.Self.Description(),
		"data":        b.Self.Data(),
	}

	for k, v := range JSONData(b.Self) {
		data[k] = v
	}

	return json.Marshal(data)
}

func (b *BaseProject) update(params map[string]interface{}) error {

	for key, value := range params {
		var err error
		switch key {
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

var ProjectForm = forms.Form{
	ErrorMsg: "invalid data encountered in the project form",
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
			Name: "data",
			Validators: []forms.Validator{
				forms.IsOptional{}, forms.IsStringMap{},
			},
		},
	},
}

func (b *BaseProject) Controller() Controller {
	return b.Controller_
}
