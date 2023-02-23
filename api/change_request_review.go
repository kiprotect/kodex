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
	"fmt"
	"github.com/kiprotect/go-helpers/forms"
	"github.com/kiprotect/kodex"
)

type ChangeRequestReviewStatus string

const (
	ReviewRequested ChangeRequestReviewStatus = "requested"
	RequestRejected ChangeRequestReviewStatus = "rejected"
	RequestApproved ChangeRequestReviewStatus = "approved"
)

type ChangeRequestReview interface {
	kodex.Model
	SetMetadata(interface{}) error
	Metadata() interface{}
	SetData(interface{}) error
	Data() interface{}
	SetStatus(ChangeRequestReviewStatus) error
	Status() ChangeRequestReviewStatus
	ChangeRequest() ChangeRequest
	Creator() User
}

type IsChangeRequestReviewStatus struct{}

func (i IsChangeRequestReviewStatus) Validate(value interface{}, values map[string]interface{}) (interface{}, error) {

	if enumValue, ok := value.(ChangeRequestReviewStatus); ok {
		// this already is an enum
		return enumValue, nil
	}

	// we expect a string
	strValue, ok := value.(string)

	if !ok {
		return nil, fmt.Errorf("expected a string")
	}

	// we convert the string...
	return ChangeRequestReviewStatus(strValue), nil
}

var ChangeRequestReviewForm = forms.Form{
	ErrorMsg: "invalid data encountered in the change request config",
	Fields: []forms.Field{
		{
			Name: "data",
			Validators: []forms.Validator{
				forms.IsOptional{},
				forms.IsStringMap{},
			},
		},
		{
			Name: "metadata",
			Validators: []forms.Validator{
				forms.IsOptional{},
				forms.IsStringMap{},
			},
		},
		{
			Name: "status",
			Validators: []forms.Validator{
				forms.IsOptional{Default: ReviewRequested},
				IsChangeRequestReviewStatus{},
				forms.IsIn{Choices: []interface{}{ReviewRequested}},
			},
		},
	},
}

/* Base Functionality */

type BaseChangeRequestReview struct {
	Self     ChangeRequestReview
	Project_ kodex.Project
	Creator_ User
}

func (b *BaseChangeRequestReview) Type() string {
	return "change-request"
}

func (b *BaseChangeRequestReview) Project() kodex.Project {
	return b.Project_
}

func (b *BaseChangeRequestReview) Creator() User {
	return b.Creator_
}

func (b *BaseChangeRequestReview) Update(values map[string]interface{}) error {

	if params, err := ChangeRequestReviewForm.ValidateUpdate(values); err != nil {
		return err
	} else {
		return b.update(params)
	}

}

func (b *BaseChangeRequestReview) Create(values map[string]interface{}) error {

	if params, err := ChangeRequestReviewForm.Validate(values); err != nil {
		return err
	} else {
		return b.update(params)
	}

}

func (b *BaseChangeRequestReview) update(params map[string]interface{}) error {

	for key, value := range params {
		var err error
		switch key {
		case "status":
			err = b.Self.SetStatus(value.(ChangeRequestReviewStatus))
		case "metadata":
			err = b.Self.SetMetadata(value)
		case "data":
			err = b.Self.SetData(value)
		}
		if err != nil {
			return err
		}
	}

	return nil

}

func (b *BaseChangeRequestReview) MarshalJSON() ([]byte, error) {

	data := map[string]interface{}{
		"data":     b.Self.Data(),
		"status":   b.Self.Status(),
		"creator":  b.Self.Creator(),
		"metadata": b.Self.Metadata(),
	}

	for k, v := range kodex.JSONData(b.Self) {
		data[k] = v
	}

	return json.Marshal(data)
}
