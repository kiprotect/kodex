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

package anonymize_test

import (
	"fmt"
	"github.com/kiprotect/kodex"
	pt "github.com/kiprotect/kodex/helpers/testing"
	pf "github.com/kiprotect/kodex/helpers/testing/fixtures"
	"testing"
	"time"
)

type AggregateTest struct {
	Config map[string]interface{}
	Items  []map[string]interface{}
	Result map[string][]map[string]interface{}
}

var tests = []AggregateTest{
	AggregateTest{
		Config: map[string]interface{}{
			"destinations": []map[string]interface{}{},
			"actions": []map[string]interface{}{
				{
					"name": "uniques-last-24h",
					"type": "anonymize",
					"config": map[string]interface{}{
						"method":   "aggregate",
						"function": "count",
						"config": map[string]interface{}{
							"epsilon": 10000,
						},
						"group-by": []map[string]interface{}{
							{
								"function":        "value",
								"always-included": true,
								"config": map[string]interface{}{
									"field":   "type",
									"is-list": true,
									"index":   0,
								},
							},
						},
						"result-name":    "count",
						"channels":       []string{"counts"},
						"finalize-after": -1,
					},
				},
			},
			"streams": []map[string]interface{}{
				{
					"name": "default",
					"configs": []map[string]interface{}{
						{
							"name":   "default",
							"status": "active",
							"actions": []map[string]interface{}{
								map[string]interface{}{
									"name": "uniques-last-24h",
								},
							},
							"destinations": []map[string]interface{}{},
						},
					},
				},
			},
		},
		Items: []map[string]interface{}{
			map[string]interface{}{
				"type": []interface{}{map[string]interface{}{"device": "desktop"}},
			},
			map[string]interface{}{
				"type": []interface{}{map[string]interface{}{"device": "desktop"}},
			},
			map[string]interface{}{
				"type": []interface{}{map[string]interface{}{"device": "mobile"}},
			},
			map[string]interface{}{
				"type": []interface{}{map[string]interface{}{"device": "mobile"}},
			},
		},
		Result: map[string][]map[string]interface{}{
			"counts": []map[string]interface{}{
				map[string]interface{}{
					"count": 2,
					"group": map[string]interface{}{"device": "desktop"},
				},
				map[string]interface{}{
					"count": 2,
					"group": map[string]interface{}{"device": "mobile"},
				},
			},
		},
	},
	AggregateTest{
		Config: map[string]interface{}{
			"destinations": []map[string]interface{}{},
			"actions": []map[string]interface{}{
				{
					"name": "uniques-last-24h",
					"type": "anonymize",
					"config": map[string]interface{}{
						"method":   "aggregate",
						"function": "count",
						"config": map[string]interface{}{
							"epsilon": 10000,
						},
						"group-by": []map[string]interface{}{
							{
								"function":        "time-window",
								"always-included": true,
								"config": map[string]interface{}{
									"field":  "created-at",
									"window": "week-by-day",
									"format": "rfc3339",
								},
							},
						},
						"result-name":    "count",
						"channels":       []string{"counts"},
						"finalize-after": -1,
					},
				},
			},
			"streams": []map[string]interface{}{
				{
					"name": "default",
					"configs": []map[string]interface{}{
						{
							"name":   "default",
							"status": "active",
							"actions": []map[string]interface{}{
								map[string]interface{}{
									"name": "uniques-last-24h",
								},
							},
							"destinations": []map[string]interface{}{},
						},
					},
				},
			},
		},
		Items: []map[string]interface{}{
			map[string]interface{}{
				"created-at": "2009-07-01T10:31:44Z",
			},
			map[string]interface{}{
				"created-at": "2009-07-02T3:33:44Z",
			},
			map[string]interface{}{
				"created-at": "2009-07-03T5:33:46Z",
			},
			map[string]interface{}{
				"created-at": "2009-07-04T9:34:44Z",
			},
		},
		Result: map[string][]map[string]interface{}{
			"counts": []map[string]interface{}{
				map[string]interface{}{
					"count": 1,
					"group": map[string]interface{}{
						"from": "2009-06-25T00:00:00Z",
						"to":   "2009-07-02T00:00:00Z",
						"tw":   "week-by-day",
					},
				},
				map[string]interface{}{
					"count": 2,
					"group": map[string]interface{}{
						"from": "2009-06-26T00:00:00Z",
						"to":   "2009-07-03T00:00:00Z",
						"tw":   "week-by-day",
					},
				},
				map[string]interface{}{
					"count": 3,
					"group": map[string]interface{}{
						"from": "2009-06-27T00:00:00Z",
						"to":   "2009-07-04T00:00:00Z",
						"tw":   "week-by-day",
					},
				},
				map[string]interface{}{
					"count": 4,
					"group": map[string]interface{}{
						"from": "2009-06-28T00:00:00Z",
						"to":   "2009-07-05T00:00:00Z",
						"tw":   "week-by-day",
					},
				},
				map[string]interface{}{
					"count": 4,
					"group": map[string]interface{}{
						"from": "2009-06-29T00:00:00Z",
						"to":   "2009-07-06T00:00:00Z",
						"tw":   "week-by-day",
					},
				},
				map[string]interface{}{
					"count": 4,
					"group": map[string]interface{}{
						"from": "2009-06-30T00:00:00Z",
						"to":   "2009-07-07T00:00:00Z",
						"tw":   "week-by-day",
					},
				},
				map[string]interface{}{
					"count": 4,
					"group": map[string]interface{}{
						"from": "2009-07-01T00:00:00Z",
						"to":   "2009-07-08T00:00:00Z",
						"tw":   "week-by-day",
					},
				},
				map[string]interface{}{
					"count": 3,
					"group": map[string]interface{}{
						"from": "2009-07-02T00:00:00Z",
						"to":   "2009-07-09T00:00:00Z",
						"tw":   "week-by-day",
					},
				},
				map[string]interface{}{
					"count": 2,
					"group": map[string]interface{}{
						"from": "2009-07-03T00:00:00Z",
						"to":   "2009-07-10T00:00:00Z",
						"tw":   "week-by-day",
					},
				},
				map[string]interface{}{
					"count": 1,
					"group": map[string]interface{}{
						"from": "2009-07-04T00:00:00Z",
						"to":   "2009-07-11T00:00:00Z",
						"tw":   "week-by-day",
					},
				},
			},
		},
	},
	AggregateTest{
		Config: map[string]interface{}{
			"destinations": []map[string]interface{}{},
			"actions": []map[string]interface{}{
				{
					"name": "count-by-minute",
					"type": "anonymize",
					"config": map[string]interface{}{
						"method":   "aggregate",
						"function": "count",
						"config": map[string]interface{}{
							"epsilon": 10000,
						},
						"group-by": []map[string]interface{}{
							{
								"function":        "time-window",
								"always-included": true,
								"config": map[string]interface{}{
									"field":  "created-at",
									"window": "minute",
									"format": "rfc3339",
								},
							},
							{
								"function":        "value",
								"always-included": false,
								"config": map[string]interface{}{
									"field": "type",
								},
							},
						},
						"channels":       []string{"counts"},
						"finalize-after": -1,
					},
				},
			},
			"streams": []map[string]interface{}{
				{
					"name": "default",
					"configs": []map[string]interface{}{
						{
							"name":   "default",
							"status": "active",
							"actions": []map[string]interface{}{
								map[string]interface{}{
									"name": "count-by-minute",
								},
							},
							"destinations": []map[string]interface{}{},
						},
					},
				},
			},
		},
		Items: []map[string]interface{}{
			map[string]interface{}{
				"created-at": "2009-07-01T10:31:44Z",
				"type":       "swipe",
			},
			map[string]interface{}{
				"created-at": "2009-07-01T10:33:44Z",
				"type":       "click",
			},
			map[string]interface{}{
				"created-at": "2009-07-01T10:33:46Z",
				"type":       "click",
			},
			map[string]interface{}{
				"created-at": "2009-07-01T10:34:44Z",
				"type":       "swipe",
			},
		},
		Result: map[string][]map[string]interface{}{
			"counts": []map[string]interface{}{
				map[string]interface{}{
					"count-by-minute": 1,
					"group": map[string]interface{}{
						"from": "2009-07-01T10:31:00Z",
						"to":   "2009-07-01T10:32:00Z",
						"tw":   "minute",
					},
				},
				map[string]interface{}{
					"count-by-minute": 2,
					"group": map[string]interface{}{
						"from": "2009-07-01T10:33:00Z",
						"to":   "2009-07-01T10:34:00Z",
						"tw":   "minute",
					},
				},
				map[string]interface{}{
					"count-by-minute": 1,
					"group": map[string]interface{}{
						"from": "2009-07-01T10:34:00Z",
						"to":   "2009-07-01T10:35:00Z",
						"tw":   "minute",
					},
				},
				map[string]interface{}{
					"count-by-minute": 1,
					"group": map[string]interface{}{
						"type": "swipe",
						"from": "2009-07-01T10:31:00Z",
						"to":   "2009-07-01T10:32:00Z",
						"tw":   "minute",
					},
				},
				map[string]interface{}{
					"count-by-minute": 2,
					"group": map[string]interface{}{
						"type": "click",
						"from": "2009-07-01T10:33:00Z",
						"to":   "2009-07-01T10:34:00Z",
						"tw":   "minute",
					},
				},
				map[string]interface{}{
					"count-by-minute": 1,
					"group": map[string]interface{}{
						"type": "swipe",
						"from": "2009-07-01T10:34:00Z",
						"to":   "2009-07-01T10:35:00Z",
						"tw":   "minute",
					},
				},
			},
		},
	},
	AggregateTest{
		Config: map[string]interface{}{
			"destinations": []map[string]interface{}{},
			"actions": []map[string]interface{}{
				{
					"name": "uniques-by-hour",
					"type": "anonymize",
					"config": map[string]interface{}{
						"method":   "aggregate",
						"function": "uniques",
						"config": map[string]interface{}{
							"id":      "ios_ifa",
							"epsilon": 100.0,
						},
						"group-by": []map[string]interface{}{
							{
								"function":        "time-window",
								"always-included": true,
								"config": map[string]interface{}{
									"field":  "time",
									"window": "hour",
									"format": "unix",
								},
							},
						},
						"channels":       []string{"counts"},
						"finalize-after": -1,
					},
				},
			},
			"streams": []map[string]interface{}{
				{
					"name": "default",
					"configs": []map[string]interface{}{
						{
							"name":   "default",
							"status": "active",
							"actions": []map[string]interface{}{
								map[string]interface{}{
									"name": "uniques-by-hour",
								},
							},
							"destinations": []map[string]interface{}{},
						},
					},
				},
			},
		},
		Items: []map[string]interface{}{
			map[string]interface{}{
				"app_release":     "2.1",
				"app_version":     "16.0",
				"carrier":         "Verizon",
				"city":            "San Francisco",
				"distinct_id":     "88A15487-EDBF-4931-ABA0-6DF1AE935AC2",
				"event":           "Enable Weather 2",
				"ios_ifa":         "88A15487-EDBF-4931-ABA0-6DF1AE935AC2",
				"lib_version":     "2.3.5",
				"manufacturer":    "Apple",
				"model":           "iPhone6,1",
				"mp_country_code": "US",
				"mp_device_model": "iPhone6,1",
				"mp_lib":          "iphone",
				"os":              "iPhone OS",
				"os_version":      "7.1.1",
				"radio":           "None",
				"region":          "California",
				"screen_height":   568,
				"screen_width":    320,
				"time":            1399551437,
				"wifi":            false,
			},
			map[string]interface{}{
				"app_release":     "2.1",
				"app_version":     "16.0",
				"carrier":         "SK Telecom",
				"city":            "Beijing",
				"distinct_id":     "D1C904B5-EA96-4BB8-A192-4AC75F7A7B09",
				"event":           "Enable Weather 2",
				"ios_ifa":         "D1C904B5-EA96-4BB8-A192-4AC75F7A7B09",
				"lib_version":     "2.3.5",
				"manufacturer":    "Apple",
				"model":           "iPhone5,2",
				"mp_country_code": "KR",
				"mp_device_model": "iPhone5,2",
				"mp_lib":          "iphone",
				"os":              "iPhone OS",
				"os_version":      "7.1.1",
				"radio":           "None",
				"region":          "China",
				"screen_height":   568,
				"screen_width":    320,
				"time":            1399556128,
				"wifi":            true,
			},
			map[string]interface{}{
				"Snooze duration": 600,
				"app_release":     "2.1",
				"app_version":     "16.0",
				"carrier":         "Verizon",
				"city":            "Houston",
				"distinct_id":     "58DD9B79-1A2B-415A-A87B-706815D0EFE4",
				"event":           "Snooze",
				"ios_ifa":         "58DD9B79-1A2B-415A-A87B-706815D0EFE4",
				"lib_version":     "2.3.5",
				"manufacturer":    "Apple",
				"model":           "iPhone3,3",
				"mp_country_code": "US",
				"mp_device_model": "iPhone3,3",
				"mp_lib":          "iphone",
				"os":              "iPhone OS",
				"os_version":      "7.1.1",
				"radio":           "None",
				"region":          "Texas",
				"screen_height":   480,
				"screen_width":    320,
				"time":            1399552244,
				"wifi":            true,
			},
			map[string]interface{}{
				"app_release":     "2.1",
				"app_version":     "16.0",
				"carrier":         "KDDI",
				"city":            "Hangzhou",
				"distinct_id":     "CE7D36F0-4F97-4951-98BC-FC7CDF54FB3D",
				"event":           "Enable Weather 2",
				"ios_ifa":         "CE7D36F0-4F97-4951-98BC-FC7CDF54FB3D",
				"lib_version":     "2.3.5",
				"manufacturer":    "Apple",
				"model":           "iPhone5,2",
				"mp_country_code": "JP",
				"mp_device_model": "iPhone5,2",
				"mp_lib":          "iphone",
				"os":              "iPhone OS",
				"os_version":      "7.1",
				"radio":           "None",
				"region":          "Wenzhou",
				"screen_height":   568,
				"screen_width":    320,
				"time":            1399552015,
				"wifi":            true,
			},
			map[string]interface{}{
				"app_release":     "2.1",
				"app_version":     "16.0",
				"carrier":         "SK Telecom",
				"city":            "Shanghai",
				"distinct_id":     "790FBB4A-954F-4691-B178-C70A0D9A1E29",
				"event":           "Edit Alarm",
				"ios_ifa":         "790FBB4A-954F-4691-B178-C70A0D9A1E29",
				"lib_version":     "2.3.5",
				"manufacturer":    "Apple",
				"model":           "iPhone4,1",
				"mp_country_code": "KR",
				"mp_device_model": "iPhone4,1",
				"mp_lib":          "iphone",
				"os":              "iPhone OS",
				"os_version":      "7.1",
				"radio":           "None",
				"region":          "China",
				"screen_height":   480,
				"screen_width":    320,
				"time":            1399556075,
				"wifi":            true,
			},
			map[string]interface{}{
				"app_release":     "2.1",
				"app_version":     "16.0",
				"carrier":         "SK Telecom",
				"city":            "Shanghai",
				"distinct_id":     "790FBB4A-954F-4691-B178-C70A0D9A1E29",
				"event":           "Edit Alarm",
				"ios_ifa":         "790FBB4A-954F-4691-B178-C70A0D9A1E29",
				"lib_version":     "2.3.5",
				"manufacturer":    "Apple",
				"model":           "iPhone4,1",
				"mp_country_code": "KR",
				"mp_device_model": "iPhone4,1",
				"mp_lib":          "iphone",
				"os":              "iPhone OS",
				"os_version":      "7.1",
				"radio":           "None",
				"region":          "China",
				"screen_height":   480,
				"screen_width":    320,
				"time":            1399556072,
				"wifi":            true,
			},
			map[string]interface{}{
				"app_release":     "2.1",
				"app_version":     "16.0",
				"carrier":         "Vodafone IT",
				"city":            "Milan",
				"distinct_id":     "4AED366C-8329-4668-8FBD-936599663CE0",
				"event":           "Enable Weather 2",
				"ios_ifa":         "4AED366C-8329-4668-8FBD-936599663CE0",
				"lib_version":     "2.3.5",
				"manufacturer":    "Apple",
				"model":           "iPhone5,2",
				"mp_country_code": "IT",
				"mp_device_model": "iPhone5,2",
				"mp_lib":          "iphone",
				"os":              "iPhone OS",
				"os_version":      "7.1.1",
				"radio":           "None",
				"region":          "Italy",
				"screen_height":   568,
				"screen_width":    320,
				"time":            1399552012,
				"wifi":            true,
			},
			map[string]interface{}{
				"app_release":     "2.1",
				"app_version":     "16.0",
				"carrier":         "Vodafone IT",
				"city":            "Milan",
				"distinct_id":     "4AED366C-8329-4668-8FBD-936599663CE0",
				"event":           "Enable Weather 2",
				"ios_ifa":         "4AED366C-8329-4668-8FBD-936599663CE0",
				"lib_version":     "2.3.5",
				"manufacturer":    "Apple",
				"model":           "iPhone5,2",
				"mp_country_code": "IT",
				"mp_device_model": "iPhone5,2",
				"mp_lib":          "iphone",
				"os":              "iPhone OS",
				"os_version":      "7.1.1",
				"radio":           "None",
				"region":          "Italy",
				"screen_height":   568,
				"screen_width":    320,
				"time":            1399551987,
				"wifi":            true,
			},
			map[string]interface{}{
				"app_release":     "2.1",
				"app_version":     "16.0",
				"carrier":         "SK Telecom",
				"city":            "Shanghai",
				"distinct_id":     "790FBB4A-954F-4691-B178-C70A0D9A1E29",
				"event":           "Edit Alarm",
				"ios_ifa":         "790FBB4A-954F-4691-B178-C70A0D9A1E29",
				"lib_version":     "2.3.5",
				"manufacturer":    "Apple",
				"model":           "iPhone4,1",
				"mp_country_code": "KR",
				"mp_device_model": "iPhone4,1",
				"mp_lib":          "iphone",
				"os":              "iPhone OS",
				"os_version":      "7.1",
				"radio":           "None",
				"region":          "China",
				"screen_height":   480,
				"screen_width":    320,
				"time":            1399556069,
				"wifi":            true,
			},
			map[string]interface{}{
				"app_release":     "2.1",
				"app_version":     "16.0",
				"carrier":         "SK Telecom",
				"city":            "Shanghai",
				"distinct_id":     "790FBB4A-954F-4691-B178-C70A0D9A1E29",
				"event":           "Edit Alarm",
				"ios_ifa":         "790FBB4A-954F-4691-B178-C70A0D9A1E29",
				"lib_version":     "2.3.5",
				"manufacturer":    "Apple",
				"model":           "iPhone4,1",
				"mp_country_code": "KR",
				"mp_device_model": "iPhone4,1",
				"mp_lib":          "iphone",
				"os":              "iPhone OS",
				"os_version":      "7.1",
				"radio":           "None",
				"region":          "China",
				"screen_height":   480,
				"screen_width":    320,
				"time":            1399556066,
				"wifi":            true,
			},
		},
		Result: map[string][]map[string]interface{}{
			"counts": []map[string]interface{}{
				map[string]interface{}{
					"uniques-by-hour": 4,
					"group": map[string]interface{}{
						"from": time.Unix(1399550400, 0).UTC().Format(time.RFC3339),
						"to":   time.Unix(1399554000, 0).UTC().Format(time.RFC3339),
						"tw":   "hour",
					},
				},
				map[string]interface{}{
					"uniques-by-hour": 2,
					"group": map[string]interface{}{
						"from": time.Unix(1399554000, 0).UTC().Format(time.RFC3339),
						"to":   time.Unix(1399557600, 0).UTC().Format(time.RFC3339),
						"tw":   "hour",
					},
				},
			},
		},
	},
}

