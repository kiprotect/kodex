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
	"fmt"
	"github.com/kiprotect/go-helpers/forms"
	"github.com/kiprotect/kodex"
)

type ChangeRequestStatus string

const (
	DraftCR     ChangeRequestStatus = "draft"
	ReadyCR     ChangeRequestStatus = "ready"
	WithdrawnCR ChangeRequestStatus = "withdrawn"
	ApprovedCR  ChangeRequestStatus = "approved"
	RejectedCR  ChangeRequestStatus = "rejected"
	MergedCR    ChangeRequestStatus = "merged"
)

type ChangeRequest interface {
	kodex.Model
	Title() string
	SetTitle(string) error
	Description() string
	SetDescription(string) error
	SetData(interface{}) error
	Data() interface{}
	Changes() []ChangeSet
	SetChanges([]ChangeSet) error
	SetStatus(ChangeRequestStatus) error
	Status() ChangeRequestStatus
	Reviews() ([]ChangeRequestReview, error)
	MakeReview(User) (ChangeRequestReview, error)
	Review([]byte) (ChangeRequestReview, error)
	Creator() User
	ObjectID() []byte
	ObjectType() string
}

type IsChangeRequestStatus struct{}

func (i IsChangeRequestStatus) Validate(value interface{}, values map[string]interface{}) (interface{}, error) {

	if enumValue, ok := value.(ChangeRequestStatus); ok {
		// this already is an enum
		return enumValue, nil
	}

	// we expect a string
	strValue, ok := value.(string)

	if !ok {
		return nil, fmt.Errorf("expected a string")
	}

	// we convert the string...
	return ChangeRequestStatus(strValue), nil
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
		{
			Name: "description",
			Validators: []forms.Validator{
				forms.IsString{MinLength: 5},
			},
		},
		{
			Name: "title",
			Validators: []forms.Validator{
				forms.IsString{MinLength: 2},
			},
		},
		{
			Name: "status",
			Validators: []forms.Validator{
				forms.IsOptional{Default: DraftCR},
				IsChangeRequestStatus{},
			},
		},
	},
}

/* Base Functionality */

type BaseChangeRequest struct {
	Self     ChangeRequest
	Creator_ User
}

func (b *BaseChangeRequest) Type() string {
	return "change-request"
}

func (b *BaseChangeRequest) Creator() User {
	return b.Creator_
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
		case "title":
			err = b.Self.SetTitle(value.(string))
		case "description":
			err = b.Self.SetDescription(value.(string))
		case "status":
			err = b.Self.SetStatus(value.(ChangeRequestStatus))
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
		"title":       b.Self.Title(),
		"description": b.Self.Description(),
		"status":      b.Self.Status(),
		"creator":     b.Self.Creator(),
		"object_id":   b.Self.ObjectID(),
		"object_type": b.Self.ObjectType(),
	}

	for k, v := range kodex.JSONData(b.Self) {
		data[k] = v
	}

	return json.Marshal(data)
}
