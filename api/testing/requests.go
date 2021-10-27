// Kodex (Enterprise Edition - EE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2021  KIProtect GmbH (HRB 208395B) - All Rights Reserved

package testing

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kiprotect/kodex/api"
	"net/http"
	"net/http/httptest"
)

func Request(controller api.Controller, user api.UserProfile, reader *bytes.Reader, method, path string) (*gin.Engine, *http.Request, error) {

	withUser := func(c *gin.Context) {
		c.Set("userProfile", user)
	}

	router, err := Router(controller, withUser)

	if err != nil {
		return nil, nil, err
	}

	req, err := http.NewRequest(method, path, reader)

	if err != nil {
		return nil, nil, err
	}

	return router, req, nil

}

func Serve(request *http.Request, router *gin.Engine) (*httptest.ResponseRecorder, error) {
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, request)
	return resp, nil
}

func PostPut(controller api.Controller, user api.UserProfile, method, path string, data interface{}) (*httptest.ResponseRecorder, error) {

	jsonData, err := json.Marshal(data)

	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(jsonData)

	router, req, err := Request(controller, user, reader, method, path)

	if err != nil {
		return nil, err
	}

	req.Header.Set("content-type", "application/json")

	return Serve(req, router)
}

func GetDelete(controller api.Controller, user api.UserProfile, method, path string, params interface{}) (*httptest.ResponseRecorder, error) {

	reader := bytes.NewReader(nil)
	router, req, err := Request(controller, user, reader, method, path)

	if err != nil {
		return nil, err
	}

	if params != nil {
		data, ok := params.(map[string]interface{})

		if !ok {
			return nil, fmt.Errorf("invalid query data")
		}

		q := req.URL.Query()

		for k, v := range data {
			strV, ok := v.(string)
			if !ok {
				return nil, fmt.Errorf("not a string")
			}
			q.Add(k, strV)
		}

		// we assign the URL back to the raw query
		req.URL.RawQuery = q.Encode()

	}

	return Serve(req, router)
}

func Get(controller api.Controller, user api.UserProfile, path string, params interface{}) (*httptest.ResponseRecorder, error) {
	return GetDelete(controller, user, "GET", path, params)
}

func Del(controller api.Controller, user api.UserProfile, path string, params interface{}) (*httptest.ResponseRecorder, error) {
	return GetDelete(controller, user, "DELETE", path, params)
}

func Post(controller api.Controller, user api.UserProfile, path string, data interface{}) (*httptest.ResponseRecorder, error) {
	return PostPut(controller, user, "POST", path, data)
}

func Put(controller api.Controller, user api.UserProfile, path string, data interface{}) (*httptest.ResponseRecorder, error) {
	return PostPut(controller, user, "PATCH", path, data)
}
