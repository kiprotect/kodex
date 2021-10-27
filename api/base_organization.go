// Kodex (Enterprise Edition - EE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2021  KIProtect GmbH (HRB 208395B) - All Rights Reserved

package api

import (
	"encoding/json"
	"github.com/kiprotect/go-helpers/forms"
	"github.com/kiprotect/kodex"
)

type BaseOrganization struct {
	Self        Organization
	Controller_ Controller
}

func (b *BaseOrganization) Type() string {
	return "organization"
}

func (b *BaseOrganization) Controller() Controller {
	return b.Controller_
}

func (b *BaseOrganization) Update(values map[string]interface{}) error {

	if params, err := OrganizationForm.ValidateUpdate(values); err != nil {
		return err
	} else {
		return b.update(params)
	}

}

func (b *BaseOrganization) Create(values map[string]interface{}) error {

	if params, err := OrganizationForm.Validate(values); err != nil {
		return err
	} else {
		return b.update(params)
	}

}

func (b *BaseOrganization) MarshalJSON() ([]byte, error) {

	data := map[string]interface{}{
		"name":        b.Self.Name(),
		"description": b.Self.Description(),
		"data":        b.Self.Data(),
	}

	for k, v := range kodex.JSONData(b.Self) {
		data[k] = v
	}

	return json.Marshal(data)
}

func (b *BaseOrganization) update(params map[string]interface{}) error {

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

var OrganizationForm = forms.Form{
	ErrorMsg: "invalid data encountered in the organization form",
	Fields: []forms.Field{
		{
			Name: "name",
			Validators: append([]forms.Validator{
				forms.IsRequired{}}, kodex.NameValidators...),
		},
		{
			Name: "description",
			Validators: append([]forms.Validator{
				forms.IsOptional{Default: ""}}, kodex.DescriptionValidators...),
		},
		{
			Name:       "data",
			Validators: []forms.Validator{forms.IsOptional{}, forms.IsStringMap{}},
		},
	},
}
