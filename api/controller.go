// Kodex (Enterprise Edition - EE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2021  KIProtect GmbH (HRB 208395B) - All Rights Reserved

package api

import (
	"github.com/kiprotect/kodex"
)

type APIControllerMaker func(
	config map[string]interface{},
	baseController kodex.Controller,
	definitions *Definitions) (Controller, error)

type APIControllerDefinitions map[string]APIControllerMaker

type Controller interface {
	kodex.Controller

	KodexController() kodex.Controller
	APIDefinitions() *Definitions
	RegisterAPIPlugin(APIPlugin) error

	// Object roles
	CanAccess(user UserProfile, object kodex.Model, objectRoles []string) (bool, error)
	ObjectRole(id []byte) (ObjectRole, error)
	RolesForObject(object kodex.Model) ([]ObjectRole, error)
	ObjectRolesForUser(objectType string, user UserProfile) ([]ObjectRole, error)
	ObjectRolesForOrganizationRoles(objectType string, organizationRoles []string, organizationID []byte) ([]ObjectRole, error)
	MakeObjectRole(object kodex.Model, organization Organization) ObjectRole

	// Organizations
	MakeOrganization() Organization
	Organization(source string, sourceID []byte) (Organization, error)
	Organizations(filters map[string]interface{}) ([]Organization, error)
}
