// Kodex (Enterprise Edition - EE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2021  KIProtect GmbH (HRB 208395B) - All Rights Reserved

package helpers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kiprotect/kodex/api"
)

func Organization(c *gin.Context) api.Organization {
	orgObj, ok := c.Get("org")

	if !ok {
		api.HandleError(c, 500, fmt.Errorf("no organization defined in context"))
		return nil
	}

	org, ok := orgObj.(api.Organization)

	if !ok {
		api.HandleError(c, 500, fmt.Errorf("no organization defined in context"))
		return nil
	}

	return org

}
