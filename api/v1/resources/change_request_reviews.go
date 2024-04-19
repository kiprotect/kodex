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

package resources

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kiprotect/go-helpers/forms"
	"github.com/kiprotect/kodex/api"
	"github.com/kiprotect/kodex/api/helpers"
)

var updateChangeRequestReviewStatusForn = forms.Form{
	ErrorMsg: "invalid data encountered in the change request update status form",
	Fields: []forms.Field{
		{
			Name: "status",
			Validators: []forms.Validator{
				forms.IsRequired{},
				api.IsChangeRequestReviewStatus{},
				forms.IsIn{Choices: []interface{}{api.ReviewRequested, api.RequestRejected, api.RequestApproved}},
			},
		},
	},
}

func UpdateChangeRequestReviewStatus(c *gin.Context) {

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

	request := changeRequest(c, controller)

	if request == nil {
		return
	}

	review := changeRequestReview(c, request)

	if review == nil {
		return
	}

	data := helpers.JSONData(c)

	if data == nil {
		return
	}

	params, err := updateChangeRequestReviewStatusForn.Validate(data)

	if err != nil {
		api.HandleError(c, 400, err)
		return
	}

	status := params["status"].(api.ChangeRequestReviewStatus)

	isReviewer, err := controller.CanAccess(user, object, []string{"reviewer", "admin", "superuser"})

	if err != nil {
		api.HandleError(c, 500, err)
		return
	}

	if !isReviewer {
		api.HandleError(c, 400, fmt.Errorf("cannot set status from %s to %s", request.Status(), status))
		return
	}

	if err := review.SetStatus(status); err != nil {
		api.HandleError(c, 500, fmt.Errorf("cannot update change request review status: %v", err))
		return
	}

	c.JSON(200, map[string]interface{}{"data": review})

}

func CreateChangeRequestReview(c *gin.Context) {

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

	request := changeRequest(c, controller)

	if request == nil {
		return
	}

	// we check that the user can actually edit the object
	if ok, err := controller.CanAccess(user, object, []string{"editor", "admin", "superuser"}); err != nil {
		api.HandleError(c, 500, err)
		return
	} else if !ok {
		api.HandleError(c, 401, fmt.Errorf("you must be an editor, admin or superuser to create a change request"))
		return
	}

	apiUser, err := user.ApiUser(controller)

	if err != nil {
		api.HandleError(c, 500, err)
		return
	}

	review, err := request.MakeReview(apiUser)

	if err != nil {
		api.HandleError(c, 500, err)
		return
	}

	data := helpers.JSONData(c)

	if data == nil {
		return
	}

	if err := review.Create(data); err != nil {
		api.HandleError(c, 400, err)
		return
	}

	if err := review.Save(); err != nil {
		api.HandleError(c, 500, err)
		return
	}

	c.JSON(200, map[string]interface{}{"data": review})

}

var changeRequestReviewIDForm = forms.Form{
	ErrorMsg: "invalid data encountered in the change request review ID form",
	Fields: []forms.Field{
		{
			Name: "review_id",
			Validators: []forms.Validator{
				forms.IsRequired{},
				forms.IsHex{ConvertToBinary: true, Strict: false},
			},
		},
	},
}

func changeRequestReview(c *gin.Context, changeRequest api.ChangeRequest) api.ChangeRequestReview {

	data := map[string]interface{}{
		"review_id": c.Param("reviewID"),
	}

	params, err := changeRequestReviewIDForm.Validate(data)

	if err != nil {
		api.HandleError(c, 400, err)
		return nil
	}

	reviewID := params["review_id"].([]byte)

	review, err := changeRequest.Review(reviewID)

	if err != nil {
		api.HandleError(c, 404, err)
		return nil
	}

	return review

}

// Update a change request
func UpdateChangeRequestReview(c *gin.Context) {

	controller := helpers.Controller(c)

	if controller == nil {
		return
	}

	request := changeRequest(c, controller)

	if request == nil {
		return
	}

	review := changeRequestReview(c, request)

	if review == nil {
		return
	}

	user := helpers.User(c)

	if user == nil {
		return
	}

	// we ensure the user is the creator of the change request
	if !bytes.Equal(review.Creator().SourceID(), user.SourceID) || review.Creator().Source() != review.Creator().Source() {
		api.HandleError(c, 401, fmt.Errorf("you cannot edit this change request review"))
		return
	}

	data := helpers.JSONData(c)

	if data == nil {
		return
	}

	if err := review.Update(data); err != nil {
		api.HandleError(c, 400, err)
		return
	}

	if err := review.Save(); err != nil {
		api.HandleError(c, 500, err)
		return
	}

	c.JSON(200, map[string]interface{}{"data": review})

}

// Delete a change request
func DeleteChangeRequestReview(c *gin.Context) {

	controller := helpers.Controller(c)

	if controller == nil {
		return
	}

	request := changeRequest(c, controller)

	if request == nil {
		return
	}

	review := changeRequestReview(c, request)

	if review == nil {
		return
	}

	user := helpers.User(c)

	if user == nil {
		return
	}

	canDelete := false

	if ok, err := controller.CanAccess(user, request, []string{"editor", "admin", "superuser"}); err != nil {
		api.HandleError(c, 500, err)
		return
	} else if ok {
		canDelete = true
	}

	if bytes.Equal(request.Creator().SourceID(), user.SourceID) && request.Creator().Source() == request.Creator().Source() {
		canDelete = true
	}

	if !canDelete {
		api.HandleError(c, 401, fmt.Errorf("you cannot delete this change request"))
		return
	}

	if err := request.Delete(); err != nil {
		api.HandleError(c, 500, err)
		return
	}

	c.JSON(200, map[string]interface{}{"message": "ok"})

}

// Get a list of change requestse
func ChangeRequestReviews(c *gin.Context) {

	controller := helpers.Controller(c)

	if controller == nil {
		return
	}

	object := GetObj(c, "objectType")

	if object == nil {
		return
	}

	request := changeRequest(c, controller)

	if request == nil {
		return
	}

	reviews, err := request.Reviews()

	if err != nil {
		api.HandleError(c, 500, err)
		return
	}

	c.JSON(200, map[string]interface{}{"data": reviews})

}