func numericValue(value interface{}) (float64, bool) {
	switch v := value.(type) {
	case float64:
		return v, true
	case int64:
		return float64(v), true
	case int:
		return float64(v), true
	case uint:
		return float64(v), true
	default:
		return 0, false
	}
}

func equal(a, b map[string]interface{}) bool {
	// two maps are equal if a contains b and b contains a...
	return contains(a, b) && contains(b, a)
}

func contains(a, b map[string]interface{}) bool {
	for keyA, valueA := range a {

		// we ignore these key values
		if keyA == "group_hash" || keyA == "action_id" || keyA == "action_name" {
			continue
		}

		valueB, ok := b[keyA]
		if !ok {
			return false
		}

		valueAMap, ok := valueA.(map[string]interface{})
		if ok {
			valueBMap, ok := valueB.(map[string]interface{})
			if !ok {
				return false
			}
			if !equal(valueAMap, valueBMap) {
				return false
			}
		} else if nvA, ok := numericValue(valueA); ok {
			nvB, ok := numericValue(valueB)
			if !ok {
				return false
			}
			if nvA != nvB {
				return false
			}
		} else if valueA != valueB {
			return false
		}
	}
	return true
}

func getMatchingItem(group interface{}, items []map[string]interface{}) map[string]interface{} {
	groupMap, ok := group.(map[string]interface{})
	if !ok {
		return nil
	}
	for _, item := range items {
		itemGroup, ok := item["group"]
		if !ok {
			continue
		}
		itemGroupMap, ok := itemGroup.(map[string]interface{})
		if !ok {
			continue
		}
		if equal(itemGroupMap, groupMap) {
			return item
		}
	}
	return nil
}

