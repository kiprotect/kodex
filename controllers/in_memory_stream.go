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

type InMemoryStream struct {
	kiprotect.BaseStream
	id          []byte
	name        string
	status      kiprotect.StreamStatus
	data        interface{}
	description string
	config      map[string]interface{}
	configs     []kiprotect.Config
	sources     map[string]kiprotect.SourceMap
}

func MakeInMemoryStream(id []byte, config map[string]interface{}, project *InMemoryProject) (kiprotect.Stream, error) {
	stream := &InMemoryStream{
		config: config,
		BaseStream: kiprotect.BaseStream{
			Project_: project,
		},
		configs: make([]kiprotect.Config, 0),
		sources: make(map[string]kiprotect.SourceMap),
		id:      id,
	}
	stream.Self = stream
	return stream, nil
}

func (c *InMemoryStream) Config(configID []byte) (kiprotect.Config, error) {
	for _, config := range c.configs {
		if bytes.Equal(config.ID(), configID) {
			return config, nil
		}
	}
	return nil, fmt.Errorf("config not found")
}

func (c *InMemoryStream) Data() interface{} {
	return c.data
}

func (c *InMemoryStream) SetData(data interface{}) error {
	c.data = data
	return nil
}

func (c *InMemoryStream) DeleteConfig(dc *InMemoryConfig) error {
	newConfigs := make([]kiprotect.Config, 0)
	for _, config := range c.configs {
		if bytes.Equal(config.ID(), dc.ID()) {
			continue
		}
		newConfigs = append(newConfigs, config)
	}
	return nil
}

func (c *InMemoryStream) MakeConfig() kiprotect.Config {
	config, err := MakeInMemoryConfig(c, kiprotect.RandomID(), map[string]interface{}{})
	if err != nil {
		panic(err)
	}
	return config
}

func (c *InMemoryStream) SaveConfig(config kiprotect.Config) error {
	for _, config := range c.configs {
		if string(config.ID()) == string(config.ID()) {
			return nil
		}
	}
	c.configs = append(c.configs, config)
	return nil
}

func (c *InMemoryStream) Configs() ([]kiprotect.Config, error) {
	return c.configs, nil
}

func (c *InMemoryStream) ID() []byte {
	return []byte(c.name)
}

func (c *InMemoryStream) InternalID() []byte {
	return c.id
}

func (i *InMemoryStream) Status() kiprotect.StreamStatus {
	return i.status
}

func (i *InMemoryStream) SetStatus(status kiprotect.StreamStatus) error {
	i.status = status
	return nil
}

func (i *InMemoryStream) Name() string {
	return i.name
}

func (i *InMemoryStream) SetName(name string) error {
	i.name = name
	return nil
}

func (i *InMemoryStream) Description() string {
	return i.description
}

func (i *InMemoryStream) SetDescription(description string) error {
	i.description = description
	return nil
}

func (i *InMemoryStream) Save() error {
	controller, ok := i.Project().Controller().(*InMemoryController)
	if !ok {
		return fmt.Errorf("not an in-memory controller")
	}
	return controller.SaveStream(i)
}

func (i *InMemoryStream) Refresh() error {
	return nil
}

func (i *InMemoryStream) CreatedAt() time.Time {
	return time.Now()
}

func (i *InMemoryStream) DeletedAt() *time.Time {
	return nil
}

func (i *InMemoryStream) UpdatedAt() time.Time {
	return time.Now()
}

func (i *InMemoryStream) Delete() error {
	return fmt.Errorf("InMemoryStream.Delete not implemented")
}

func (c *InMemoryStream) AddSource(source kiprotect.Source, status kiprotect.SourceStatus) error {
	inMemorySource, ok := source.(*InMemorySource)
	if !ok {
		return fmt.Errorf("not an in-memory source")
	}
	c.sources[string(source.ID())] = MakeInMemorySourceMap(kiprotect.RandomID(), c, inMemorySource, status)
	return nil
}

func (c *InMemoryStream) RemoveSource(source kiprotect.Source) error {
	for id, _ := range c.sources {
		if id == string(source.ID()) {
			delete(c.sources, id)
			return nil
		}
	}
	return fmt.Errorf("source not found")
}

func (c *InMemoryStream) Sources() (map[string]kiprotect.SourceMap, error) {
	return c.sources, nil
}
