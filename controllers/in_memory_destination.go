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

type InMemoryDestination struct {
	kodex.BaseDestination
	name            string
	description     string
	createdAt       time.Time
	updatedAt       time.Time
	deletedAt       *time.Time
	data            interface{}
	destinationType string
	configData      map[string]interface{}
	id              []byte
}

func MakeInMemoryDestination(id []byte, project kodex.Project) *InMemoryDestination {
	destination := &InMemoryDestination{
		id: id,
		BaseDestination: kodex.BaseDestination{
			Project_: project,
		},
	}
	destination.Self = destination
	return destination
}

func (i *InMemoryDestination) ID() []byte {
	return []byte(i.name)
}

func (i *InMemoryDestination) InternalID() []byte {
	return i.id
}

func (i *InMemoryDestination) Delete() error {
	controller, ok := i.Project().Controller().(*InMemoryController)
	if !ok {
		return fmt.Errorf("not an in-memory controller")
	}
	return controller.DeleteDestination(i)
}

func (i *InMemoryDestination) CreatedAt() time.Time {
	return i.createdAt
}

func (i *InMemoryDestination) DeletedAt() *time.Time {
	return i.deletedAt
}

func (i *InMemoryDestination) UpdatedAt() time.Time {
	return i.updatedAt
}

func (i *InMemoryDestination) SetCreatedAt(t time.Time) error {
	i.createdAt = t
	return nil
}

func (i *InMemoryDestination) SetUpdatedAt(t time.Time) error {
	i.updatedAt = t
	return nil
}

func (i *InMemoryDestination) SetDeletedAt(t *time.Time) error {
	i.deletedAt = t
	return nil
}

func (i *InMemoryDestination) Save() error {
	controller, ok := i.Project().Controller().(*InMemoryController)
	if !ok {
		return fmt.Errorf("not an in-memory controller")
	}
	return controller.SaveDestination(i)
}

func (i *InMemoryDestination) Refresh() error {
	return nil
}

func (i *InMemoryDestination) Data() interface{} {
	return i.data
}

func (i *InMemoryDestination) SetData(data interface{}) error {
	i.data = data
	return nil
}

func (i *InMemoryDestination) ConfigData() map[string]interface{} {
	return i.configData
}

func (i *InMemoryDestination) SetConfigData(configData map[string]interface{}) error {
	i.configData = configData
	return nil
}

func (i *InMemoryDestination) Name() string {
	return i.name
}

func (i *InMemoryDestination) SetName(name string) error {
	i.name = name
	return nil
}

func (i *InMemoryDestination) DestinationType() string {
	return i.destinationType
}

func (i *InMemoryDestination) SetDestinationType(destinationType string) error {
	i.destinationType = destinationType
	return nil
}

func (i *InMemoryDestination) Description() string {
	return i.description
}

func (i *InMemoryDestination) SetDescription(description string) error {
	i.description = description
	return nil
}
