// Kodex (Enterprise Edition - EE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2021  KIProtect GmbH (HRB 208395B) - All Rights Reserved

package resources

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kiprotect/kodex"
	"github.com/kiprotect/kodex/api"
	"github.com/kiprotect/kodex/api/helpers"
)

func TransformConfigEndpoint(meter kodex.Meter) func(*gin.Context) {
	return func(c *gin.Context) {

		// generate payload from POST Body

		data := helpers.JSONData(c)

		if data == nil {
			return
		}

		params, err := TransformActionConfigForm.Validate(data)

		if err != nil {
			api.HandleError(c, 400, err)
			return
		}

		items := params["items"].([]*kodex.Item)

		configObj, ok := c.Get("config")
		if !ok {
			api.HandleError(c, 500, fmt.Errorf("invalid config"))
			return
		}

		psConfig, ok := configObj.(kodex.Config)

		if !ok {
			api.HandleError(c, 500, fmt.Errorf("invalid config"))
			return
		}

		writer := kodex.MakeInMemoryChannelWriter()

		processor, err := psConfig.Processor(false)

		if err != nil {
			api.HandleError(c, 500, err)
			return
		}

		processor.SetWriter(writer)

		if newItems, err := processor.Process(items, nil); err != nil {
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
			c.JSON(200, data)
			return
		}

		meterUsage(meter, c, items)
	}

}
