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

type AssociateConfigDestinationAdaptor struct{}

func (a AssociateConfigDestinationAdaptor) LeftType() string {
	return "config"
}

func (a AssociateConfigDestinationAdaptor) RightType() string {
	return "destination"
}

func (a AssociateConfigDestinationAdaptor) Associate(c *gin.Context, left, right kodex.Model) bool {

	config, ok := left.(kodex.Config)

	if !ok {
		api.HandleError(c, 500, fmt.Errorf("config missing"))
		return false
	}

	destination, ok := right.(kodex.Destination)

	if !ok {
		api.HandleError(c, 500, fmt.Errorf("destination missing"))
		return false
	}

	data := helpers.JSONData(c)

	if data == nil {
		return false
	}

	params, err := AddConfigDestinationForm.Validate(data)

	if err != nil {
		api.HandleError(c, 400, err)
		return false
	}

	if err := config.AddDestination(destination, params["name"].(string), kodex.DestinationStatus(params["status"].(string))); err != nil {
		api.HandleError(c, 500, err)
		return false
	}

	c.JSON(200, map[string]interface{}{"message": "success"})

	return true

}

func (a AssociateConfigDestinationAdaptor) Dissociate(c *gin.Context, left, right kodex.Model) bool {

	config, ok := left.(kodex.Config)

	if !ok {
		api.HandleError(c, 500, fmt.Errorf("config missing"))
		return false
	}

	destination, ok := right.(kodex.Destination)

	if !ok {
		api.HandleError(c, 500, fmt.Errorf("destination missing"))
		return false
	}

	if err := config.RemoveDestination(destination); err != nil {
		api.HandleError(c, 500, err)
		return false
	}

	return true
}

func (a AssociateConfigDestinationAdaptor) Get(c *gin.Context, left kodex.Model) interface{} {

	config, ok := left.(kodex.Config)

	if !ok {
		api.HandleError(c, 500, fmt.Errorf("config missing"))
		return nil
	}

	if destinations, err := config.Destinations(); err != nil {
		api.HandleError(c, 500, err)
		return nil
	} else {
		return destinations
	}
}
