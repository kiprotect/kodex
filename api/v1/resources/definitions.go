// Kodex (Enterprise Edition - EE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2021  KIProtect GmbH (HRB 208395B) - All Rights Reserved

package resources

import (
	"github.com/gin-gonic/gin"
	"github.com/kiprotect/kodex/api/helpers"
)

// Get the definitions
func Definitions(c *gin.Context) {

	controller := helpers.Controller(c)

	if controller == nil {
		return
	}

	c.JSON(200, map[string]interface{}{"data": controller.APIDefinitions()})

}
