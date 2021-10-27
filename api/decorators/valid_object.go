// Kodex (Enterprise Edition - EE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2021  KIProtect GmbH (HRB 208395B) - All Rights Reserved

package decorators

import (
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

		var userProfile api.UserProfile

		if len(objectRoles) > 0 || len(scopes) > 0 {
			up, ok := c.Get("userProfile")

			if !ok {
				api.HandleError(c, 401, fmt.Errorf("unauthorized"))
				return
			}

			userProfile, ok = up.(api.UserProfile)

			if !ok {
				api.HandleError(c, 500, fmt.Errorf("corrupt user profile"))
				return
			}

			if len(scopes) > 0 && !CheckScopes(scopes, userProfile.AccessToken().Scopes()) {
				api.HandleError(c, 403, errors.MakeExternalError("access denied", "ACCESS-DENIED", map[string]interface{}{"user_scopes": userProfile.AccessToken().Scopes(), "required_scopes": scopes}, nil))
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

		object, roleObject, err = adaptor.Get(apiController, c, objectID)

		if err != nil {
			if cErr, ok := err.(errors.ChainableError); ok && cErr.Code() == "NOT-FOUND" {
				api.HandleError(c, 404, fmt.Errorf("object not found"))
			} else {
				api.HandleError(c, 500, err)
			}
			return
		}

		if len(objectRoles) > 0 {
			if ok, err := apiController.CanAccess(userProfile, roleObject, objectRoles); !ok || err != nil {
				if err != nil {
					api.HandleError(c, 500, err)
				} else {
					api.HandleError(c, 404, fmt.Errorf("object not found"))
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
