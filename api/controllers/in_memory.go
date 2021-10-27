// Kodex (Enterprise Edition - EE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2021  KIProtect GmbH (HRB 208395B) - All Rights Reserved

package controllers

import (
	"bytes"
	"fmt"
	"github.com/kiprotect/kodex"
	"github.com/kiprotect/kodex/api"
	kiprotectControllers "github.com/kiprotect/kodex/controllers"
)

type InMemoryController struct {
	api.BaseController
	objectRoles   map[string]api.ObjectRole
	organizations map[string]api.Organization
	*kiprotectControllers.InMemoryController
}

func MakeInMemoryController(config map[string]interface{}, controller kodex.Controller, definitions *api.Definitions) (api.Controller, error) {
	inMemoryController, ok := controller.(*kiprotectControllers.InMemoryController)
	if !ok {
		return nil, fmt.Errorf("not an InMemory controller")
	}
	apiController := &InMemoryController{
		organizations:      make(map[string]api.Organization),
		objectRoles:        make(map[string]api.ObjectRole),
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

/* Object Role Management */

func (m *InMemoryController) ObjectRole(id []byte) (api.ObjectRole, error) {
	for _, objectRole := range m.objectRoles {
		if bytes.Equal(objectRole.ID(), id) {
			return objectRole, nil
		}
	}
	return nil, fmt.Errorf("not found")
}

func (m *InMemoryController) SaveObjectRole(objectRole *InMemoryObjectRole) error {
	kodex.Log.Info("Saving object role...")
	m.objectRoles[string(objectRole.ID())] = objectRole
	return nil
}

func (m *InMemoryController) MakeObjectRole(object kodex.Model, organization api.Organization) api.ObjectRole {
	return MakeInMemoryObjectRole(kodex.RandomID(), organization.ID(), object.ID(), object.Type(), m)
}

func (m *InMemoryController) RolesForObject(object kodex.Model) ([]api.ObjectRole, error) {
	osrs := make([]api.ObjectRole, 0)

	kodex.Log.Info(m.objectRoles)

	for _, objectRole := range m.objectRoles {
		if bytes.Equal(objectRole.ObjectID(), object.ID()) {
			osrs = append(osrs, objectRole)
		}
	}

	return osrs, nil
}

func (m *InMemoryController) ObjectRolesForOrganizationRoles(objectType string, organizationRoles []string, organizationID []byte) ([]api.ObjectRole, error) {

	osrs := make([]api.ObjectRole, 0)

	for _, objectRole := range m.objectRoles {
		if objectRole.ObjectType() != objectType {
			continue
		}
		if !bytes.Equal(objectRole.OrganizationID(), organizationID) {
			continue
		}

		found := false
		for _, role := range organizationRoles {
			if role == objectRole.OrganizationRole() {
				found = true
				break
			}
		}
		if !found {
			continue
		}
		osrs = append(osrs, objectRole)
	}

	return osrs, nil

}

/* Organization Management */

func (c *InMemoryController) Organizations(filters map[string]interface{}) ([]api.Organization, error) {
	organizations := make([]api.Organization, 0)
outer:
	for _, organization := range c.organizations {
		for key, value := range filters {
			switch key {
			case "name":
				strValue, ok := value.(string)
				if !ok {
					return nil, fmt.Errorf("expected a name")
				}
				if organization.Name() != strValue {
					continue outer
				}
			default:
				return nil, fmt.Errorf("unknown filter key: %s", key)
			}
		}
		organizations = append(organizations, organization)
	}
	return organizations, nil
}

func (c *InMemoryController) Organization(source string, sourceID []byte) (api.Organization, error) {
	for _, organization := range c.organizations {
		if organization.Source() == source && bytes.Equal(organization.SourceID(), sourceID) {
			return organization, nil
		}
	}
	return nil, nil
}

func (c *InMemoryController) SaveOrganization(organization *InMemoryOrganization) error {
	c.organizations[string(organization.ID())] = organization
	return nil
}

func (c *InMemoryController) MakeOrganization() api.Organization {
	return MakeInMemoryOrganization(kodex.RandomID(), c)
}
