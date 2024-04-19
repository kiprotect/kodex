// Kodex (Community Edition - CE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2024  KIProtect GmbH (HRB 208395B) - Germany
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

package controllers

import (
	"github.com/kiprotect/kodex/api"
	"time"
)

type InMemoryDefaultObjectRole struct {
	api.BaseDefaultObjectRole
	createdAt        time.Time
	updatedAt        time.Time
	deletedAt        *time.Time
	organizationID   []byte
	id               []byte
	objectType       string
	objectRole       string
	organizationRole string
	controller       *InMemoryController
}

func MakeInMemoryDefaultObjectRole(id, organizationID []byte, objectType string, controller *InMemoryController) api.DefaultObjectRole {
	inMemoryRole := &InMemoryDefaultObjectRole{
		id:             id,
		organizationID: organizationID,
		objectType:     objectType,
		controller:     controller,
	}
	inMemoryRole.Self = inMemoryRole
	return inMemoryRole
}

func (c *InMemoryDefaultObjectRole) Save() error {
	return c.controller.SaveDefaultObjectRole(c)
}

func (c *InMemoryDefaultObjectRole) Delete() error {
	return c.controller.DeleteDefaultObjectRole(c)
}

func (c *InMemoryDefaultObjectRole) Refresh() error {
	return nil
}

func (c *InMemoryDefaultObjectRole) CreatedAt() time.Time {
	return c.createdAt
}

func (c *InMemoryDefaultObjectRole) DeletedAt() *time.Time {
	return c.deletedAt
}

func (c *InMemoryDefaultObjectRole) UpdatedAt() time.Time {
	return c.updatedAt
}

func (c *InMemoryDefaultObjectRole) OrganizationID() []byte {
	return c.organizationID
}

func (c *InMemoryDefaultObjectRole) ObjectType() string {
	return c.objectType
}

func (c *InMemoryDefaultObjectRole) SetObjectRole(role string) error {
	c.objectRole = role
	return nil
}

func (c *InMemoryDefaultObjectRole) SetOrganizationRole(role string) error {
	c.organizationRole = role
	return nil
}

func (c *InMemoryDefaultObjectRole) ID() []byte {
	return c.id
}

func (c *InMemoryDefaultObjectRole) ObjectRole() string {
	return c.objectRole
}

func (c *InMemoryDefaultObjectRole) OrganizationRole() string {
	return c.organizationRole
}
