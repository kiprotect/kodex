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
	"github.com/kiprotect/kiprotect"
	"github.com/kiprotect/kiprotect/plugins"
)

func RegisterPlugins(controller kiprotect.Controller, settings kiprotect.Settings) error {
	pluginsSetting, err := settings.Get("plugins")

	if err == nil {
		pluginsList, ok := pluginsSetting.([]interface{})
		if ok {
			for _, pluginName := range pluginsList {
				pluginNameStr, ok := pluginName.(string)
				if !ok {
					return fmt.Errorf("expected a string")
				}
				if definition, ok := plugins.Plugins[pluginNameStr]; ok {
					plugin, err := definition.Maker(nil)
					if err != nil {
						return err
					}
					if err := controller.RegisterPlugin(plugin); err != nil {
						return err
					} else {
						kiprotect.Log.Infof("Successfully registered plugin '%s'", pluginName)
					}
				}
			}
		}
	}
	return nil
}
