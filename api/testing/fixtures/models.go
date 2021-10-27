// Kodex (Enterprise Edition - EE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2021  KIProtect GmbH (HRB 208395B) - All Rights Reserved

package fixtures

import (
	"fmt"
	"github.com/kiprotect/kodex"
	"github.com/kiprotect/kodex/api"
)

type ObjectRole struct {
	ObjectName       string
	OrganizationRole string
	ObjectRole       string
	Organization     string
}

func (o ObjectRole) Setup(fixtures map[string]interface{}) (interface{}, error) {
	controller, err := GetController(fixtures)
	if err != nil {
		return nil, err
	}

	org, ok := fixtures[o.Organization].(api.Organization)

	if !ok {
		return nil, fmt.Errorf("organization %s not found", o.Organization)
	}

	object, ok := fixtures[o.ObjectName].(kodex.Model)

	if !ok {
		return nil, fmt.Errorf("object %s not found", o.ObjectName)
	}

	objectRole := controller.MakeObjectRole(object, org)

	values := map[string]interface{}{
		"organization_role": o.OrganizationRole,
		"role":              o.ObjectRole,
	}

	if err := objectRole.Create(values); err != nil {
		return nil, err
	}

	if err := objectRole.Save(); err != nil {
		return nil, err
	}

	return objectRole, nil

}

func (o ObjectRole) Teardown(fixture interface{}) error {
	return nil
}
