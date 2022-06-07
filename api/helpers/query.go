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

package helpers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kiprotect/go-helpers/maps"
	"github.com/kiprotect/kodex/api"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func QueryToConfig(query url.Values) map[string]interface{} {
	config := make(map[string]interface{})
	for key, value := range query {
		if len(value) > 1 {
			config[key] = value
		} else if len(value) == 0 || len(value[0]) == 0 {
			config[key] = nil
		} else {
			config[key] = value[0]
		}
	}
	return config
}

type ContentType struct {
	MediaType string
	Charset   string
	Boundary  string
}

func getContentType(header string) (ContentType, error) {
	ct := ContentType{}
	if header == "" {
		return ContentType{}, fmt.Errorf("invalid content type header")
	}
	split := strings.Split(header, ";")
	ct.MediaType = strings.ToLower(strings.Trim(split[0], " "))
	if len(split) > 1 {
		for _, s := range split[1:len(split)] {
			kv := strings.Split(s, "=")
			if len(kv) != 2 {
				return ct, fmt.Errorf("invalid key-value: %s", kv)
			}
			key := strings.ToLower(strings.Trim(kv[0], " "))
			if key == "charset" {
				ct.Charset = strings.ToLower(strings.Trim(kv[1], " "))
			} else if key == "boundary" {
				ct.Boundary = strings.Trim(kv[1], " ")
			} else {
				return ct, fmt.Errorf("invalid key: %s", key)
			}
		}
	}
	return ct, nil
}

func JSON(request *http.Request, checkContentType bool) (interface{}, error) {
	if checkContentType {
		ct, err := getContentType(request.Header.Get("content-type"))
		if err != nil {
			return nil, err
		}
		// we also accept 'text/plain' for JSON since it's sometimes useful to send
		// as JSON body as this type since it enables avoiding preflight requests...
		if ct.MediaType != "application/json" && ct.MediaType != "text/plain" {
			return nil, fmt.Errorf("invalid content-type: expected application/json or text/plain")
		}

	}
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return nil, err
	}
	var iv interface{}
	return iv, json.Unmarshal(body, &iv)
}

func PlainTextJSONData(c *gin.Context) map[string]interface{} {

	data, err := JSON(c.Request, false)
	if err != nil {
		api.HandleError(c, 400, err)
		return nil
	}

	return parseJSONData(c, data)
}

func parseJSONData(c *gin.Context, data interface{}) map[string]interface{} {
	mapData, ok := maps.ToStringMap(data)

	if !ok {
		api.HandleError(c, 400, fmt.Errorf("invalid data format"))
		return nil
	}

	return mapData

}

// Return the JSON data submitted in the request body
func JSONData(c *gin.Context) map[string]interface{} {
	data, err := JSON(c.Request, true)
	if err != nil {
		api.HandleError(c, 400, err)
		return nil
	}

	return parseJSONData(c, data)

}
