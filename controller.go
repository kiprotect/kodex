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
	"fmt"
)

type ControllerMaker func(map[string]interface{}, Settings, *Definitions) (Controller, error)
type ControllerDefinitions map[string]ControllerMaker

var NotFound = fmt.Errorf("object not found")

type Controller interface {
	SetVar(key string, value interface{}) error
	GetVar(key string) (interface{}, bool)

	// Transaction Helpers
	Begin() error
	Commit() error
	Rollback() error

	// Retrieve Settings
	Settings() Settings

	// Initialize a plugin
	InitializePlugin(Plugin) error
	// Initialize all plugins as defined in the settings
	InitializePlugins() error

	// Streams
	Streams(filters map[string]interface{}) ([]Stream, error)
	Stream(streamID []byte) (Stream, error)

	// Sources
	Sources(filters map[string]interface{}) ([]Source, error)
	Source(sourceID []byte) (Source, error)

	// Destinations
	Destinations(filters map[string]interface{}) ([]Destination, error)
	Destination(destinationID []byte) (Destination, error)

	// Configs
	Config(configID []byte) (Config, error)

	// Action Configs
	ActionConfigs(filters map[string]interface{}) ([]ActionConfig, error)
	ActionConfig(configID []byte) (ActionConfig, error)

	Definitions() *Definitions

	// Retrieve a list of streams by urgency
	StreamsByUrgency(n int) ([]Stream, error)
	// Retrieve a list of sources by urgency
	SourcesByUrgency(n int) ([]SourceMap, error)
	// Retrieve a list of destinations by urgency
	DestinationsByUrgency(n int) ([]DestinationMap, error)

	// Acquire a processable entity
	Acquire(Processable, []byte) (bool, error)
	// Release a processable entity
	Release(Processable, []byte) (bool, error)
	// Send a pingback with stats for a processable entity
	Ping(Processable, ProcessingStats) error

	// Projects
	MakeProject(id []byte) Project
	Project(projectID []byte) (Project, error)
	Projects(filters map[string]interface{}) ([]Project, error)

	// Resets the database (warning, this is a destructive action...)
	ResetDB() error

	// Parameter store
	ParameterStore() ParameterStore

	// Run all hooks of the given name
	RunHooks(name string, data interface{}) (interface{}, error)
}

/* Base Functionality */

type BaseController struct {
	definitions    *Definitions
	parameterStore ParameterStore
	vars           map[string]interface{}
	settings       Settings
}

func MakeBaseController(settings Settings, definitions *Definitions) (BaseController, error) {
	parameterStore, err := MakeParameterStore(settings, definitions)

	if err != nil {
		return BaseController{}, err
	}

	return BaseController{
		parameterStore: parameterStore,
		definitions:    definitions,
		settings:       settings,
		vars:           map[string]interface{}{},
	}, nil
}

func (b *BaseController) RunHooks(name string, data interface{}) (interface{}, error) {
	var err error

	hooks := b.definitions.HookDefinitions[name]
	currentData := data

	for _, hook := range hooks {
		if currentData, err = hook.Hook(data); err != nil {
			return nil, err
		}
	}
	return currentData, nil
}

func (b *BaseController) ParameterStore() ParameterStore {
	return b.parameterStore
}

func (b *BaseController) SetVar(key string, value interface{}) error {
	b.vars[key] = value
	return nil
}

func (b *BaseController) GetVar(key string) (interface{}, bool) {
	value, ok := b.vars[key]
	return value, ok
}

func (b *BaseController) Definitions() *Definitions {
	return b.definitions
}

func (b *BaseController) InitializePlugin(plugin Plugin) error {
	return plugin.Initialize(b.definitions)
}

func (b *BaseController) Settings() Settings {
	return b.settings
}

func (b *BaseController) InitializePlugins() error {
	pluginsSetting, err := b.settings.Get("plugins")

	if err == nil {
		pluginsList, ok := pluginsSetting.([]interface{})
		if ok {
			for _, pluginName := range pluginsList {
				pluginNameStr, ok := pluginName.(string)
				if !ok {
					return fmt.Errorf("expected a string")
				}
				if definition, ok := b.definitions.PluginDefinitions[pluginNameStr]; ok {
					plugin, err := definition.Maker(nil)
					if err != nil {
						return err
					}
					if err := b.InitializePlugin(plugin); err != nil {
						return err
					} else {
						Log.Infof("Successfully initialized plugin '%s'", pluginName)
					}
				} else {
					return fmt.Errorf("plugin '%s' is not registered", pluginNameStr)
				}
			}
		}
	}
	return nil
}
