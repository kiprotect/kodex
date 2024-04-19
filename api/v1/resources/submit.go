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
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kiprotect/kodex"
	"github.com/kiprotect/kodex/api"
	"github.com/kiprotect/kodex/api/helpers"
)

func Submit(c *gin.Context) {

	// generate payload from POST Body

	data := helpers.JSONData(c)

	if data == nil {
		return
	}

	params, err := TransformForm.Validate(data)

	if err != nil {
		api.HandleError(c, 400, err)
		return
	}

	controllerObj, ok := c.Get("controller")

	if !ok {
		api.HandleError(c, 500, fmt.Errorf("invalid controller"))
	}

	ctrl, ok := controllerObj.(api.Controller)

	if !ok {
		api.HandleError(c, 500, fmt.Errorf("invalid controller"))
	}

	streamObj, ok := c.Get("stream")
	if !ok {
		api.HandleError(c, 500, fmt.Errorf("invalid stream"))
		return
	}

	stream, ok := streamObj.(kodex.Stream)
	if !ok {
		api.HandleError(c, 500, fmt.Errorf("invalid stream"))
		return
	}

	channel := kodex.MakeInternalChannel()

	if err := channel.Setup(ctrl, stream); err != nil {
		api.HandleError(c, 500, err)
		return
	}

	// we get the items that should be submitted to the source
	payload := kodex.MakeBasicPayload(params["items"].([]*kodex.Item), map[string]interface{}{}, true)

	// we write the items to the internal API writer
	if err := channel.Write(payload); err != nil {
		api.HandleError(c, 500, err)
		return
	}

	if err := channel.Teardown(); err != nil {
		api.HandleError(c, 500, err)
		return
	}

	// we return a success message
	c.JSON(200, map[string]interface{}{"message": "success"})

}
