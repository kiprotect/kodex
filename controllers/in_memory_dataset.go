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
	"fmt"
	"github.com/kiprotect/kodex"
	"time"
)

type InMemoryDataset struct {
	kodex.BaseDataset
	id          []byte
	index       int
	actionType  string
	description string
	name        string
	items       []map[string]any
	createdAt   time.Time
	updatedAt   time.Time
	deletedAt   *time.Time
	data        interface{}
	configData  map[string]interface{}
}

func MakeInMemoryDataset(id []byte, project kodex.Project) *InMemoryDataset {
	inMemoryDataset := &InMemoryDataset{
		id: id,
		BaseDataset: kodex.BaseDataset{
			Project_: project,
		},
	}
	inMemoryDataset.Self = inMemoryDataset
	return inMemoryDataset
}

func (c *InMemoryDataset) ID() []byte {
	return c.id
}

func (c *InMemoryDataset) InternalID() []byte {
	return c.id
}

func (c *InMemoryDataset) Items() []map[string]any {
	return c.items
}

func (c *InMemoryDataset) SetItems(items []map[string]any) error {
	c.items = items
	return nil
}

func (c *InMemoryDataset) Index() int {
	return c.index
}

func (c *InMemoryDataset) SetIndex(index int) error {
	c.index = index
	return nil
}

func (c *InMemoryDataset) Name() string {
	return c.name
}

func (c *InMemoryDataset) Data() interface{} {
	return c.data
}

func (c *InMemoryDataset) SetData(data interface{}) error {
	c.data = data
	return nil
}

func (c *InMemoryDataset) Description() string {
	return c.description
}

func (c *InMemoryDataset) SetName(name string) error {
	c.name = name
	return nil
}

func (c *InMemoryDataset) SetDescription(description string) error {
	c.description = description
	return nil
}

func (c *InMemoryDataset) SetUpdatedAt(t time.Time) error {
	c.updatedAt = t
	return nil
}

func (c *InMemoryDataset) SetCreatedAt(t time.Time) error {
	c.createdAt = t
	return nil
}

func (c *InMemoryDataset) SetDeletedAt(t *time.Time) error {
	c.deletedAt = t
	return nil
}

func (c *InMemoryDataset) UpdatedAt() time.Time {
	return c.updatedAt
}

func (c *InMemoryDataset) CreatedAt() time.Time {
	return c.createdAt
}

func (c *InMemoryDataset) DeletedAt() *time.Time {
	return c.deletedAt
}

func (c *InMemoryDataset) Refresh() error {
	// we don't need to do anything here
	return nil
}

func (c *InMemoryDataset) Save() error {
	controller, ok := c.Project().Controller().(*InMemoryController)
	if !ok {
		return fmt.Errorf("not an in-memory controller")
	}
	return controller.SaveDataset(c)
}

func (c *InMemoryDataset) Delete() error {
	return fmt.Errorf("InMemoryDataset.Delete not implemented")
}
