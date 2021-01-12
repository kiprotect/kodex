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
	"github.com/kiprotect/go-helpers/maps"
)

type Blueprint struct {
	config map[string]interface{}
}

func initSources(project Project, config map[string]interface{}) error {
	sourcesConfig, ok := maps.ToStringMapList(config["sources"])
	if !ok {
		return nil
	}
	Log.Debug("Initializing sources...")
	for _, sourceConfig := range sourcesConfig {
		name, ok := sourceConfig["name"].(string)
		if !ok {
			return fmt.Errorf("name is missing")
		}
		sourceMapConfig, ok := maps.ToStringMap(sourceConfig)
		if !ok {
			return fmt.Errorf("invalid config for source %s", name)
		}
		Log.Debugf("Creating source: %s", name)
		source := project.MakeSource()

		if err := source.Create(sourceMapConfig); err != nil {
			return err
		}

		if err := source.Save(); err != nil {
			return err
		}
	}
	return nil

}

func initActionConfigs(project Project, config map[string]interface{}) error {
	actionsConfig, ok := maps.ToStringMapList(config["actions"])
	if !ok {
		Log.Debug("actions config does not exist")
		return nil
	}
	Log.Debug("Initializing actions...")
	for _, actionConfig := range actionsConfig {
		name, ok := actionConfig["name"].(string)
		if !ok {
			return fmt.Errorf("name is missing")
		}
		actionMapConfig, ok := maps.ToStringMap(actionConfig)
		if !ok {
			return fmt.Errorf("invalid config for action %s", name)
		}
		Log.Debugf("Creating action: %s", name)
		action := project.MakeActionConfig()

		if err := action.Create(actionMapConfig); err != nil {
			return err
		}

		if err := action.Save(); err != nil {
			return err
		}
	}
	return nil

}

func initDestinations(project Project, config map[string]interface{}) error {
	destinationsConfig, ok := maps.ToStringMapList(config["destinations"])
	if !ok {
		Log.Debug("destinations config does not exist")
		return nil
	}
	Log.Debug("Initializing destinations...")
	for _, destinationConfig := range destinationsConfig {
		name, ok := destinationConfig["name"].(string)
		if !ok {
			return fmt.Errorf("name is missing")
		}
		destinationMapConfig, ok := maps.ToStringMap(destinationConfig)
		if !ok {
			return fmt.Errorf("invalid config for destination %s", name)
		}
		Log.Debugf("Creating destination: %s", name)
		destination := project.MakeDestination()

		if err := destination.Create(destinationMapConfig); err != nil {
			return err
		}

		if err := destination.Save(); err != nil {
			return err
		}

	}
	return nil

}

func initStreams(project Project, config map[string]interface{}) error {
	streamsConfig, ok := maps.ToStringMapList(config["streams"])
	if !ok {
		return nil
	}
	Log.Debug("Initializing streams...")
	for _, streamConfig := range streamsConfig {
		name, ok := streamConfig["name"].(string)
		if !ok {
			return fmt.Errorf("name is missing")
		}
		streamConfigMap, ok := maps.ToStringMap(streamConfig)
		if !ok {
			return fmt.Errorf("stream config missing")
		}
		Log.Debugf("Creating stream: %s", name)
		stream := project.MakeStream()
		values := map[string]interface{}{
			"name":        name,
			"description": "",
			"status":      string(ActiveStream),
		}

		if err := stream.Create(values); err != nil {
			return err
		}

		if err := stream.Save(); err != nil {
			return err
		}

		if err := initStreamSources(stream, streamConfigMap); err != nil {
			return err
		}

		if err := initStreamConfigs(stream, streamConfigMap); err != nil {
			return err
		}
	}
	return nil
}

func initStreamSources(stream Stream, config map[string]interface{}) error {
	sourceConfigs, ok := maps.ToStringMapList(config["sources"])

	if !ok {
		return nil
	}

	allSources, err := stream.Project().Controller().Sources(map[string]interface{}{})

	if err != nil {
		return err
	}

	sourcesByName := map[string]Source{}

	for _, source := range allSources {
		sourcesByName[source.Name()] = source
	}

	for i, sourceConfig := range sourceConfigs {
		sourceName, ok := sourceConfig["source"].(string)
		if !ok {
			return fmt.Errorf("name missing for source %d", i)
		}
		sourceStatus, ok := sourceConfig["status"].(string)
		if !ok {
			sourceStatus = "active"
		}
		source, ok := sourcesByName[sourceName]
		if !ok {
			return fmt.Errorf("source '%s' does not exist", sourceName)
		}
		if err := stream.AddSource(source, SourceStatus(sourceStatus)); err != nil {
			return err
		}
	}
	return nil
}

