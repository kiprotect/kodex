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

package web

import (
	"fmt"
	"github.com/gospel-sh/gospel"
	"github.com/kiprotect/kodex/api"
)

// Returns all requested web plugins
func GetPlugins(controller api.Controller) ([]WebPlugin, error) {

	plugins := make([]WebPlugin, 0)

	pluginSettings, err := controller.Settings().Get("web.plugins")

	if err == nil {
		pluginsList, ok := pluginSettings.([]interface{})
		if ok {
			for _, pluginName := range pluginsList {
				pluginNameStr, ok := pluginName.(string)
				if !ok {
					return nil, fmt.Errorf("expected a string")
				}
				if definition, ok := controller.Definitions().PluginDefinitions[pluginNameStr]; ok {

					plugin, err := definition.Maker(nil)

					if err != nil {
						return nil, err
					}

					webPluginMaker, ok := plugin.(WebPluginMaker)

					if !ok {
						return nil, fmt.Errorf("plugin '%s' is not a web plugin", pluginNameStr)
					}

					webPlugin, err := webPluginMaker.InitializeWebPlugin(controller)

					if err != nil {
						return nil, fmt.Errorf("cannot make plugin '%s': %v", pluginNameStr, err)
					}

					plugins = append(plugins, webPlugin)

				} else {
					return nil, fmt.Errorf("plugin '%s' not found", pluginNameStr)
				}
			}
		}
	}
	return plugins, nil
}

// Returns the app server for the Kodex UI
func AppServer(controller api.Controller) (*gospel.Server, error) {

	plugins, err := GetPlugins(controller)

	if err != nil {
		return nil, fmt.Errorf("cannot get plugins: %v", err)
	}

	root, err := Root(controller, plugins)

	if err != nil {
		return nil, fmt.Errorf("cannot get root: %v", err)
	}

	fs := StaticFiles

	for _, plugin := range plugins {
		if fsPlugin, ok := plugin.(StaticFilesPlugin); ok {
			fs = append(fs, fsPlugin.StaticFiles())
		}
	}

	return gospel.MakeServer(&gospel.App{
		Root:         root,
		StaticFiles:  fs,
		StaticPrefix: "/static",
	}), nil
}
