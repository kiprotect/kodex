// Kodex (Enterprise Edition - EE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2021  KIProtect GmbH (HRB 208395B) - All Rights Reserved

package helpers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kiprotect/kodex/api"
)

// Get the controller and user profile (as created by the decorators)
func Controller(c *gin.Context) api.Controller {

	controller, ok := c.Get("controller")
	if !ok {
		api.HandleError(c, 500, fmt.Errorf("no controller defined (API controller check)"))
		return nil
	}

	apiController, ok := controller.(api.Controller)
	if !ok {
		api.HandleError(c, 500, fmt.Errorf("not an API controller"))
		return nil
	}

	return apiController

}
