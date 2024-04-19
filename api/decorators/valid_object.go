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

package decorators

import (
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kiprotect/go-helpers/errors"
	"github.com/kiprotect/go-helpers/forms"
	"github.com/kiprotect/kodex"
	"github.com/kiprotect/kodex/api"
)

var ObjectIDForm = forms.Form{
	ErrorMsg: "invalid data encountered in the object ID form",
	Fields: []forms.Field{
		{
			Name: "object_id",
			Validators: []forms.Validator{
				forms.IsRequired{},
				forms.IsHex{ConvertToBinary: true},
			},
		},
	},
}

func ValidObject(settings kodex.Settings, objectType string, objectRoles []string, scopes []string) gin.HandlerFunc {

	return func(c *gin.Context) {

		var user *api.ExternalUser

		if len(objectRoles) > 0 || len(scopes) > 0 {
			up, ok := c.Get("user")

			if !ok {
				api.HandleError(c, 401, fmt.Errorf("unauthorized"))
				return
			}

			user, ok = up.(*api.ExternalUser)

			if !ok {
				api.HandleError(c, 500, fmt.Errorf("corrupt user"))
				return
			}
		}

		params, err := ObjectIDForm.Validate(map[string]interface{}{
			"object_id": c.Param(fmt.Sprintf("%sID", objectType)),
		})

		if err != nil {
			api.HandleError(c, 400, err)
			return
		}

		objectID := params["object_id"].([]byte)

		controller, ok := c.Get("controller")
		if !ok {
			api.HandleError(c, 500, fmt.Errorf("no controller defined (valid object check)"))
			return
		}

		apiController, ok := controller.(api.Controller)
		if !ok {
			api.HandleError(c, 500, fmt.Errorf("not an API controller"))
			return
		}

		var object kodex.Model
		var roleObject kodex.Model

		adaptors := apiController.APIDefinitions().ObjectAdaptors

		adaptor, ok := adaptors[objectType]

		if !ok {
			api.HandleError(c, 500, fmt.Errorf("unknown object type"))
			c.Abort()
			return
		}

		notFound := func() {
			api.HandleError(c, 404, fmt.Errorf("object not found"))
		}

		object, roleObject, err = adaptor.Get(apiController, c, objectID)

		if err != nil {
			if err == kodex.NotFound {
				notFound()
			} else {
				api.HandleError(c, 500, err)
			}
			return
		}

		if len(scopes) > 0 {

			fullScopes := []string{}

			// we append the object ID to the scopes, such that it's possible to restrict access to
			// a specific object not only an object type...
			for _, scope := range scopes {
				fullScopes = append(fullScopes, fmt.Sprintf("%s:%s", scope, hex.EncodeToString(objectID)))
			}

			// we check that the access token has the required scopes to access the object
			if !CheckScopes(fullScopes, user.AccessToken.Scopes) {
				api.HandleError(c, 403, errors.MakeExternalError("access denied", "ACCESS-DENIED", map[string]interface{}{"user_scopes": user.AccessToken.Scopes, "required_scopes": fullScopes}, nil))
				return
			}

		}

		if len(objectRoles) > 0 {
			if ok, err := apiController.CanAccess(user, roleObject, objectRoles); !ok || err != nil {
				if err != nil {
					api.HandleError(c, 500, err)
				} else {
					notFound()
				}
				return
			}
		}

		// we set the object
		c.Set(objectType, object)

		// we also set the role object as an object (for the roles endpoints)
		c.Set("roleObject", roleObject)

	}
}
