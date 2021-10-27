// Kodex (Community Edition - CE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2021  KIProtect GmbH (HRB 208395B) - Germany
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

package decorators

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kiprotect/go-helpers/errors"
	"github.com/kiprotect/kodex"
	"github.com/kiprotect/kodex/api"
	"github.com/kiprotect/kodex/api/user_provider"
	"regexp"
	"strings"
)

func extractAccessToken(c *gin.Context) (string, bool) {
	authorizationHeader := c.Request.Header.Get("Authorization")

	if authorizationHeader == "" {
		return "", false
	}

	regex, _ := regexp.Compile("(?i)\\s*Bearer\\s+([\\w\\d-]+)")
	result := regex.FindStringSubmatch(authorizationHeader)
	if result == nil {
		return "", false
	}
	return result[1], true
}

func CheckScopes(requiredScopes, userScopes []string) bool {
	for _, scope := range requiredScopes {
		if len(scope) > 0 && scope[len(scope)-1:] == "*" {
			// this is a wildcard, we check for the given prefix
			// e.g. kiprotect:api
			prefix := scope[:len(scope)-1]
			for _, userScope := range userScopes {
				// e.g. a literal match (kiprotect:api) or a prefix match
				// with a colon, e.g. kiprotect:api:read
				if userScope == prefix || strings.HasPrefix(userScope, prefix+":") {
					return true
				}
			}
		} else {
			// this is a full scope, e.g. kiprotect:api:privacy-managers:read
			for _, userScope := range userScopes {
				// either we get a literal match (kiprotect:api:read) or one
				// of the user tokens includes the necessary scope.
				// For example, if the user has a 'kiprotect:api' scope and
				// we look for a 'kiprotect:api:privacy-managers:read' scope
				// then that scope has a prefix 'kiprotect:api:' so the given
				// user token matches. If the user has a token 'kiprotect:apic'
				// then 'kiprotect:apic:' will not match.
				if userScope == scope || strings.HasPrefix(scope, userScope+":") {
					return true
				}
			}
		}
	}
	return false
}

//Makes sure that the user has provided a valid access token.
//Stores the token, user ID and user profile in the context.
func ValidUser(settings kodex.Settings, scopes []string, superUser bool) gin.HandlerFunc {

	testDecorator := func(c *gin.Context) {
		c.Set("userId", "test")
		c.Set("accessTokenId", "foobar")
	}

	decorator := func(c *gin.Context) {

		ch, ok := c.Get("profileProvider")

		if !ok {
			api.HandleError(c, 500, fmt.Errorf("internal server error"))
			return
		}

		profileProvider, ok := ch.(provider.UserProfileProvider)

		if !ok {
			api.HandleError(c, 500, fmt.Errorf("internal server error"))
			return
		}

		accessToken, ok := extractAccessToken(c)

		if !ok {
			api.HandleError(c, 401, fmt.Errorf("malformed/missing authorization header"))
			return
		}

		userProfile, err := profileProvider.Get(accessToken)

		if err != nil {
			api.HandleError(c, 401, fmt.Errorf("invalid access token"))
			return
		}

		if superUser && !userProfile.SuperUser() {
			api.HandleError(c, 403, fmt.Errorf("access denied"))
			return
		}

		if !CheckScopes(scopes, userProfile.AccessToken().Scopes()) {
			api.HandleError(c, 403, errors.MakeExternalError("access denied", "ACCESS-DENIED", map[string]interface{}{"user_scopes": userProfile.AccessToken().Scopes(), "required_scopes": scopes}, nil))
			return
		}

		//if successful, we set the userId to the given value
		c.Set("userId", userProfile.SourceID())
		c.Set("userSource", userProfile.Source())
		c.Set("userProfile", userProfile)
		c.Set("accessTokenId", accessToken)
	}

	if test, _ := settings.Bool("test"); test {
		return testDecorator
	}
	return decorator
}
