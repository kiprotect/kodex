// Kodex (Enterprise Edition - EE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2021  KIProtect GmbH (HRB 208395B) - All Rights Reserved

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

	if disable, ok := controller.Settings().Bool("user-profile-provider.disable"); !(ok && disable) {
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
