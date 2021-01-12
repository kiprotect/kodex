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
	"fmt"
	"github.com/kiprotect/kodex"
	"time"
)

type InMemoryProject struct {
	kodex.BaseProject
	name        string
	description string
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
	return []byte(i.name)
}

func (i *InMemoryProject) InternalID() []byte {
	return i.id
}

func (i *InMemoryProject) Delete() error {
	return fmt.Errorf("InMemoryProject.Delete not implemented")
}

func (i *InMemoryProject) CreatedAt() time.Time {
	return time.Now()
}

func (i *InMemoryProject) DeletedAt() *time.Time {
	return nil
}

func (i *InMemoryProject) UpdatedAt() time.Time {
	return time.Now()
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

func (c *InMemoryProject) MakeStream() kodex.Stream {
	id := kodex.RandomID()
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

func (c *InMemoryProject) MakeActionConfig() kodex.ActionConfig {
	return MakeInMemoryActionConfig(kodex.RandomID(), c)
}

func (c *InMemoryProject) MakeSource() kodex.Source {
	return MakeInMemorySource(kodex.RandomID(), c)
}

func (c *InMemoryProject) MakeDestination() kodex.Destination {
	return MakeInMemoryDestination(kodex.RandomID(), c)
}
