package web

import (
	. "github.com/gospel-sh/gospel"
	"github.com/kiprotect/kodex"
	"github.com/kiprotect/kodex/web/ui"
)

func ProjectRolesRoutes(project kodex.Project) ElementFunction {
	return func(c Context) Element {
		return UseRouter(c).Match(
			c,
			Route("^/?$", ProjectRoles(project)),
			Route("/details/(?P<roleId>[^/]+)", ProjectRoleDetails(project)),
			Route("/new", NewProjectRole(project)),
		)
	}
}

func ProjectRoles(project kodex.Project) ElementFunction {

	return func(c Context) Element {
		controller := UseController(c)
		// organization := UseDefaultOrganization(c)
		roles, err := controller.RolesForObject(project)

		if err != nil {
			return Div("cannot get organization")
		}

		roleItems := make([]Element, len(roles))

		for i, role := range roles {
			roleItems[i] = A(
				Href(Fmt("/flows/projects/%s/settings/roles/details/%s", Hex(project.ID()), Hex(role.ID()))),
				ui.ListItem(
					ui.ListColumn("md", role.OrganizationRole()),
					ui.ListColumn("md", role.ObjectRole()),
				),
			)
		}

		return F(
			ui.List(
				ui.ListHeader(
					ui.ListColumn("md", "Organization Role"),
					ui.ListColumn("md", "Object Role"),
				),
				roleItems,
			),
			A(Href(Fmt("/flows/projects/%s/settings/roles/new", Hex(project.ID()))), Class("bulma-button", "bulma-is-success"), "New Role"),
		)

		return Ul(
			roleItems,
		)

	}
}
