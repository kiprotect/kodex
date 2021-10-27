// Kodex (Enterprise Edition - EE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2021  KIProtect GmbH (HRB 208395B) - All Rights Reserved

package definitions

import (
	"github.com/kiprotect/kodex/api"
	"github.com/kiprotect/kodex/api/controllers"
	"github.com/kiprotect/kodex/api/v1"
)

var DefaultDefinitions = api.Definitions{
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
}
