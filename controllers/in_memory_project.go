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

type InMemoryProject struct {
	kodex.BaseProject
	name        string
	description string
	createdAt   time.Time
	updatedAt   time.Time
	deletedAt   *time.Time
	data        interface{}
	id          []byte
}

func MakeInMemoryProject(id []byte, controller kodex.Controller) *InMemoryProject {
	destination := &InMemoryProject{
		id: id,
		BaseProject: kodex.BaseProject{
			Controller_: controller,
		},
	}
	destination.Self = destination
	return destination
}

func (i *InMemoryProject) ID() []byte {
	return i.id
}

func (i *InMemoryProject) InternalID() []byte {
	return i.id
}

func (i *InMemoryProject) Delete() error {

	if err := i.DeleteRelated(); err != nil {
		return err
	}

	controller, ok := i.Controller().(*InMemoryController)
	if !ok {
		return fmt.Errorf("not an in-memory controller")
	}
	return controller.DeleteProject(i)
}

func (i *InMemoryProject) CreatedAt() time.Time {
	return i.createdAt
}

func (i *InMemoryProject) DeletedAt() *time.Time {
	return i.deletedAt
}

func (i *InMemoryProject) UpdatedAt() time.Time {
	return i.updatedAt
}

func (i *InMemoryProject) SetCreatedAt(t time.Time) error {
	i.createdAt = t
	return nil
}

func (i *InMemoryProject) SetUpdatedAt(t time.Time) error {
	i.updatedAt = t
	return nil
}

func (i *InMemoryProject) SetDeletedAt(t *time.Time) error {
	i.deletedAt = t
	return nil
}

func (i *InMemoryProject) Save() error {
	controller, ok := i.Controller().(*InMemoryController)
	if !ok {
		return fmt.Errorf("not an in-memory controller")
	}
	return controller.SaveProject(i)
}

func (i *InMemoryProject) Refresh() error {
	return nil
}

func (i *InMemoryProject) Data() interface{} {
	return i.data
}

func (i *InMemoryProject) SetData(data interface{}) error {
	i.data = data
	return nil
}

func (i *InMemoryProject) Name() string {
	return i.name
}

func (i *InMemoryProject) SetName(name string) error {
	i.name = name
	return nil
}

func (i *InMemoryProject) Description() string {
	return i.description
}

func (i *InMemoryProject) SetDescription(description string) error {
	i.description = description
	return nil
}

func (c *InMemoryProject) MakeStream(id []byte) kodex.Stream {
	if id == nil {
		id = kodex.RandomID()
	}
	stream, err := MakeInMemoryStream(id, map[string]interface{}{
		"configs": []map[string]interface{}{},
		"params":  []map[string]interface{}{},
	}, c)
	if err != nil {
		// this should never happen
		panic(err)
	}
	return stream
}

func (c *InMemoryProject) MakeDataset(id []byte) kodex.Dataset {
	if id == nil {
		id = kodex.RandomID()
	}
	return MakeInMemoryDataset(id, c)
}

func (c *InMemoryProject) MakeActionConfig(id []byte) kodex.ActionConfig {
	if id == nil {
		id = kodex.RandomID()
	}
	return MakeInMemoryActionConfig(id, c)
}

func (c *InMemoryProject) MakeSource(id []byte) kodex.Source {
	if id == nil {
		id = kodex.RandomID()
	}
	return MakeInMemorySource(id, c)
}

func (c *InMemoryProject) MakeDestination(id []byte) kodex.Destination {
	if id == nil {
		id = kodex.RandomID()
	}
	return MakeInMemoryDestination(id, c)
}
