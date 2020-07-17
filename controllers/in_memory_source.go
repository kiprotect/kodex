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
	"bytes"
	"fmt"
	"github.com/kiprotect/kiprotect"
	"time"
)

type InMemorySource struct {
	kiprotect.BaseSource
	name        string
	description string
	sourceType  string
	data        interface{}
	reader      kiprotect.Reader
	configData  map[string]interface{}
	id          []byte
}

func MakeInMemorySource(id []byte,
	project kiprotect.Project) *InMemorySource {

	source := &InMemorySource{
		id: id,
		BaseSource: kiprotect.BaseSource{
			Project_: project,
		},
	}

	source.Self = source
	return source
}

func (i *InMemorySource) Streams(status kiprotect.SourceStatus) ([]kiprotect.Stream, error) {

	streams := make([]kiprotect.Stream, 0)

	allStreams, err := i.Project().Controller().Streams(map[string]interface{}{})

	if err != nil {
		return nil, err
	}

	for _, stream := range allStreams {
		sources, err := stream.Sources()
		if err != nil {
			return nil, err
		}
		for _, sourceMap := range sources {
			if sourceMap.Status() == status && bytes.Equal(sourceMap.Source().ID(), i.ID()) {
				streams = append(streams, stream)
				break
			}
		}
	}
	return streams, nil
}

func (i *InMemorySource) Delete() error {
	return nil
}

func (i *InMemorySource) ID() []byte {
	return []byte(i.name)
}

func (c *InMemorySource) InternalID() []byte {
	return c.id
}

func (i *InMemorySource) CreatedAt() time.Time {
	return time.Now()
}

func (i *InMemorySource) DeletedAt() *time.Time {
	return nil
}

func (i *InMemorySource) UpdatedAt() time.Time {
	return time.Now()
}

func (i *InMemorySource) ConfigData() map[string]interface{} {
	return i.configData
}

func (i *InMemorySource) Data() interface{} {
	return i.data
}

func (i *InMemorySource) SetData(data interface{}) error {
	i.data = data
	return nil
}

func (i *InMemorySource) SetConfigData(configData map[string]interface{}) error {
	i.configData = configData
	return nil
}

func (i *InMemorySource) Name() string {
	return i.name
}

func (i *InMemorySource) SetName(name string) error {
	i.name = name
	return nil
}

func (i *InMemorySource) SourceType() string {
	return i.sourceType
}

func (i *InMemorySource) SetSourceType(sourceType string) error {
	i.sourceType = sourceType
	return nil
}

func (i *InMemorySource) Description() string {
	return i.description
}

func (i *InMemorySource) SetDescription(description string) error {
	i.description = description
	return nil
}

func (i *InMemorySource) Refresh() error {
	return nil
}

func (i *InMemorySource) Save() error {
	controller, ok := i.Project().Controller().(*InMemoryController)
	if !ok {
		return fmt.Errorf("not an in-memory controller")
	}
	return controller.SaveSource(i)
}

func (i *InMemorySource) Service() kiprotect.Service {
	return nil
}

func (i *InMemorySource) SetService(kiprotect.Service) error {
	return nil
}
