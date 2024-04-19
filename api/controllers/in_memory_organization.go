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
	"github.com/kiprotect/kodex/api"
	"time"
)

type InMemoryOrganization struct {
	api.BaseOrganization
	name        string
	description string
	source      string
	sourceID    []byte
	data        interface{}
	id          []byte
}

func MakeInMemoryOrganization(id []byte,
	controller api.Controller) *InMemoryOrganization {

	organization := &InMemoryOrganization{
		id: id,
		BaseOrganization: api.BaseOrganization{
			Controller_: controller,
		},
	}
	organization.Self = organization
	return organization
}

func (i *InMemoryOrganization) Delete() error {
	return fmt.Errorf("not implemented")
}

func (i *InMemoryOrganization) ID() []byte {
	return i.id
}

func (i *InMemoryOrganization) CreatedAt() time.Time {
	return time.Now().UTC()
}

func (i *InMemoryOrganization) DeletedAt() *time.Time {
	return nil
}

func (i *InMemoryOrganization) UpdatedAt() time.Time {
	return time.Now().UTC()
}

func (i *InMemoryOrganization) Data() interface{} {
	return i.data
}

func (i *InMemoryOrganization) SetData(data interface{}) error {
	i.data = data
	return nil
}

func (i *InMemoryOrganization) Name() string {
	return i.name
}

func (i *InMemoryOrganization) SetName(name string) error {
	i.name = name
	return nil
}

func (i *InMemoryOrganization) Source() string {
	return i.source
}

func (i *InMemoryOrganization) SetSource(source string) error {
	i.source = source
	return nil
}

func (i *InMemoryOrganization) SourceID() []byte {
	return i.sourceID
}

func (i *InMemoryOrganization) SetSourceID(sourceID []byte) error {
	i.sourceID = sourceID
	return nil
}

func (i *InMemoryOrganization) Description() string {
	return i.description
}

func (i *InMemoryOrganization) SetDescription(description string) error {
	i.description = description
	return nil
}

func (i *InMemoryOrganization) Refresh() error {
	return nil
}

func (i *InMemoryOrganization) Save() error {
	inMemoryController, ok := i.Controller().(*InMemoryController)
	if !ok {
		return fmt.Errorf("expected an InMemory controller")
	}
	return inMemoryController.SaveOrganization(i)
}
