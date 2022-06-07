// Kodex (Community Edition - CE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2022  KIProtect GmbH (HRB 208395B) - Germany
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
		"project.id": project.ID(),
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
