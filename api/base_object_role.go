// Kodex (Enterprise Edition - EE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2021  KIProtect GmbH (HRB 208395B) - All Rights Reserved

package api

import (
	"github.com/kiprotect/go-helpers/forms"
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
				forms.MatchesRegex{Regex: regexp.MustCompile(`^[\w\d\-\:\.]{2,100}$`)},
			},
		},
		{
			Name: "role",
			Validators: []forms.Validator{
				forms.IsRequired{},
				forms.IsString{},
				forms.IsIn{Choices: []interface{}{"superuser", "admin", "viewer"}},
			},
		},
	},
}
