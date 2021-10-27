// Kodex (Enterprise Edition - EE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2021  KIProtect GmbH (HRB 208395B) - All Rights Reserved

package helpers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kiprotect/kodex/api"
)

func UserProfile(c *gin.Context) api.UserProfile {
	userProfileObj, ok := c.Get("userProfile")

	if !ok {
		api.HandleError(c, 500, fmt.Errorf("no user profile defined in context"))
		return nil
	}

	userProfile, ok := userProfileObj.(api.UserProfile)

	if !ok {
		api.HandleError(c, 500, fmt.Errorf("no user profile defined in context"))
		return nil
	}

	return userProfile

}
