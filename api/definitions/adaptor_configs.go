// Kodex (Enterprise Edition - EE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2021  KIProtect GmbH (HRB 208395B) - All Rights Reserved

package definitions

import (
	"github.com/gin-gonic/gin"
	"github.com/kiprotect/kodex"
	"github.com/kiprotect/kodex/api"
	"github.com/kiprotect/kodex/api/helpers"
)

type ConfigAdaptor struct{}

func getStream(c *gin.Context) kodex.Stream {
	stream, ok := c.Get("stream")

	handleError := func() {
	}

	if !ok {
		handleError()
		return nil
	}

	kiprotectStream, ok := stream.(kodex.Stream)

	if !ok {
		handleError()
		return nil
	}

	return kiprotectStream

}

func (f ConfigAdaptor) Type() string {
	return "config"
}

func (f ConfigAdaptor) DependsOn() string {
	return "stream"
}

func (f ConfigAdaptor) Initialize(controller api.Controller, g *gin.RouterGroup) error {
	return nil
}

func (f ConfigAdaptor) Get(controller api.Controller, c *gin.Context, id []byte) (kodex.Model, kodex.Model, error) {
	object, err := controller.Config(id)
	if err == nil {
		return object, object.Stream().Project(), nil
	}
	return nil, nil, err
}

func (a ConfigAdaptor) Objects(c *gin.Context) []kodex.Model {

	stream := getStream(c)

	if stream == nil {
		return nil
	}

	configs, err := stream.Configs()

	if err != nil {
		api.HandleError(c, 500, err)
		return nil
	}

	objects := make([]kodex.Model, len(configs))
	for i, config := range configs {
		objects[i] = config
	}
	return objects

}

func (a ConfigAdaptor) MakeObject(c *gin.Context) kodex.Model {

	controller := helpers.Controller(c)

	if controller == nil {
		return nil
	}

	stream := getStream(c)

	if stream == nil {
		return nil
	}

	return stream.MakeConfig()
}
