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

package controller

import (
	"fmt"
	"github.com/kiprotect/go-helpers/maps"
	"github.com/kiprotect/kodex"
	"github.com/kiprotect/kodex/api"
	"github.com/kiprotect/kodex/helpers"
)

func InMemoryController(settings kodex.Settings, config map[string]interface{}, definitions *api.Definitions) (api.Controller, error) {
	kodexController, err := helpers.InMemoryController(settings, &definitions.Definitions, config)
	if err != nil {
		return nil, err
	}
	return ControllerType("inMemory", config, kodexController, definitions)
}

func RegisterPlugins(controller api.Controller) error {
	pluginSettings, err := controller.Settings().Get("api.plugins")

	if err == nil {
		pluginsList, ok := pluginSettings.([]interface{})
		if ok {
			for _, pluginName := range pluginsList {
				pluginNameStr, ok := pluginName.(string)
				if !ok {
					return fmt.Errorf("expected a string")
				}
				if definition, ok := controller.Definitions().PluginDefinitions[pluginNameStr]; ok {
					plugin, err := definition.Maker(nil)
					if err != nil {
						return err
					}
					apiPlugin, ok := plugin.(api.APIPlugin)
					if ok {
						if err := controller.RegisterAPIPlugin(apiPlugin); err != nil {
							return err
						} else {
							kodex.Log.Infof("Successfully registered plugin '%s'", pluginName)
						}
					}
				} else {
					kodex.Log.Errorf("plugin '%s' not found", pluginName)
				}
			}
		}
	}
	return nil
}

func ApiController(kodexController kodex.Controller, definitions *api.Definitions) (api.Controller, error) {

	apiControllerType, ok := kodexController.Settings().String("controller.type")

	if !ok {
		return nil, fmt.Errorf("No controller type given (controller.type)!")
	}

	config, err := kodexController.Settings().Get("controller")

	if err != nil {
		return nil, err
	}

	strMapConfig, ok := maps.ToStringMap(config)

	if !ok {
		return nil, fmt.Errorf("Invalid config")
	}

	if apiController, err := ControllerType(apiControllerType, strMapConfig, kodexController, definitions); err != nil {
		return nil, err
	} else if err := RegisterPlugins(apiController); err != nil {
		return nil, err
	} else {
		return apiController, nil
	}

}

func Controller(settings kodex.Settings, definitions *api.Definitions) (api.Controller, error) {

	kodexController, err := helpers.Controller(settings, &definitions.Definitions)

	if err != nil {
		return nil, err
	}

	return ApiController(kodexController, definitions)

}

func ControllerType(apiControllerType string, config map[string]interface{}, kodexController kodex.Controller, definitions *api.Definitions) (api.Controller, error) {

	maker, ok := definitions.APIControllerDefinitions[apiControllerType]

	if !ok {
		return nil, fmt.Errorf("unknown API controller type: %s", apiControllerType)
	}

	return maker(config, kodexController, definitions)
}
