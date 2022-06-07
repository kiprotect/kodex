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

package kodex

import (
	"github.com/kiprotect/go-helpers/forms"
)

type PluginDefinition struct {
	Name        string
	Description string
	Maker       PluginMaker `json:"-"`
	Form        forms.Form  `json:"form"`
}

type PluginMaker func(map[string]interface{}) (Plugin, error)
type PluginDefinitions map[string]PluginDefinition

type Plugin interface {
	Initialize(*Definitions) error
}
