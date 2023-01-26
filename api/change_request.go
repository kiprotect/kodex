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

type ChangeRequest interface {
	kodex.Model
	SetData(interface{}) error
	Data() interface{}
	ObjectID() []byte
	ObjectType() string
}

var ChangeRequestForm = forms.Form{
	ErrorMsg: "invalid data encountered in the change request config",
	Fields: []forms.Field{
		{
			Name: "data",
			Validators: []forms.Validator{
				forms.IsOptional{},
				forms.IsStringMap{},
			},
		},
	},
}

/* Base Functionality */

type BaseChangeRequest struct {
	Self     ChangeRequest
	Project_ kodex.Project
}

func (b *BaseChangeRequest) Type() string {
	return "change-request"
}

func (b *BaseChangeRequest) Project() kodex.Project {
	return b.Project_
}

func (b *BaseChangeRequest) Update(values map[string]interface{}) error {

	if params, err := ChangeRequestForm.ValidateUpdate(values); err != nil {
		return err
	} else {
		return b.update(params)
	}

}

func (b *BaseChangeRequest) Create(values map[string]interface{}) error {

	if params, err := ChangeRequestForm.Validate(values); err != nil {
		return err
	} else {
		return b.update(params)
	}

}

func (b *BaseChangeRequest) update(params map[string]interface{}) error {

	for key, value := range params {
		var err error
		switch key {
		case "data":
			err = b.Self.SetData(value)
		}
		if err != nil {
			return err
		}
	}

	return nil

}

func (b *BaseChangeRequest) MarshalJSON() ([]byte, error) {

	data := map[string]interface{}{
		"data":        b.Self.Data(),
		"object_id":   b.Self.ObjectID(),
		"object_type": b.Self.ObjectType(),
	}

	for k, v := range kodex.JSONData(b.Self) {
		data[k] = v
	}

	return json.Marshal(data)
}