func initStreamConfigs(stream Stream, config map[string]interface{}) error {

	configConfigs, ok := maps.ToStringMapList(config["configs"])

	if !ok {
		return nil
	}

	for _, configConfig := range configConfigs {
		name, ok := configConfig["name"].(string)
		if !ok {
			return fmt.Errorf("name is missing")
		}

		mapConfigConfig, ok := maps.ToStringMap(configConfig)
		if !ok {
			return fmt.Errorf("invalid config: %s", name)
		}
		config := stream.MakeConfig()

		status, ok := mapConfigConfig["status"]

		if !ok {
			status = "active"
		}

		values := map[string]interface{}{
			"name":   name,
			"status": status,
		}

		if err := config.Create(values); err != nil {
			return err
		}

		if err := config.Save(); err != nil {
			return err
		}

		if err := initConfigDestinations(config, mapConfigConfig); err != nil {
			return err
		}

		if err := initConfigActions(config, mapConfigConfig); err != nil {
			return err
		}

		Log.Debugf("Created config '%s'", name)
	}

	return nil

}

func initConfigDestinations(config Config, configData map[string]interface{}) error {

	destinationConfigs, ok := maps.ToStringMapList(configData["destinations"])

	if !ok {
		// not destinations specified, we skip initialization
		return nil
	}

	allDestinations, err := config.Stream().Project().Controller().Destinations(map[string]interface{}{})

	if err != nil {
		return err
	}

	destinationsByName := map[string]Destination{}

	for _, destination := range allDestinations {
		destinationsByName[destination.Name()] = destination
	}

	for _, destinationConfig := range destinationConfigs {

		nameStr, ok := destinationConfig["name"].(string)

		if !ok {
			return fmt.Errorf("destination name missing")
		}

		destinationNameStr, ok := destinationConfig["destination"].(string)

		if !ok {
			destinationNameStr = nameStr
		}

		status, ok := destinationConfig["status"].(string)

		if !ok {
			status = "active"
		}

		destination, ok := destinationsByName[destinationNameStr]

		if !ok {
			return fmt.Errorf("destination '%s' does not exist (%d destinations)", nameStr, len(destinationsByName))
		}

		if err := config.AddDestination(destination, nameStr, DestinationStatus(status)); err != nil {
			return err
		}
	}

	return nil
}

func initConfigActions(config Config, configData map[string]interface{}) error {

	actionConfigConfigs, ok := maps.ToStringMapList(configData["actions"])

	if !ok {
		// not actionConfigs specified, we skip initialization
		return nil
	}

	allActionConfigs, err := config.Stream().Project().Controller().ActionConfigs(map[string]interface{}{})

	if err != nil {
		return err
	}

	actionConfigsByName := map[string]ActionConfig{}

	for _, actionConfig := range allActionConfigs {
		actionConfigsByName[actionConfig.Name()] = actionConfig
	}

	for i, actionConfigConfig := range actionConfigConfigs {

		nameStr, ok := actionConfigConfig["name"].(string)

		if !ok {
			return fmt.Errorf("actionConfig name missing")
		}

		actionConfigNameStr, ok := actionConfigConfig["actionConfig"].(string)

		if !ok {
			actionConfigNameStr = nameStr
		}

		actionConfig, ok := actionConfigsByName[actionConfigNameStr]

		if !ok {
			return fmt.Errorf("action config '%s' does not exist", nameStr)
		}

		Log.Debugf("Adding action %s", nameStr)

		if err := config.AddActionConfig(actionConfig, i); err != nil {
			return err
		}
	}

	return nil

}

func initKeys(project Project, config map[string]interface{}) error {
	settings := project.Controller().Settings()
	if salt, ok := config["salt"].(string); ok {
		settings.Set("salt", salt)
	}
	if key, ok := config["key"].(string); ok {
		settings.Set("key", key)
	}
	return nil
}

func MakeBlueprint(config map[string]interface{}) *Blueprint {
	return &Blueprint{
		config: config,
	}
}

func (b *Blueprint) Create(project Project) error {
	if err := initSources(project, b.config); err != nil {
		return err
	}
	if err := initDestinations(project, b.config); err != nil {
		return err
	}
	if err := initActionConfigs(project, b.config); err != nil {
		return err
	}
	if err := initStreams(project, b.config); err != nil {
		return err
	}
	if err := initKeys(project, b.config); err != nil {
		return err
	}
	return nil
}
