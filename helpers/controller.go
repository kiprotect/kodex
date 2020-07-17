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
	"github.com/kiprotect/go-helpers/maps"
	"github.com/kiprotect/kiprotect"
)

func ControllerType(controllerType string, config map[string]interface{}, settingsObj kiprotect.Settings, definitions kiprotect.Definitions) (kiprotect.Controller, error) {

	controllerMaker, ok := definitions.ControllerDefinitions[controllerType]

	if !ok {
		return nil, fmt.Errorf("Unknown controller type: %s", controllerType)
	}

	return controllerMaker(config, settingsObj, definitions)

}

func InMemoryController(settingsObj kiprotect.Settings, definitions kiprotect.Definitions, config map[string]interface{}) (kiprotect.Controller, error) {
	return ControllerType("inMemory", config, settingsObj, definitions)
}

func Controller(settingsObj kiprotect.Settings, definitions kiprotect.Definitions) (kiprotect.Controller, error) {
	controllerType, ok := settingsObj.String("controller.type")

	if !ok {
		return nil, fmt.Errorf("No controller type given (controller.type)!")
	}

	config, err := settingsObj.Get("controller")

	if err != nil {
		return nil, err
	}

	strMapConfig, ok := maps.ToStringMap(config)
	if !ok {
		return nil, fmt.Errorf("Invalid controller config")
	}

	return ControllerType(controllerType, strMapConfig, settingsObj, definitions)

}
