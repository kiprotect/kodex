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
	"bytes"
	"fmt"
	"github.com/kiprotect/kodex"
	"time"
)

type ActionConfigMap struct {
	ActionConfig kodex.ActionConfig
}

type InMemoryConfig struct {
	kodex.BaseConfig
	id            []byte
	status        kodex.ConfigStatus
	name          string
	description   string
	version       string
	source        string
	createdAt     time.Time
	updatedAt     time.Time
	deletedAt     *time.Time
	destinations  map[string][]kodex.DestinationMap
	config        map[string]interface{}
	data          interface{}
	actionConfigs []*ActionConfigMap
}

func MakeInMemoryConfig(stream *InMemoryStream, id []byte, config map[string]interface{}) (*InMemoryConfig, error) {
	inMemoryConfig := &InMemoryConfig{
		id: id,
		BaseConfig: kodex.BaseConfig{
			Stream_: stream,
		},
		destinations:  make(map[string][]kodex.DestinationMap),
		actionConfigs: make([]*ActionConfigMap, 0, 10),
		config:        config,
	}
	inMemoryConfig.Self = inMemoryConfig
	return inMemoryConfig, nil
}

func (c *InMemoryConfig) Data() interface{} {
	return c.data
}

func (c *InMemoryConfig) SetData(data interface{}) error {
	c.data = data
	return nil
}

func (c *InMemoryConfig) ID() []byte {
	return c.id
}

func (c *InMemoryConfig) Destinations() (map[string][]kodex.DestinationMap, error) {
	return c.destinations, nil
}

func (i *InMemoryConfig) CreatedAt() time.Time {
	return i.createdAt
}

func (i *InMemoryConfig) DeletedAt() *time.Time {
	return i.deletedAt
}

func (i *InMemoryConfig) UpdatedAt() time.Time {
	return i.updatedAt
}

func (i *InMemoryConfig) SetCreatedAt(t time.Time) error {
	i.createdAt = t
	return nil
}

func (i *InMemoryConfig) SetUpdatedAt(t time.Time) error {
	i.updatedAt = t
	return nil
}

func (i *InMemoryConfig) SetDeletedAt(t *time.Time) error {
	i.deletedAt = t
	return nil
}

func (c *InMemoryConfig) Refresh() error {
	return nil
}

func (c *InMemoryConfig) Save() error {
	stream, ok := c.Stream().(*InMemoryStream)
	if !ok {
		return fmt.Errorf("not a inMemory stream")
	}
	return stream.SaveConfig(c)
}

func (c *InMemoryConfig) Delete() error {
	stream, ok := c.Stream().(*InMemoryStream)
	if !ok {
		return fmt.Errorf("not a inMemory stream")
	}
	return stream.DeleteConfig(c)
}

func (c *InMemoryConfig) ActionConfigs() ([]kodex.ActionConfig, error) {
	actionConfigs := make([]kodex.ActionConfig, len(c.actionConfigs))
	for i, actionConfig := range c.actionConfigs {
		actionConfigs[i] = actionConfig.ActionConfig
	}
	return actionConfigs, nil
}

func (c *InMemoryConfig) AddActionConfig(actionConfig kodex.ActionConfig, index int) error {
	if index > len(c.actionConfigs) || index < 0 {
		return fmt.Errorf("invalid index: out of bounds")
	}
	actionConfigMap := &ActionConfigMap{
		ActionConfig: actionConfig,
	}
	newActionConfigs := make([]*ActionConfigMap, 0, len(c.actionConfigs)+1)
	for i, existingActionConfig := range c.actionConfigs {
		if i == index {
			newActionConfigs = append(newActionConfigs, actionConfigMap)
		}
		newActionConfigs = append(newActionConfigs, existingActionConfig)
	}
	if index == len(c.actionConfigs) {
		newActionConfigs = append(newActionConfigs, actionConfigMap)
	}
	c.actionConfigs = newActionConfigs
	return nil
}

func (c *InMemoryConfig) RemoveActionConfig(actionConfig kodex.ActionConfig) error {
	newActionConfigs := make([]*ActionConfigMap, 0, len(c.actionConfigs)-1)
	for _, actionConfigMap := range c.actionConfigs {
		if bytes.Equal(actionConfigMap.ActionConfig.ID(), actionConfig.ID()) {
			continue
		}
		newActionConfigs = append(newActionConfigs, actionConfigMap)
	}
	c.actionConfigs = newActionConfigs
	return nil
}

func (c *InMemoryConfig) Status() kodex.ConfigStatus {
	return c.status
}

func (c *InMemoryConfig) SetStatus(status kodex.ConfigStatus) error {
	c.status = status
	return nil
}

func (c *InMemoryConfig) Description() string {
	return c.description
}

func (c *InMemoryConfig) SetDescription(description string) error {
	c.description = description
	return nil
}

func (c *InMemoryConfig) Name() string {
	return c.name
}

func (c *InMemoryConfig) SetName(name string) error {
	c.name = name
	return nil
}

func (c *InMemoryConfig) Version() string {
	return c.version
}

func (c *InMemoryConfig) SetVersion(version string) error {
	c.version = version
	return nil
}

func (c *InMemoryConfig) Source() string {
	return c.source
}

func (c *InMemoryConfig) SetSource(source string) error {
	c.source = source
	return nil
}

func (c *InMemoryConfig) AddDestination(destination kodex.Destination, name string, status kodex.DestinationStatus) error {
	inMemoryDestination, ok := destination.(*InMemoryDestination)
	if !ok {
		return fmt.Errorf("not an in-memory source")
	}
	if _, ok := c.destinations[name]; !ok {
		c.destinations[name] = make([]kodex.DestinationMap, 0, 1)
	}
	for _, destinationMap := range c.destinations[name] {
		if bytes.Equal(destinationMap.Destination().ID(), destination.ID()) && destinationMap.Name() == name {
			// this destination map already exists, we just update it...
			destinationMap.SetStatus(status)
			return nil
		}
	}
	c.destinations[name] = append(c.destinations[name], MakeInMemoryDestinationMap(kodex.RandomID(), name, c, inMemoryDestination, status))
	return nil
}

func (c *InMemoryConfig) RemoveDestination(destination kodex.Destination) error {
	for key, destinationMaps := range c.destinations {
		newDestinationMaps := make([]kodex.DestinationMap, 0, len(destinationMaps))
		found := false
		for _, existingDestinationMap := range destinationMaps {
			if string(existingDestinationMap.Destination().ID()) == string(destination.ID()) {
				found = true
				continue
			}
			newDestinationMaps = append(newDestinationMaps, existingDestinationMap)
		}
		if len(newDestinationMaps) == 0 {
			delete(c.destinations, key)
		} else {
			c.destinations[key] = newDestinationMaps
		}
		if found {
			return nil
		}
	}
	return fmt.Errorf("destination not found")
}
