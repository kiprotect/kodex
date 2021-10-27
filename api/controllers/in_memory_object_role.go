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

package controllers

import (
	"github.com/kiprotect/kodex/api"
	"time"
)

type InMemoryObjectRole struct {
	api.BaseObjectRole
	createdAt        time.Time
	updatedAt        time.Time
	deletedAt        *time.Time
	organizationID   []byte
	objectID         []byte
	id               []byte
	objectType       string
	objectRole       string
	organizationRole string
	controller       *InMemoryController
}

func MakeInMemoryObjectRole(id, organizationID, objectID []byte, objectType string, controller *InMemoryController) api.ObjectRole {
	inMemoryRole := &InMemoryObjectRole{
		id:             id,
		organizationID: organizationID,
		objectID:       objectID,
		objectType:     objectType,
		controller:     controller,
	}
	inMemoryRole.Self = inMemoryRole
	return inMemoryRole
}

func (c *InMemoryObjectRole) Save() error {
	return c.controller.SaveObjectRole(c)
}

func (c *InMemoryObjectRole) Delete() error {
	return c.controller.DeleteObjectRole(c)
}

func (c *InMemoryObjectRole) Refresh() error {
	return nil
}

func (c *InMemoryObjectRole) CreatedAt() time.Time {
	return c.createdAt
}

func (c *InMemoryObjectRole) DeletedAt() *time.Time {
	return c.deletedAt
}

func (c *InMemoryObjectRole) UpdatedAt() time.Time {
	return c.updatedAt
}

func (c *InMemoryObjectRole) OrganizationID() []byte {
	return c.organizationID
}

func (c *InMemoryObjectRole) ObjectID() []byte {
	return c.objectID
}

func (c *InMemoryObjectRole) ObjectType() string {
	return c.objectType
}

func (c *InMemoryObjectRole) SetObjectRole(role string) error {
	c.objectRole = role
	return nil
}

func (c *InMemoryObjectRole) SetOrganizationRole(role string) error {
	c.organizationRole = role
	return nil
}

func (c *InMemoryObjectRole) ID() []byte {
	return c.id
}

func (c *InMemoryObjectRole) SetObjectID(id []byte) error {
	c.objectID = id
	return nil
}

func (c *InMemoryObjectRole) ObjectRole() string {
	return c.objectRole
}

func (c *InMemoryObjectRole) OrganizationRole() string {
	return c.organizationRole
}
