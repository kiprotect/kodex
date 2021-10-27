// Kodex (Enterprise Edition - EE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2021  KIProtect GmbH (HRB 208395B) - All Rights Reserved

package definitions

import (
	"github.com/gin-gonic/gin"
	"github.com/kiprotect/kodex"
	"github.com/kiprotect/kodex/api"
	"github.com/kiprotect/kodex/api/helpers"
)

type DestinationAdaptor struct{}

func (f DestinationAdaptor) Type() string {
	return "destination"
}

func (f DestinationAdaptor) DependsOn() string {
	return "project"
}

func (f DestinationAdaptor) Initialize(controller api.Controller, g *gin.RouterGroup) error {
	return nil
}

func (f DestinationAdaptor) Get(controller api.Controller, c *gin.Context, id []byte) (kodex.Model, kodex.Model, error) {
	object, err := controller.Destination(id)

	if err != nil {
		return nil, nil, err
	}

	return object, object.Project(), err
}

func (a DestinationAdaptor) Objects(c *gin.Context) []kodex.Model {

	controller := helpers.Controller(c)

	if controller == nil {
		return nil
	}

	project := helpers.GetProject(c)

	if project == nil {
		return nil
	}

	destinations, err := controller.Destinations(map[string]interface{}{
		"ProjectID": project.ID(),
	})

	if err != nil {
		api.HandleError(c, 400, err)
		return nil
	}

	objects := make([]kodex.Model, len(destinations))
	for i, destination := range destinations {
		objects[i] = destination
	}
	return objects

}

func (a DestinationAdaptor) MakeObject(c *gin.Context) kodex.Model {

	project := helpers.GetProject(c)

	if project == nil {
		return nil
	}

	return project.MakeDestination()
}
