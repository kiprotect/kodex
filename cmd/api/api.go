// Kodex (Community Edition - CE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2024  KIProtect GmbH (HRB 208395B) - Germany
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
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/kiprotect/kodex"
	"github.com/kiprotect/kodex/api"
	controllerHelpers "github.com/kiprotect/kodex/api/helpers/controller"
	ginHelpers "github.com/kiprotect/kodex/api/helpers/gin"
	"github.com/urfave/cli"
	"sync"
)

func API(controller kodex.Controller, definitions interface{}) ([]cli.Command, error) {

	apiDefinitions := definitions.(*api.Definitions)

	return []cli.Command{
		{
			Name:  "api",
			Usage: "API related commands.",
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
						cli.IntFlag{
							Name:  "port",
							Usage: "The port to bind to",
						},
						cli.StringFlag{
							Name:  "host",
							Usage: "The host to bind to",
						},
					},
					Usage: "Run the Kodex API.",
					Action: func(c *cli.Context) error {

						blueprintName := ""

						if c.NArg() > 0 {
							blueprintName = c.Args().Get(0)
						}

						return RunAPI(controller, apiDefinitions, c.String("host"), c.Int("port"), "", nil, blueprintName)
					},
				},
			},
		},
	}, nil

}

func RunAPI(controller kodex.Controller, definitions *api.Definitions, host string, port int, prefix string, handlerMaker func(api.Controller, http.Handler) (http.Handler, error), blueprintName string) error {
	kodex.Log.Infof("Running Kodex - API %s", kodex.Version)

	var wg sync.WaitGroup

	apiController, err := controllerHelpers.ApiController(controller, definitions)

	if err != nil {
		return err
	}

	var ok bool

	if port == 0 {

		port, ok = controller.Settings().Int("port")

		if !ok {
			port = 8000
		}

	}

	if host == "" {

		host, ok = controller.Settings().String("host")

		if !ok {
			host = "0.0.0.0"
		}

	}

	bindAddress := fmt.Sprintf("%s:%d", host, port)

	if blueprintName != "" {

		blueprintConfig, err := kodex.LoadBlueprintConfig(apiController.Settings(), blueprintName, "")

		if err != nil {
			return fmt.Errorf("Cannot load blueprint: %v", err)
		}

		blueprint := kodex.MakeBlueprint(blueprintConfig)
		apiBlueprint := api.MakeBlueprint(blueprintConfig)

		if _, err := blueprint.Create(apiController, true); err != nil {
			return fmt.Errorf("Cannot create blueprint: %v", err)
		}

		if err := apiBlueprint.Create(apiController); err != nil {
			return fmt.Errorf("Cannot create API blueprint: %v", err)
		}

	}

	srv, _, err := ginHelpers.RunApi(apiController, bindAddress, prefix, handlerMaker, &wg)

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
