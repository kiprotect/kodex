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
	"github.com/kiprotect/kodex/api"
	"sync"
)

type ProcessorStats struct {
	ProcessorID    []byte
	IdleFraction   float64
	ItemsProcessed int64
	Capacity       float64
}

type Stats struct {
	ProcessorStats []ProcessorStats
}

type InMemoryController struct {
	kodex.BaseController
	mutex            sync.Mutex
	streams          map[string]kodex.Stream
	sources          map[string]kodex.Source
	destinations     map[string]kodex.Destination
	actionConfigs    map[string]kodex.ActionConfig
	projects         map[string]kodex.Project
	streamStats      map[string]Stats
	sourceStats      map[string]Stats
	destinationStats map[string]Stats
}

func MakeInMemoryController(config map[string]interface{}, settings kodex.Settings, definitions *kodex.Definitions) (kodex.Controller, error) {
	var err error
	var baseController kodex.BaseController

	if baseController, err = kodex.MakeBaseController(settings, definitions); err != nil {
		return nil, err
	}

	controller := InMemoryController{
		BaseController:   baseController,
		streamStats:      make(map[string]Stats),
		sourceStats:      make(map[string]Stats),
		destinationStats: make(map[string]Stats),
		destinations:     make(map[string]kodex.Destination),
		actionConfigs:    make(map[string]kodex.ActionConfig),
		projects:         make(map[string]kodex.Project),
		streams:          make(map[string]kodex.Stream),
		sources:          make(map[string]kodex.Source),
	}

	return &controller, nil
}

func (c *InMemoryController) SaveActionConfig(actionConfig kodex.ActionConfig) error {
	inMemoryActionConfig, ok := actionConfig.(*InMemoryActionConfig)
	if !ok {
		return fmt.Errorf("not an in-memory action config")
	}
	if existingActionConfig, ok := c.actionConfigs[string(actionConfig.ID())].(*InMemoryActionConfig); ok {
		if bytes.Equal(existingActionConfig.InternalID(), inMemoryActionConfig.InternalID()) && existingActionConfig != inMemoryActionConfig {
			return fmt.Errorf("ID conflict")
		}
	}
	c.actionConfigs[string(actionConfig.ID())] = actionConfig
	return nil
}

func (c *InMemoryController) SaveSource(source kodex.Source) error {
	inMemorySource, ok := source.(*InMemorySource)
	if !ok {
		return fmt.Errorf("not an in-memory action config")
	}
	if existingSource, ok := c.sources[string(source.ID())].(*InMemorySource); ok {
		if bytes.Equal(existingSource.InternalID(), inMemorySource.InternalID()) && existingSource != inMemorySource {
			return fmt.Errorf("ID conflict")
		}
	}
	c.sources[string(source.ID())] = source
	return nil
}

func (c *InMemoryController) SaveDestination(destination kodex.Destination) error {
	inMemoryDestination, ok := destination.(*InMemoryDestination)
	if !ok {
		return fmt.Errorf("not an in-memory action config")
	}
	if existingDestination, ok := c.destinations[string(destination.ID())].(*InMemoryDestination); ok {
		if bytes.Equal(existingDestination.InternalID(), inMemoryDestination.InternalID()) && existingDestination != inMemoryDestination {
			return fmt.Errorf("ID conflict")
		}
	}
	c.destinations[string(destination.ID())] = destination
	return nil
}

func (c *InMemoryController) DeleteStream(stream *InMemoryStream) error {
	delete(c.streams, string(stream.ID()))
	return nil
}

func (c *InMemoryController) SaveStream(stream kodex.Stream) error {
	inMemoryStream, ok := stream.(*InMemoryStream)
	if !ok {
		return fmt.Errorf("not an in-memory action config")
	}
	if existingStream, ok := c.streams[string(stream.ID())].(*InMemoryStream); ok {
		if bytes.Equal(existingStream.InternalID(), inMemoryStream.InternalID()) && existingStream != inMemoryStream {
			return fmt.Errorf("ID conflict")
		}
	}
	c.streams[string(stream.ID())] = stream
	return nil
}

func (c *InMemoryController) SaveProject(project kodex.Project) error {
	inMemoryProject, ok := project.(*InMemoryProject)
	if !ok {
		return fmt.Errorf("not an in-memory action config")
	}
	if existingProject, ok := c.projects[string(project.ID())].(*InMemoryProject); ok {
		if bytes.Equal(existingProject.InternalID(), inMemoryProject.InternalID()) && existingProject != inMemoryProject {
			return fmt.Errorf("ID conflict")
		}
	}
	c.projects[string(project.ID())] = project
	return nil
}

func (c *InMemoryController) Begin() error {
	return nil
}

func (c *InMemoryController) Commit() error {
	return nil
}

func (c *InMemoryController) Rollback() error {
	return nil
}

/* Stream Management */

