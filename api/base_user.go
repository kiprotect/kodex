// Kodex (Community Edition - CE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2022  KIProtect GmbH (HRB 208395B) - Germany
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

package api

import (
	"encoding/json"
	"github.com/kiprotect/go-helpers/forms"
	"github.com/kiprotect/kodex"
)

type BaseUser struct {
	Self        User
	Controller_ Controller
}

func (b *BaseUser) Type() string {
	return "user"
}

func (b *BaseUser) Controller() Controller {
	return b.Controller_
}

func (b *BaseUser) Update(values map[string]interface{}) error {

	if params, err := UserForm.ValidateUpdate(values); err != nil {
		return err
	} else {
		return b.update(params)
	}

}

func (b *BaseUser) Create(values map[string]interface{}) error {

	if params, err := UserForm.Validate(values); err != nil {
		return err
	} else {
		return b.update(params)
	}

}

func (b *BaseUser) MarshalJSON() ([]byte, error) {

	data := map[string]interface{}{
		"displayName": b.Self.DisplayName(),
		"email":       b.Self.Email(),
		"data":        b.Self.Data(),
	}

	for k, v := range kodex.JSONData(b.Self) {
		data[k] = v
	}

	return json.Marshal(data)
}

func (b *BaseUser) update(params map[string]interface{}) error {

	for key, value := range params {
		var err error
		switch key {
		case "displayName":
			err = b.Self.SetDisplayName(value.(string))
		case "email":
			err = b.Self.SetEmail(value.(string))
		case "data":
			err = b.Self.SetData(value)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

var UserForm = forms.Form{
	ErrorMsg: "invalid data encountered in the user form",
	Fields: []forms.Field{
		{
			Name: "displayName",
			Validators: append([]forms.Validator{
				forms.IsRequired{}}, kodex.NameValidators...),
		},
		{
			Name: "email",
			Validators: append([]forms.Validator{
				forms.IsRequired{}}, kodex.NameValidators...),
		},
		{
			Name:       "data",
			Validators: []forms.Validator{forms.IsOptional{}, forms.IsStringMap{}},
		},
	},
}
