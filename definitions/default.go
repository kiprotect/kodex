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

package definitions

import (
	"github.com/kiprotect/kiprotect"
	"github.com/kiprotect/kiprotect/actions"
	"github.com/kiprotect/kiprotect/cmd"
	"github.com/kiprotect/kiprotect/controllers"
	"github.com/kiprotect/kiprotect/parameters"
	"github.com/kiprotect/kiprotect/plugins"
	"github.com/kiprotect/kiprotect/readers"
	"github.com/kiprotect/kiprotect/writers"
)

func Merge(a, b kiprotect.Definitions) kiprotect.Definitions {
	c := kiprotect.Definitions{
		CommandsDefinitions:       kiprotect.CommandsDefinitions{},
		ParameterStoreDefinitions: kiprotect.ParameterStoreDefinitions{},
		PluginDefinitions:         kiprotect.PluginDefinitions{},
		ActionDefinitions:         kiprotect.ActionDefinitions{},
		WriterDefinitions:         kiprotect.WriterDefinitions{},
		ReaderDefinitions:         kiprotect.ReaderDefinitions{},
		ControllerDefinitions:     kiprotect.ControllerDefinitions{},
	}
	for _, obj := range []kiprotect.Definitions{a, b} {
		for _, v := range obj.CommandsDefinitions {
			c.CommandsDefinitions = append(c.CommandsDefinitions, v)
		}
		for k, v := range obj.PluginDefinitions {
			c.PluginDefinitions[k] = v
		}
		for k, v := range obj.ActionDefinitions {
			c.ActionDefinitions[k] = v
		}
		for k, v := range obj.WriterDefinitions {
			c.WriterDefinitions[k] = v
		}
		for k, v := range obj.ReaderDefinitions {
			c.ReaderDefinitions[k] = v
		}
		for k, v := range obj.ControllerDefinitions {
			c.ControllerDefinitions[k] = v
		}
		for k, v := range obj.ParameterStoreDefinitions {
			c.ParameterStoreDefinitions[k] = v
		}
	}
	return c
}

var DefaultDefinitions = kiprotect.Definitions{
	ParameterStoreDefinitions: parameters.ParameterStores,
	CommandsDefinitions:       cmd.Commands,
	PluginDefinitions:         plugins.Plugins,
	ActionDefinitions:         actions.Actions,
	WriterDefinitions:         writers.Writers,
	ReaderDefinitions:         readers.Readers,
	ControllerDefinitions:     controllers.Controllers,
}