// Return a list of streams identified by the list of IDs and in addition
// filtered by the given arguments
func (c *InMemoryController) Streams(filters map[string]interface{}) ([]kodex.Stream, error) {
	streams := make([]kodex.Stream, 0)
outer:
	for _, stream := range c.streams {
		for key, value := range filters {
			switch key {
			case "ProjectID":
				bytesValue, ok := value.([]byte)
				if !ok {
					return nil, fmt.Errorf("expected at name")
				}
				if !bytes.Equal(stream.Project().ID(), bytesValue) {
					continue outer
				}
			case "name":
				strValue, ok := value.(string)
				if !ok {
					return nil, fmt.Errorf("expected a name")
				}
				if stream.Name() != strValue {
					continue outer
				}
			default:
				return nil, fmt.Errorf("unknown filter key: %s", key)
			}
		}
		streams = append(streams, stream)
	}
	return streams, nil
}

func (c *InMemoryController) Stream(streamID []byte) (kodex.Stream, error) {

	for _, stream := range c.streams {
		if bytes.Equal(stream.ID(), streamID) {
			return stream, nil
		}
	}

	return nil, kodex.NotFound
}

func (c *InMemoryController) Config(configID []byte) (kodex.Config, error) {
	for _, stream := range c.streams {
		configs, err := stream.Configs()
		if err != nil {
			return nil, err
		}
		for _, config := range configs {
			if bytes.Equal(config.ID(), configID) {
				return config, nil
			}
		}
	}
	return nil, kodex.NotFound
}

func (c *InMemoryController) ActionConfig(actionConfigID []byte) (kodex.ActionConfig, error) {
	for _, actionConfig := range c.actionConfigs {
		if bytes.Equal(actionConfig.ID(), actionConfigID) {
			return actionConfig, nil
		}
	}
	return nil, kodex.NotFound
}

/* Action Config Management */

func (c *InMemoryController) ActionConfigs(filters map[string]interface{}) ([]kodex.ActionConfig, error) {

	actionConfigs := make([]kodex.ActionConfig, 0, len(c.actionConfigs))

outer:
	for _, actionConfig := range c.actionConfigs {
		for key, value := range filters {
			switch key {
			case "ProjectID":
				bytesValue, ok := value.([]byte)
				if !ok {
					return nil, fmt.Errorf("expected at name")
				}
				if !bytes.Equal(actionConfig.Project().ID(), bytesValue) {
					continue outer
				}
			case "name":
				strValue, ok := value.(string)
				if !ok {
					return nil, fmt.Errorf("expected a name")
				}
				if actionConfig.Name() != strValue {
					continue outer
				}
			default:
				return nil, fmt.Errorf("unknown filter key: %s", key)
			}
		}
		actionConfigs = append(actionConfigs, actionConfig)
	}

	return actionConfigs, nil
}

/* Source Management */

func (c *InMemoryController) Sources(filters map[string]interface{}) ([]kodex.Source, error) {
	sources := make([]kodex.Source, 0, len(c.sources))

outer:
	for _, source := range c.sources {
		for key, value := range filters {
			switch key {
			case "ProjectID":
				bytesValue, ok := value.([]byte)
				if !ok {
					return nil, fmt.Errorf("expected at name")
				}
				if !bytes.Equal(source.Project().ID(), bytesValue) {
					continue outer
				}
			case "name":
				strValue, ok := value.(string)
				if !ok {
					return nil, fmt.Errorf("expected a name")
				}
				if source.Name() != strValue {
					continue outer
				}
			default:
				return nil, fmt.Errorf("unknown filter key: %s", key)
			}
		}
		sources = append(sources, source)
	}

	return sources, nil
}

func (c *InMemoryController) Source(sourceID []byte) (kodex.Source, error) {
	for _, source := range c.sources {
		if bytes.Equal(source.ID(), sourceID) {
			return source, nil
		}
	}
	return nil, kodex.NotFound
}

/* Destination Management */

func (c *InMemoryController) Destinations(filters map[string]interface{}) ([]kodex.Destination, error) {
	destinations := make([]kodex.Destination, 0, len(c.sources))

outer:
	for _, destination := range c.destinations {
		for key, value := range filters {
			switch key {
			case "ProjectID":
				bytesValue, ok := value.([]byte)
				if !ok {
					return nil, fmt.Errorf("expected at name")
				}
				if !bytes.Equal(destination.Project().ID(), bytesValue) {
					continue outer
				}
			case "name":
				strValue, ok := value.(string)
				if !ok {
					return nil, fmt.Errorf("expected a name")
				}
				if destination.Name() != strValue {
					continue outer
				}
			default:
				return nil, fmt.Errorf("unknown filter key: %s", key)
			}
		}
		destinations = append(destinations, destination)
	}

	return destinations, nil
}

