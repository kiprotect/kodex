// Kodex (Enterprise Edition - EE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2021  KIProtect GmbH (HRB 208395B) - All Rights Reserved

package helpers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kiprotect/kodex"
	"github.com/kiprotect/kodex/api"
)

func GetProject(c *gin.Context) kodex.Project {
	project, ok := c.Get("project")

	handleError := func() {
		api.HandleError(c, 500, fmt.Errorf("cannot load project"))
	}

	if !ok {
		handleError()
		return nil
	}

	kodexProject, ok := project.(kodex.Project)

	if !ok {
		handleError()
		return nil
	}

	return kodexProject

}
