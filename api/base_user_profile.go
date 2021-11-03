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
)

type BaseUser struct {
	Self User
}

type BaseOrganizationRoles struct {
	Self OrganizationRoles
}

type BaseUserOrganization struct {
	apiOrg Organization
	Self   UserOrganization
}

type BaseAccessToken struct {
	Self AccessToken
}

func (b *BaseAccessToken) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"scopes": b.Self.Scopes(),
	})
}

func (b *BaseOrganizationRoles) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"roles":        b.Self.Roles(),
		"organization": b.Self.Organization(),
	})
}

func (b *BaseUserOrganization) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"name":        b.Self.Name(),
		"source":      b.Self.Source(),
		"default":     b.Self.Default(),
		"description": b.Self.Description(),
		"source_id":   hex.EncodeToString(b.Self.ID()),
	})
}

func (b *BaseUserOrganization) ApiOrganization(controller Controller) (Organization, error) {
	if b.apiOrg == nil {
		org, err := controller.Organization(b.Self.Source(), b.Self.ID())
		if err == nil {
			b.apiOrg = org
		} else {
			org := controller.MakeOrganization()
			org.SetName(b.Self.Name())
			org.SetDescription(b.Self.Description())
			org.SetSourceID(b.Self.ID())
			org.SetSource(b.Self.Source())
			if err := org.Save(); err != nil {
				return nil, err
			}
			b.apiOrg = org
		}
	}
	return b.apiOrg, nil
}

func (b *BaseUser) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"source":       b.Self.Source(),
		"source_id":    hex.EncodeToString(b.Self.SourceID()),
		"email":        b.Self.EMail(),
		"display_name": b.Self.DisplayName(),
		"access_token": b.Self.AccessToken(),
		"roles":        b.Self.Roles(),
	})
}
