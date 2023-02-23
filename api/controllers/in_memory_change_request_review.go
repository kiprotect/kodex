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

type InMemoryChangeRequestReview struct {
	api.BaseChangeRequestReview
	createdAt     time.Time
	updatedAt     time.Time
	deletedAt     *time.Time
	status        api.ChangeRequestReviewStatus
	id            []byte
	data          interface{}
	metadata      interface{}
	changeRequest *InMemoryChangeRequest
	controller    *InMemoryController
}

func MakeInMemoryChangeRequestReview(changeRequest *InMemoryChangeRequest) api.ChangeRequestReview {
	inMemoryChangeRequestReview := &InMemoryChangeRequestReview{
		changeRequest: changeRequest,
	}
	inMemoryChangeRequestReview.Self = inMemoryChangeRequestReview
	return inMemoryChangeRequestReview
}

func (c *InMemoryChangeRequestReview) Save() error {
	return c.changeRequest.SaveChangeRequestReview(c)
}

func (c *InMemoryChangeRequestReview) Delete() error {
	return c.changeRequest.DeleteChangeRequestReview(c)
}

func (c *InMemoryChangeRequestReview) Refresh() error {
	return nil
}

func (c *InMemoryChangeRequestReview) ChangeRequest() api.ChangeRequest {
	return c.changeRequest
}

func (c *InMemoryChangeRequestReview) CreatedAt() time.Time {
	return c.createdAt
}

func (c *InMemoryChangeRequestReview) SetStatus(status api.ChangeRequestReviewStatus) error {
	c.status = status
	return nil
}

func (c *InMemoryChangeRequestReview) Status() api.ChangeRequestReviewStatus {
	return c.status
}

func (c *InMemoryChangeRequestReview) DeletedAt() *time.Time {
	return c.deletedAt
}

func (c *InMemoryChangeRequestReview) UpdatedAt() time.Time {
	return c.updatedAt
}

func (c *InMemoryChangeRequestReview) ID() []byte {
	return c.id
}

func (c *InMemoryChangeRequestReview) Data() interface{} {
	return c.data
}

func (c *InMemoryChangeRequestReview) SetData(data interface{}) error {
	c.data = data
	return nil
}

func (c *InMemoryChangeRequestReview) Metadata() interface{} {
	return c.metadata
}

func (c *InMemoryChangeRequestReview) SetMetadata(data interface{}) error {
	c.metadata = data
	return nil
}
