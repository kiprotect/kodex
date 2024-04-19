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

type InMemoryUser struct {
	api.BaseUser
	displayName string
	email       string
	superuser   bool
	source      string
	sourceID    []byte
	data        interface{}
	id          []byte
}

func MakeInMemoryUser(id []byte,
	controller api.Controller) *InMemoryUser {

	organization := &InMemoryUser{
		id: id,
		BaseUser: api.BaseUser{
			Controller_: controller,
		},
	}
	organization.Self = organization
	return organization
}

func (i *InMemoryUser) Delete() error {
	return fmt.Errorf("not implemented")
}

func (i *InMemoryUser) ID() []byte {
	return i.id
}

func (i *InMemoryUser) CreatedAt() time.Time {
	return time.Now().UTC()
}

func (i *InMemoryUser) DeletedAt() *time.Time {
	return nil
}

func (i *InMemoryUser) UpdatedAt() time.Time {
	return time.Now().UTC()
}

func (i *InMemoryUser) Data() interface{} {
	return i.data
}

func (i *InMemoryUser) SetData(data interface{}) error {
	i.data = data
	return nil
}

func (i *InMemoryUser) Email() string {
	return i.email
}

func (i *InMemoryUser) SetEmail(email string) error {
	i.email = email
	return nil
}

func (i *InMemoryUser) Superuser() bool {
	return i.superuser
}

func (i *InMemoryUser) SetSuperuser(superuser bool) error {
	i.superuser = superuser
	return nil
}

func (i *InMemoryUser) DisplayName() string {
	return i.displayName
}

func (i *InMemoryUser) SetDisplayName(displayName string) error {
	i.displayName = displayName
	return nil
}

func (i *InMemoryUser) Source() string {
	return i.source
}

func (i *InMemoryUser) SetSource(source string) error {
	i.source = source
	return nil
}

func (i *InMemoryUser) SourceID() []byte {
	return i.sourceID
}

func (i *InMemoryUser) SetSourceID(sourceID []byte) error {
	i.sourceID = sourceID
	return nil
}

func (i *InMemoryUser) Refresh() error {
	return nil
}

func (i *InMemoryUser) Save() error {
	inMemoryController, ok := i.Controller().(*InMemoryController)
	if !ok {
		return fmt.Errorf("expected an InMemory controller")
	}
	return inMemoryController.SaveUser(i)
}
