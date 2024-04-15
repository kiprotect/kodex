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
