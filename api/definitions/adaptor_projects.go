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

type ProjectAdaptor struct{}

func (f ProjectAdaptor) Type() string {
	return "project"
}

func (f ProjectAdaptor) DependsOn() string {
	return ""
}

func (f ProjectAdaptor) Initialize(controller api.Controller, g *gin.RouterGroup) error {
	return nil
}

func (f ProjectAdaptor) Get(controller api.Controller, c *gin.Context, id []byte) (kodex.Model, kodex.Model, error) {
	object, err := controller.Project(id)
	if err == nil {
		return object, object, nil
	}
	return nil, nil, err
}

func (a ProjectAdaptor) Objects(c *gin.Context) []kodex.Model {

	controller := helpers.Controller(c)
	user := helpers.User(c)

	if controller == nil || user == nil {
		return nil
	}

	objectRoles, err := controller.ObjectRolesForUser("project", user)

	if err != nil {
		api.HandleError(c, 500, err)
		return nil
	}

	ids := make([]interface{}, len(objectRoles))

	for i, role := range objectRoles {
		ids[i] = role.ObjectID()
	}

	projects, err := controller.Projects(map[string]interface{}{
		"ID": api.In{Values: ids},
	})

	if err != nil {
		api.HandleError(c, 400, err)
		return nil
	}

	objects := make([]kodex.Model, len(projects))
	for i, project := range projects {
		objects[i] = project
	}
	return objects

}

func (a ProjectAdaptor) MakeObject(c *gin.Context) kodex.Model {

	controller := helpers.Controller(c)

	if controller == nil {
		return nil
	}

	return controller.MakeProject(nil)
}
