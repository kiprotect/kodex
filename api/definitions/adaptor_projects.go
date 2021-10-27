// Kodex (Enterprise Edition - EE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2021  KIProtect GmbH (HRB 208395B) - All Rights Reserved

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
	user := helpers.UserProfile(c)

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

	return controller.MakeProject()
}
