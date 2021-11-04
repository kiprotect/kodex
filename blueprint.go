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
	"github.com/kiprotect/go-helpers/forms"
	"github.com/kiprotect/go-helpers/maps"
	"github.com/kiprotect/go-helpers/settings"
	kipStrings "github.com/kiprotect/go-helpers/strings"
	"github.com/kiprotect/go-helpers/yaml"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type Blueprint struct {
	config map[string]interface{}
}

func GetBlueprintsPaths(settings Settings) ([]string, error) {
	blueprintsPaths, err := settings.Get("blueprints.paths")
	if err != nil {
		return []string{}, nil
	}
	blueprintsPathsList, err := kipStrings.ToListOfStr(blueprintsPaths)
	if err != nil {
		return nil, err
	}
	return blueprintsPathsList, nil
}

func blueprintSettings(path string) (map[string]interface{}, error) {
	var settings map[string]interface{}
	if f, err := os.OpenFile(path, os.O_RDONLY, 0700); err != nil {
		return nil, err
	} else if data, err := ioutil.ReadAll(f); err != nil {
		return nil, err
	} else if err := yaml.Unmarshal(data, &settings); err != nil {
		return nil, err
	}
	return settings, nil
}

var versionRegexp = regexp.MustCompile(`^(\d+)\.(\d+)\.(\d+)(?:(\-|\+)(.*))?$`)

func compareVersions(a, b string) (int, error) {
	matchA := versionRegexp.FindStringSubmatch(a)
	matchB := versionRegexp.FindStringSubmatch(b)
	if matchA == nil || matchB == nil {
		return 0, fmt.Errorf("not a valid semantic version")
	}
	i := 0
	for i = 0; i < 3; i++ {
		vA, err := strconv.Atoi(matchA[i+1])
		if err != nil {
			// should never happen
			return 0, err
		}
		vB, err := strconv.Atoi(matchB[i+1])
		if err != nil {
			// should never happen
			return 0, err
		}
		if vA > vB {
			return 1, nil
		} else if vB != vA {
			return -1, nil
		}
	}
	// if the extra part does not match we return an error (better safe than sorry)
	if matchA[4] != matchB[4] {
		return 0, fmt.Errorf("cannot determine highest version between '%s' and '%s', please specify", a, b)
	}
	// versions match
	return 0, nil
}

func findBlueprint(settings Settings, name string, version string) (string, error) {
	blueprintsPaths, err := GetBlueprintsPaths(settings)
	if err != nil {
		return "", err
	}
	var bestPath string
	var bestVersion string
outer:
	for _, path := range blueprintsPaths {
		var err error
		if path, err = NormalizePath(path); err != nil {
			return "", err
		}
		files, err := ioutil.ReadDir(path)
		if err != nil {
			return "", err
		}
		for _, file := range files {
			if file.IsDir() {
				subfiles, err := ioutil.ReadDir(filepath.Join(path, file.Name()))
				if err != nil {
					return "", err
				}
				for _, subfile := range subfiles {
					if subfile.Name() == ".blueprints.yml" {
						Log.Debugf("found blueprints directory: '%s'", filepath.Join(path, file.Name()))
						settingsPath := filepath.Join(path, file.Name(), subfile.Name())
						trialPath := filepath.Join(path, file.Name(), name)
						if settings, err := blueprintSettings(settingsPath); err != nil {
							return "", err
						} else {
							if versionInfo, ok := settings["version"].(string); !ok {
								return "", fmt.Errorf("version information missing in settings file '%s'", settingsPath)
							} else if version != "" && version == versionInfo {
								// this is not the version we're looking for
								bestVersion = versionInfo
								bestPath = trialPath
								break outer
							} else if version == "" {
								if bestVersion == "" {
									bestVersion = versionInfo
									bestPath = trialPath
								} else if cp, err := compareVersions(versionInfo, bestVersion); err != nil {
									return "", err
								} else if cp > 0 {
									bestVersion = versionInfo
									bestPath = trialPath
								}
							}
						}
					}
				}
			}
		}
	}
	if bestPath != "" {
		if _, err := os.Stat(bestPath); err == nil {
			Log.Debugf("found blueprint '%s' at path '%s' (version: '%s')", name, bestPath, bestVersion)
			return bestPath, nil
		}
	}
	return "", fmt.Errorf("blueprint '%s' with version '%s' not found", name, version)
}

func NormalizePath(path string) (string, error) {
	if strings.HasPrefix(path, "~") {
		usr, err := user.Current()
		if err != nil {
			return "", err
		}
		return usr.HomeDir + path[1:len(path)], nil
	}
	return path, nil
}

func LoadBlueprintConfig(settingsObj Settings, filename, version string) (map[string]interface{}, error) {
	if filename == "" {
		filename = ".yml"
	} else {
		if !strings.HasSuffix(filename, ".yml") {
			filename = filename + ".yml"
		}
		var err error
		if filename, err = NormalizePath(filename); err != nil {
			return nil, err
		}
		// we check if we can directly locate the blueprint. If not, we try to
		// find it using the blueprints directories.
		if _, err := os.Stat(filename); err != nil {
			var err error
			if filename, err = findBlueprint(settingsObj, filename, version); err != nil {
				return nil, err
			}
		}
	}
	if _, err := os.Stat(filename); err != nil {
		return nil, fmt.Errorf("blueprint '%s' not found", filename)
	}
	if config, err := settings.LoadYaml(filename); err != nil {
		return nil, err
	} else if configMap, ok := config.(map[string]interface{}); !ok {
		return nil, fmt.Errorf("expected a map")
	} else {
		if values, err := settings.ParseVars(configMap); err != nil {
			return nil, err
		} else {
			if err := settings.InsertVars(configMap, values); err != nil {
				return nil, err
			}
		}
		return configMap, nil
	}
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

var BlueprintProjectForm = forms.Form{
	Fields: []forms.Field{
		{
			Name: "id",
			Validators: []forms.Validator{
				forms.IsBytes{Encoding: "hex"},
			},
		},
		{
			Name: "name",
			Validators: []forms.Validator{
				forms.IsString{},
			},
		},
		{
			Name: "description",
			Validators: []forms.Validator{
				forms.IsOptional{Default: ""},
				forms.IsString{},
			},
		},
	},
}

func initProject(controller Controller, configData map[string]interface{}) (Project, error) {
	projectConfigData, ok := configData["project"]
	if !ok {
		project := controller.MakeProject(nil)
		project.SetName("default")
		return project, project.Save()
	}
	projectConfig, ok := maps.ToStringMap(projectConfigData)

	if params, err := BlueprintProjectForm.Validate(projectConfig); err != nil {
		return nil, err
	} else {
		project := controller.MakeProject(params["id"].([]byte))
		if err := project.Create(params); err != nil {
			return nil, err
		}
		return project, project.Save()
	}
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

func (b *Blueprint) Create(controller Controller) error {

	project, err := initProject(controller, b.config)

	if err != nil {
		return err
	}

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
