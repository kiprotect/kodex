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

package helpers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kiprotect/kodex/api"
)

// Get the controller and user (as created by the decorators)
func Controller(c *gin.Context) api.Controller {

	controller, ok := c.Get("controller")
	if !ok {
		api.HandleError(c, 500, fmt.Errorf("no controller defined (API controller check)"))
		return nil
	}

	apiController, ok := controller.(api.Controller)

	if !ok {
		api.HandleError(c, 500, fmt.Errorf("not an API controller"))
		return nil
	}

	apiController, err := apiController.ApiClone()

	if err != nil {
		api.HandleError(c, 500, fmt.Errorf("not an API controller"))
		return nil
	}

	return apiController

}
