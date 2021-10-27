// Kodex (Enterprise Edition - EE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2021  KIProtect GmbH (HRB 208395B) - All Rights Reserved

package definitions

import (
	"github.com/gin-gonic/gin"
	"github.com/kiprotect/go-helpers/forms"
	"github.com/kiprotect/kodex"
	"github.com/kiprotect/kodex/api"
	"github.com/kiprotect/kodex/api/helpers"
)

type SourceAdaptor struct{}

func (f SourceAdaptor) Type() string {
	return "source"
}

func (f SourceAdaptor) DependsOn() string {
	return "project"
}

func (f SourceAdaptor) Form() forms.Form {
	return kodex.SourceForm
}

func (f SourceAdaptor) Initialize(controller api.Controller, g *gin.RouterGroup) error {
	return nil
}

func (f SourceAdaptor) Get(controller api.Controller, c *gin.Context, id []byte) (kodex.Model, kodex.Model, error) {
	object, err := controller.Source(id)

	if err != nil {
		return nil, nil, err
	}

	return object, object.Project(), err
}

func (a SourceAdaptor) Objects(c *gin.Context) []kodex.Model {

	controller := helpers.Controller(c)

	if controller == nil {
		return nil
	}

	project := helpers.GetProject(c)

	if project == nil {
		return nil
	}

	sources, err := controller.Sources(map[string]interface{}{
		"ProjectID": project.ID(),
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

func (a SourceAdaptor) MakeObject(c *gin.Context) kodex.Model {

	project := helpers.GetProject(c)

	if project == nil {
		return nil
	}

	return project.MakeSource()
}
