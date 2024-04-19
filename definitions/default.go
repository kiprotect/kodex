// Kodex (Community Edition - CE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2024  KIProtect GmbH (HRB 208395B) - Germany
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
	"github.com/kiprotect/kodex"
	"github.com/kiprotect/kodex/actions"
	"github.com/kiprotect/kodex/cmd"
	"github.com/kiprotect/kodex/controllers"
	"github.com/kiprotect/kodex/parameters"
	"github.com/kiprotect/kodex/plugins"
	"github.com/kiprotect/kodex/readers"
	"github.com/kiprotect/kodex/writers"
)

var DefaultDefinitions = kodex.Definitions{
	ParameterStoreDefinitions: parameters.ParameterStores,
	CommandsDefinitions:       cmd.Commands,
	PluginDefinitions:         plugins.Plugins,
	ActionDefinitions:         actions.Actions,
	WriterDefinitions:         writers.Writers,
	ReaderDefinitions:         readers.Readers,
	ControllerDefinitions:     controllers.Controllers,
}
