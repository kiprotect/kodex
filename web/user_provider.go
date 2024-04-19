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
	"github.com/gospel-sh/gospel"
	"github.com/kiprotect/kodex/api"
)

func SetUserProvider(c gospel.Context, userProvider api.UserProvider) {
	gospel.GlobalVar(c, "userProvider", userProvider)
}

func UseUserProvider(c gospel.Context) api.UserProvider {
	return gospel.UseGlobal[api.UserProvider](c, "userProvider")
}

func SetApiUser(c gospel.Context, user api.User) {
	gospel.GlobalVar(c, "apiUser", user)
}

func UseApiUser(c gospel.Context) api.User {
	return gospel.UseGlobal[api.User](c, "apiUser")
}

func SetExternalUser(c gospel.Context, user *api.ExternalUser) {
	gospel.GlobalVar(c, "externalUser", user)
}

func UseExternalUser(c gospel.Context) *api.ExternalUser {
	return gospel.UseGlobal[*api.ExternalUser](c, "externalUser")
}

func UseDefaultOrganization(c gospel.Context) *api.UserOrganization {
	user := UseExternalUser(c)

	for _, role := range user.Roles {
		if role.Organization.Default {
			return role.Organization
		}
	}
	return nil
}
