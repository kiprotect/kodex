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

package resources

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kiprotect/kodex"
	"github.com/kiprotect/kodex/api"
	"github.com/kiprotect/kodex/api/helpers"
)

func TransformEndpoint(meter kodex.Meter) func(c *gin.Context) {
	return func(c *gin.Context) {

		controllerObj, ok := c.Get("controller")

		if !ok {
			api.HandleError(c, 500, fmt.Errorf("invalid controller"))
			return
		}

		apiController, ok := controllerObj.(api.Controller)

		if !ok {
			api.HandleError(c, 500, fmt.Errorf("invalid API controllers"))
			return
		}

		data := helpers.JSONData(c)

		if data == nil {
			return
		}

		params, err := TransformForm.ValidateWithContext(data, map[string]interface{}{"definitions": apiController.Definitions()})

		if err != nil {
			api.HandleError(c, 400, err)
			return
		}

		definitions := apiController.Definitions()

		items := params["items"].([]*kodex.Item)
		actionSpecs := params["actions"].([]kodex.ActionSpecification)
		var key, salt []byte
		var returnKey bool
		if params["key"] != nil {
			key = params["key"].([]byte)
		} else {
			if key, err = kodex.RandomBytes(32); err != nil {
				api.HandleError(c, 500, err)
				return
			}
			returnKey = true
		}
		if params["salt"] != nil {
			salt = params["salt"].([]byte)
		}
		actions, err := kodex.MakeActions(actionSpecs, definitions)

		if err != nil {
			api.HandleError(c, 500, err)
			return
		}

		parameterSet, err := kodex.MakeParameterSet(actions, nil)

		if err != nil {
			api.HandleError(c, 500, err)
			return
		}

		writer := kodex.MakeInMemoryChannelWriter()

		processor, err := kodex.MakeProcessor(parameterSet, writer, nil)

		if err != nil {
			api.HandleError(c, 500, err)
			return
		}

		processor.SetKey(key)
		processor.SetSalt(salt)

		process := func() ([]*kodex.Item, error) {
			if params["undo"].(bool) {
				return processor.Undo(items, nil)
			} else {
				return processor.Process(items, nil)
			}
		}

		if newItems, err := process(); err != nil {
			api.HandleError(c, 500, err)
			return
		} else {
			channels := make(map[string]interface{})
			channels["items"] = newItems
			for k, v := range writer.Items {
				channels[k] = v
			}
			data := map[string]interface{}{
				"data": channels,
			}
			if returnKey {
				data["key"] = key
			}
			c.JSON(200, data)
			return
		}

		meterUsage(meter, c, items)

	}

}

var tws = []kodex.TimeWindowFunc{
	kodex.Minute,
	kodex.Hour,
	kodex.Day,
	kodex.Week,
	kodex.Month,
}

func meterUsage(meter kodex.Meter, c *gin.Context, items []*kodex.Item) {
	idObj, ok := c.Get("organizationMeterId")
	if !ok {
		kodex.Log.Error("ID object is missing")
		return
	}
	id, ok := idObj.(string)
	if !ok {
		kodex.Log.Error("ID is not a string")
		return
	}

	data := map[string]string{}

	for _, twt := range tws {

		tw := twt(time.Now().UTC().UnixNano())

		values := map[string]int64{
			"source-items":       int64(len(items)),
			"source-volume":      int64(c.Request.ContentLength),
			"destination-volume": int64(c.Writer.Size()),
		}

		for key, value := range values {
			err := meter.Add(id, key, data, tw, value)
			if err != nil {
				kodex.Log.Error(fmt.Sprintf("Error storing metric %s for ID %s...", key, id))
				kodex.Log.Error(err)
			}
		}

	}
}
