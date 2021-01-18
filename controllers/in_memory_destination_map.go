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

type InMemoryDestinationMap struct {
	kodex.BaseDestinationMap
	name        string
	status      kodex.DestinationStatus
	destination *InMemoryDestination
	config      *InMemoryConfig
	id          []byte
}

func MakeInMemoryDestinationMap(id []byte, name string, config *InMemoryConfig, destination *InMemoryDestination, status kodex.DestinationStatus) *InMemoryDestinationMap {
	destinationMap := &InMemoryDestinationMap{
		id:                 id,
		name:               name,
		destination:        destination,
		config:             config,
		status:             status,
		BaseDestinationMap: kodex.BaseDestinationMap{},
	}
	destinationMap.Self = destinationMap
	return destinationMap
}

func (i *InMemoryDestinationMap) ID() []byte {
	return i.id
}

func (i *InMemoryDestinationMap) Delete() error {
	return fmt.Errorf("InMemoryDestinationMap.Delete not implemented")
}

func (i *InMemoryDestinationMap) Destination() kodex.Destination {
	return i.destination
}

func (i *InMemoryDestinationMap) Config() kodex.Config {
	return i.config
}

func (i *InMemoryDestinationMap) Status() kodex.DestinationStatus {
	return i.status
}

func (i *InMemoryDestinationMap) SetStatus(status kodex.DestinationStatus) error {
	i.status = status
	return nil
}

func (i *InMemoryDestinationMap) SetConfig(config kodex.Config) error {
	inMemoryConfig, ok := config.(*InMemoryConfig)
	if !ok {
		return fmt.Errorf("not a inMemory config")
	}
	i.config = inMemoryConfig
	return nil
}

func (i *InMemoryDestinationMap) SetDestination(destination kodex.Destination) error {
	inMemoryDestination, ok := destination.(*InMemoryDestination)
	if !ok {
		return fmt.Errorf("not a inMemory destination")
	}
	i.destination = inMemoryDestination
	return nil
}

func (i *InMemoryDestinationMap) CreatedAt() time.Time {
	return time.Now().UTC()
}

func (i *InMemoryDestinationMap) DeletedAt() *time.Time {
	return nil
}

func (i *InMemoryDestinationMap) UpdatedAt() time.Time {
	return time.Now().UTC()
}

func (i *InMemoryDestinationMap) Save() error {
	return nil
}

func (i *InMemoryDestinationMap) Refresh() error {
	return nil
}

func (i *InMemoryDestinationMap) Name() string {
	return i.name
}

func (i *InMemoryDestinationMap) SetName(name string) error {
	i.name = name
	return nil
}
