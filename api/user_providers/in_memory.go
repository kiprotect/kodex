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
	"github.com/kiprotect/go-helpers/maps"
	"github.com/kiprotect/kodex"
	"github.com/kiprotect/kodex/api"
	"net/http"
	"regexp"
)

var InMemoryUserProviderForm = api.BlueprintConfigForm

type InMemoryUserProviderSettings struct {
	api.UsersAndRoles
}

func ValidateInMemoryUserProviderSettings(settings map[string]interface{}) (interface{}, error) {
	if params, err := InMemoryUserProviderForm.Validate(settings); err != nil {
		return nil, err
	} else {
		providerSettings := &InMemoryUserProviderSettings{}
		if err := InMemoryUserProviderForm.Coerce(providerSettings, params); err != nil {
			return nil, err
		}
		return providerSettings, nil
	}
}

type InMemoryUserProvider struct {
	settings *InMemoryUserProviderSettings
}

func MakeInMemoryUserProvider(settings kodex.Settings) (api.UserProvider, error) {

	providerSettings, err := settings.Get("user-provider.settings")

	if err != nil {
		return nil, err
	}

	providerSettingsMap, ok := maps.ToStringMap(providerSettings)

	if !ok {
		return nil, fmt.Errorf("invalid config")
	}

	if params, err := InMemoryUserProviderForm.Validate(providerSettingsMap); err != nil {
		return nil, err
	} else {

		settingsStruct := &InMemoryUserProviderSettings{}

		if err := InMemoryUserProviderForm.Coerce(settingsStruct, params); err != nil {
			return nil, err
		}

		return &InMemoryUserProvider{
			settings: settingsStruct,
		}, nil

	}
}

func (i *InMemoryUserProvider) Initialize(group *gin.RouterGroup) error {
	return nil
}

func (i *InMemoryUserProvider) Create(user *api.ExternalUser) error {
	i.settings.Users = append(i.settings.Users, user)
	return nil
}

func extractAccessToken(request *http.Request) (string, bool) {
	authorizationHeader := request.Header.Get("Authorization")

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
func (i *InMemoryUserProvider) Get(controller api.Controller, request *http.Request) (*api.ExternalUser, error) {

	accessToken, ok := extractAccessToken(request)

	if !ok {
		return nil, fmt.Errorf("malformed/missing authorization header")
	}

	if token, err := hex.DecodeString(accessToken); err != nil {
		return nil, fmt.Errorf("malformed access token")
	} else {
		for _, user := range i.settings.Users {
			if bytes.Equal(user.AccessToken.Token, token) {
				return user, nil
			}
		}
	}

	return nil, fmt.Errorf("invalid access token")

}
