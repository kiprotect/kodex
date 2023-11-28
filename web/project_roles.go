package web

import (
	. "github.com/gospel-sh/gospel"
	"github.com/kiprotect/kodex/web/ui"
)

func ProjectRoles(c Context) Element {

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
