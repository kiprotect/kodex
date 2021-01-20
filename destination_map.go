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

package kodex

import (
	"encoding/json"
	"fmt"
)

type DestinationMap interface {
	Processable
	PriorityModel
	Destination() Destination
	Config() Config
	Name() string
	SetName(string) error
	Status() DestinationStatus
	SetStatus(DestinationStatus) error
	InternalWriter() (Writer, error)
}

/* Base Functionality */

type BaseDestinationMap struct {
	Self DestinationMap
}

func (b *BaseDestinationMap) Type() string {
	return "destination_map"
}

func (b *BaseDestinationMap) Update(values map[string]interface{}) error {
	return fmt.Errorf("BaseDestinationMap.Update implemented")
}

func (b *BaseDestinationMap) Create(values map[string]interface{}) error {
	return fmt.Errorf("BaseDestinationMap.Create not implemented")
}

func (b *BaseDestinationMap) createWriter() (Writer, error) {
	channel := MakeInternalChannel()
	if err := channel.Setup(b.Self.Destination().Project().Controller(), b.Self); err != nil {
		return nil, err
	}
	return MakeInternalWriter(channel), nil
}

func (b *BaseDestinationMap) InternalWriter() (Writer, error) {
	channel := MakeInternalChannel()
	if err := channel.Setup(b.Self.Destination().Project().Controller(), b.Self); err != nil {
		return nil, err
	}
	return MakeInternalWriter(channel), nil
}

func (b *BaseDestinationMap) MarshalJSON() ([]byte, error) {

	data := map[string]interface{}{
		"name":        b.Self.Name(),
		"status":      b.Self.Status(),
		"destination": b.Self.Destination(),
		"config":      b.Self.Config(),
	}

	for k, v := range JSONData(b.Self) {
		data[k] = v
	}

	return json.Marshal(data)
}
