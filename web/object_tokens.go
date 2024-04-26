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
	"bytes"
	"fmt"
	. "github.com/gospel-sh/gospel"
	"github.com/kiprotect/kodex"
	"github.com/kiprotect/kodex/api"
	"github.com/kiprotect/kodex/web/ui"
	"regexp"
	"strings"
	"time"
)

func ObjectTokens(c Context, object kodex.Model, actions []string) Element {
	router := UseRouter(c)
	return router.Match(c,
		Route("/new", func(c Context) Element { return NewToken(c, object, actions) }),
		Route("/details/([a-z-0-9]+)", func(c Context, tokenId string) Element { return TokenDetails(c, object, tokenId) }),
		Route("", func(c Context) Element { return ObjectTokensList(c, object) }),
	)
}

func TokenDetails(c Context, object kodex.Model, tokenId string) Element {

	router := UseRouter(c)
	controller := UseController(c)
	externalUser := UseExternalUser(c)
	userProvider := UseUserProvider(c)
	tokenProvider := userProvider.(api.AuthTokenUserProvider)
	tokenController, err := tokenProvider.TokenController(controller)

	if err != nil {
		return Div("not a token provider")
	}

	tokenValue := PersistentGlobalVar(c, "sso-token", "")

	apiOrg, err := externalUser.Roles[0].Organization.ApiOrganization(controller)

	if err != nil {
		// to do: improve
		return Div("error")
	}

	token, err := tokenController.Token(apiOrg, Unhex(tokenId))

	if err != nil {
		return Div("error")
	}

	deleteForm := MakeFormData(c, "deleteToken", POST)

	onSubmit := func() {
		token.Delete()
		router.RedirectTo(router.LastPath())
	}

	deleteForm.OnSubmit(onSubmit)

	return F(
		H2(Class("bulma-subtitle"), "Access Token Details"),
		P(
			"Scopes: ",
			Strong(strings.Join(token.Scopes(), ", ")),
		),
		If(
			tokenValue.Get() == Hex(token.ID()),
			F(
				Hr(),
				ui.Message(
					"success",
					Strong(Hex(token.Token())),
				),
			),
		),
		Hr(),
		P(
			token.Description(),
		),
		Hr(),
		A(
			Class("bulma-button", "bulma-is-danger"),
			Href(router.CurrentPath()+"/delete"),
			"delete token",
		),
		router.Match(
			c,
			Route("/delete$",
				func(c Context) Element {
					return ui.Modal(
						c,
						"Do you really want to delete this token?",
						Span(
							"Do you really want to delete this token?",
						),
						F(
							A(
								Class("bulma-button"),
								Href(router.LastPath()),
								"Cancel",
							),
							Span(Style("flex-grow: 1")),
							Span(
								deleteForm.Form(
									Class("bulma-is-inline"),
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
						router.LastPath(),
					)
				},
			),
		),
	)
}

func NewToken(c Context, object kodex.Model, actions []string) Element {

	form := MakeFormData(c, "newToken", POST)
	description := form.Var("description", "")
	action := form.Var("action", "")
	error := Var(c, "")
	router := UseRouter(c)

	controller := UseController(c)
	userProvider := UseUserProvider(c)
	tokenProvider := userProvider.(api.AuthTokenUserProvider)
	tokenController, err := tokenProvider.TokenController(controller)

	if err != nil {
		return Div("not a token controller")
	}

	externalUser := UseExternalUser(c)
	tokenValue := PersistentGlobalVar(c, "sso-token", "")
	apiUser, err := externalUser.ApiUser(controller)

	if err != nil {
		// to do: improve
		return Div("error")
	}

	apiOrg, err := externalUser.Roles[0].Organization.ApiOrganization(controller)

	if err != nil {
		// to do: improve
		return Div("error")
	}

	onSubmit := func() {

		if description.Get() == "" {
			error.Set("Please enter a description")
			return
		}

		found := false
		for _, act := range actions {
			if act == action.Get() {
				found = true
				break
			}
		}

		if !found {
			error.Set("invalid action value")
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

		token, err := tokenController.MakeToken(apiOrg, apiUser)

		if err != nil {
			error.Set("Cannot create token")
			return
		}

		token.SetDescription(description.Get())

		// we set the scopes
		token.SetScopes([]string{Fmt("kiprotect:api:%s:%s:%s", object.Type(), action.Get(), Hex(object.ID()))})

		for _, role := range externalUser.Roles {

			if role.Organization.Source == apiOrg.Source() && bytes.Equal(role.Organization.ID, apiOrg.SourceID()) {
				token.SetRoles(role.Roles)
				break
			}
		}

		if rv, err := kodex.RandomBytes(16); err != nil {
			error.Set("Cannot create random value")
			return
		} else if err := token.SetToken(rv); err != nil {
			error.Set("Cannot set token value")
			return
		}

		if err := token.Save(); err != nil {
			error.Set(Fmt("Cannot save token: %v", err))
			return
		}

		success = true

		tokenValue.Set(Hex(token.ID()))

		router.RedirectTo(Fmt(router.LastPath()+"/details/%s", Hex(token.ID())))
	}

	form.OnSubmit(onSubmit)

	values := []any{}

	for _, action := range actions {
		values = append(values, Option(Value(action), action))
	}

	var errorNotice Element

	if error.Get() != "" {
		errorNotice = P(
			Class("bulma-help", "bulma-is-danger"),
			error.Get(),
		)
	}

	return form.Form(
		H1(Class("bulma-subtitle"), "New Token"),
		Div(
			Class("bulma-field"),
			Label(
				Class("bulma-label"),
				"Action",
			),
			Div(
				Class("bulma-select", "bulma-is-fullwidth"),
				Select(
					Class("bulma-select"),
					values,
					Value(action),
				),
			),
		),
		Div(
			Class("bulma-field"),
			errorNotice,
			Label(
				Class("bulma-label"),
				"Description",
				Input(
					Class("bulma-input", If(error.Get() != "", "bulma-is-danger")),
					Type("text"),
					Value(description),
					Placeholder("token description"),
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
					"Create Token",
				),
			),
		),
	)
}

func ObjectTokensList(c Context, object kodex.Model) Element {

	user := UseExternalUser(c)
	router := UseRouter(c)
	controller := UseController(c)
	userProvider := UseUserProvider(c)
	tokenProvider := userProvider.(api.AuthTokenUserProvider)
	tokenController, err := tokenProvider.TokenController(controller)

	apiOrg, err := user.Roles[0].Organization.ApiOrganization(controller)

	if err != nil {
		return Div("cannot load API organization")
	}

	if err != nil {
		return Div("cannot load controller")
	}

	tokens, err := tokenController.Tokens(apiOrg, map[string]any{})

	if err != nil {
		return Div("error")
	}

	objectTokens := make([]api.Token, 0, len(tokens))
	objectScopeRegexp, err := regexp.Compile(fmt.Sprintf(`^kiprotect:api:%s(?::[a-z0-9]+)?:%s$`, object.Type(), Hex(object.ID())))

	if err != nil {
		return Div("error")
	}

	for _, token := range tokens {
		for _, scope := range token.Scopes() {
			if objectScopeRegexp.FindString(scope) != "" {
				objectTokens = append(objectTokens, token)
			}
		}
	}

	pis := make([]any, 0, len(objectTokens))

	for _, token := range objectTokens {

		if token.Description() == "" {
			// this is an SSO token, we skip it...
			continue
		}

		var expiresAt Element = ui.ListColumn("sm", Span("never"))

		if token.ExpiresAt() != nil {
			expiresAt = ui.ListColumn("sm", HumanDuration(time.Now().Sub(*token.ExpiresAt())))
		}

		tokenItem := A(
			Href(Fmt(router.CurrentPath()+"/details/%s", Hex(token.ID()))),
			ui.ListItem(
				ui.ListColumn("md", token.Description()),
				ui.ListColumn("sm", strings.Join(token.Scopes(), ", ")),
				ui.ListColumn("sm", HumanDuration(time.Now().Sub(token.CreatedAt()))),
				expiresAt,
			),
		)
		pis = append(pis, tokenItem)
	}

	return F(
		IfElse(
			len(pis) > 0,
			ui.List(
				ui.ListHeader(
					ui.ListColumn("md", "Name"),
					ui.ListColumn("sm", "Scopes"),
					ui.ListColumn("sm", "Created At"),
					ui.ListColumn("sm", "Expires At"),
				),
				pis,
			),
			F(Div("No access tokens defined yet."), Hr()),
		),
		A(Href(router.CurrentPath()+"/new"), Class("bulma-button", "bulma-is-success"), "New Token"),
	)
}
