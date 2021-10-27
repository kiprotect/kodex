// Kodex (Enterprise Edition - EE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2021  KIProtect GmbH (HRB 208395B) - All Rights Reserved

package testing

import (
	"github.com/gin-gonic/gin"
	"github.com/kiprotect/kodex/api"
	ginHelpers "github.com/kiprotect/kodex/api/helpers/gin"
)

func Router(controller api.Controller, decorator gin.HandlerFunc) (*gin.Engine, error) {
	return ginHelpers.Router(controller, decorator)
}
