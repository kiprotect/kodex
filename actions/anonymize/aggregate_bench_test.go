// KIProtect (Community Edition - CE) - Privacy & Security Engineering Platform
// Copyright (C) 2020  KIProtect GmbH (HRB 208395B) - Germany
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
	"github.com/kiprotect/kodex"
	pt "github.com/kiprotect/kodex/helpers/testing"
	pf "github.com/kiprotect/kodex/helpers/testing/fixtures"
	"testing"
)

type AggregateBenchmark struct {
	Config map[string]interface{}
	Items  []map[string]interface{}
}

var simpleBenchmark = AggregateBenchmark{
	Config: map[string]interface{}{
		"destinations": []map[string]interface{}{
			map[string]interface{}{
				"type":        "in-memory",
				"name":        "counts",
				"status":      "active",
				"description": "counts",
			},
		},
		"actions": map[string]interface{}{
			"uniques-last-24h": map[string]interface{}{
				"type": "anonymize",
				"config": map[string]interface{}{
					"method":   "aggregate",
					"function": "count",
					"config": map[string]interface{}{
						"k":     0,
						"sigma": 0,
					},
					"time-window": map[string]interface{}{
						"field":  "created-at",
						"window": "day",
						"format": "rfc3339",
					},
					"group-by":       []string{},
					"result-name":    "counts",
					"channels":       []string{"counts"},
					"finalize-after": -1,
				},
			},
		},
		"streams": map[string]interface{}{
			"default": map[string]interface{}{
				"configs": map[string]interface{}{
					"default": map[string]interface{}{
						"status": "active",
						"actions": []map[string]interface{}{
							{
								"name": "uniques-last-24h",
							},
						},
					},
				},
				"sources": map[string]interface{}{},
			},
		},
	},
	Items: []map[string]interface{}{
		map[string]interface{}{
			"_ignore": true,
			"_reset":  true,
		},
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
		map[string]interface{}{
			"created-at": "2009-07-05T10:31:44Z",
		},
		map[string]interface{}{
			"created-at": "2009-07-06T3:33:44Z",
		},
		map[string]interface{}{
			"created-at": "2009-07-07T5:33:46Z",
		},
		map[string]interface{}{
			"created-at": "2009-07-08T9:34:44Z",
		},
		map[string]interface{}{
			"created-at": "2009-07-09T5:33:46Z",
		},
		map[string]interface{}{
			"created-at": "2009-07-10T9:34:44Z",
		},
		map[string]interface{}{
			"_ignore": true,
			"_flush":  true,
		},
	},
}

func BenchmarkAggregate(b *testing.B) {

	var fixtureConfig = []pt.FC{
		pt.FC{pf.Settings{}, "settings"},
		pt.FC{pf.Controller{}, "controller"},
		pt.FC{pf.Project{Name: "test"}, "project"},
		pt.FC{pf.Blueprint{Config: simpleBenchmark.Config, Project: "project"}, "blueprint"},
	}

	fixtures, err := pt.SetupFixtures(fixtureConfig)
	defer pt.TeardownFixtures(fixtureConfig, fixtures)

	if err != nil {
		b.Fatal(err)
	}

	controller := fixtures["controller"].(kodex.Controller)
	streams, err := controller.Streams(map[string]interface{}{"name": "default"})

	if err != nil {
		b.Fatal(err)
	}

	if len(streams) != 1 {
		b.Fatalf("Stream not found")
	}

	stream := streams[0]

	configs, err := stream.Configs()

	if err != nil {
		b.Fatal(err)
	}

	if len(configs) == 0 {
		b.Fatal("No active configs!")
	}

	config := configs[0]

	processor, err := config.Processor()
	if err != nil {
		b.Fatal(err)
	}

	if err := processor.Setup(); err != nil {
		b.Fatal(err)
	}

	sourceItems := make([]*kodex.Item, 0)

	for _, item := range simpleBenchmark.Items {
		sourceItems = append(sourceItems, kodex.MakeItem(item))
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err = processor.Process(sourceItems, nil)
		if err != nil {
			b.Fatal(err)
		}
		if _, err := processor.Finalize(); err != nil {
			b.Fatal(err)
		}
	}

}
