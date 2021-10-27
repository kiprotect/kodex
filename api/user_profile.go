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

type UserProfile interface {
	Source() string
	SourceID() []byte
	EMail() string
	SuperUser() bool
	DisplayName() string
	AccessToken() AccessToken
	Roles() []OrganizationRoles
	Limits() map[string]interface{}
}

type AccessToken interface {
	Scopes() []string
}

type OrganizationRoles interface {
	Roles() []string
	Organization() UserOrganization
}

type UserOrganization interface {
	Name() string
	Source() string
	Default() bool
	Description() string
	ID() []byte
	ApiOrganization(Controller) (Organization, error)
}

type UserProfileProvider interface {
	Get(string) (UserProfile, error)
	Start()
	Stop()
}
