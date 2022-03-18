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

package fixtures

import (
	"fmt"
	"github.com/kiprotect/kodex"
	"github.com/kiprotect/kodex/api"
	apiDefinitions "github.com/kiprotect/kodex/api/definitions"
	controllerHelpers "github.com/kiprotect/kodex/api/helpers/controller"
	"github.com/kiprotect/kodex/definitions"
	"github.com/kiprotect/kodex/helpers"
)

type Controller struct{}

func (c Controller) Setup(fixtures map[string]interface{}) (interface{}, error) {

	defs, ok := fixtures["definitions"].(api.Definitions)

	if !ok {
		defs = apiDefinitions.DefaultDefinitions
		defs.Definitions = definitions.DefaultDefinitions
	}

	settings, ok := fixtures["settings"].(kodex.Settings)

	if !ok {
		return nil, fmt.Errorf("No settings present")
	}

	if ctrl, err := helpers.Controller(settings, &defs.Definitions); err != nil {
		return nil, err
	} else {
		if err := ctrl.ResetDB(); err != nil {
			return nil, err
		}
		return controllerHelpers.ApiController(ctrl, &defs)
	}

}

func (c Controller) Teardown(fixture interface{}) error {
	return nil
}

func GetController(fixtures map[string]interface{}) (api.Controller, error) {

	controllerObj := fixtures["controller"]

	if controllerObj == nil {
		return nil, fmt.Errorf("A controller is required")
	}

	controller, ok := controllerObj.(api.Controller)

	if !ok {
		return nil, fmt.Errorf("controller should be an API controller")
	}

	return controller, nil

}
