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

package api

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/kiprotect/kodex"
	"github.com/kiprotect/kodex/api"
	apiDefinitions "github.com/kiprotect/kodex/api/definitions"
	controllerHelpers "github.com/kiprotect/kodex/api/helpers/controller"
	ginHelpers "github.com/kiprotect/kodex/api/helpers/gin"
	"github.com/urfave/cli"
	"strconv"
	"sync"
)

func API(controller kodex.Controller) ([]cli.Command, error) {

	return []cli.Command{
		{
			Name: "api",
			Subcommands: []cli.Command{
				{
					Name:    "run",
					Aliases: []string{"r"},
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "profile",
							Value: "",
							Usage: "enable profiler and store results to given filename",
						},
					},
					Usage: "Run the KIProtect API.",
					Action: func(c *cli.Context) error {

						blueprintName := ""

						if c.NArg() > 0 {
							blueprintName = c.Args().Get(0)
						}

						apiBlueprintName := ""

						if c.NArg() > 1 {
							apiBlueprintName = c.Args().Get(1)
						}

						return runAPI(controller, blueprintName, apiBlueprintName)
					},
				},
			},
		},
	}, nil

}

func runAPI(controller kodex.Controller, blueprintName, apiBlueprintName string) error {
	kodex.Log.Info("KIProtect - API", ginHelpers.ApiVersion)

	var wg sync.WaitGroup

	port, ok := controller.Settings().Int("port")
	if !ok {
		port = 8000
	}
	host, ok := controller.Settings().String("host")
	if !ok {
		host = "0.0.0.0"
	}

	// we load the default definitions and merge them with the given definitions
	definitions := api.MergeDefinitions(api.Definitions{}, apiDefinitions.DefaultDefinitions)
	definitions.Definitions = kodex.MergeDefinitions(kodex.Definitions{}, *controller.Definitions())

	apiController, err := controllerHelpers.Controller(controller.Settings(), &definitions)

	if err != nil {
		return err
	}

	if blueprintName != "" {

		project := controller.MakeProject()
		project.SetName("default")

		if err := project.Save(); err != nil {
			return err
		}

		blueprintConfig, err := kodex.LoadBlueprintConfig(controller.Settings(), blueprintName, "")

		if err != nil {
			return err
		}

		blueprint := kodex.MakeBlueprint(blueprintConfig)

		if err := blueprint.Create(project); err != nil {
			return err
		}

		if apiBlueprintName != "" {

			blueprintConfig, err := kodex.LoadBlueprintConfig(controller.Settings(), apiBlueprintName, "")

			if err != nil {
				return err
			}

			blueprint := api.MakeBlueprint(blueprintConfig)

			if err := blueprint.Create(apiController); err != nil {
				return err
			}

		}

	}

	var addr = host + ":" + strconv.Itoa(port)
	srv, _, err := ginHelpers.RunApi(apiController, addr, &wg)

	if err != nil {
		return err
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			fmt.Println("\nCTRL-C pressed, shutting down server...")
			ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
			srv.Shutdown(ctx)
		}
	}()

	wg.Wait()

	return nil
}
