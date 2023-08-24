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

package app

import (
	"github.com/kiprotect/kodex"
	"github.com/kiprotect/kodex/api"
	apiCmd "github.com/kiprotect/kodex/cmd/api"
	"github.com/kiprotect/kodex/web"
	"github.com/urfave/cli"
	"net/http"
	"strings"
)

func App(controller kodex.Controller, definitions interface{}) ([]cli.Command, error) {

	apiDefinitions := definitions.(*api.Definitions)

	return []cli.Command{
		{
			Name: "app",
			Subcommands: []cli.Command{
				{
					Name:    "run",
					Aliases: []string{"r"},
					Flags:   []cli.Flag{},
					Usage:   "Run the Kodex App.",
					Action: func(c *cli.Context) error {

						blueprintName := ""

						if c.NArg() > 0 {
							blueprintName = c.Args().Get(0)
						}

						return RunApp(controller, apiDefinitions, blueprintName)
					},
				},
			},
		},
	}, nil

}

type AppWithApiHandler struct {
	App http.Handler
	API http.Handler
}

func (a *AppWithApiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/api") {
		// we serve the API
		a.API.ServeHTTP(w, r)
	} else {
		// we serve the app
		a.App.ServeHTTP(w, r)
	}
}

func RunApp(controller kodex.Controller, definitions *api.Definitions, blueprintName string) error {
	kodex.Log.Infof("Running Kodex - App %s", kodex.Version)

	handlerMaker := func(controller api.Controller, api http.Handler) (http.Handler, error) {

		app, err := web.AppServer(controller)

		if err != nil {
			return nil, err
		}

		return &AppWithApiHandler{
			App: app,
			API: api,
		}, nil
	}

	return apiCmd.RunAPI(controller, definitions, "", 0, "/api", handlerMaker, blueprintName)
}
