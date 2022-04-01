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
	"github.com/kiprotect/kodex"
	"github.com/kiprotect/kodex/api"
	kipHelpers "github.com/kiprotect/kodex/helpers"
	"github.com/kiprotect/kodex/processing"
	"github.com/urfave/cli"
	"io/ioutil"
	"os"
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
	if bytes, err := ioutil.ReadFile(path); err != nil {
		return err
	} else if err := json.Unmarshal(bytes, &parametersStruct); err != nil {
		return err
	} else {
		for _, parametersData := range parametersStruct.Parameters {
			if parameters, err := kodex.RestoreParameters(parametersData, parameterStore); err != nil {
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

func Settings() (kodex.Settings, error) {
	if settingsPaths, fS, err := kipHelpers.SettingsPaths(); err != nil {
		return nil, err
	} else if settings, err := kipHelpers.Settings(settingsPaths, fS); err != nil {
		return nil, err
	} else {
		return settings, nil
	}
}

func Kodex(definitions *api.Definitions) {

	var controller kodex.Controller
	var settings kodex.Settings
	var err error

	if settings, err = Settings(); err != nil {
		kodex.Log.Fatal(err)
	}

	if controller, err = kipHelpers.Controller(settings, &definitions.Definitions); err != nil {
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
						blueprintsPaths, err := kodex.GetBlueprintsPaths(controller.Settings())
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
						if blueprintsPath, err := kodex.NormalizePath(blueprintsPaths[0]); err != nil {
							return err
						} else {
							return downloadBlueprints(blueprintsPath, blueprintsUrl)
						}
					},
				},
			},
		},
		cli.Command{
			Name: "export",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "version",
					Value: "",
					Usage: "optional: the version of the blueprint to load",
				},
			},
			Action: func(c *cli.Context) error {

				blueprintName := ""

				if c.NArg() > 0 {
					blueprintName = c.Args().Get(0)
				}

				blueprintConfig, err := kodex.LoadBlueprintConfig(controller.Settings(), blueprintName, c.String("version"))

				if err != nil {
					return err
				}

				bytes, err := json.MarshalIndent(blueprintConfig, "", "  ")

				if err != nil {
					return err
				}

				fmt.Println(string(bytes))

				return nil

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

				blueprintName := ""

				if c.NArg() > 0 {
					blueprintName = c.Args().Get(0)
				}

				blueprintConfig, err := kodex.LoadBlueprintConfig(controller.Settings(), blueprintName, c.String("version"))

				if err != nil {
					return err
				}

				blueprint := kodex.MakeBlueprint(blueprintConfig)

				project, err := blueprint.Create(controller)

				if err != nil {
					return err
				}

				streams, err := controller.Streams(map[string]interface{}{"name": "default", "stream_project_id_project.ext_id": project.ID()})

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
		if commands, err := commandsDefinition.Maker(controller, definitions); err != nil {
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
