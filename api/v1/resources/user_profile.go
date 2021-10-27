// Kodex (Enterprise Edition - EE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2021  KIProtect GmbH (HRB 208395B) - All Rights Reserved

package resources

import (
	"github.com/gin-gonic/gin"
	"github.com/kiprotect/kodex/api/helpers"
)

// Return the profile of the logged in user
func UserProfile(c *gin.Context) {

	controller := helpers.Controller(c)

	if controller == nil {
		return
	}

	userProfile := helpers.UserProfile(c)

	if userProfile == nil {
		return
	}

	c.JSON(200, map[string]interface{}{"data": userProfile})

}
