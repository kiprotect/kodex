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

package api

import (
	"encoding/hex"
	"encoding/json"
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
	Get(string) (*User, error)
	Start()
	Stop()
}

type AccessToken struct {
	scopes []string `json:"scopes"`
	token  []byte
}

func (i *AccessToken) Scopes() []string {
	return i.scopes
}

func (i *AccessToken) Token() []byte {
	return i.token
}

type OrganizationRoles struct {
	roles        []string          `json:"roles"`
	organization *UserOrganization `json:"organization"`
}

func (i *OrganizationRoles) Roles() []string {
	return i.roles
}

func (i *OrganizationRoles) Organization() *UserOrganization {
	return i.organization
}

type UserOrganization struct {
	source      string       `json:"source"`
	name        string       `json:"name"`
	isDefault   bool         `json:"default"`
	description string       `json:"description"`
	id          []byte       `json:"id"`
	apiOrg      Organization `json:"-"`
}

func (i *UserOrganization) Name() string {
	return i.name
}

func (i *UserOrganization) Source() string {
	return i.source
}

func (i *UserOrganization) Default() bool {
	return i.isDefault
}

func (i *UserOrganization) Description() string {
	return i.description
}

func (i *UserOrganization) ID() []byte {
	return i.id
}

type User struct {
	source      string                 `json:"source"`
	sourceID    []byte                 `json:"sourceID"`
	email       string                 `json:"email"`
	displayName string                 `json:"displayName"`
	superuser   bool                   `json:"superuser"`
	accessToken *AccessToken           `json:"accessToken"`
	roles       []*OrganizationRoles   `json:"roles"`
	limits      map[string]interface{} `json:"limits"`
}

func (i *User) Source() string {
	return i.source
}

func (i *User) SourceID() []byte {
	return i.sourceID
}

func (i *User) EMail() string {
	return i.email
}

func (i *User) SuperUser() bool {
	return i.superuser
}

func (i *User) DisplayName() string {
	return i.displayName
}

func (i *User) AccessToken() *AccessToken {
	return i.accessToken
}

func (i *User) Roles() []*OrganizationRoles {
	return i.roles
}

func (i *User) Limits() map[string]interface{} {
	return i.limits
}

func (i *AccessToken) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"scopes": i.scopes,
	})
}

func (i *OrganizationRoles) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"roles":        i.roles,
		"organization": i.organization,
	})
}

func (i *UserOrganization) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"name":        i.name,
		"source":      i.source,
		"default":     i.isDefault,
		"description": i.description,
		"source_id":   hex.EncodeToString(i.id),
	})
}

func (i *UserOrganization) ApiOrganization(controller Controller) (Organization, error) {
	if i.apiOrg == nil {
		org, err := controller.Organization(i.source, i.id)
		if err == nil {
			i.apiOrg = org
		} else {
			org := controller.MakeOrganization()
			org.SetName(i.name)
			org.SetDescription(i.description)
			org.SetSourceID(i.id)
			org.SetSource(i.source)
			if err := org.Save(); err != nil {
				return nil, err
			}
			i.apiOrg = org
		}
	}
	return i.apiOrg, nil
}

func (i *User) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"source":       i.source,
		"source_id":    hex.EncodeToString(i.sourceID),
		"email":        i.email,
		"display_name": i.displayName,
		"access_token": i.accessToken,
		"roles":        i.roles,
	})
}

func MakeUserOrganization(source, name, description string, id []byte) *UserOrganization {
	return &UserOrganization{
		source:      source,
		name:        name,
		description: description,
		id:          id,
	}
}

func MakeAccessToken(scopes []string) *AccessToken {
	return &AccessToken{
		scopes: scopes,
	}
}

func MakeOrganizationRoles(org *UserOrganization, roles []string) *OrganizationRoles {
	return &OrganizationRoles{
		organization: org,
		roles:        roles,
	}
}

func MakeUser(source, email string, superuser bool, roles []*OrganizationRoles, limits map[string]interface{}, token *AccessToken) *User {

	return &User{
		source:      source,
		limits:      limits,
		email:       email,
		superuser:   superuser,
		roles:       roles,
		accessToken: token,
	}

}
