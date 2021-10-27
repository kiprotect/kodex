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

func TransformActionConfigEndpoint(meter kodex.Meter) func(*gin.Context) {
	return func(c *gin.Context) {

		// generate payload from POST Body

		data := helpers.JSONData(c)

		if data == nil {
			return
		}

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

		params, err := TransformActionConfigForm.Validate(data)

		if err != nil {
			api.HandleError(c, 400, err)
			return
		}

		items := params["items"].([]*kodex.Item)

		actionConfigObj, ok := c.Get("action")
		if !ok {
			api.HandleError(c, 500, fmt.Errorf("invalid action config"))
			return
		}

		psActionConfig, ok := actionConfigObj.(kodex.ActionConfig)

		if !ok {
			api.HandleError(c, 500, fmt.Errorf("invalid config"))
			return
		}

		action, err := psActionConfig.Action()

		if err != nil {
			api.HandleError(c, 500, err)
			return
		}

		parameterSet, err := kodex.MakeParameterSet([]kodex.Action{action}, apiController.ParameterStore())

		writer := kodex.MakeInMemoryChannelWriter()

		processor, err := kodex.MakeProcessor(parameterSet, writer, nil)

		if err != nil {
			api.HandleError(c, 500, err)
			return
		}

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
