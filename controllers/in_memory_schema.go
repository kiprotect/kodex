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
	"github.com/kiprotect/kiprotect"
	"time"
)

type InMemorySchema struct {
	kiprotect.BaseSchema
	data       interface{}
	dataSchema *kiprotect.DataSchema
	id         []byte
}

func MakeInMemorySchema(id []byte,
	dataSchema *kiprotect.DataSchema, project kiprotect.Project) *InMemorySchema {
	schema := &InMemorySchema{
		id:         id,
		dataSchema: dataSchema,
		BaseSchema: kiprotect.BaseSchema{
			Project_: project,
		},
	}
	schema.Self = schema
	return schema
}

func (i *InMemorySchema) Delete() error {
	return nil
}

func (i *InMemorySchema) ID() []byte {
	return i.id
}

func (i *InMemorySchema) CreatedAt() time.Time {
	return time.Now()
}

func (i *InMemorySchema) DeletedAt() *time.Time {
	return nil
}

func (i *InMemorySchema) UpdatedAt() time.Time {
	return time.Now()
}

func (i *InMemorySchema) Schema() *kiprotect.DataSchema {
	return i.dataSchema
}

func (i *InMemorySchema) Data() interface{} {
	return i.data
}

func (i *InMemorySchema) SetData(data interface{}) error {
	i.data = data
	return nil
}

func (i *InMemorySchema) SetSchema(dataSchema *kiprotect.DataSchema) error {
	i.dataSchema = dataSchema
	return nil
}

func (i *InMemorySchema) Refresh() error {
	return nil
}

func (i *InMemorySchema) Save() error {
	return nil
}
