// Kodex (Enterprise Edition - EE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2021  KIProtect GmbH (HRB 208395B) - All Rights Reserved

package decorators

import (
	"github.com/gin-gonic/gin"
	"github.com/kiprotect/kodex"
)

func WithSettings(settings kodex.Settings) gin.HandlerFunc {

	/*
	   This decorator adds a reference to the settings object to the request context.
	*/

	decorator := func(c *gin.Context) {
		c.Set("settings", settings)
	}
	return decorator
}
