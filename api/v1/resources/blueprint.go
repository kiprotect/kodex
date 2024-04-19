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

package resources

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kiprotect/go-helpers/forms"
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

	blueprint, err := kodex.ExportBlueprint(project)

	if err != nil {
		api.HandleError(c, 500, err)
		return
	}

	// to do: convert the project to a blueprint....

	c.JSON(200, map[string]interface{}{"message": "success", "data": blueprint})

}

func UploadBlueprint(c *gin.Context) {

	// generate payload from POST Body

	// to do: require a valid project

	data := helpers.JSONData(c)

	if data == nil {
		return
	}

	blueprint := kodex.MakeBlueprint(data)

	controllerObj, ok := c.Get("controller")

	if !ok {
		api.HandleError(c, 500, fmt.Errorf("invalid controller"))
		return
	}

	ctrl, ok := controllerObj.(api.Controller)

	if !ok {
		api.HandleError(c, 500, fmt.Errorf("invalid controller"))
		return
	}

	organization := helpers.Organization(c)

	if organization == nil {
		return
	}

	// we only create a blueprint for an existing project
	project, err := blueprint.Create(ctrl, false)

	if err != nil {
		if _, ok := err.(*forms.FormError); ok {
			api.HandleError(c, 400, err)
		} else {
			api.HandleError(c, 500, err)
		}
		return
	}

	c.JSON(200, map[string]interface{}{"message": "success", "data": project})

}
