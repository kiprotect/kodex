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
	"bytes"
	"fmt"
	"github.com/kiprotect/kodex"
	"time"
)

type InMemoryStream struct {
	kodex.BaseStream
	id          []byte
	name        string
	status      kodex.StreamStatus
	data        interface{}
	description string
	config      map[string]interface{}
	configs     []kodex.Config
	sources     map[string]kodex.SourceMap
}

func MakeInMemoryStream(id []byte, config map[string]interface{}, project *InMemoryProject) (kodex.Stream, error) {
	stream := &InMemoryStream{
		config: config,
		BaseStream: kodex.BaseStream{
			Project_: project,
		},
		configs: make([]kodex.Config, 0),
		sources: make(map[string]kodex.SourceMap),
		id:      id,
	}
	stream.Self = stream
	return stream, nil
}

func (c *InMemoryStream) Config(configID []byte) (kodex.Config, error) {
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
	newConfigs := make([]kodex.Config, 0)
	for _, config := range c.configs {
		if bytes.Equal(config.ID(), dc.ID()) {
			continue
		}
		newConfigs = append(newConfigs, config)
	}
	return nil
}

func (c *InMemoryStream) MakeConfig() kodex.Config {
	config, err := MakeInMemoryConfig(c, kodex.RandomID(), map[string]interface{}{})
	if err != nil {
		panic(err)
	}
	return config
}

func (c *InMemoryStream) SaveConfig(config kodex.Config) error {
	for _, config := range c.configs {
		if string(config.ID()) == string(config.ID()) {
			return nil
		}
	}
	c.configs = append(c.configs, config)
	return nil
}

func (c *InMemoryStream) Configs() ([]kodex.Config, error) {
	return c.configs, nil
}

func (c *InMemoryStream) ID() []byte {
	return []byte(c.name)
}

func (c *InMemoryStream) InternalID() []byte {
	return c.id
}

func (i *InMemoryStream) Status() kodex.StreamStatus {
	return i.status
}

func (i *InMemoryStream) SetStatus(status kodex.StreamStatus) error {
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
	// we call the deletion hooks
	if _, err := i.Project().Controller().RunHooks("stream.delete", i); err != nil {
		return err
	}
	return fmt.Errorf("InMemoryStream.Delete not implemented")
}

func (c *InMemoryStream) AddSource(source kodex.Source, status kodex.SourceStatus) error {
	inMemorySource, ok := source.(*InMemorySource)
	if !ok {
		return fmt.Errorf("not an in-memory source")
	}
	c.sources[string(source.ID())] = MakeInMemorySourceMap(kodex.RandomID(), c, inMemorySource, status)
	return nil
}

func (c *InMemoryStream) RemoveSource(source kodex.Source) error {
	for id, _ := range c.sources {
		if id == string(source.ID()) {
			delete(c.sources, id)
			return nil
		}
	}
	return fmt.Errorf("source not found")
}

func (c *InMemoryStream) Sources() (map[string]kodex.SourceMap, error) {
	return c.sources, nil
}
