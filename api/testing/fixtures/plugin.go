// Kodex (Enterprise Edition - EE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2021  KIProtect GmbH (HRB 208395B) - All Rights Reserved

package fixtures

import (
	"github.com/kiprotect/kodex/api"
)

type Plugin struct {
	Plugin api.APIPlugin
}

func (o Plugin) Setup(fixtures map[string]interface{}) (interface{}, error) {

	controller, err := GetController(fixtures)

	if err != nil {
		return nil, err
	}

	if err := controller.RegisterAPIPlugin(o.Plugin); err != nil {
		return nil, err
	}

	return nil, nil
}

func (o Plugin) Teardown(fixture interface{}) error {
	return nil
}
