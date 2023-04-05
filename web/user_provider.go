package web

import (
	"github.com/kiprotect/gospel"
	"github.com/kiprotect/kodex/api"
)

func SetUserProvider(c gospel.Context, userProvider api.UserProvider) {
	gospel.GlobalVar(c, "userProvider", userProvider)
}

func UseUserProvider(c gospel.Context) api.UserProvider {
	return gospel.UseGlobal[api.UserProvider](c, "userProvider")
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
