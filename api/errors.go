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
