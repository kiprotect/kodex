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

package web

import (
	. "github.com/gospel-sh/gospel"
	"github.com/kiprotect/kodex/web/ui"
)

func DefaultRoles(c Context) Element {

	controller := UseController(c)
	organization := UseDefaultOrganization(c)
	apiOrg, err := organization.ApiOrganization(controller)

	if err != nil {
		return Div("cannot get organization")
	}

	roles, err := controller.DefaultObjectRoles(apiOrg.ID())

	if err != nil {
		return Div("cannot load default object roles")
	}

	roleItems := make([]Element, len(roles))

	for i, role := range roles {
		roleItems[i] = A(
			Href(Fmt("/admin/roles/details/%s", Hex(role.ID()))),
			ui.ListItem(
				ui.ListColumn("md", role.OrganizationRole()),
				ui.ListColumn("md", role.ObjectRole()),
				ui.ListColumn("md", role.ObjectType()),
			),
		)
	}

	return F(
		ui.List(
			ui.ListHeader(
				ui.ListColumn("md", "Organization Role"),
				ui.ListColumn("md", "Object Role"),
				ui.ListColumn("md", "Object Type"),
			),
			roleItems,
		),
		A(Href("/admin/roles/new"), Class("bulma-button", "bulma-is-success"), "New Role"),
	)

	return Ul(
		roleItems,
	)
}
