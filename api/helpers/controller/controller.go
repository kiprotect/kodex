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

package controller

import (
	"fmt"
	"github.com/kiprotect/go-helpers/maps"
	"github.com/kiprotect/kodex"
	"github.com/kiprotect/kodex/api"
	"github.com/kiprotect/kodex/helpers"
)

func InMemoryController(settings kodex.Settings, config map[string]interface{}, definitions *api.Definitions) (api.Controller, error) {
	kiprotectController, err := helpers.InMemoryController(settings, &definitions.Definitions, config)
	if err != nil {
		return nil, err
	}
	return ControllerType("inMemory", config, kiprotectController, definitions)
}

func ApiController(kiprotectController kodex.Controller, definitions *api.Definitions) (api.Controller, error) {

	apiControllerType, ok := kiprotectController.Settings().String("controller.type")

	if !ok {
		return nil, fmt.Errorf("No controller type given (controller.type)!")
	}

	config, err := kiprotectController.Settings().Get("controller")

	if err != nil {
		return nil, err
	}

	strMapConfig, ok := maps.ToStringMap(config)

	if !ok {
		return nil, fmt.Errorf("Invalid config")
	}

	return ControllerType(apiControllerType, strMapConfig, kiprotectController, definitions)

}

func Controller(settings kodex.Settings, definitions *api.Definitions) (api.Controller, error) {

	kiprotectController, err := helpers.Controller(settings, &definitions.Definitions)

	if err != nil {
		return nil, err
	}

	return ApiController(kiprotectController, definitions)

}

func ControllerType(apiControllerType string, config map[string]interface{}, kiprotectController kodex.Controller, definitions *api.Definitions) (api.Controller, error) {

	maker, ok := definitions.APIControllerDefinitions[apiControllerType]

	if !ok {
		return nil, fmt.Errorf("unknown API controller type: %s", apiControllerType)
	}

	return maker(config, kiprotectController, definitions)
}
