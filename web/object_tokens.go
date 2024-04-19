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
	"github.com/kiprotect/kodex/web/ui"
	"strings"
	"time"
)

func ObjectTokens(c Context, object kodex.Model) Element {

	user := UseExternalUser(c)
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

	pis := make([]any, 0, len(tokens))

	for _, token := range tokens {

		if token.Description() == "" {
			// this is an SSO token, we skip it...
			continue
		}

		var expiresAt Element = ui.ListColumn("sm", Span("never"))

		if token.ExpiresAt() != nil {
			expiresAt = ui.ListColumn("sm", HumanDuration(time.Now().Sub(*token.ExpiresAt())))
		}

		tokenItem := A(
			Href(Fmt("/sso/tokens/details/%s", Hex(token.ID()))),
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
		ui.List(
			ui.ListHeader(
				ui.ListColumn("md", "Name"),
				ui.ListColumn("sm", "Scopes"),
				ui.ListColumn("sm", "Created At"),
				ui.ListColumn("sm", "Expires At"),
			),
			pis,
		),
		A(Href("/sso/tokens/new"), Class("bulma-button", "bulma-is-success"), "New Token"),
	)
}
