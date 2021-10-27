// Kodex (Enterprise Edition - EE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2021  KIProtect GmbH (HRB 208395B) - All Rights Reserved

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
	}

	allDefinitions := api.MergeDefinitions(api.Definitions{}, defs)
	allDefinitions.Definitions = kodex.MergeDefinitions(kodex.Definitions{}, definitions.DefaultDefinitions)

	settings, ok := fixtures["settings"].(kodex.Settings)

	if !ok {
		return nil, fmt.Errorf("No settings present")
	}

	if ctrl, err := helpers.Controller(settings, &allDefinitions.Definitions); err != nil {
		return nil, err
	} else {
		if err := ctrl.ResetDB(); err != nil {
			return nil, err
		}
		return controllerHelpers.ApiController(ctrl, &allDefinitions)
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
