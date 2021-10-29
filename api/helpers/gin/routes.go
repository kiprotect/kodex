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

package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/kiprotect/kodex/api"
	"github.com/kiprotect/kodex/api/decorators"
	"github.com/kiprotect/kodex/api/user_provider"
)

func InitializeRouterGroup(engine *gin.Engine,
	controller api.Controller) (*gin.RouterGroup, error) {
	/*
	   Here we define the routes for the V1 of the API
	*/

	endpoints := engine.Group("")

	//attach settings to all handlers
	endpoints.Use(decorators.WithSettings(controller.Settings()))

	var err error
	var profileProvider provider.UserProfileProvider

	if _, ok := controller.Settings().Bool("user-profile-provider.type"); ok {
		if profileProvider, err = provider.MakeUserProfileProvider(controller.Settings()); err != nil {
			return nil, err
		}

		profileProvider.Start()

		// we add the profile provider so the WithUser decorator can use it
		endpoints.Use(decorators.WithValue("profileProvider", profileProvider))

	}

	// we add the CORS handler
	engine.NoRoute(decorators.Cors(controller.Settings(), true))

	// we add the CORS handler to all existing endpoints as well
	endpoints.Use(decorators.Cors(controller.Settings(), false))

	// we provide the API controller to all endpoints
	endpoints.Use(decorators.WithController(controller))

	return endpoints, nil
}
