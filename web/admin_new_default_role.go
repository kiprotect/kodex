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

func NewDefaultRole(c Context) Element {

	orgRole := Var(c, "")
	objectRole := Var(c, "")
	objectType := Var(c, "")
	error := Var(c, "")
	router := UseRouter(c)

	controller := UseController(c)
	organization := UseDefaultOrganization(c)
	apiOrg, err := organization.ApiOrganization(controller)

	if err != nil {
		// to do: improve
		return Div("error")
	}

	onSubmit := Func[any](c, func() {

		if orgRole.Get() == "" {
			error.Set("Please enter a SSO group")
			return
		}

		if objectRole.Get() == "" {
			error.Set("Please enter a Kodex group")
			return
		}

		if objectType.Get() == "" {
			error.Set("Please enter an object type")
		}

		controller.Begin()

		success := false

		defer func() {
			if success {
				controller.Commit()
			}
			controller.Rollback()
		}()

		defaultRole := controller.MakeDefaultObjectRole(objectType.Get(), apiOrg)

		defaultRole.SetOrganizationRole(orgRole.Get())
		defaultRole.SetObjectRole(objectRole.Get())

		if err := defaultRole.Save(); err != nil {
			error.Set(Fmt("Cannot save role: %v", err))
			return
		}

		success = true

		router.RedirectTo(Fmt("/admin/roles/details/%s", Hex(defaultRole.ID())))
	})

	var errorNotice Element

	if error.Get() != "" {
		errorNotice = P(
			Class("bulma-help", "bulma-is-danger"),
			error.Get(),
		)
	}

	return Form(
		Method("POST"),
		OnSubmit(onSubmit),
		H1(Class("bulma-subtitle"), "New Default Role"),
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
				Input(
					Class("bulma-input", If(error.Get() != "", "bulma-is-danger")),
					Type("text"),
					Value(objectRole),
					Placeholder("object role"),
				),
			),
			Label(
				Class("bulma-label"),
				"Object Type",
				Input(
					Class("bulma-input", If(error.Get() != "", "bulma-is-danger")),
					Type("text"),
					Value(objectType),
					Placeholder("object type"),
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
					"Create Default Role",
				),
			),
		),
	)
}
