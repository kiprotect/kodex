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

package fixtures

import (
	"fmt"
	"github.com/kiprotect/kiprotect"
	"github.com/kiprotect/kiprotect/definitions"
	"github.com/kiprotect/kiprotect/helpers"
)

type Controller struct {
}

func (c Controller) Setup(fixtures map[string]interface{}) (interface{}, error) {

	defs, ok := fixtures["definitions"].(kiprotect.Definitions)

	if !ok {
		defs = definitions.DefaultDefinitions
	}

	sett := fixtures["settings"]

	if sett == nil {
		return nil, fmt.Errorf("no settings found")
	}

	settingsObj, ok := sett.(kiprotect.Settings)

	if !ok {
		return nil, fmt.Errorf("not a real settings object")
	}

	if ctrl, err := helpers.Controller(settingsObj, defs); err != nil {
		return nil, err
	} else {
		if err := ctrl.ResetDB(); err != nil {
			return nil, err
		} else if err := ctrl.InitializePlugins(); err != nil {
			return nil, err
		}
		return ctrl, nil
	}

}

func (c Controller) Teardown(fixture interface{}) error {
	return nil
}
