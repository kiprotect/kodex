package web

import (
	"github.com/kiprotect/gospel"
	"github.com/kiprotect/kodex/api"
)

func SetUserProvider(c gospel.Context, userProvider api.UserProvider) {
	gospel.SetVar(c, "userProvider", userProvider)
}

func UseUserProvider(c gospel.Context) api.UserProvider {
	return gospel.UseVar[api.UserProvider](c, "userProvider")
}

func SetUser(c gospel.Context, user api.User) {
	gospel.SetVar(c, "user", user)
}

func UseUser(c gospel.Context) api.User {
	return gospel.UseVar[api.User](c, "user")
}

func SetExternalUser(c gospel.Context, user *api.ExternalUser) {
	gospel.SetVar(c, "externalUser", user)
}

func UseExternalUser(c gospel.Context) *api.ExternalUser {
	return gospel.UseVar[*api.ExternalUser](c, "externalUser")
}
