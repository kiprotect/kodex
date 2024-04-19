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
	"github.com/kiprotect/kodex"
	"github.com/kiprotect/kodex/api"
)

func NewProjectRole(project kodex.Project) ElementFunction {

	return func(c Context) Element {
		orgRole := Var(c, "")
		objectRole := Var(c, "")
		error := Var(c, "")
		router := UseRouter(c)
		controller := UseController(c)
		organization := UseDefaultOrganization(c)
		apiOrg, err := organization.ApiOrganization(controller)

		if err != nil {
			return Div("cannot load API organization")
		}

		onSubmit := Func[any](c, func() {

			if orgRole.Get() == "" {
				error.Set("Please enter an organization role")
				return
			}

			if objectRole.Get() == "" {
				error.Set("Please enter an object role")
				return
			}

			controller.Begin()

			success := false

			defer func() {
				if success {
					controller.Commit()
				}
				controller.Rollback()
			}()

			role := controller.MakeObjectRole(project, apiOrg)
			role.SetOrganizationRole(orgRole.Get())
			role.SetObjectRole(objectRole.Get())

			if err := role.Save(); err != nil {
				error.Set(Fmt("Cannot save role: %v", err))
				return
			}

			success = true

			router.RedirectTo(Fmt("/flows/projects/%s/settings/roles/details/%s", Hex(project.ID()), Hex(role.ID())))
		})

		var errorNotice Element

		if error.Get() != "" {
			errorNotice = P(
				Class("bulma-help", "bulma-is-danger"),
				error.Get(),
			)
		}

		roles := []Element{}

		for _, item := range api.ObjectRoleValues {
			roles = append(roles, Option(If(item == objectRole.Get(), BooleanAttrib("selected")()), Value(item), item))
		}

		return Form(
			Method("POST"),
			OnSubmit(onSubmit),
			H1(Class("bulma-subtitle"), "New Object Role"),
			Div(
				Class("bulma-field"),
				errorNotice,
				Label(
					Class("bulma-label"),
					"Organization Role",
					Input(
						Class("bulma-input", If(error.Get() != "", "bulma-is-danger")),
						Type("text"),
						Value(orgRole),
						Placeholder("organization role"),
					),
				),
				Label(
					Class("bulma-label"),
					"Object Role",
					Div(
						Class("bulma-select", "bulma-is-fullwidth"),
						Select(
							roles,
							Value(objectRole),
							Attrib("autocomplete")("off"),
							Id("objRoleSelect"),
						),
					),
				),
			),
			Div(
				Class("bulma-field"),
				P(
					Class("bulma-control"),
					Button(
						Class("bulma-button", "bulma-is-success"),
						Type("submit"),
						"Create Object Role",
					),
				),
			),
		)

	}
}
