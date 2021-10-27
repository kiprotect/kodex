// Kodex (Enterprise Edition - EE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2021  KIProtect GmbH (HRB 208395B) - All Rights Reserved

package api

import (
	"bytes"
	"github.com/kiprotect/kodex"
)

type BaseController struct {
	Definitions_ *Definitions
	Self         Controller
}

func (b *BaseController) APIDefinitions() *Definitions {
	return b.Definitions_
}

func (b *BaseController) ObjectRolesForUser(objectType string, user UserProfile) ([]ObjectRole, error) {
	objectRoles := make([]ObjectRole, 0)
	for _, organizationRoles := range user.Roles() {
		// we retrieve the organization of the user
		apiOrg, err := organizationRoles.Organization().ApiOrganization(b.Self)
		if err != nil {
			return nil, err
		}
		newRoles, err := b.Self.ObjectRolesForOrganizationRoles(objectType, organizationRoles.Roles(), apiOrg.ID())
		if err != nil {
			return nil, err
		}
		objectRoles = append(objectRoles, newRoles...)
	}
	return objectRoles, nil
}

func (b *BaseController) CanAccess(user UserProfile, object kodex.Model, objectRoles []string) (bool, error) {

	// we retrive all organization roles for this object
	roles, err := b.Self.RolesForObject(object)

	kodex.Log.Info(roles)

	if err != nil {
		return false, err
	}

	for _, organizationRoles := range user.Roles() {

		apiOrg, err := organizationRoles.Organization().ApiOrganization(b.Self)
		if err != nil {
			return false, err
		}
		organizationID := apiOrg.ID()
		userRoles := organizationRoles.Roles()

		for _, role := range roles {
			if !bytes.Equal(organizationID, role.OrganizationID()) {
				// this role is for another organization
				continue
			}
			for _, userRole := range userRoles {
				// users with the "superuser" role always have access to the object
				if userRole == "superuser" {
					return true, nil
				}
				if userRole == role.OrganizationRole() {
					// the user has a matching role
					for _, objectRole := range objectRoles {
						if objectRole == role.ObjectRole() {
							// there's a matching object role
							return true, nil
						}
					}
				}
			}
		}
	}

	return false, nil
}

func (b *BaseController) RegisterAPIPlugin(plugin APIPlugin) error {
	b.Definitions_.Routes = append(b.Definitions_.Routes, plugin.InitializeAPI)
	if err := plugin.InitializeAdaptors(b.Definitions_.ObjectAdaptors); err != nil {
		return err
	}
	return nil
}
