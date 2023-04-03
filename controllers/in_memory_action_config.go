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
	"fmt"
	"github.com/kiprotect/kodex"
	"time"
)

type InMemoryActionConfig struct {
	kodex.BaseActionConfig
	id          []byte
	index       int
	actionType  string
	description string
	name        string
	data        interface{}
	configData  map[string]interface{}
}

func MakeInMemoryActionConfig(id []byte, project kodex.Project) *InMemoryActionConfig {
	inMemoryActionConfig := &InMemoryActionConfig{
		id: id,
		BaseActionConfig: kodex.BaseActionConfig{
			Project_: project,
		},
	}
	inMemoryActionConfig.Self = inMemoryActionConfig
	return inMemoryActionConfig
}

func (c *InMemoryActionConfig) ID() []byte {
	return c.id
}

func (c *InMemoryActionConfig) InternalID() []byte {
	return c.id
}

func (c *InMemoryActionConfig) Index() int {
	return c.index
}

func (c *InMemoryActionConfig) SetIndex(index int) error {
	c.index = index
	return nil
}

func (c *InMemoryActionConfig) Name() string {
	return c.name
}

func (c *InMemoryActionConfig) ConfigData() map[string]interface{} {
	return c.configData
}

func (c *InMemoryActionConfig) SetConfigData(configData map[string]interface{}) error {
	c.configData = configData
	return nil
}

func (c *InMemoryActionConfig) Data() interface{} {
	return c.data
}

func (c *InMemoryActionConfig) SetData(data interface{}) error {
	c.data = data
	return nil
}

func (c *InMemoryActionConfig) Description() string {
	return c.description
}

func (c *InMemoryActionConfig) ActionType() string {
	return c.actionType
}

func (c *InMemoryActionConfig) SetActionType(actionType string) error {
	c.actionType = actionType
	return nil
}

func (c *InMemoryActionConfig) SetName(name string) error {
	c.name = name
	return nil
}

func (c *InMemoryActionConfig) SetDescription(description string) error {
	c.description = description
	return nil
}

func (c *InMemoryActionConfig) UpdatedAt() time.Time {
	return time.Now().UTC()
}

func (c *InMemoryActionConfig) CreatedAt() time.Time {
	return time.Now().UTC()
}

func (c *InMemoryActionConfig) DeletedAt() *time.Time {
	return nil
}

func (c *InMemoryActionConfig) Refresh() error {
	// we don't need to do anything here
	return nil
}

func (c *InMemoryActionConfig) Save() error {
	controller, ok := c.Project().Controller().(*InMemoryController)
	if !ok {
		return fmt.Errorf("not an in-memory controller")
	}
	return controller.SaveActionConfig(c)
}

func (c *InMemoryActionConfig) Delete() error {
	return fmt.Errorf("InMemoryActionConfig.Delete not implemented")
}
