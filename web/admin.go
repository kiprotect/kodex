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

	superuser := false

	// we check if the user is a superuser
outer:
	for _, role := range user.Roles {
		if role.Organization.Default {
			for _, userRole := range role.Roles {
				if userRole == "superuser" {
					superuser = true
					break outer
				}
			}
			break outer
		}
	}

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
		Route("", F(
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
