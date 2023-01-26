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
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kiprotect/go-helpers/forms"
	"github.com/kiprotect/kodex/api"
	"github.com/kiprotect/kodex/api/helpers"
)

var createDefaultObjectRoleForm = forms.Form{
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

// Create a new object role
func CreateDefaultObjectRole(c *gin.Context) {

	controller := helpers.Controller(c)

	if controller == nil {
		return
	}

	organization := helpers.Organization(c)

	if organization == nil {
		return
	}

	data := map[string]interface{}{
		"object_type": c.Param("objectType"),
	}

	params, err := createDefaultObjectRoleForm.Validate(data)

	if err != nil {
		api.HandleError(c, 400, err)
		return
	}

	objectType := params["object_type"].(string)

	data = helpers.JSONData(c)

	if data == nil {
		return
	}

	role := controller.MakeDefaultObjectRole(objectType, organization)

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

var deleteDefaultObjectRoleForm = forms.Form{
	ErrorMsg: "invalid data encountered in the default role deletion form",
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
func DeleteDefaultObjectRole(c *gin.Context) {
	controller := helpers.Controller(c)
	organization := helpers.Organization(c)

	if controller == nil || organization == nil {
		return
	}

	data := map[string]interface{}{
		"role_id": c.Param("roleID"),
	}

	params, err := deleteDefaultObjectRoleForm.Validate(data)

	if err != nil {
		api.HandleError(c, 400, err)
		return
	}

	roleID := params["role_id"].([]byte)

	role, err := controller.DefaultObjectRole(roleID)

	if err != nil {
		api.HandleError(c, 404, err)
		return
	}

	if !bytes.Equal(organization.ID(), role.OrganizationID()) {
		api.HandleError(c, 404, fmt.Errorf("not found"))
		return
	}

	if err := role.Delete(); err != nil {
		api.HandleError(c, 500, err)
		return
	}

	c.JSON(200, map[string]interface{}{"message": "deleted"})

}

// Get a list of object roles
func DefaultObjectRoles(c *gin.Context) {

	controller := helpers.Controller(c)
	organization := helpers.Organization(c)

	if controller == nil || organization == nil {
		return
	}

	roles, err := controller.DefaultObjectRoles(organization.ID())

	if err != nil {
		api.HandleError(c, 500, err)
		return
	}

	c.JSON(200, map[string]interface{}{"data": roles})

}