func (c *InMemoryController) Destination(destinationID []byte) (kodex.Destination, error) {
	for _, destination := range c.destinations {
		if bytes.Equal(destination.ID(), destinationID) {
			return destination, nil
		}
	}
	return nil, kodex.NotFound
}

func (c *InMemoryController) StreamsByUrgency(n int) ([]kodex.Stream, error) {

	streams := make([]kodex.Stream, 0)
	for _, stream := range c.streams {
		streams = append(streams, stream)
		if len(streams) >= n {
			break
		}
	}
	return streams, nil
}

func (c *InMemoryController) SourcesByUrgency(n int) ([]kodex.SourceMap, error) {

	sources := make([]kodex.SourceMap, 0)
OUTER:
	for _, stream := range c.streams {
		streamSources, err := stream.Sources()
		if err != nil {
			return nil, err
		}

		for _, source := range streamSources {
			sources = append(sources, source)
			if len(sources) >= n {
				break OUTER
			}
		}
	}
	return sources, nil
}

func (c *InMemoryController) DestinationsByUrgency(n int) ([]kodex.DestinationMap, error) {
	destinations := make([]kodex.DestinationMap, 0)
OUTER:
	for _, stream := range c.streams {
		streamConfigs, err := stream.Configs()
		if err != nil {
			return nil, err
		}

		for _, config := range streamConfigs {
			configDestinations, err := config.Destinations()
			if err != nil {
				return nil, err
			}
			for _, destinationMaps := range configDestinations {
				for _, destinationMap := range destinationMaps {
					destinations = append(destinations, destinationMap)
					if len(destinations) >= n {
						break OUTER
					}
				}
			}
		}
	}
	return destinations, nil
}

func (c *InMemoryController) getTable(processable kodex.Processable) (map[string]Stats, error) {
	switch processable.Type() {
	case "stream":
		return c.streamStats, nil
	case "source":
		return c.sourceStats, nil
	case "destination":
		return c.destinationStats, nil
	default:
		return nil, fmt.Errorf("invalid type: %s", processable.Type())
	}
}

// Acquire a processable entity
func (c *InMemoryController) Acquire(processable kodex.Processable, processorID []byte) (bool, error) {

	c.mutex.Lock()
	defer c.mutex.Unlock()

	table, err := c.getTable(processable)

	if err != nil {
		return false, err
	}

	pId := string(processable.ID())
	if _, ok := table[pId]; ok {
		return false, nil
	}

	table[pId] = Stats{
		ProcessorStats: []ProcessorStats{
			ProcessorStats{
				ProcessorID:    processorID,
				IdleFraction:   0,
				ItemsProcessed: 0,
				Capacity:       0,
			},
		},
	}

	return true, nil
}

// Release a processable entity
func (c *InMemoryController) Release(processable kodex.Processable, processorID []byte) (bool, error) {

	c.mutex.Lock()
	defer c.mutex.Unlock()

	table, err := c.getTable(processable)

	if err != nil {
		return false, err
	}

	pId := string(processable.ID())
	if _, ok := table[pId]; ok {
		delete(table, pId)
		return true, nil
	}
	return false, nil
}

// Send a pingback with stats for a processable entity
func (c *InMemoryController) Ping(processable kodex.Processable, stats kodex.ProcessingStats) error {
	return nil
}

/* Project Management */

func (c *InMemoryController) Project(id []byte) (kodex.Project, error) {
	for _, project := range c.projects {
		if bytes.Equal(project.ID(), id) {
			return project, nil
		}
	}
	return nil, kodex.NotFound
}

func (c *InMemoryController) Projects(filters map[string]interface{}) ([]kodex.Project, error) {
	projects := make([]kodex.Project, 0, len(c.projects))
outer:
	for _, project := range c.projects {
		for key, value := range filters {
			switch key {
			case "ID":
				switch tv := value.(type) {
				case []byte:
					if !bytes.Equal(project.ID(), tv) {
						continue outer
					}
				case api.In:
					found := false
					for _, id := range tv.Values {
						if bytesID, ok := id.([]byte); !ok {
							return nil, fmt.Errorf("invalid ID")
						} else if bytes.Equal(project.ID(), bytesID) {
							found = true
						}
					}
					if !found {
						continue outer
					}
				default:
					return nil, fmt.Errorf("invalid type")
				}
			case "name":
				strValue, ok := value.(string)
				if !ok {
					return nil, fmt.Errorf("expected a name")
				}
				if project.Name() != strValue {
					continue outer
				}
			default:
				return nil, fmt.Errorf("unknown filter key: %s", key)
			}
		}
		projects = append(projects, project)
	}

	return projects, nil
}

func (c *InMemoryController) MakeProject(id []byte) kodex.Project {
	if id == nil {
		id = kodex.RandomID()
	}
	return MakeInMemoryProject(id, c)
}

func (c *InMemoryController) ResetDB() error {
	return nil
}
