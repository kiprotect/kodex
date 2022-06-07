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

func Organization(c *gin.Context) api.Organization {
	orgObj, ok := c.Get("org")

	if !ok {
		api.HandleError(c, 500, fmt.Errorf("no organization defined in context"))
		return nil
	}

	org, ok := orgObj.(api.Organization)

	if !ok {
		api.HandleError(c, 500, fmt.Errorf("no organization defined in context"))
		return nil
	}

	return org

}
