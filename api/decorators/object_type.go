// Kodex (Enterprise Edition - EE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2021  KIProtect GmbH (HRB 208395B) - All Rights Reserved

package decorators

import (
	"github.com/gin-gonic/gin"
)

func ObjectType(objectType string) gin.HandlerFunc {

	return func(c *gin.Context) {
		c.Set("objectType", objectType)
	}
}