// Tests aggregation in serial mode
func TestAggregate(t *testing.T) {
	testAggregate(t, false)
}

// Tests aggregation in parallel mode
func TestParallelAggregate(t *testing.T) {
	testAggregate(t, true)
}

func testAggregate(t *testing.T, parallel bool) {

	// we set the logging level to "debug"
	kodex.Log.SetLevel(kodex.DebugLogLevel)

	for testI, test := range tests {

		var fixtureConfig = []pt.FC{
			pt.FC{pf.Settings{}, "settings"},
			pt.FC{pf.Controller{}, "controller"},
			pt.FC{pf.Project{Name: "test"}, "project"},
			pt.FC{pf.Blueprint{Config: test.Config, Project: "project"}, "blueprint"},
		}

		fixtures, err := pt.SetupFixtures(fixtureConfig)
		defer pt.TeardownFixtures(fixtureConfig, fixtures)

		if err != nil {
			t.Fatal(err)
		}

		controller := fixtures["controller"].(kodex.Controller)
		streams, err := controller.Streams(map[string]interface{}{"name": "default"})

		if err != nil {
			t.Fatal(err)
		}

		if len(streams) != 1 {
			t.Fatalf("Stream not found")
		}

		stream := streams[0]

		configs, err := stream.Configs()

		if err != nil {
			t.Fatal(err)
		}

		if len(configs) == 0 {
			t.Fatal("No configs defined")
		}

		config := configs[0]

		sourceItems := make([]*kodex.Item, 0)

		for _, item := range test.Items {
			sourceItems = append(sourceItems, kodex.MakeItem(item))
		}

		c := make(chan error, len(sourceItems)-1)

		processor, err := config.Processor(false)

		if err != nil {
			t.Fatal(err)
		}

		channelWriter := kodex.MakeInMemoryChannelWriter()
		processor.SetWriter(channelWriter)

		if err := processor.Setup(); err != nil {
			t.Fatal(err)
		}

		if err := processor.Reset(); err != nil {
			t.Fatal(err)
		}

		process := func(items []*kodex.Item) error {
			processor, err := config.Processor(false)
			if err != nil {
				return err
			}
			if err := processor.Setup(); err != nil {
				return err
			}
			if _, err := processor.Process(items, nil); err != nil {
				return err
			}
			if err := processor.Teardown(); err != nil {
				return err
			}
			return nil
		}

		if parallel {
			// we always reset the group store state at the beginning of a test
			process([]*kodex.Item{sourceItems[0]})
			for i, item := range sourceItems {
				if i == len(sourceItems)-1 || i == 0 {
					continue
				} else {
					go func(item *kodex.Item) {
						c <- process([]*kodex.Item{item})
					}(item)
				}
			}
			for i := 0; i < len(sourceItems)-2; i++ {
				err := <-c
				if err != nil {
					t.Fatal(err)
				}
			}
			// we process the last item only after all others are
			// finished (as this item will flush the groups)
			process([]*kodex.Item{sourceItems[len(sourceItems)-1]})
		} else {
			if err := process(sourceItems); err != nil {
				t.Fatal(err)
			}
		}

		if _, err := processor.Finalize(); err != nil {
			t.Fatal(err)
		}

		for key, expectedDestinationItems := range test.Result {

			destinationItems, ok := channelWriter.Items[key]

			if !ok {
				t.Errorf("No items for key %s found for test %d", key, testI)
				continue
			}

			if len(destinationItems) != len(expectedDestinationItems) {
				t.Errorf("Expected %d destination items, got %d for test %d",
					len(expectedDestinationItems),
					len(destinationItems),
					testI)
				continue
			}

			for i, destinationItem := range destinationItems {
				group, ok := destinationItem.Get("group")
				if !ok {
					t.Fatalf("Group information missing")
				}
				expectedDestinationItem := getMatchingItem(group, expectedDestinationItems)
				if expectedDestinationItem == nil {
					t.Error(expectedDestinationItem, destinationItem)
					t.Fatalf("Could not find an item with a matching group for test %d and item %d", testI, i)
				}
				if !equal(destinationItem.All(), expectedDestinationItem) {
					o := ""
					for k, v := range expectedDestinationItem {
						vg, _ := destinationItem.Get(k)
						mapV, ok := v.(map[string]interface{})
						if ok {
							mapVg, ok := vg.(map[string]interface{})
							if !ok {
								o += fmt.Sprintf("%s: expected %v, got %v", k, mapV, vg)
							} else {
								for mk, mv := range mapV {
									mgv, _ := mapVg[mk]
									if mgv != mv {
										o += fmt.Sprintf("%s-%s: expected: %v, received: %v\n", k, mk, mv, mgv)
									}
								}
							}
						} else if v != vg {
							o += fmt.Sprintf("%s: expected: %v, received: %v\n", k, v, vg)
						}
					}
					t.Fatalf("Items do not match for test %d and item %d:\n%s", testI, i, o)
				}
			}

		}
	}
}
