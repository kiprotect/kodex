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

func DefaultRoleDetails(c Context, roleId string) Element {
	router := UseRouter(c)

	controller := UseController(c)

	role, err := controller.DefaultObjectRole(Unhex(roleId))

	if err != nil {
		return Div("cannot load default role")
	}

	onSubmit := Func[any](c, func() {

		role.Delete()

		router.RedirectTo("/admin/roles")
	})

	return F(
		H1(
			Class("bulma-subtitle"),
			Fmt("Mapping Details - %s", roleId),
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
				Tr(
					Td("Object type"),
					Td(role.ObjectType()),
				),
			),
		),
		Hr(),
		A(
			Class("bulma-button", "bulma-is-danger"),
			Href(Fmt("/admin/roles/details/%s/delete", roleId)),
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
								Href(Fmt("/admin/roles/details/%s", roleId)),
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
						Fmt("/admin/roles/details/%s", roleId),
					)
				},
			),
		),
	)

}
