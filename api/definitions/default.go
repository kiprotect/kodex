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

package definitions

import (
	"github.com/kiprotect/kodex/api"
	"github.com/kiprotect/kodex/api/controllers"
	"github.com/kiprotect/kodex/api/user_providers"
	"github.com/kiprotect/kodex/api/v1"
	"github.com/kiprotect/kodex/definitions"
)

var DefaultDefinitions = api.Definitions{
	Definitions: definitions.DefaultDefinitions,
	APIControllerDefinitions: map[string]api.APIControllerMaker{
		"inMemory": controllers.MakeInMemoryController,
	},
	Routes: []api.Routes{v1.Initialize},
	ObjectAdaptors: map[string]api.ObjectAdaptor{
		"stream":      StreamAdaptor{},
		"config":      ConfigAdaptor{},
		"source":      SourceAdaptor{},
		"destination": DestinationAdaptor{},
		"action":      ActionConfigAdaptor{},
		"project":     ProjectAdaptor{},
	},
	AssociateAdaptors: map[string]api.AssociateAdaptor{
		"config-action":      AssociateConfigActionConfigAdaptor{},
		"stream-source":      AssociateStreamSourceAdaptor{},
		"config-destination": AssociateConfigDestinationAdaptor{},
	},
	UserProviders: providers.Definitions,
}
