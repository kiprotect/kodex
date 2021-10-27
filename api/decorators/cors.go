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

package decorators

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kiprotect/go-helpers/forms"
	"github.com/kiprotect/go-helpers/maps"
	"github.com/kiprotect/kodex"
	"regexp"
	"strings"
)

var CorsForm = forms.Form{
	ErrorMsg: "invalid data encountered in the Cors form",
	Fields: []forms.Field{
		{
			Name: "allowed-hosts",
			Validators: []forms.Validator{
				forms.IsRequired{},
				forms.IsStringList{},
			},
		},
		{
			Name: "allowed-headers",
			Validators: []forms.Validator{
				forms.IsOptional{Default: []string{}},
				forms.IsStringList{},
			},
		},
		{
			Name: "allowed-methods",
			Validators: []forms.Validator{
				forms.IsOptional{Default: []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"}},
				forms.IsStringList{},
			},
		},
		{
			Name: "disable",
			Validators: []forms.Validator{
				forms.IsOptional{Default: false},
				forms.IsBoolean{},
			},
		},
	},
}

func uniques(list []string) []string {
	us := make([]string, 0)
	found := make(map[string]bool)
	for _, s := range list {
		s = strings.ToLower(s)
		if _, ok := found[s]; ok {
			continue
		}
		us = append(us, s)
	}
	return us
}

func Cors(settings kodex.Settings, defaultRoute bool) gin.HandlerFunc {

	testDecorator := func(c *gin.Context) {}

	corsConfig, err := settings.Get("cors")

	if err != nil {
		return testDecorator
	}

	corsConfigMap, ok := maps.ToStringMap(corsConfig)

	if !ok {
		return testDecorator
	}

	corsParams, err := CorsForm.Validate(corsConfigMap)

	if err != nil {
		// to do: proper error handling
		panic(err)
	}

	disabled := corsParams["disable"].(bool)
	allowedHosts := corsParams["allowed-hosts"].([]string)
	allowedHeaders := corsParams["allowed-headers"].([]string)
	allowedMethods := corsParams["allowed-methods"].([]string)

	allowedHostPatterns := make([]*regexp.Regexp, len(allowedHosts))

	for i, allowedHost := range allowedHosts {
		if pattern, err := regexp.Compile(allowedHost); err != nil {
			panic(err)
		} else {
			allowedHostPatterns[i] = pattern
		}
	}

	decorator := func(c *gin.Context) {

		if disabled {
			return
		}

		allAllowedHeaders := strings.Join(
			uniques(append([]string{c.Request.Header.Get("Access-Control-Request-Headers")},
				allowedHeaders...)), ", ")

		origin := c.Request.Header.Get("Origin")
		found := false
		for _, pattern := range allowedHostPatterns {
			if pattern.MatchString(origin) {
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
				found = true
				break
			}
		}

		if found {
			c.Writer.Header().Set("Access-Control-Max-Age", fmt.Sprintf("%d", 60))
			c.Writer.Header().Set("Access-Control-Allow-Headers", allAllowedHeaders)
			c.Writer.Header().Set("Access-Control-Allow-Methods", strings.Join(allowedMethods, ", "))

			// for OPTIONS calls we set the status code explicitly
			if c.Request.Method == "OPTIONS" {
				c.AbortWithStatus(200)
				c.Abort()
				return
			}

		}

		if defaultRoute {
			c.JSON(404, gin.H{"message": "route not found"})
			c.Abort()
			return
		}

		c.Next()

	}

	if test, _ := settings.Bool("test"); test {
		return testDecorator
	}
	return decorator
}

func CorsFromEverywhere(settings kodex.Settings) gin.HandlerFunc {

	testDecorator := func(c *gin.Context) {}

	decorator := func(c *gin.Context) {

		origin := c.Request.Header.Get("Origin")
		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		c.Writer.Header().Set("Access-Control-Max-Age", fmt.Sprintf("%d", 60))
		c.Writer.Header().Set("Access-Control-Allow-Headers", c.Request.Header.Get("Access-Control-Request-Headers"))
		c.Writer.Header().Set("Access-Control-Allow-Methods", strings.Join([]string{"POST", "GET"}, ", "))

		// for OPTIONS calls we set the status code explicitly
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			c.Abort()
			return
		}

		c.Next()

	}

	if test, _ := settings.Bool("test"); test {
		return testDecorator
	}
	return decorator
}
