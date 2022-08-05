// Kodex (Community Edition - CE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2022  KIProtect GmbH (HRB 208395B) - Germany
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

package providers

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kiprotect/go-helpers/forms"
	"github.com/kiprotect/kodex"
	"github.com/kiprotect/kodex/api"
	"regexp"
)

var InMemoryUserForm = forms.Form{
	Fields: []forms.Field{
		{
			Name: "displayName",
			Validators: []forms.Validator{
				forms.IsOptional{Default: ""},
				forms.IsString{},
			},
		},
	},
}

type InMemoryUserProviderSettings struct {
}

func ValidateInMemoryUserProviderSettings(settings map[string]interface{}) (interface{}, error) {
	if params, err := InMemoryUserForm.Validate(settings); err != nil {
		return nil, err
	} else {
		providerSettings := &InMemoryUserProviderSettings{}
		if err := InMemoryUserForm.Coerce(providerSettings, params); err != nil {
			return nil, err
		}
		return providerSettings, nil
	}
}

type InMemoryUserProvider struct {
	users []*api.User
}

func MakeInMemoryUserProvider(settings kodex.Settings) (api.UserProvider, error) {
	return &InMemoryUserProvider{
		users: make([]*api.User, 0),
	}, nil
}

func (i *InMemoryUserProvider) Initialize(group *gin.RouterGroup) error {
	return nil
}

func (i *InMemoryUserProvider) Create(user *api.User) error {
	i.users = append(i.users, user)
	return nil
}

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

// Return a user with the given access token
func (i *InMemoryUserProvider) Get(c *gin.Context) (*api.User, error) {

	accessToken, ok := extractAccessToken(c)

	if !ok {

		err := fmt.Errorf("malformed/missing authorization header")
		api.HandleError(c, 401, err)

		return nil, err
	}

	if token, err := hex.DecodeString(accessToken); err != nil {
		err := fmt.Errorf("malformed access token")
		api.HandleError(c, 401, err)
		return nil, err
	} else {
		for _, user := range i.users {
			if bytes.Equal(user.AccessToken.Token, token) {
				return user, nil
			}
		}
	}

	err := fmt.Errorf("invalid access token")
	api.HandleError(c, 401, err)
	return nil, err

}
