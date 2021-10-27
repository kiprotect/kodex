// Kodex (Enterprise Edition - EE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2021  KIProtect GmbH (HRB 208395B) - All Rights Reserved

package api

import (
	"github.com/gin-gonic/gin"
	"github.com/kiprotect/kodex"
)

type APIPlugin interface {
	InitializeAPI(*gin.RouterGroup, Controller, kodex.Meter) error
	InitializeAdaptors(map[string]ObjectAdaptor) error
}
