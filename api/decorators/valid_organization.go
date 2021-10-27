// Kodex (Enterprise Edition - EE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2021  KIProtect GmbH (HRB 208395B) - All Rights Reserved

package decorators

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kiprotect/kodex/api"
)

func ValidOrganization(orgRoles []string) gin.HandlerFunc {

	return func(c *gin.Context) {

		controller, ok := c.Get("controller")
		if !ok {
			api.HandleError(c, 500, fmt.Errorf("no controller defined"))
			return
		}

		apiController, ok := controller.(api.Controller)
		if !ok {
			api.HandleError(c, 500, fmt.Errorf("not an API controller"))
			return
		}

		up, ok := c.Get("userProfile")

		if !ok {
			api.HandleError(c, 401, fmt.Errorf("unauthorized"))
			return
		}

		userProfile, ok := up.(api.UserProfile)

		if !ok {
			api.HandleError(c, 500, fmt.Errorf("corrupt user profile"))
			return
		}

		found := false
		useDefault := false
		var orgID []byte

		orgIDParam := c.Param("organizationID")

		if orgIDParam == "default" {
			// we just use the default organization for the user
			useDefault = true
		} else {
			// we parse the organization ID as a hex value
			params, err := ObjectIDForm.Validate(map[string]interface{}{
				"object_id": orgIDParam,
			})
			if err != nil {
				api.HandleError(c, 400, err)
				return
			}
			orgID = params["object_id"].([]byte)
		}

		var org api.UserOrganization
		var err error

		// to do: take into account the given organization
		for _, organizationRoles := range userProfile.Roles() {
			org = organizationRoles.Organization()
			if err != nil {
				api.HandleError(c, 500, err)
				return
			}
			if useDefault {
				if !org.Default() {
					continue
				}
			} else if !bytes.Equal(org.ID(), orgID) {
				continue
			}
			if len(orgRoles) == 0 {
				found = true
			} else {
				for _, userRole := range organizationRoles.Roles() {
					for _, orgRole := range orgRoles {
						if orgRole == userRole || orgRole == "*" {
							found = true
							break
						}
					}
					if found {
						break
					}
				}
			}
			if found {
				break
			}
		}

		if !found {
			api.HandleError(c, 403, fmt.Errorf("not allowed"))
			return
		} else {
			if apiOrg, err := org.ApiOrganization(apiController); err != nil {
				api.HandleError(c, 500, err)
				return
			} else {
				c.Set("org", apiOrg)
			}
		}

	}
}
