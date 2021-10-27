// Kodex (Enterprise Edition - EE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2021  KIProtect GmbH (HRB 208395B) - All Rights Reserved

package controllers

import (
	"github.com/kiprotect/kodex/api"
)

var Controllers = map[string]api.APIControllerMaker{
	"inMemory": MakeInMemoryController,
}
