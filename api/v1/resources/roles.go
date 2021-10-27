// Kodex (Community Edition - CE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2021  KIProtect GmbH (HRB 208395B) - Germany
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
	"github.com/kiprotect/kodex"
	"github.com/kiprotect/kodex/api"
	"github.com/kiprotect/kodex/api/helpers"
)

// Return the object that was initialized via the ValidObject decorator
func GetObject(c *gin.Context) kodex.Model {
	objObj, ok := c.Get("roleObject")

	if !ok {
		api.HandleError(c, 404, fmt.Errorf("object not found"))
		return nil
	}

	obj, ok := objObj.(kodex.Model)

	if !ok {
		api.HandleError(c, 500, fmt.Errorf("invalid object"))
		return nil
	}

	return obj

}

// Create a new object role
func CreateObjectRole(c *gin.Context) {

	controller := helpers.Controller(c)

	if controller == nil {
		return
	}

	obj := GetObject(c)

	if obj == nil {
		return
	}

	organization := helpers.Organization(c)

	if organization == nil {
		return
	}

	data := helpers.JSONData(c)

	if data == nil {
		return
	}

	role := controller.MakeObjectRole(obj, organization)

	if err := role.Create(data); err != nil {
		api.HandleError(c, 400, err)
		return
	}

	if err := role.Save(); err != nil {
		api.HandleError(c, 500, err)
		return
	}

	c.JSON(200, map[string]interface{}{"data": role})

}

var deleteObjectRoleForm = forms.Form{
	ErrorMsg: "invalid data encountered in the role deletion form",
	Fields: []forms.Field{
		{
			Name: "role_id",
			Validators: []forms.Validator{
				forms.IsRequired{},
				forms.IsHex{ConvertToBinary: true, Strict: false},
			},
		},
	},
}

// Delete a object role
func DeleteObjectRole(c *gin.Context) {
	controller := helpers.Controller(c)
	user := helpers.UserProfile(c)

	if controller == nil || user == nil {
		return
	}

	obj := GetObject(c)

	if obj == nil {
		return
	}

	data := map[string]interface{}{
		"role_id": c.Param("roleID"),
	}

	params, err := deleteObjectRoleForm.Validate(data)

	if err != nil {
		api.HandleError(c, 400, err)
		return
	}

	roleID := params["role_id"].([]byte)

	roles, err := controller.RolesForObject(obj)

	if err != nil {
		api.HandleError(c, 500, err)
		return
	}

	// we ensure that the user does not cut off his/her own access to the object
	userRoleCount := 0
	roleAffected := false
	for _, role := range roles {
		for _, organizationRoles := range user.Roles() {
			apiOrg, err := organizationRoles.Organization().ApiOrganization(controller)
			if err != nil {
				api.HandleError(c, 500, err)
				return
			}
			if !bytes.Equal(apiOrg.ID(), role.OrganizationID()) {
				continue
			}
			for _, userRole := range organizationRoles.Roles() {
				if role.OrganizationRole() == userRole && role.ObjectRole() == "superuser" {
					userRoleCount++
					if bytes.Equal(role.ID(), roleID) {
						roleAffected = true
					}
				}
			}
		}
	}

	if roleAffected && userRoleCount == 1 {
		api.HandleError(c, 400, fmt.Errorf("you cannot delete this role as this would cut off your access to this object"))
		return
	}

	for _, role := range roles {
		if bytes.Equal(role.ID(), roleID) {
			if err := role.Delete(); err != nil {
				api.HandleError(c, 500, err)
				return
			}
			break
		}
	}

}

// Get a list of object roles
func ObjectRoles(c *gin.Context) {

	controller := helpers.Controller(c)

	if controller == nil {
		return
	}

	obj := GetObject(c)

	if obj == nil {
		return
	}

	roles, err := controller.RolesForObject(obj)

	if err != nil {
		api.HandleError(c, 500, err)
		return
	}

	c.JSON(200, map[string]interface{}{"data": roles})

}
