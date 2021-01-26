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

package helpers

import (
	"encoding/json"
	"fmt"
	"github.com/kiprotect/go-helpers/settings"
	kipStrings "github.com/kiprotect/go-helpers/strings"
	"github.com/kiprotect/go-helpers/yaml"
	"github.com/kiprotect/kodex"
	kipHelpers "github.com/kiprotect/kodex/helpers"
	"github.com/kiprotect/kodex/processing"
	"github.com/urfave/cli"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type decorator func(f func(c *cli.Context) error) func(c *cli.Context) error

func decorate(commands []cli.Command, decorator decorator) []cli.Command {
	newCommands := make([]cli.Command, len(commands))
	for i, command := range commands {
		if command.Action != nil {
			command.Action = decorator(command.Action.(func(c *cli.Context) error))
		}
		if command.Subcommands != nil {
			command.Subcommands = decorate(command.Subcommands, decorator)
		}
		newCommands[i] = command
	}
	return newCommands
}

type ParametersStruct struct {
	Parameters    []map[string]interface{} `json:"parameters"`
	ParameterSets []map[string]interface{} `json:"parameter-sets"`
}

func importParameters(controller kodex.Controller, path string) error {
	var parametersStruct ParametersStruct
	parameterStore := controller.ParameterStore()
	definitions := controller.Definitions()
	if bytes, err := ioutil.ReadFile(path); err != nil {
		return err
	} else if err := json.Unmarshal(bytes, &parametersStruct); err != nil {
		return err
	} else {
		for _, parametersData := range parametersStruct.Parameters {
			if parameters, err := kodex.RestoreParameters(parametersData, parameterStore, definitions); err != nil {
				return err
			} else {
				if err := parameters.Save(); err != nil {
					kodex.Log.Error(err)
					continue
				}
			}
		}
		for _, parameterSetData := range parametersStruct.ParameterSets {
			if parameterSet, err := kodex.RestoreParameterSet(parameterSetData, parameterStore); err != nil {
				return err
			} else if err := parameterSet.Save(); err != nil {
				kodex.Log.Error(err)
				continue
			}
		}
	}
	return nil
}

func exportParameters(controller kodex.Controller, path string) error {
	parameterStore := controller.ParameterStore()
	allParameterSets, err := parameterStore.AllParameterSets()
	if err != nil {
		return err
	}
	allParameters, err := parameterStore.AllParameters()
	if err != nil {
		return err
	}
	if bytes, err := json.Marshal(map[string]interface{}{
		"parameter-sets": allParameterSets,
		"parameters":     allParameters,
	}); err != nil {
		return err
	} else {
		return ioutil.WriteFile(path, bytes, 0644)
	}
}

var defaultSettings = map[string]interface{}{
	"controller": map[string]interface{}{
		"type": "inMemory",
	},
	"parameter-store": map[string]interface{}{
		"type":     "file",
		"filename": "~/.kiprotect/parameters.kip",
	},
	"blueprints": map[string]interface{}{
		"paths": []interface{}{
			"~/.kiprotect/blueprints",
		},
	},
}

func downloadBlueprints(path, url string) error {
	if data, err := Download(url); err != nil {
		return err
	} else {
		if err := ExtractBlueprints(data, path); err != nil {
			return err
		}
	}
	return nil
}

func normalizePath(path string) (string, error) {
	if strings.HasPrefix(path, "~") {
		usr, err := user.Current()
		if err != nil {
			return "", err
		}
		return usr.HomeDir + path[1:len(path)], nil
	}
	return path, nil
}

func LoadBlueprint(settingsObj kodex.Settings, filename, version string) (*kodex.Blueprint, error) {
	if filename == "" {
		filename = ".kodex.yml"
	} else {
		if !strings.HasSuffix(filename, ".yml") {
			filename = filename + ".yml"
		}
		var err error
		if filename, err = normalizePath(filename); err != nil {
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
		return kodex.MakeBlueprint(configMap), nil
	}
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

func findBlueprint(settings kodex.Settings, name string, version string) (string, error) {
	blueprintsPaths, err := getBlueprintsPaths(settings)
	if err != nil {
		return "", err
	}
	var bestPath string
	var bestVersion string
outer:
	for _, path := range blueprintsPaths {
		var err error
		if path, err = normalizePath(path); err != nil {
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
						kodex.Log.Debugf("found blueprints directory: '%s'", filepath.Join(path, file.Name()))
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
			kodex.Log.Debugf("found blueprint '%s' at path '%s' (version: '%s')", name, bestPath, bestVersion)
			return bestPath, nil
		}
	}
	return "", fmt.Errorf("blueprint '%s' with version '%s' not found", name, version)
}

func getBlueprintsPaths(settings kodex.Settings) ([]string, error) {
	blueprintsPaths, err := settings.Get("blueprints.paths")
	if err != nil {
		return nil, err
	}
	blueprintsPathsList, err := kipStrings.ToListOfStr(blueprintsPaths)
	if err != nil {
		return nil, err
	}
	return blueprintsPathsList, nil
}

func Settings() (kodex.Settings, error) {
	settingsPaths := kipHelpers.SettingsPaths()
	if settings, err := kipHelpers.Settings(settingsPaths); err != nil {
		return nil, err
	} else {
		// if no settings were given we use the default settings above
		if len(settingsPaths) == 0 {
			settings.Update(defaultSettings)
		}
		return settings, nil
	}
}

func Kodex(definitions *kodex.Definitions) {

	var controller kodex.Controller
	var settings kodex.Settings
	var err error

	if settings, err = Settings(); err != nil {
		kodex.Log.Fatal(err)
	}

	if controller, err = kipHelpers.Controller(settings, definitions); err != nil {
		kodex.Log.Fatal(err)
	}

	if err := controller.InitializePlugins(); err != nil {
		kodex.Log.Fatal(err)
	}

	init := func(f func(c *cli.Context) error) func(c *cli.Context) error {
		return func(c *cli.Context) error {

			level := c.GlobalString("level")
			logLevel, err := kodex.ParseLevel(level)
			if err != nil {
				return err
			}
			kodex.Log.SetLevel(logLevel)

			runner := func() error { return f(c) }
			profiler := c.GlobalString("profile")
			if profiler != "" {
				return runWithProfiler(profiler, runner)
			}

			return f(c)
		}
	}

	app := cli.NewApp()
	app.Name = "KIProtect"
	app.Usage = "Run all KIProtect commands"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "level",
			Value: "info",
			Usage: "The desired log level",
		},
		cli.StringFlag{
			Name:  "profile",
			Value: "",
			Usage: "enable profiler and store results to given filename",
		},
	}

	bareCommands := []cli.Command{
		cli.Command{
			Name: "parameters",
			Subcommands: []cli.Command{
				cli.Command{
					Name: "export",
					Action: func(c *cli.Context) error {
						if c.NArg() != 1 {
							return fmt.Errorf("usage: export [filename]")
						}
						return exportParameters(controller, c.Args().Get(0))
					},
				},
				cli.Command{
					Name: "import",
					Action: func(c *cli.Context) error {
						if c.NArg() != 1 {
							return fmt.Errorf("usage: import [filename]")
						}
						return importParameters(controller, c.Args().Get(0))
					},
				},
			},
		},
		cli.Command{
			Name: "blueprints",
			Subcommands: []cli.Command{
				cli.Command{
					Name: "download",
					Action: func(c *cli.Context) error {
						blueprintsPaths, err := getBlueprintsPaths(controller.Settings())
						if err != nil {
							return err
						}
						if len(blueprintsPaths) == 0 {
							return fmt.Errorf("no blueprint paths specified")
						}
						blueprintsUrl := "https://github.com/kiprotect/blueprints/archive/master.zip"
						if c.NArg() == 1 {
							blueprintsUrl = c.Args().Get(0)
						}
						if blueprintsPath, err := normalizePath(blueprintsPaths[0]); err != nil {
							return err
						} else {
							return downloadBlueprints(blueprintsPath, blueprintsUrl)
						}
					},
				},
			},
		},
		cli.Command{
			Name: "run",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "version",
					Value: "",
					Usage: "optional: the version of the blueprint to load",
				},
			},
			Action: func(c *cli.Context) error {
				project := controller.MakeProject()
				project.SetName("default")
				if err := project.Save(); err != nil {
					return err
				}
				blueprintName := ""
				if c.NArg() > 0 {
					blueprintName = c.Args().Get(0)
				}
				blueprint, err := LoadBlueprint(controller.Settings(), blueprintName, c.String("version"))
				if err != nil {
					return err
				}
				if err := blueprint.Create(project); err != nil {
					return err
				}
				streams, err := controller.Streams(map[string]interface{}{"name": "default"})
				if err != nil {
					return err
				}
				if len(streams) != 1 {
					return fmt.Errorf("expected one stream")
				}
				stream := streams[0]
				return processing.ProcessStream(stream, 0)
			},
		},
	}

	// we add commands from the definitions
	for _, commandsDefinition := range controller.Definitions().CommandsDefinitions {
		if commands, err := commandsDefinition.Maker(controller); err != nil {
			kodex.Log.Fatal(err)
		} else {
			bareCommands = append(bareCommands, commands...)
		}
	}

	app.Commands = decorate(bareCommands, init)

	err = app.Run(os.Args)

	if err != nil {
		kodex.Log.Error(err)
	}

}
