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

package api

import (
	"github.com/gin-gonic/gin"
	"github.com/kiprotect/kodex"
)

type UserProviderDefinition struct {
	Name              string            `json:"name"`
	Description       string            `json:"description"`
	Maker             UserProviderMaker `json:"-"`
	SettingsValidator SettingsValidator `json:"-"`
}

type UserProviderDefinitions map[string]UserProviderDefinition
type UserProviderMaker func(settings kodex.Settings) (UserProvider, error)

type UserProvider interface {
	Initialize(group *gin.RouterGroup) error
	Get(context *gin.Context) (*ExternalUser, error)
}

type CreateUserProvider interface {
	Create(user *ExternalUser) error
}

type AccessToken struct {
	Scopes []string `json:"scopes"`
	Token  []byte   `json:"-" coerce:"name:token"`
}

type OrganizationRoles struct {
	Roles        []string          `json:"roles"`
	Organization *UserOrganization `json:"organization"`
}

type UserOrganization struct {
	Source      string       `json:"source"`
	Name        string       `json:"name"`
	Default     bool         `json:"default"`
	Description string       `json:"description"`
	ID          []byte       `json:"id"`
	apiOrg      Organization `json:"-"`
}

type ExternalUser struct {
	Source      string                 `json:"source"`
	SourceID    []byte                 `json:"sourceID"`
	EMail       string                 `json:"email"`
	DisplayName string                 `json:"displayName"`
	Superuser   bool                   `json:"superuser"`
	AccessToken *AccessToken           `json:"accessToken"`
	Roles       []*OrganizationRoles   `json:"roles"`
	Limits      map[string]interface{} `json:"limits"`
}

func (i *UserOrganization) ApiOrganization(controller Controller) (Organization, error) {
	if i.apiOrg == nil {
		org, err := controller.Organization(i.Source, i.ID)
		if err == nil {
			i.apiOrg = org
		} else {
			org := controller.MakeOrganization()
			org.SetName(i.Name)
			org.SetDescription(i.Description)
			org.SetSourceID(i.ID)
			org.SetSource(i.Source)
			if err := org.Save(); err != nil {
				return nil, err
			}
			i.apiOrg = org
		}
	}
	return i.apiOrg, nil
}

func MakeUserOrganization(source, name, description string, id []byte) *UserOrganization {
	return &UserOrganization{
		Source:      source,
		Name:        name,
		Description: description,
		ID:          id,
	}
}

func MakeAccessToken(scopes []string) *AccessToken {
	return &AccessToken{
		Scopes: scopes,
	}
}

func MakeOrganizationRoles(org *UserOrganization, roles []string) *OrganizationRoles {
	return &OrganizationRoles{
		Organization: org,
		Roles:        roles,
	}
}

func MakeUser(source, email string, superuser bool, roles []*OrganizationRoles, limits map[string]interface{}, token *AccessToken) *ExternalUser {

	return &ExternalUser{
		Source:      source,
		Limits:      limits,
		EMail:       email,
		Superuser:   superuser,
		Roles:       roles,
		AccessToken: token,
	}

}
