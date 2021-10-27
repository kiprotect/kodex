// Kodex (Enterprise Edition - EE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2021  KIProtect GmbH (HRB 208395B) - All Rights Reserved

package resources

import (
	"github.com/gin-gonic/gin"
	"github.com/kiprotect/kodex/api/helpers"
)

// Simply return a 200 if the access token works
func SayHello(c *gin.Context) {

	user := helpers.UserProfile(c)

	if user == nil {
		return
	}

	c.JSON(200, map[string]interface{}{"data": map[string]interface{}{
		"user": user,
	}})

}
