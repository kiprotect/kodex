// Kodex (Enterprise Edition - EE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2021  KIProtect GmbH (HRB 208395B) - All Rights Reserved

package definitions

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kiprotect/kodex"
	"github.com/kiprotect/kodex/api"
	"github.com/kiprotect/kodex/api/helpers"
)

type ActionConfigAdaptor struct{}

func (f ActionConfigAdaptor) Type() string {
	return "action"
}

func (f ActionConfigAdaptor) DependsOn() string {
	return "project"
}

func (f ActionConfigAdaptor) Initialize(controller api.Controller, g *gin.RouterGroup) error {
	return nil
}

func (f ActionConfigAdaptor) Get(controller api.Controller, c *gin.Context, id []byte) (kodex.Model, kodex.Model, error) {
	object, err := controller.ActionConfig(id)
	if err == nil {
		return object, object.Project(), nil
	}
	return nil, nil, err
}

func getConfig(c *gin.Context) kodex.Config {
	config, ok := c.Get("config")

	handleError := func() {
		api.HandleError(c, 500, fmt.Errorf("cannot load config"))
	}

	if !ok {
		handleError()
		return nil
	}

	kiprotectConfig, ok := config.(kodex.Config)

	if !ok {
		handleError()
		return nil
	}

	return kiprotectConfig

}

func (a ActionConfigAdaptor) Objects(c *gin.Context) []kodex.Model {

	controller := helpers.Controller(c)

	if controller == nil {
		return nil
	}

	project := helpers.GetProject(c)

	if project == nil {
		return nil
	}

	sources, err := controller.ActionConfigs(map[string]interface{}{
		"action_project_id_project.ext_id": project.ID(),
	})

	if err != nil {
		api.HandleError(c, 400, err)
		return nil
	}

	objects := make([]kodex.Model, len(sources))
	for i, source := range sources {
		objects[i] = source
	}
	return objects

}

func (a ActionConfigAdaptor) MakeObject(c *gin.Context) kodex.Model {
	project := helpers.GetProject(c)

	if project == nil {
		return nil
	}

	return project.MakeActionConfig()
}
