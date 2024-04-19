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
		"project.id": project.ID(),
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

	return project.MakeActionConfig(nil)
}
