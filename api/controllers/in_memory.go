// Kodex (Enterprise Edition - EE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2021  KIProtect GmbH (HRB 208395B) - All Rights Reserved

package controllers

import (
	"fmt"
	"github.com/kiprotect/kodex"
	"github.com/kiprotect/kodex/api"
	kiprotectControllers "github.com/kiprotect/kodex/controllers"
)

type InMemoryController struct {
	api.BaseController
	*kiprotectControllers.InMemoryController
}

func MakeInMemoryController(config map[string]interface{}, controller kodex.Controller, definitions *api.Definitions) (api.Controller, error) {
	inMemoryController, ok := controller.(*kiprotectControllers.InMemoryController)
	if !ok {
		return nil, fmt.Errorf("not an InMemory controller")
	}
	apiController := &InMemoryController{
		InMemoryController: inMemoryController,
		BaseController: api.BaseController{
			Definitions_: definitions,
		},
	}

	apiController.Self = apiController

	return apiController, nil
}

func (m *InMemoryController) KodexController() kodex.Controller {
	return m.InMemoryController
}

func (m *InMemoryController) CanAccess(user api.UserProfile, object kodex.Model, objectRoles []string) (bool, error) {
	return true, nil
}

/* Object Role Management */

func (m *InMemoryController) ObjectRole(id []byte) (api.ObjectRole, error) {
	return nil, nil
}

func (m *InMemoryController) MakeObjectRole(object kodex.Model, organization api.Organization) api.ObjectRole {
	return MakeInMemoryObjectRole(kodex.RandomID(), organization.ID(), object.ID(), object.Type())
}

func (m *InMemoryController) ObjectRolesForUser(objectType string, user api.UserProfile) ([]api.ObjectRole, error) {
	return nil, nil
}

func (m *InMemoryController) RolesForObject(object kodex.Model) ([]api.ObjectRole, error) {
	return nil, nil
}

func (m *InMemoryController) ObjectRolesForOrganizationRoles(objectType string, organizationRoles []string, organizationID []byte) ([]api.ObjectRole, error) {
	return nil, nil
}

/* Organization Management */

func (c *InMemoryController) Organizations(filters map[string]interface{}) ([]api.Organization, error) {
	return nil, fmt.Errorf("not implemented")
}

func (c *InMemoryController) Organization(source string, sourceID []byte) (api.Organization, error) {
	return nil, nil
}

func (c *InMemoryController) MakeOrganization() api.Organization {
	return MakeInMemoryOrganization(kodex.RandomID(), c)
}
