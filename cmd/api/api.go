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
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/kiprotect/kodex"
	"github.com/kiprotect/kodex/api"
	controllerHelpers "github.com/kiprotect/kodex/api/helpers/controller"
	ginHelpers "github.com/kiprotect/kodex/api/helpers/gin"
	"github.com/urfave/cli"
	"strconv"
	"sync"
)

func API(controller kodex.Controller, definitions interface{}) ([]cli.Command, error) {

	apiDefinitions := definitions.(*api.Definitions)

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

						return RunAPI(controller, apiDefinitions, "", nil, blueprintName)
					},
				},
			},
		},
	}, nil

}

func RunAPI(controller kodex.Controller, definitions *api.Definitions, prefix string, handlerMaker func(http.Handler) http.Handler, blueprintName string) error {
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

	apiController, err := controllerHelpers.ApiController(controller, definitions)

	if err != nil {
		return err
	}

	if blueprintName != "" {

		apiBlueprintConfig, err := kodex.LoadBlueprintConfig(apiController.Settings(), blueprintName, "")

		if err != nil {
			return err
		}

		apiBlueprint := api.MakeBlueprint(apiBlueprintConfig)

		if err := apiBlueprint.Create(apiController); err != nil {
			return err
		}

	}

	var addr = host + ":" + strconv.Itoa(port)
	srv, _, err := ginHelpers.RunApi(apiController, addr, prefix, handlerMaker, &wg)

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
