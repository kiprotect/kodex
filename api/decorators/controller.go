// Kodex (Enterprise Edition - EE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2021  KIProtect GmbH (HRB 208395B) - All Rights Reserved

package decorators

import (
	"github.com/gin-gonic/gin"
	"github.com/kiprotect/kodex/api"
)

func WithController(controller api.Controller) gin.HandlerFunc {

	decorator := func(c *gin.Context) {
		c.Set("controller", controller)
	}

	return decorator
}
