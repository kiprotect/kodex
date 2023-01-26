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

package resources

import (
	"github.com/gin-gonic/gin"
	"github.com/kiprotect/go-helpers/forms"
	"github.com/kiprotect/kodex/api"
	"github.com/kiprotect/kodex/api/helpers"
)

var createChangeRequestForm = forms.Form{
	ErrorMsg: "invalid data encountered in the default role creation form",
	Fields: []forms.Field{
		{
			Name: "object_type",
			Validators: []forms.Validator{
				forms.IsRequired{},
				forms.IsIn{Choices: []interface{}{"project"}},
			},
		},
	},
}

func CreateChangeRequest(c *gin.Context) {

	controller := helpers.Controller(c)

	if controller == nil {
		return
	}

	object := GetObj(c, "objectType")

	if object == nil {
		return
	}

	data := helpers.JSONData(c)

	if data == nil {
		return
	}

	request := controller.MakeChangeRequest(object)

	if err := request.Create(data); err != nil {
		api.HandleError(c, 400, err)
		return
	}

	if err := request.Save(); err != nil {
		api.HandleError(c, 500, err)
		return
	}

	c.JSON(200, map[string]interface{}{"data": request})

}

var deleteChangeRequestForm = forms.Form{
	ErrorMsg: "invalid data encountered in the change request deletion form",
	Fields: []forms.Field{
		{
			Name: "request_id",
			Validators: []forms.Validator{
				forms.IsRequired{},
				forms.IsHex{ConvertToBinary: true, Strict: false},
			},
		},
	},
}

// Delete a change request
func DeleteChangeRequest(c *gin.Context) {

	controller := helpers.Controller(c)

	if controller == nil {
		return
	}

	data := map[string]interface{}{
		"request_id": c.Param("requestID"),
	}

	params, err := deleteChangeRequestForm.Validate(data)

	if err != nil {
		api.HandleError(c, 400, err)
		return
	}

	requestID := params["request_id"].([]byte)

	request, err := controller.ChangeRequest(requestID)

	if err != nil {
		api.HandleError(c, 404, err)
		return
	}

	if err := request.Delete(); err != nil {
		api.HandleError(c, 500, err)
		return
	}

	c.JSON(200, map[string]interface{}{"message": "ok"})

}

// Get a list of change requestse
func ChangeRequests(c *gin.Context) {

	controller := helpers.Controller(c)

	if controller == nil {
		return
	}

	object := GetObj(c, "objectType")

	if object == nil {
		return
	}

	requests, err := controller.ChangeRequests(object)

	if err != nil {
		api.HandleError(c, 500, err)
		return
	}

	c.JSON(200, map[string]interface{}{"data": requests})

}
