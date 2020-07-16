// KIProtect (Community Edition - CE) - Privacy & Security Engineering Platform
// Copyright (C) 2020  KIProtect GmbH (HRB 208395B) - Germany
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
	"github.com/kiprotect/kiprotect"
	"time"
)

type InMemoryProject struct {
	kiprotect.BaseProject
	name        string
	description string
	data        interface{}
	id          []byte
}

func MakeInMemoryProject(id []byte, controller kiprotect.Controller) *InMemoryProject {
	destination := &InMemoryProject{
		id: id,
		BaseProject: kiprotect.BaseProject{
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

func (i *InMemoryProject) MakeSchema() kiprotect.Schema {
	return MakeInMemorySchema(kiprotect.RandomID(), &kiprotect.DataSchema{}, i)
}

func (c *InMemoryProject) MakeStream() kiprotect.Stream {
	id := kiprotect.RandomID()
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

func (c *InMemoryProject) MakeActionConfig() kiprotect.ActionConfig {
	return MakeInMemoryActionConfig(kiprotect.RandomID(), c)
}

func (c *InMemoryProject) MakeSource() kiprotect.Source {
	return MakeInMemorySource(kiprotect.RandomID(), c)
}

func (c *InMemoryProject) MakeDestination() kiprotect.Destination {
	return MakeInMemoryDestination(kiprotect.RandomID(), c)
}
