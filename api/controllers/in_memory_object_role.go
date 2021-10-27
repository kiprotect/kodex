// Kodex (Enterprise Edition - EE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2021  KIProtect GmbH (HRB 208395B) - All Rights Reserved

package controllers

import (
	"fmt"
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
}

func MakeInMemoryObjectRole(id, organizationID, objectID []byte, objectType string) api.ObjectRole {
	inMemoryRole := &InMemoryObjectRole{
		id:             id,
		organizationID: organizationID,
		objectID:       objectID,
		objectType:     objectType,
	}
	inMemoryRole.Self = inMemoryRole
	return inMemoryRole
}

func (c *InMemoryObjectRole) Save() error {
	return nil
}

func (c *InMemoryObjectRole) Delete() error {
	return fmt.Errorf("not implemented")
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

func (c *InMemoryObjectRole) MarshalJSON() ([]byte, error) {
	return nil, fmt.Errorf("not implemented")
}
