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

package decorators

import (
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kiprotect/kodex"
	"github.com/kiprotect/kodex/api"
	"strings"
	"time"
)

var tws = []kodex.TimeWindowFunc{
	kodex.Minute,
	kodex.Hour,
	kodex.Day,
	kodex.Week,
	kodex.Month,
}

func MeterEndpointCalls(meter kodex.Meter, meterId string) func(*gin.Context) {

	return func(c *gin.Context) {

		var id string

		if meterId != "global" {
			idObj, ok := c.Get(meterId)

			if !ok {
				api.HandleError(c, 500, fmt.Errorf("meter ID is undefined"))
				return
			}

			id, ok = idObj.(string)
			if !ok {
				api.HandleError(c, 500, fmt.Errorf("invalid meter ID"))
				return
			}
		} else {
			id = "global"
		}

		path := c.Request.URL.Path

		pathComponents := strings.Split(path, "/")

		now := time.Now().UTC().UnixNano()
		for _, twt := range tws {

			// we submit statistics for the full path
			tw := twt(now)
			if err := meter.Add(id, "endpoints", map[string]string{"path": path}, tw, 1); err != nil {
				continue
			}

			// we submit statistics for partial paths as well
			partialPaths := []string{}
			for i, pathComponent := range pathComponents {
				if i == len(pathComponents)-1 {
					break
				}
				partialPaths = append(partialPaths, pathComponent)
				if err := meter.Add(id, "endpoints", map[string]string{"path": strings.Join(append(partialPaths, "*"), "/")}, tw, 1); err != nil {
					continue
				}
			}
		}

	}

}

func OrganizationMeterId(settings kodex.Settings) gin.HandlerFunc {
	decorator := func(c *gin.Context) {

		up, ok := c.Get("user")
		if !ok {
			c.Set("organizationMeterId", "org:anonymous")
			return
		}

		user, ok := up.(*api.User)

		if !ok {
			api.HandleError(c, 500, fmt.Errorf("cannot get user"))
			return
		}

		var id string

		if len(user.Roles()) == 0 {
			api.HandleError(c, 400, fmt.Errorf("you need to be associated with an organization to use this endpoint"))
			return
		}

		orgId := hex.EncodeToString(user.Roles()[0].Organization().ID())

		// to do: select a given organization based on the access token
		id = "org:" + orgId

		c.Set("organizationMeterId", id)

	}
	return decorator

}

func Metered(settings kodex.Settings, meter kodex.Meter) gin.HandlerFunc {

	testDecorator := func(c *gin.Context) {}

	decorator := func(c *gin.Context) {

		if disabled, ok := settings.Bool("meter.disable"); ok && disabled {
			kodex.Log.Info("Metering is disabled...")
			return
		}

		c.Set("meter", meter)

	}

	if test, _ := settings.Bool("test"); test {
		return testDecorator
	}
	return decorator
}
