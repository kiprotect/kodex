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
	"fmt"
	"github.com/kiprotect/go-helpers/settings"
	"github.com/kiprotect/kiprotect"
	"github.com/kiprotect/kiprotect/definitions"
	kipHelpers "github.com/kiprotect/kiprotect/helpers"
	"github.com/kiprotect/kiprotect/processing"
	"github.com/urfave/cli"
	"os"
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

func loadBlueprint(filename string) (*kiprotect.Blueprint, error) {
	if filename == "" {
		filename = ".kiprotect.yml"
	}
	if !strings.HasSuffix(filename,".yml") {
		filename = filename + ".yml"
	}
	if config, err := settings.LoadYaml(filename); err != nil {
		return nil, err
	} else {
		if values, err := settings.ParseVars(config); err != nil {
			return nil, err
		} else {
			if err := settings.InsertVars(config, values); err != nil {
				return nil, err
			}
		}
		return kiprotect.MakeBlueprint(config), nil
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
			Name: "run",
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
				blueprint, err := loadBlueprint(blueprintName)
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
