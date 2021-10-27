// Kodex (Enterprise Edition - EE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2021  KIProtect GmbH (HRB 208395B) - All Rights Reserved

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
