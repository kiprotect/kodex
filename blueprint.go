// Kodex (Community Edition - CE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2022  KIProtect GmbH (HRB 208395B) - Germany
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
	"bytes"
	"encoding/hex"
	"encoding/json"
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
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
)

var BlueprintDateForm = forms.Form{
	Fields: []forms.Field{
		{
			Name: "created_at",
			Validators: []forms.Validator{
				forms.IsOptional{},
				forms.IsTime{Format: "rfc3339"},
			},
		},
		{
			Name: "updated_at",
			Validators: []forms.Validator{
				forms.IsOptional{},
				forms.IsTime{Format: "rfc3339"},
			},
		},
		{
			Name: "deleted_at",
			Validators: []forms.Validator{
				forms.IsOptional{},
				forms.IsTime{Format: "rfc3339"},
			},
		},
	},
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

var BlueprintStreamForm = forms.Form{
	Fields: []forms.Field{
		{
			Name: "id",
			Validators: []forms.Validator{
				forms.IsOptional{},
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
		{
			Name: "data",
			Validators: []forms.Validator{
				forms.IsOptional{},
				forms.IsStringMap{},
			},
		},
		{
			Name: "status",
			Validators: []forms.Validator{
				forms.IsOptional{Default: "active"},
				forms.IsIn{Choices: []interface{}{"active", "inactive", "testing"}},
			},
		},
	},
}

var BlueprintConfigForm = forms.Form{
	Fields: []forms.Field{
		{
			Name: "id",
			Validators: []forms.Validator{
				forms.IsOptional{},
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
		{
			Name: "status",
			Validators: []forms.Validator{
				forms.IsOptional{Default: "active"},
				forms.IsIn{Choices: []interface{}{"active", "inactive", "testing"}},
			},
		},
		{
			Name: "data",
			Validators: []forms.Validator{
				forms.IsOptional{},
				forms.IsStringMap{},
			},
		},
	},
}

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
	filename = filepath.ToSlash(filename)
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
			Log.Error(err)
			var err error
			if filename, err = findBlueprint(settingsObj, filename, version); err != nil {
				return nil, err
			}
		}
	}
	if _, err := os.Stat(filename); err != nil {
		return nil, fmt.Errorf("blueprint '%s' not found", filename)
	}
	if absFilename, err := filepath.Abs(filename); err != nil {
		return nil, err
	} else {
		// we remove the leading '/' as that's illegal for the FS interface
		// filename = absFilename[1:]
		if runtime.GOOS == "windows" {
			filename = filepath.ToSlash(absFilename)[3:] // we remove the drive letter and first slash
		} else {
			filename = filepath.ToSlash(absFilename)[1:] // we remove the first slash
		}
	}
	if config, err := settings.LoadYaml(filename, os.DirFS("")); err != nil {
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

		var id []byte

		if strId, ok := sourceMapConfig["id"].(string); ok {
			var err error
			if id, err = hex.DecodeString(strId); err != nil {
				id = nil
			}
		}

		Log.Debugf("Creating source: %s", name)
		source := project.MakeSource(id)

		if err := source.Create(sourceMapConfig); err != nil {
			return err
		}

		if err := FixDates(source, sourceMapConfig); err != nil {
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

		var id []byte

		if strId, ok := actionMapConfig["id"].(string); ok {
			var err error
			if id, err = hex.DecodeString(strId); err != nil {
				id = nil
			}
		}

		Log.Debugf("Creating action: %s", name)

		action := project.MakeActionConfig(id)

		if err := action.Create(actionMapConfig); err != nil {
			return fmt.Errorf("error creating action '%s': %v", name, err)
		}

		if err := FixDates(action, actionMapConfig); err != nil {
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

		var id []byte

		if strId, ok := destinationMapConfig["id"].(string); ok {
			var err error
			if id, err = hex.DecodeString(strId); err != nil {
				id = nil
			}
		}

		Log.Debugf("Creating destination: %s", name)
		destination := project.MakeDestination(id)

		if err := destination.Create(destinationMapConfig); err != nil {
			return err
		}

		if err := FixDates(destination, destinationMapConfig); err != nil {
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

		if params, err := BlueprintStreamForm.Validate(streamConfigMap); err != nil {
			return err
		} else {

			var stream Stream
			var err error

			Log.Debugf("Creating stream: %s", name)

			streamID, _ := params["id"].([]byte)

			if stream, err = project.Controller().Stream(streamID); err != nil {

				if err != NotFound {
					return err
				}

				stream = project.MakeStream(streamID)

				if err := stream.Create(params); err != nil {
					return err
				}

			} else if err := stream.Update(params); err != nil {
				return err
			}

			if err := FixDates(stream, streamConfigMap); err != nil {
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

		if params, err := BlueprintConfigForm.Validate(mapConfigConfig); err != nil {
			return err
		} else {

			var config Config
			var err error

			configID, _ := params["id"].([]byte)

			if config, err = stream.Config(configID); err != nil {

				if err != NotFound {
					return err
				}

				config = stream.MakeConfig(configID)

				if err := config.Create(params); err != nil {
					return err
				}

			} else if err := config.Update(params); err != nil {
				return err
			}

			if err := FixDates(config, mapConfigConfig); err != nil {
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

type Dates struct {
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type DatesSettable interface {
	SetCreatedAt(time.Time) error
	SetUpdatedAt(time.Time) error
	SetDeletedAt(*time.Time) error
}

func FixDates(object Model, data map[string]interface{}) error {

	settable, ok := object.(DatesSettable)

	if !ok {
		// can't set dates on this object...
		return nil
	}

	dates := &Dates{}
	if params, err := BlueprintDateForm.Validate(data); err != nil {
		return err
	} else if err := BlueprintDateForm.Coerce(dates, params); err != nil {
		return err
	} else {

		if !dates.CreatedAt.IsZero() {
			settable.SetCreatedAt(dates.CreatedAt)
		} else {
			settable.SetCreatedAt(time.Now())
		}
		if !dates.UpdatedAt.IsZero() {
			settable.SetUpdatedAt(dates.UpdatedAt)
		} else {
			settable.SetUpdatedAt(time.Now())
		}
		if dates.DeletedAt != nil {
			settable.SetDeletedAt(dates.DeletedAt)
		}

		return nil
	}
}

func initProject(controller Controller, configData map[string]interface{}, createProject bool) (Project, error) {
	projectConfigData, ok := configData["project"]

	var projectConfig map[string]interface{}
	var project Project

	if !ok {
		projectConfig = map[string]interface{}{
			"id":   []byte("default"),
			"name": "default",
		}

	} else if projectConfig, ok = maps.ToStringMap(projectConfigData); !ok {
		return nil, fmt.Errorf("expected a map")
	}

	if params, err := BlueprintProjectForm.Validate(projectConfig); err != nil {
		return nil, err
	} else {
		id := params["id"].([]byte)

		// if the project already exists we delete it
		if project, err = controller.Project(id); err != nil {

			if err != NotFound {
				return nil, err
			}

			if !createProject {
				return nil, fmt.Errorf("project does not exist, will not create it...")
			}

		} else {
			// the project exists, we delete it and recreate it
			if err := project.Delete(); err != nil {
				return nil, err
			}
		}

		project = controller.MakeProject(id)

		if err := project.Create(params); err != nil {
			return nil, err
		}

		if err := FixDates(project, projectConfig); err != nil {
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

type ByID[T Model] struct {
	Entries []T
}

func (b ByID[T]) Len() int {
	return len(b.Entries)
}

func (b ByID[T]) Swap(i, j int) {
	b.Entries[i], b.Entries[j] = b.Entries[j], b.Entries[i]
}

func (b ByID[T]) Less(i, j int) bool {
	if bytes.Compare(b.Entries[i].ID(), b.Entries[j].ID()) <= 0 {
		return true
	}
	return false
}

func Sort[T Model](models []T) {

	sortedByID := &ByID[T]{
		Entries: models,
	}
	sort.Sort(sortedByID)
}

func ExportBlueprint(project Project) (map[string]interface{}, error) {
	blueprint := make(map[string]interface{})

	blueprint["project"] = project

	// we sort all dependent objects to make diffing possible...

	if actions, err := project.Controller().ActionConfigs(map[string]interface{}{"project.id": project.ID()}); err != nil {
		return nil, err
	} else {

		Sort(actions)

		blueprint["actions"] = actions
	}

	if sources, err := project.Controller().Sources(map[string]interface{}{"project.id": project.ID()}); err != nil {
		return nil, err
	} else {

		Sort(sources)

		blueprint["sources"] = sources
	}

	if destinations, err := project.Controller().Destinations(map[string]interface{}{"project.id": project.ID()}); err != nil {
		return nil, err
	} else {

		Sort(destinations)

		blueprint["destinations"] = destinations
	}

	if streams, err := project.Controller().Streams(map[string]interface{}{"project.id": project.ID()}); err != nil {
		return nil, err
	} else {

		Sort(streams)

		blueprint["streams"] = streams
	}

	if jsonData, err := json.Marshal(blueprint); err != nil {
		return nil, err
	} else {
		var mapData map[string]interface{}

		if err := json.Unmarshal(jsonData, &mapData); err != nil {
			return nil, err
		}

		return mapData, nil
	}

}

func MakeBlueprint(config map[string]interface{}) *Blueprint {
	return &Blueprint{
		config: config,
	}
}

func (b *Blueprint) Create(controller Controller, createProject bool) (Project, error) {

	succeeded := false

	if err := controller.Begin(); err != nil {
		return nil, err
	}

	defer func() {
		// we roll back the transaction, but only if it hasn't been successful
		if !succeeded {
			if err := controller.Rollback(); err != nil {
				Log.Error(err)
			}
		}
	}()

	project, err := initProject(controller, b.config, createProject)

	if err != nil {
		return nil, err
	}

	if err := b.createWithProject(controller, project); err != nil {
		return nil, err
	}

	succeeded = true

	return project, nil
}

func (b *Blueprint) CreateWithProject(controller Controller, project Project) error {

	succeeded := false

	if err := controller.Begin(); err != nil {
		return err
	}

	defer func() {
		// we roll back the transaction, but only if it hasn't been successful
		if !succeeded {
			if err := controller.Rollback(); err != nil {
				Log.Error(err)
			}
		}
	}()

	if err := b.createWithProject(controller, project); err != nil {
		return err
	}

	succeeded = true

	return nil

}

func (b *Blueprint) createWithProject(controller Controller, project Project) error {

	if err := initSources(project, b.config); err != nil {
		return fmt.Errorf("error creating sources: %v", err)
	}
	if err := initDestinations(project, b.config); err != nil {
		return fmt.Errorf("error creating destinations: %v", err)
	}
	if err := initActionConfigs(project, b.config); err != nil {
		return fmt.Errorf("error creating action configs: %v", err)
	}
	if err := initStreams(project, b.config); err != nil {
		return fmt.Errorf("error creating streams: %v", err)
	}
	if err := initKeys(project, b.config); err != nil {
		return fmt.Errorf("error creating keys: %v", err)
	}
	if err := controller.Commit(); err != nil {
		return fmt.Errorf("error committing changes: %v", err)
	}

	return nil
}
