// Kodex (Community Edition - CE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2021  KIProtect GmbH (HRB 208395B) - Germany
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

type StreamAdaptor struct{}

func (f StreamAdaptor) Type() string {
	return "stream"
}

func (f StreamAdaptor) DependsOn() string {
	return "project"
}

func (f StreamAdaptor) Initialize(controller api.Controller, g *gin.RouterGroup) error {
	return nil
}

func (f StreamAdaptor) Get(controller api.Controller, c *gin.Context, id []byte) (kodex.Model, kodex.Model, error) {
	object, err := controller.Stream(id)

	if err != nil {
		return nil, nil, err
	}

	return object, object.Project(), err
}

func (a StreamAdaptor) Objects(c *gin.Context) []kodex.Model {

	controller := helpers.Controller(c)

	if controller == nil {
		return nil
	}

	project := helpers.GetProject(c)

	if project == nil {
		return nil
	}

	streams, err := controller.Streams(map[string]interface{}{
		"ProjectID": project.ID(),
	})

	if err != nil {
		api.HandleError(c, 400, err)
		return nil
	}

	objects := make([]kodex.Model, len(streams))
	for i, stream := range streams {
		objects[i] = stream
	}
	return objects

}

func (a StreamAdaptor) MakeObject(c *gin.Context) kodex.Model {

	project := helpers.GetProject(c)

	if project == nil {
		return nil
	}

	return project.MakeStream(nil)
}
