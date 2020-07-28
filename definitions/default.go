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

var DefaultDefinitions = kiprotect.Definitions{
	ParameterStoreDefinitions: parameters.ParameterStores,
	CommandsDefinitions:       cmd.Commands,
	PluginDefinitions:         plugins.Plugins,
	ActionDefinitions:         actions.Actions,
	WriterDefinitions:         writers.Writers,
	ReaderDefinitions:         readers.Readers,
	ControllerDefinitions:     controllers.Controllers,
}
