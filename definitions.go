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

package kodex

import (
	"encoding/json"
)

type Definitions struct {
	CommandsDefinitions
	ParameterStoreDefinitions
	PluginDefinitions
	ActionDefinitions
	WriterDefinitions
	ReaderDefinitions
	ControllerDefinitions
	HookDefinitions
}

func (d Definitions) Marshal() map[string]interface{} {
	return map[string]interface{}{
		"commands":   d.CommandsDefinitions,
		"parameters": d.ParameterStoreDefinitions,
		"plugins":    d.PluginDefinitions,
		"actions":    d.ActionDefinitions,
		"writers":    d.WriterDefinitions,
		"readers":    d.ReaderDefinitions,
		"hooks":      d.HookDefinitions,
	}
}

// We perform JSON marshalling manually to gain more flexibility...
func (d Definitions) MarshalJSON() ([]byte, error) {
	ed := d.Marshal()
	return json.Marshal(ed)
}

func MergeDefinitions(a, b Definitions) Definitions {
	c := Definitions{
		CommandsDefinitions:       CommandsDefinitions{},
		ParameterStoreDefinitions: ParameterStoreDefinitions{},
		PluginDefinitions:         PluginDefinitions{},
		ActionDefinitions:         ActionDefinitions{},
		WriterDefinitions:         WriterDefinitions{},
		ReaderDefinitions:         ReaderDefinitions{},
		ControllerDefinitions:     ControllerDefinitions{},
		HookDefinitions:           make(HookDefinitions, 0),
	}
	for _, obj := range []Definitions{a, b} {
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
		for k, v := range obj.HookDefinitions {
			c.HookDefinitions[k] = append(c.HookDefinitions[k], v...)
		}
	}
	return c
}
