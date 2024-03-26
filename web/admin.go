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
