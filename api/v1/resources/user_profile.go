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

package resources

import (
	"github.com/gin-gonic/gin"
	"github.com/kiprotect/kodex/api/helpers"
)

// Return the profile of the logged in user
func User(c *gin.Context) {

	controller := helpers.Controller(c)

	if controller == nil {
		return
	}

	user := helpers.User(c)

	if user == nil {
		return
	}

	c.JSON(200, map[string]interface{}{"data": user})

}
