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
)

func DefaultRolesRoutes(c Context) Element {
	return UseRouter(c).Match(
		c,
		Route("^/?$", DefaultRoles),
		Route("/details/(?P<mappingId>[^/]+)", DefaultRoleDetails),
		Route("/new", NewDefaultRole),
	)
}

func Admin(c Context) Element {

	user := UseExternalUser(c)

	superuser := user.HasRole(nil, "superuser")

	if !superuser {
		return Div("Administration is for superusers only")
	}

	AddBreadcrumb(c, "Admin", "/admin")

	AddSidebarItem(c, &SidebarItem{Title: "Admin", Path: "/admin", Icon: "hammer", Submenu: []*SidebarItem{
		{
			Title: "Default Roles",
			Path:  "/admin/roles",
			Icon:  "users",
		},
	}})

	router := UseRouter(c)

	return router.Match(c,
		Route("/roles", DefaultRolesRoutes),
		Route("^$", F(
			H1(
				Class("bulma-title"),
				"Administrative Settings",
			),
			P(
				Class("bulma-text"),
				"Here you can manage administrative settings. Please select a menu point from the left to continue.",
			),
		),
		),
	)
}
