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

	return stream.MakeConfig(nil)
}
