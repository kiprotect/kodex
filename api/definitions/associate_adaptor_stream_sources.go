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

package definitions

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kiprotect/kodex"
	"github.com/kiprotect/kodex/api"
	"github.com/kiprotect/kodex/api/helpers"
)

type AssociateStreamSourceAdaptor struct{}

func (a AssociateStreamSourceAdaptor) LeftType() string {
	return "stream"
}

func (a AssociateStreamSourceAdaptor) RightType() string {
	return "source"
}

func (a AssociateStreamSourceAdaptor) Associate(c *gin.Context, left, right kodex.Model) bool {

	stream, ok := left.(kodex.Stream)

	if !ok {
		api.HandleError(c, 500, fmt.Errorf("stream missing"))
		return false
	}

	source, ok := right.(kodex.Source)

	if !ok {
		api.HandleError(c, 500, fmt.Errorf("source missing"))
		return false
	}

	data := helpers.JSONData(c)

	if data == nil {
		return false
	}

	params, err := AddStreamSourceForm.Validate(data)

	if err != nil {
		api.HandleError(c, 400, err)
		return false
	}

	if err := stream.AddSource(source, kodex.SourceStatus(params["status"].(string))); err != nil {
		api.HandleError(c, 500, err)
		return false
	}

	c.JSON(200, map[string]interface{}{"message": "success"})

	return true

}

func (a AssociateStreamSourceAdaptor) Dissociate(c *gin.Context, left, right kodex.Model) bool {

	stream, ok := left.(kodex.Stream)

	if !ok {
		api.HandleError(c, 500, fmt.Errorf("stream missing"))
		return false
	}

	source, ok := right.(kodex.Source)

	if !ok {
		api.HandleError(c, 500, fmt.Errorf("source missing"))
		return false
	}

	if err := stream.RemoveSource(source); err != nil {
		api.HandleError(c, 500, err)
		return false
	}

	return true
}

func (a AssociateStreamSourceAdaptor) Get(c *gin.Context, left kodex.Model) interface{} {

	stream, ok := left.(kodex.Stream)

	if !ok {
		api.HandleError(c, 500, fmt.Errorf("stream missing"))
		return nil
	}

	if sources, err := stream.Sources(); err != nil {
		api.HandleError(c, 500, err)
		return nil
	} else {
		return sources
	}
}
