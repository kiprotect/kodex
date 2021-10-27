// Kodex (Enterprise Edition - EE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2021  KIProtect GmbH (HRB 208395B) - All Rights Reserved

package api

import (
	"github.com/gin-gonic/gin"
	"github.com/kiprotect/go-helpers/errors"
	"github.com/kiprotect/kodex"
)

func FormatError(err error) interface{} {
	if chainableErr, ok := err.(errors.ChainableError); ok {
		return errors.MakeStructuredErrorWithTraceback(chainableErr, errors.ExternalError)
	}
	return map[string]interface{}{"message": err.Error()}
}

func HandleError(c *gin.Context, code int, err error) {
	unexpectedError := map[string]interface{}{
		"message": "an unexpected error occurred, please contact tech support",
	}
	if code >= 400 && code < 500 {
		c.JSON(code, FormatError(err))
	} else {
		// we log all errors that are not 400s
		kodex.Log.Error(err)
		// we do not return other error messages to clients
		// (as they could contain sensitive information)
		c.JSON(code, unexpectedError)
	}
	c.Abort()
}
