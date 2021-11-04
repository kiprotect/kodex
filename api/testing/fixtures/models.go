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
