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

type CommentStatus string

const (
	DraftComment     CommentStatus = "draft"
	PublishedComment CommentStatus = "published"
	DeletedComment   CommentStatus = "deleted"
)

type Comment interface {
	kodex.Model
	Text() string
	SetText(string) error
	SetData(interface{}) error
	Data() interface{}
	SetStatus(CommentStatus) error
	Status() CommentStatus
	Creator() User
	ObjectID() []byte
	ObjectType() string
}

type IsCommentStatus struct{}

func (i IsCommentStatus) Validate(value interface{}, values map[string]interface{}) (interface{}, error) {

	if enumValue, ok := value.(CommentStatus); ok {
		// this already is an enum
		return enumValue, nil
	}

	// we expect a string
	strValue, ok := value.(string)

	if !ok {
		return nil, fmt.Errorf("expected a string")
	}

	// we convert the string...
	return CommentStatus(strValue), nil
}

var CommentForm = forms.Form{
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
			Name: "text",
			Validators: []forms.Validator{
				forms.IsString{MinLength: 5},
			},
		},
		{
			Name: "status",
			Validators: []forms.Validator{
				forms.IsOptional{Default: DraftComment},
				IsCommentStatus{},
				forms.IsIn{Choices: []interface{}{DraftComment, PublishedComment, DeletedComment}},
			},
		},
	},
}

/* Base Functionality */

type BaseComment struct {
	Self     Comment
	Creator_ User
}

func (b *BaseComment) Type() string {
	return "comment"
}

func (b *BaseComment) Creator() User {
	return b.Creator_
}

func (b *BaseComment) Update(values map[string]interface{}) error {

	if params, err := CommentForm.ValidateUpdate(values); err != nil {
		return err
	} else {
		return b.update(params)
	}

}

func (b *BaseComment) Create(values map[string]interface{}) error {

	if params, err := CommentForm.Validate(values); err != nil {
		return err
	} else {
		return b.update(params)
	}

}

func (b *BaseComment) update(params map[string]interface{}) error {

	for key, value := range params {
		var err error
		switch key {
		case "text":
			err = b.Self.SetText(value.(string))
		case "status":
			err = b.Self.SetStatus(value.(CommentStatus))
		case "data":
			err = b.Self.SetData(value)
		}
		if err != nil {
			return err
		}
	}

	return nil

}

func (b *BaseComment) MarshalJSON() ([]byte, error) {

	data := map[string]interface{}{
		"data":        b.Self.Data(),
		"text":        b.Self.Text(),
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
