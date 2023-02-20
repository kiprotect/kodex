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

package controllers

import (
	"github.com/kiprotect/kodex/api"
	"time"
)

type InMemoryChangeRequest struct {
	api.BaseChangeRequest
	createdAt  time.Time
	updatedAt  time.Time
	deletedAt  *time.Time
	objectType string
	status     api.ChangeRequestStatus
	objectID   []byte
	id         []byte
	data       interface{}
	metadata       interface{}
	controller *InMemoryController
}

func MakeInMemoryChangeRequest(objectType string, objectID []byte, controller *InMemoryController) api.ChangeRequest {
	inMemoryChangeRequest := &InMemoryChangeRequest{
		objectID:   objectID,
		objectType: objectType,
	}
	inMemoryChangeRequest.Self = inMemoryChangeRequest
	return inMemoryChangeRequest
}

func (c *InMemoryChangeRequest) Save() error {
	return c.controller.SaveChangeRequest(c)
}

func (c *InMemoryChangeRequest) Delete() error {
	return c.controller.DeleteChangeRequest(c)
}

func (c *InMemoryChangeRequest) Refresh() error {
	return nil
}

func (c *InMemoryChangeRequest) CreatedAt() time.Time {
	return c.createdAt
}

func (c *InMemoryChangeRequest) SetReviewer(user api.User) error {
	c.Reviewer_ = user
	return nil
}

func (c *InMemoryChangeRequest) SetStatus(status api.ChangeRequestStatus) error {
	c.status = status
	return nil
}

func (c *InMemoryChangeRequest) Status() api.ChangeRequestStatus {
	return c.status
}

func (c *InMemoryChangeRequest) DeletedAt() *time.Time {
	return c.deletedAt
}

func (c *InMemoryChangeRequest) UpdatedAt() time.Time {
	return c.updatedAt
}

func (c *InMemoryChangeRequest) ObjectID() []byte {
	return c.objectID
}

func (c *InMemoryChangeRequest) ID() []byte {
	return c.id
}

func (c *InMemoryChangeRequest) ObjectType() string {
	return c.objectType
}

func (c *InMemoryChangeRequest) Data() interface{} {
	return c.data
}

func (c *InMemoryChangeRequest) SetData(data interface{}) error {
	c.data = data
	return nil
}

func (c *InMemoryChangeRequest) Metadata() interface{} {
	return c.metadata
}

func (c *InMemoryChangeRequest) SetMetadata(data interface{}) error {
	c.metadata = data
	return nil
}
