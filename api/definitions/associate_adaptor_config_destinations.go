// Kodex (Community Edition - CE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2024  KIProtect GmbH (HRB 208395B) - Germany
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

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
