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
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kiprotect/kodex"
	"github.com/kiprotect/kodex/api"
	"time"
)

func BaseRateLimitSetup(settings kodex.Settings, endpoints *gin.RouterGroup, meter kodex.Meter) {

	// we only perform rate limiting if a meter is defined
	if meter == nil {
		return
	}

	endpoints.Use(Metered(settings, meter))
	endpoints.Use(OrganizationMeterId(settings))
	// we meter all endpoint calls globally
	endpoints.Use(MeterEndpointCalls(meter, "global"))
	// we also meter endpoint calls for each organization
	endpoints.Use(MeterEndpointCalls(meter, "organizationMeterId"))
	endpoints.Use(RateLimited(meter, settings, "organizationMeterId", "dataVolume",
		kodex.Hour, VolumeRateLimit, 1e12))
	endpoints.Use(RateLimited(meter, settings, "organizationMeterId", "callsPerHour",
		kodex.Hour, CallRateLimit, 10*60*60))
	endpoints.Use(RateLimited(meter, settings, "organizationMeterId", "callsPerMinute",
		kodex.Minute, CallRateLimit, 100*60))
	endpoints.Use(RateLimited(meter, settings, "organizationMeterId", "callsPerSecond",
		kodex.Second, CallRateLimit, 100))

}

func getMeterId(c *gin.Context, idName string) string {

	idObj, ok := c.Get(idName)
	if !ok {
		api.HandleError(c, 404, fmt.Errorf("no meter ID defined"))
		return ""
	}
	id, ok := idObj.(string)
	if !ok {
		api.HandleError(c, 404, fmt.Errorf("meter ID is not a string"))
		return ""
	}
	return id
}

func IPMeterID(settings kodex.Settings) gin.HandlerFunc {

	decorator := func(c *gin.Context) {

		var id, key, header string
		var ok bool

		// ClientIP already takes care of X-Forwarded-For, but we probably
		// shouldn't always trust this information...
		ip := c.ClientIP()

		header, ok = settings.String("meter.ip-header")

		if !ok {
			header = "X-Real-Ip"
		}

		if headerIp := c.Request.Header.Get(header); headerIp != "" {
			ip = headerIp
		}

		if key, ok = settings.String("meter.ip-key"); !ok {
			api.HandleError(c, 500, fmt.Errorf("Hashing key not defined"))
		}

		mac := hmac.New(sha256.New, []byte(key))
		mac.Write([]byte(ip))

		id = "ip:" + hex.EncodeToString(mac.Sum(nil))

		c.Set("ipMeterId", id)

	}
	return decorator
}

type RateLimitType int

const (
	VolumeRateLimit RateLimitType = iota
	CallRateLimit
)

func RateLimited(meter kodex.Meter, settings kodex.Settings, idName, rateLimitName string, tw kodex.TimeWindowFunc, rateLimitType RateLimitType, max int64) gin.HandlerFunc {

	testDecorator := func(c *gin.Context) {}

	decorator := func(c *gin.Context) {

		ts := time.Now().UTC().UnixNano()
		c.Set("meterTimestamp", ts)

		tw := tw(ts)

		var id string
		if idName == "global" {
			id = "global"
		} else {
			id = getMeterId(c, idName)
		}

		if id == "" {
			return
		}

		metric, err := meter.Get(id, rateLimitName, nil, tw)
		if err != nil {
			api.HandleError(c, 404, err)
			return
		}

		c.Writer.Header().Set("X-Quota-Before-"+rateLimitName, fmt.Sprintf("%d", metric.Value))
		c.Writer.Header().Set("X-Quota-Maximum-"+rateLimitName, fmt.Sprintf("%d", max))
		c.Writer.Header().Set("X-Quota-From-"+rateLimitName, fmt.Sprintf("%d", tw.From))
		c.Writer.Header().Set("X-Quota-To-"+rateLimitName, fmt.Sprintf("%d", tw.To))

		if metric.Value >= max {
			api.HandleError(c, 429, fmt.Errorf("sorry, your quota is exceeded"))
			return
		}

		c.Next()

		var q int64
		switch rateLimitType {
		case VolumeRateLimit:
			q = int64(c.Request.ContentLength + int64(c.Writer.Size()))
		case CallRateLimit:
			q = 1
		}

		if err := meter.Add(id, rateLimitName, nil, tw, q); err != nil {
			api.HandleError(c, 404, err)
			return
		}

	}

	if disabled, ok := settings.Bool("meter.disable"); ok && disabled {
		return testDecorator
	}

	if test, _ := settings.Bool("test"); test {
		return testDecorator
	}
	return decorator
}
