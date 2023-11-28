package web

import (
	"bytes"
	. "github.com/gospel-sh/gospel"
	"github.com/kiprotect/kodex"
	"github.com/kiprotect/kodex/web/ui"
)

func ProjectRoleDetails(project kodex.Project) any {

	return func(c Context, roleId string) Element {
		router := UseRouter(c)

		controller := UseController(c)
		role, err := controller.ObjectRole(Unhex(roleId))

		if err != nil {
			return Div("cannot load default role")
		}

		if !bytes.Equal(role.ObjectID(), project.ID()) {
			return Div("illegal object")
		}

		onSubmit := Func[any](c, func() {
			role.Delete()
			router.RedirectTo(Fmt("/flows/projects/%s/settings/roles", Hex(project.ID())))
		})

		return F(
			H1(
				Class("bulma-subtitle"),
				Fmt("Role Details - %s", roleId),
			),
			Table(
				Class("bulma-table"),
				Thead(
					Tr(
						Th("Key"),
						Th("Value"),
					),
				),
				Tbody(
					Tr(
						Td("Organization role"),
						Td(role.OrganizationRole()),
					),
					Tr(
						Td("Object role"),
						Td(role.ObjectRole()),
					),
				),
			),
			Hr(),
			A(
				Class("bulma-button", "bulma-is-danger"),
				Href(Fmt("/flows/projects/%s/settings/roles/details/%s/delete", Hex(project.ID()), Hex(role.ID()))),
				"delete role",
			),
			router.Match(
				c,
				Route("/delete$",
					func(c Context) Element {
						return ui.Modal(
							c,
							"Do you really want to delete this role?",
							Span(
								"Do you really want to delete this role?",
							),
							F(
								A(
									Class("bulma-button"),
									Href(Fmt("/flows/projects/%s/settings/roles/details/%s", Hex(project.ID()), Hex(role.ID()))),
									"Cancel",
								),
								Span(Style("flex-grow: 1")),
								Span(
									Form(
										Class("bulma-is-inline"),
										Method("POST"),
										OnSubmit(onSubmit),
										Button(
											Name("action"),
											Value("edit"),
											Class("bulma-button", "bulma-is-danger"),
											Type("submit"),
											"Yes, delete",
										),
									),
								),
							),
							Fmt("/flows/projects/%s/settings/roles/details/%s", Hex(project.ID()), Hex(role.ID())),
						)
					},
				),
			),
		)

	}
}
