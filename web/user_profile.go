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
	"strings"
)

func UserProfile(c Context) Element {

	user := UseExternalUser(c)

	roles := []Element{}

	for _, userRoles := range user.Roles {
		roles = append(roles, P(
			Fmt("In organization '%s', you have roles '%s'.", userRoles.Organization.Name, strings.Join(userRoles.Roles, ", ")),
		))
	}

	return F(
		H1(Class("bulma-title"), user.Email),
		roles,
	)
}
