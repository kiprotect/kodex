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
	"github.com/gin-gonic/gin"
	"github.com/kiprotect/kodex"
	"github.com/kiprotect/kodex/api"
	"github.com/kiprotect/kodex/api/helpers"
)

// Delete a given object.
func DeleteObject(c *gin.Context) {

	controller := helpers.Controller(c)

	if controller == nil {
		return
	}

	object := GetObj(c, "objectType")

	if object == nil {
		return
	}

	objectRoles, err := controller.RolesForObject(object)

	if err != nil {
		api.HandleError(c, 500, err)
		return
	}

	if err := object.Delete(); err != nil {
		api.HandleError(c, 500, err)
		return
	}

	for _, role := range objectRoles {
		if err := role.Delete(); err != nil {
			api.HandleError(c, 500, err)
			return
		}
	}

	c.JSON(200, map[string]interface{}{"message": "success"})
}

func CreateObject(c *gin.Context) {

	controller := helpers.Controller(c)

	if controller == nil {
		return
	}

	roleObj, ok := getRoleObj(c)

	if !ok {
		return
	}

	var organization api.Organization

	if roleObj == nil {
		organization = helpers.Organization(c)

		if organization == nil {
			return
		}
	}

	data := helpers.JSONData(c)

	if data == nil {
		return
	}

	object := makeObject(c)

	if object == nil {
		return
	}

	if err := object.Create(data); err != nil {
		api.HandleError(c, 400, err)
		return
	}

	if err := object.Save(); err != nil {
		api.HandleError(c, 500, err)
		return
	}

	handleError := func(err error) {
		if err := object.Delete(); err != nil {
			kodex.Log.Error("cannot delete object")
		}
		api.HandleError(c, 500, err)
		return
	}

	if roleObj == nil {

		// we always add admin and superuser roles
		for _, orgRole := range []string{"admin", "superuser"} {
			role := controller.MakeObjectRole(object, organization)
			values := map[string]interface{}{
				"organization_role": orgRole,
				"role":              "superuser",
			}

			if err := role.Create(values); err != nil {
				handleError(err)
				return
			}
			if err := role.Save(); err != nil {
				handleError(err)
				return
			}
		}

		// we try to add default roles as well
		if defaultRoles, err := controller.DefaultObjectRoles(organization.ID()); err != nil {
			kodex.Log.Errorf("Cannot load default roles: %v", err)
		} else {
			for _, defaultRole := range defaultRoles {
				if defaultRole.ObjectType() != object.Type() {
					continue
				}

				role := controller.MakeObjectRole(object, organization)

				values := map[string]interface{}{
					"organization_role": defaultRole.OrganizationRole(),
					"role":              defaultRole.ObjectRole(),
				}

				if err := role.Create(values); err != nil {
					handleError(err)
					return
				}
				if err := role.Save(); err != nil {
					handleError(err)
					return
				}

			}
		}

	}

	c.JSON(200, map[string]interface{}{"data": object})

}

func ObjectDetails(c *gin.Context) {
	object := GetObj(c, "objectType")

	if object == nil {
		return
	}

	c.JSON(200, map[string]interface{}{"message": "success", "data": object})
}

// Update a object
func UpdateObject(c *gin.Context) {

	object := GetObj(c, "objectType")

	if object == nil {
		return
	}

	data := helpers.JSONData(c)

	if data == nil {
		return
	}

	adaptor := getAdaptor(c)

	if adaptor == nil {
		return
	}

	if updateAdaptor, ok := adaptor.(api.UpdateObjectAdaptor); ok {
		// this object implements its own update flow
		if updateObject, err := updateAdaptor.UpdateObject(object, data); err != nil {
			api.HandleError(c, 400, err)
			return
		} else {
			if updateObject == nil {
				if err := object.Save(); err != nil {
					api.HandleError(c, 500, err)
					return
				}
				c.JSON(200, map[string]interface{}{"message": "success", "data": object})
			} else {
				if err := updateAdaptor.SaveUpdated(updateObject, object); err != nil {
					api.HandleError(c, 500, err)
					return
				}
				c.JSON(200, map[string]interface{}{"message": "success", "data": updateObject})
			}
		}
	} else {
		err := object.Update(data)

		if err != nil {
			api.HandleError(c, 400, err)
			return
		}

		if err := object.Save(); err != nil {
			api.HandleError(c, 500, err)
			return
		}
		c.JSON(200, map[string]interface{}{"message": "success", "data": object})
	}
}

type JSONDataObject interface {
	ListJSONData() map[string]interface{}
}

// Get a list of objects
func Objects(c *gin.Context) {

	objects := GetObjs(c)

	if objects == nil {
		return
	}

	if len(objects) > 0 {
		if _, ok := objects[0].(JSONDataObject); ok {
			listObjects := make([]map[string]interface{}, len(objects))
			for i, obj := range objects {
				listObjects[i] = obj.(JSONDataObject).ListJSONData()
			}
			c.JSON(200, map[string]interface{}{"data": listObjects})
			return
		}
	}

	c.JSON(200, map[string]interface{}{"data": objects})

}

// Get a list of objects
func AllObjects(c *gin.Context) {

	objects := getAllObjs(c)

	if objects == nil {
		return
	}

	c.JSON(200, map[string]interface{}{"data": objects})

}

func AssociatedObjects(c *gin.Context) {
	left := GetObj(c, "leftType")

	if left == nil {
		return
	}

	adaptor := getAssociateAdaptor(c)

	if adaptor == nil {
		return
	}

	if objs := adaptor.Get(c, left); objs != nil {
		c.JSON(200, map[string]interface{}{"message": "success", "data": objs})
	}

}

func AssociateObjects(c *gin.Context) {
	left := GetObj(c, "leftType")
	right := GetObj(c, "rightType")

	if left == nil || right == nil {
		return
	}

	adaptor := getAssociateAdaptor(c)

	if adaptor == nil {
		return
	}

	if ok := adaptor.Associate(c, left, right); ok {
		c.JSON(200, map[string]interface{}{"message": "success"})
	}

}

func DissociateObjects(c *gin.Context) {

	left := GetObj(c, "leftType")
	right := GetObj(c, "rightType")

	if left == nil || right == nil {
		return
	}

	adaptor := getAssociateAdaptor(c)

	if adaptor == nil {
		return
	}

	if ok := adaptor.Dissociate(c, left, right); ok {
		c.JSON(200, map[string]interface{}{"message": "success"})
	} else {
		c.JSON(404, map[string]interface{}{"message": "objects not associated"})
	}

}
