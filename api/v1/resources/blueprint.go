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

package resources

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kiprotect/kodex"
	"github.com/kiprotect/kodex/api"
	"github.com/kiprotect/kodex/api/helpers"
)

func GetBlueprint(c *gin.Context) {

	projectObj, ok := c.Get("project")

	if !ok {
		api.HandleError(c, 500, fmt.Errorf("invalid project"))
		return
	}

	project, ok := projectObj.(kodex.Project)
	if !ok {
		api.HandleError(c, 500, fmt.Errorf("invalid project"))
		return
	}

	// to do: convert the project to a blueprint....

	c.JSON(200, map[string]interface{}{"message": "success", "project": project})

}

func UploadBlueprint(c *gin.Context) {

	// generate payload from POST Body

	data := helpers.JSONData(c)

	if data == nil {
		return
	}

	blueprint := kodex.MakeBlueprint(data)

	controllerObj, ok := c.Get("controller")

	if !ok {
		api.HandleError(c, 500, fmt.Errorf("invalid controller"))
	}

	ctrl, ok := controllerObj.(api.Controller)

	if !ok {
		api.HandleError(c, 500, fmt.Errorf("invalid controller"))
	}

	project, err := blueprint.Create(ctrl)

	if err != nil {
		api.HandleError(c, 500, err)
	}

	c.JSON(200, map[string]interface{}{"message": "success", "project": project})

}
