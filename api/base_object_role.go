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
	"encoding/hex"
	"encoding/json"
	"github.com/kiprotect/go-helpers/forms"
	"github.com/kiprotect/kodex"
	"regexp"
)

// BaseObjectRole contains useful common functionality that should be shared by
// all implementations of the interface, such as validation.
type BaseObjectRole struct {
	Self ObjectRole
}

func (b *BaseObjectRole) Type() string {
	return "object_role"
}

func (b *BaseObjectRole) MarshalJSON() ([]byte, error) {

	data := map[string]interface{}{
		"organization_id":   hex.EncodeToString(b.Self.OrganizationID()),
		"object_id":         hex.EncodeToString(b.Self.ObjectID()),
		"organization_role": b.Self.OrganizationRole(),
		"object_role":       b.Self.ObjectRole(),
		"object_type":       b.Self.ObjectType(),
	}

	for k, v := range kodex.JSONData(b.Self) {
		data[k] = v
	}

	return json.Marshal(data)
}

func (b *BaseObjectRole) Update(values map[string]interface{}) error {

	if params, err := ObjectRoleForm.ValidateUpdate(values); err != nil {
		return err
	} else {
		return b.update(params)
	}

}

func (b *BaseObjectRole) Create(values map[string]interface{}) error {

	if params, err := ObjectRoleForm.Validate(values); err != nil {
		return err
	} else {
		return b.update(params)
	}

}

func (b *BaseObjectRole) update(params map[string]interface{}) error {

	for key, value := range params {
		var err error
		switch key {
		case "organization_role":
			err = b.Self.SetOrganizationRole(value.(string))
		case "role":
			err = b.Self.SetObjectRole(value.(string))
		}
		if err != nil {
			return err
		}
	}
	return nil

}

var ObjectRoleForm = forms.Form{
	ErrorMsg: "invalid data encountered in the object role form",
	Fields: []forms.Field{
		{
			Name: "organization_role",
			Validators: []forms.Validator{
				forms.IsRequired{},
				forms.IsString{MinLength: 2, MaxLength: 100},
				forms.MatchesRegex{Regexp: regexp.MustCompile(`^[\w\d\-\:\.]{2,100}$`)},
			},
		},
		{
			Name: "role",
			Validators: []forms.Validator{
				forms.IsRequired{},
				forms.IsString{},
				forms.IsIn{Choices: []interface{}{"superuser", "admin", "viewer", "reviewer", "editor", "legal-reviewer", "technical-reviewer"}},
			},
		},
	},
}
