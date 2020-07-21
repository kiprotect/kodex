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

package helpers

import (
	"encoding/json"
	"fmt"
	"github.com/kiprotect/go-helpers/settings"
	kipStrings "github.com/kiprotect/go-helpers/strings"
	"github.com/kiprotect/go-helpers/yaml"
	"github.com/kiprotect/kiprotect"
	"github.com/kiprotect/kiprotect/definitions"
	kipHelpers "github.com/kiprotect/kiprotect/helpers"
	"github.com/kiprotect/kiprotect/processing"
	"github.com/urfave/cli"
	"io/ioutil"
	"os"
	"path/filepath"
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

func importParameters(controller kiprotect.Controller, path string) error {
	var parametersStruct ParametersStruct
	parameterStore := controller.ParameterStore()
	definitions := controller.Definitions()
	if bytes, err := ioutil.ReadFile(path); err != nil {
		return err
	} else if err := json.Unmarshal(bytes, &parametersStruct); err != nil {
		return err
	} else {
		for _, parametersData := range parametersStruct.Parameters {
			if parameters, err := kiprotect.RestoreParameters(parametersData, parameterStore, definitions); err != nil {
				return err
			} else {
				if err := parameters.Save(); err != nil {
					kiprotect.Log.Error(err)
					continue
				}
			}
		}
		for _, parameterSetData := range parametersStruct.ParameterSets {
			if parameterSet, err := kiprotect.RestoreParameterSet(parameterSetData, parameterStore); err != nil {
				return err
			} else if err := parameterSet.Save(); err != nil {
				kiprotect.Log.Error(err)
				continue
			}
		}
	}
	return nil
}

func exportParameters(controller kiprotect.Controller, path string) error {
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

func loadBlueprint(settingsObj kiprotect.Settings, filename, version string) (*kiprotect.Blueprint, error) {
	if filename == "" {
		filename = ".kiprotect.yml"
	} else {
		if !strings.HasSuffix(filename, ".yml") {
			filename = filename + ".yml"
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
		return kiprotect.MakeBlueprint(configMap), nil
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

func findBlueprint(settings kiprotect.Settings, name string, version string) (string, error) {
	blueprintsPaths, err := getBlueprintsPaths(settings)
	if err != nil {
		return "", err
	}
	for _, path := range blueprintsPaths {
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
						settingsPath := filepath.Join(path, file.Name(), subfile.Name())
						if settings, err := blueprintSettings(settingsPath); err != nil {
							return "", err
						} else {
							if versionInfo, ok := settings["version"].(string); !ok {
								return "", fmt.Errorf("version information missing in settings file '%s'", settingsPath)
							} else if version != "" && version != versionInfo {
								// this is not the version we're looking for
								continue
							}
						}
						kiprotect.Log.Debugf("found blueprints directory: '%s'", filepath.Join(path, file.Name()))
						trialPath := filepath.Join(path, file.Name(), name)
						if _, err := os.Stat(trialPath); err == nil {
							kiprotect.Log.Debugf("found blueprint '%s' at path '%s'", name, trialPath)
							return trialPath, nil
						}
					}
				}
			}
		}
	}
	return "", fmt.Errorf("blueprint '%s' with version '%s' not found", name, version)
}

func getBlueprintsPaths(settings kiprotect.Settings) ([]string, error) {
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

func KIProtect() {

	var controller kiprotect.Controller
	var settings kiprotect.Settings
	var err error

	init := func(f func(c *cli.Context) error) func(c *cli.Context) error {
		return func(c *cli.Context) error {

			level := c.GlobalString("level")
			logLevel, err := kiprotect.ParseLevel(level)
			if err != nil {
				return err
			}
			kiprotect.Log.SetLevel(logLevel)

			kiprotect.Log.Debug("Initializing settings and controller...")

			settingsArg := c.GlobalString("settings")
			var settingsPaths []string
			if settingsArg != "" {
				settingsPaths = strings.Split(settingsArg, ":")
			} else {
				settingsPaths = kipHelpers.SettingsPaths()
			}
			if settings, err = kipHelpers.Settings(settingsPaths); err != nil {
				kiprotect.Log.Error("An error occurred when loadings the settings.")
				kiprotect.Log.Fatal(err)
			}
			// if no settings were given we use the default settings above
			if len(settingsPaths) == 0 {
				settings.Update(defaultSettings)
			}
			if controller, err = kipHelpers.Controller(settings, definitions.DefaultDefinitions); err != nil {
				kiprotect.Log.Error("An error occurred when creating the controller.")
				kiprotect.Log.Fatal(err)
			}
			kiprotect.Log.Debug("Initialization successful...")

			runner := func() error { return f(c) }

			profiler := c.String("profile")
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
			Name:  "settings",
			Value: "",
			Usage: "path(s) to load settings from (separated by ':' e.g. 'foo:bar:baz')",
		},
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
						blueprintsPaths, err := getBlueprintsPaths(settings)
						if err != nil {
							return err
						}
						if len(blueprintsPaths) == 0 {
							return fmt.Errorf("no blueprint paths specified")
						}
						blueprintsPath := blueprintsPaths[0]
						blueprintsUrl := "https://github.com/kiprotect/blueprints/archive/master.zip"
						if c.NArg() == 1 {
							blueprintsUrl = c.Args().Get(0)
						}
						return downloadBlueprints(blueprintsPath, blueprintsUrl)
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
				blueprint, err := loadBlueprint(settings, blueprintName, c.String("version"))
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

	app.Commands = decorate(bareCommands, init)

	err = app.Run(os.Args)

	if err != nil {
		kiprotect.Log.Error(err)
	}

}
