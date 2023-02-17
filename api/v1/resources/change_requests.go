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
	"fmt"
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

var updateChangeRequestStatusForn = forms.Form{
	ErrorMsg: "invalid data encountered in the change request update status form",
	Fields: []forms.Field{
		{
			Name: "status",
			Validators: []forms.Validator{
				forms.IsRequired{},
				forms.IsIn{Choices: []interface{}{api.Draft, api.Ready, api.Withdrawn, api.Rejected, api.Approved}},
			},
		},
	},
}

func UpdateChangeRequestStatus(c *gin.Context) {

	controller := helpers.Controller(c)

	if controller == nil {
		return
	}

	user := helpers.User(c)

	if user == nil {
		return
	}

	object := GetObj(c, "objectType")

	if object == nil {
		return
	}

	changeRequestObj := GetObj(c, "changeRequest")

	if changeRequestObj == nil {
		return
	}

	changeRequest, ok := changeRequestObj.(api.ChangeRequest)

	if !ok {
		api.HandleError(c, 500, fmt.Errorf("invalid change request"))
		return
	}

	data := helpers.JSONData(c)

	if data == nil {
		return
	}

	params, err := updateChangeRequestStatusForn.Validate(data)

	if err != nil {
		api.HandleError(c, 400, err)
		return
	}

	status := params["status"].(api.ChangeRequestStatus)

	isReviewer, err := controller.CanAccess(user, object, []string{"reviewer", "admin", "superuser"})

	if err != nil {
		api.HandleError(c, 500, err)
		return
	}

	if !isReviewer {
		// requests can only go from "draft" to "ready" or "withdrawn"
		if changeRequest.Status() == api.Draft && (status != api.Ready && status != api.Withdrawn) {
			api.HandleError(c, 400, fmt.Errorf("cannot set status from %s to %s", changeRequest.Status(), status))
			return
		} else if changeRequest.Status() == api.Ready && (status != api.Withdrawn && status != api.Draft) {
			api.HandleError(c, 400, fmt.Errorf("cannot set status from %s to %s", changeRequest.Status(), status))
			return
		}
	}

	if err := changeRequest.SetStatus(status); err != nil {
		api.HandleError(c, 500, fmt.Errorf("cannot update change request status: %v", err))
		return
	}

	c.JSON(200, map[string]interface{}{"data": changeRequest})

}

func CreateChangeRequest(c *gin.Context) {

	controller := helpers.Controller(c)

	if controller == nil {
		return
	}

	user := helpers.User(c)

	if user == nil {
		return
	}

	object := GetObj(c, "objectType")

	if object == nil {
		return
	}

	// we check that the user can actually access the change request
	if ok, err := controller.CanAccess(user, object, []string{"editor", "admin", "superuser"}); err != nil {
		api.HandleError(c, 500, err)
		return
	} else if !ok {
		api.HandleError(c, 401, fmt.Errorf("you must be an editor, admin or superuser to create a change request"))
		return
	}

	data := helpers.JSONData(c)

	if data == nil {
		return
	}

	apiUser, err := user.ApiUser(controller)

	if err != nil {
		api.HandleError(c, 500, err)
		return
	}

	request := controller.MakeChangeRequest(object, apiUser)

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
func UpdateChangeRequest(c *gin.Context) {

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
