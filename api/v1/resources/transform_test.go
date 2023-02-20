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

package resources_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kiprotect/go-helpers/maps"
	"github.com/kiprotect/kodex/api"
	at "github.com/kiprotect/kodex/api/testing"
	af "github.com/kiprotect/kodex/api/testing/fixtures"
	pt "github.com/kiprotect/kodex/helpers/testing"
	pf "github.com/kiprotect/kodex/helpers/testing/fixtures"
	"net/http"
	"net/http/httptest"
	"testing"
)

type TransformTestStruct struct {
	Name         string
	InputParams  map[string]string
	URLParams    map[string]string
	ResponseCode int
	Actions      []map[string]interface{}
	Items        []map[string]interface{}
	Result       map[string][]map[string]interface{}
}

var tests = []TransformTestStruct{
	TransformTestStruct{
		Name:         "Pseudonymize",
		ResponseCode: 200,
		InputParams: map[string]string{
			"key": "test",
		},
		Items: []map[string]interface{}{
			map[string]interface{}{
				"foo": "foobar",
			},
		},
		Result: map[string][]map[string]interface{}{
			"items": []map[string]interface{}{
				map[string]interface{}{
					"foo": "19Lozle/",
				},
			},
		},
		Actions: []map[string]interface{}{
			{
				"name": "pseudonymize-foo",
				"type": "pseudonymize",
				"config": map[string]interface{}{
					"key":    "foo",
					"method": "merengue",
					"config": map[string]interface{}{},
				},
			},
		},
	},
}

func equal(a, b map[string]interface{}) bool {
	for keyA, valueA := range a {
		valueB, ok := b[keyA]
		if !ok {
			fmt.Println("not a map")
			return false
		}
		valueAMap, ok := valueA.(map[string]interface{})
		if ok {
			valueBMap, ok := valueB.(map[string]interface{})
			if !ok {
				fmt.Println("B not a map")
				return false
			}
			if !equal(valueAMap, valueBMap) {
				fmt.Println("not equal maps")
				return false
			}
		} else if valueA != valueB {
			return false
		}
	}
	return true
}

func TestTransform(t *testing.T) {

	var fixturesConfig = []pt.FC{

		// we create the settings
		pt.FC{pf.Settings{}, "settings"},
		// we create an SQL controller
		pt.FC{af.Controller{}, "controller"},
		pt.FC{af.Organization{Name: "test"}, "org"},

		pt.FC{af.User{Email: "max@mustermann.de", Organization: "org", Roles: []string{"admin"}, Scopes: []string{"kiprotect:api:transform"}}, "user"},
	}

	fixtures, err := pt.SetupFixtures(fixturesConfig)

	user := fixtures["user"].(*api.ExternalUser)
	controller := fixtures["controller"].(api.Controller)

	withUser := func(c *gin.Context) {
		c.Set("user", user)
	}

	router, err := at.Router(controller, withUser)

	if err != nil {
		t.Fatal(err)
	}

	for _, test := range tests {

		data := map[string]interface{}{
			"actions": test.Actions,
			"items":   test.Items,
		}

		for key, value := range test.InputParams {
			data[key] = value
		}

		res, err := json.Marshal(data)
		if err != nil {
			t.Fatal("Invalid test data")
		}

		reader := bytes.NewReader(res)
		req, _ := http.NewRequest("POST", "/v1/transform", reader)
		req.Header.Set("content-type", "application/json")
		q := req.URL.Query()

		req.URL.RawQuery = q.Encode()
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if resp.Code != test.ResponseCode {
			t.Error(fmt.Errorf("Response code should be %d but is %d", test.ResponseCode, resp.Code))
			fmt.Println(resp.Body.String())
			continue
		}

		if test.ResponseCode != 200 {
			continue
		}

		var values map[string]interface{}
		if err = json.Unmarshal(resp.Body.Bytes(), &values); err != nil {
			t.Error("invalid JSON")
			continue
		}

		t.Log(values)

		data, ok := maps.ToStringMap(values["data"])

		if !ok {
			t.Error("no result found")
			continue
		}

		for key, results := range test.Result {

			items, ok := data[key]
			if !ok {
				t.Errorf("Items for key %s in test %s are missing", key, test.Name)
				continue
			}

			itemsList, ok := items.([]interface{})

			if !ok {
				t.Error("Items should be a list of map[string]interface{} elements")
				continue
			}

			if len(itemsList) != len(results) {
				t.Errorf("The same number of items should be returned (key: %s, test: %s)", key, test.Name)
				continue
			}

			for i, item := range itemsList {
				expectedItem := results[i]
				itemMap, ok := item.(map[string]interface{})
				if !ok {
					t.Error("Item should be of type map[string]interface{}")
					continue
				}
				// we remove auto-generated fields
				for _, key := range []string{"_kip", "action_name", "action_id", "group_hash", "params_hash"} {
					if _, ok := itemMap[key]; ok {
						delete(itemMap, key)
					}
				}
				if len(itemMap) != len(expectedItem) {
					t.Log(itemMap)
					t.Errorf("Items should have the same number of keys (plus ID)")
					continue
				}
				if !equal(itemMap, expectedItem) {
					t.Errorf("Items are not equal")
					fmt.Println("A:", itemMap)
					fmt.Println("B:", expectedItem)
				}
			}

		}

	}

}
