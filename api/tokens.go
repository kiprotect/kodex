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

package api

import (
	"encoding/json"
	"github.com/kiprotect/go-helpers/forms"
	"github.com/kiprotect/go-helpers/strings"
	"github.com/kiprotect/kodex"
	"time"
)

type TokenController interface {
	Tokens(org Organization, filters map[string]interface{}) ([]Token, error)
	Token(org Organization, id []byte) (Token, error)
	TokenByValue(token []byte) (Token, error)
	MakeToken(org Organization, user User) (Token, error)
}

type Token interface {
	kodex.Model
	Description() string
	SetDescription(string) error
	Scopes() []string
	SetScopes([]string) error
	Roles() []string
	SetRoles([]string) error
	Token() []byte
	User() User
	SetToken([]byte) error
	ExpiresAt() *time.Time
	SetExpiresAt(time.Time) error
	Data() interface{}
	SetData(interface{}) error
	Organization() Organization
}

// BaseToken contains useful common functionality that should be shared by
// all implementations of the interface, such as validation.
type BaseToken struct {
	Self Token
}

func (b *BaseToken) Type() string {
	return "token"
}

func (b *BaseToken) Update(values map[string]interface{}) error {
	if params, err := TokenForm.ValidateUpdate(values); err != nil {
		return err
	} else {
		return b.update(params)
	}
}

func (b *BaseToken) Create(values map[string]interface{}) error {

	if params, err := TokenForm.Validate(values); err != nil {
		return err
	} else {
		return b.update(params)
	}

}

func (b *BaseToken) update(params map[string]interface{}) error {

	for key, value := range params {
		var err error
		switch key {
		case "description":
			err = b.Self.SetDescription(value.(string))
		case "roles":
			if sl, slErr := strings.ToListOfStr(value); err != nil {
				err = slErr
			} else {
				err = b.Self.SetRoles(sl)
			}
		case "scopes":
			if sl, slErr := strings.ToListOfStr(value); err != nil {
				err = slErr
			} else {
				err = b.Self.SetScopes(sl)
			}
		case "expires_at":
			err = b.Self.SetExpiresAt(value.(time.Time))
		case "data":
			err = b.Self.SetData(value)
		}
		if err != nil {
			return err
		}
	}
	return nil

}

func (b *BaseToken) MarshalJSON() ([]byte, error) {

	data := map[string]interface{}{
		"description":  b.Self.Description(),
		"data":         b.Self.Data(),
		"roles":        b.Self.Roles(),
		"scopes":       b.Self.Scopes(),
		"organization": b.Self.Organization(),
		"expires_at":   b.Self.ExpiresAt(),
	}

	for k, v := range kodex.JSONData(b.Self) {
		data[k] = v
	}

	return json.Marshal(data)
}

var TokenForm = forms.Form{
	ErrorMsg: "invalid data encountered in the consent config form",
	Fields: []forms.Field{
		{
			Name: "description",
			Validators: append([]forms.Validator{
				forms.IsOptional{Default: ""}},
				kodex.DescriptionValidators...),
		},
		{
			Name: "scopes",
			Validators: []forms.Validator{
				forms.IsOptional{Default: []string{"kiprotect:api"}},
				forms.IsList{
					Validators: []forms.Validator{
						forms.IsString{},
					},
				},
			},
		},
		{
			Name: "expires_at",
			Validators: []forms.Validator{
				forms.IsOptional{},
				forms.IsTime{},
			},
		},
		{
			Name: "data",
			Validators: []forms.Validator{
				forms.IsOptional{},
				forms.IsStringMap{},
			},
		},
	},
}
