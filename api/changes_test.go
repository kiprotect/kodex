// Kodex (Community Edition - CE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2024  KIProtect GmbH (HRB 208395B) - Germany
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
	"testing"
)

type TestCase struct {
	Index    int
	Expected []any
}

func TestByteArrays(t *testing.T) {
	a := map[string]any{
		"foo": []byte("test"),
		"bar": []byte("bar"),
	}

	b := map[string]any{
		"foo": []byte("test"),
		"bar": []byte("bar"),
	}

	changes := Diff(a, b)

	if len(changes) != 0 {
		t.Fatalf("expected 0 change, got %d - %v", len(changes), changes)
	}

}

func TestDifferentTypes(t *testing.T) {
	a := map[string]any{
		"foo": int64(4),
		"bar": float64(3),
	}

	b := map[string]any{
		"foo": float64(4),
		"bar": int64(3),
	}

	changes := Diff(a, b)

	if len(changes) != 0 {
		t.Fatalf("expected 0 change, got %d - %v", len(changes), changes)
	}

	if err := ApplyChanges(a, changes); err != nil {
		t.Fatal(err, changes)
	}

	if newChanges := Diff(a, b); len(newChanges) != 0 {
		t.Fatalf("should be identical: %v - %v - %v", a, b, changes)
	}

}

func TestDifferentTypesDifferentValues(t *testing.T) {
	a := map[string]any{
		"foo": int64(4),
		"bar": float64(3.3),
	}

	b := map[string]any{
		"foo": float64(4.1),
		"bar": int64(3),
	}

	changes := Diff(a, b)

	if len(changes) != 2 {
		t.Fatalf("expected 2 changes, got %d - %v", len(changes), changes)
	}

	if err := ApplyChanges(a, changes); err != nil {
		t.Fatal(err, changes)
	}

	if newChanges := Diff(a, b); len(newChanges) != 0 {
		t.Fatalf("should be identical: %v - %v - %v", a, b, changes)
	}

}

func TestDeepMapDiff(t *testing.T) {
	a := map[string]any{
		"fields": []any{
			map[string]any{
				"name": "bar",
				"validators": []any{
					map[string]any{
						"type": "IsStringMap",
					},
				},
			},
		},
	}

	b := map[string]any{
		"fields": []any{
			map[string]any{
				"name": "bar",
				"validators": []any{
					map[string]any{
						"type": "IsStringMap",
					},
					map[string]any{
						"type": "IsString",
					},
				},
			},
		},
	}

	changes := Diff(a, b)

	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d - %v", len(changes), changes)
	}

	if err := ApplyChanges(a, changes); err != nil {
		t.Fatal(err, changes)
	}

	if newChanges := Diff(a, b); len(newChanges) != 0 {
		t.Fatalf("should be identical: %v - %v - %v", a, b, changes)
	}
}

func TestMapDiff(t *testing.T) {
	a := map[string]any{
		"foo": map[string]any{"bum": 1},
		"baz": []any{1, 2, 3, "test"},
	}

	b := map[string]any{
		"foo":    map[string]any{"bum": 2},
		"baz":    []any{"test", 1, 2, 3, "foo"},
		"number": 4,
	}

	changes := Diff(a, b)

	if len(changes) != 6 {
		t.Fatalf("expected 6 changes, got %d - %v", len(changes), changes)
	}

	if err := ApplyChanges(a, changes); err != nil {
		t.Fatal(err)
	}

	if newChanges := Diff(a, b); len(newChanges) != 0 {
		t.Fatalf("should be identical: %v - %v - %v", a, b, changes)
	}
}

func TestDiffWithDuplicateIds(t *testing.T) {

	// the diff library does not check for ID uniqueness...

	a := map[string]any{
		"foo": []any{map[string]any{"id": "a"}, map[string]any{"id": "a"}, map[string]any{"id": "c"}, map[string]any{"id": "c"}},
	}

	b := map[string]any{
		"foo": []any{map[string]any{"id": "a", "baz": "bam"}, map[string]any{"id": "a", "buz": "bar"}, map[string]any{"id": "c"}},
	}

	changes := Diff(a, b)

	if len(changes) != 3 {
		t.Fatalf("expected 3 changes, got %d - %v", len(changes), changes)
	}

	if err := ApplyChanges(a, changes); err != nil {
		t.Fatal(err)
	}

	if changes := Diff(a, b); len(changes) == 0 {
		// we expect this to fail...
		t.Fatalf("should not be identical - %v vs %v", a, b)
	}

}

func TestDiffWithIds(t *testing.T) {
	a := map[string]any{
		"foo": []any{map[string]any{"id": "a"}, map[string]any{"id": "b"}, map[string]any{"id": "c"}},
	}

	b := map[string]any{
		"foo": []any{map[string]any{"id": "a", "baz": "bam"}, map[string]any{"id": "c"}, map[string]any{"id": "d"}},
	}

	changes := Diff(a, b)

	if len(changes) != 3 {
		t.Fatalf("expected 3 changes, got %d - %v", len(changes), changes)
	}

	if err := ApplyChanges(a, changes); err != nil {
		t.Fatal(err)
	}

	if newChanges := Diff(a, b); len(newChanges) != 0 {
		t.Fatalf("should be identical - %v vs %v - %v", a, b, newChanges)
	}

}

func TestDiffWithIdsAndSwap(t *testing.T) {
	a := []any{
		map[string]any{
			"id": "a",
		},
		map[string]any{
			"id": "b",
		},
		map[string]any{
			"id": "c",
		},
	}

	b := []any{
		map[string]any{
			"id": "c",
		},
		map[string]any{
			"id": "a",
		},
		map[string]any{
			"id": "d",
		},
		map[string]any{
			"id": "b",
		},
	}

	changes := Diff(a, b)

	if len(changes) != 3 {
		t.Fatalf("expected 3 changes, got %d - %v", len(changes), changes)
	}

	ap, err := ApplyChangesWithObject(a, changes)

	if err != nil {
		t.Fatal(err)
	}

	if newChanges := Diff(ap, b); len(newChanges) != 0 {
		t.Fatalf("should be identical - %v vs %v - %v (%v)", ap, b, newChanges, changes)
	}

}

func TestDiffWithIdsAndSwapAndRemoves(t *testing.T) {
	a := []any{
		map[string]any{
			"id": "a",
		},
		map[string]any{
			"id": "b",
		},
		map[string]any{
			"id": "c",
		},
	}

	b := []any{
		map[string]any{
			"id": "d",
		},
	}

	changes := Diff(a, b)

	if len(changes) != 4 {
		t.Fatalf("expected 3 changes, got %d - %v", len(changes), changes)
	}

	ap, err := ApplyChangesWithObject(a, changes)

	if err != nil {
		t.Fatal(err)
	}

	if newChanges := Diff(ap, b); len(newChanges) != 0 {
		t.Fatalf("should be identical - %v vs %v - %v (%v)", ap, b, newChanges, changes)
	}

}

func TestDiffWithIdsAndSwapAndRemove(t *testing.T) {
	a := []any{
		map[string]any{
			"id": "a",
		},
		map[string]any{
			"id": "b",
		},
		map[string]any{
			"id": "c",
		},
	}

	b := []any{
		map[string]any{
			"id": "d",
		},
		map[string]any{
			"id": "a",
		},
	}

	changes := Diff(a, b)

	if len(changes) != 3 {
		t.Fatalf("expected 3 changes, got %d - %v", len(changes), changes)
	}

	ap, err := ApplyChangesWithObject(a, changes)

	if err != nil {
		t.Fatal(err)
	}

	if newChanges := Diff(ap, b); len(newChanges) != 0 {
		t.Fatalf("should be identical - %v vs %v - %v (%v)", ap, b, newChanges, changes)
	}

}

func TestAdvancedEquality(t *testing.T) {
	a := map[string]any{
		"foo": []any{1, map[string]any{"blub": "blab"}, map[string]any{"blip": "blop"}, "aka", 5, 6},
	}

	b := map[string]any{
		"foo": []any{"aka", "aka", map[string]any{"blom": "blab", "bleb": "blob"}, map[string]any{"fooz": "bar"}, 4, "mama"},
	}

	changes := Diff(a, b)

	if len(changes) != 10 {
		t.Fatalf("expected 10 changes, got %d - %v", len(changes), changes)
	}

	if err := ApplyChanges(a, changes); err != nil {
		t.Fatalf("%v, %v", err, changes)
	}

	newChanges := Diff(a, b)

	if len(newChanges) != 0 {
		t.Fatalf("should be identical - %v vs %v . %v", a, b, changes)
	}

}
func TestDiffWithoutIds(t *testing.T) {
	a := map[string]any{
		"foo": []any{1, map[string]any{"blub": "blab"}, map[string]any{"blip": "blop"}, "aka", 5, 6},
	}

	b := map[string]any{
		"foo": []any{"aka", "aka", map[string]any{"blub": "blab"}, map[string]any{"foo": "bar"}, 4, "mama"},
	}

	changes := Diff(a, b)

	if len(changes) != 9 {
		t.Fatalf("expected 9 changes, got %d - %v", len(changes), changes)
	}

	if err := ApplyChanges(a, changes); err != nil {
		t.Fatalf("%v, %v", err, changes)
	}

	newChanges := Diff(a, b)

	if len(newChanges) != 0 {
		t.Fatalf("should be identical - %v vs %v . %v", a, b, changes)
	}

}

func TestSwap(t *testing.T) {

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {

			l := []any{map[string]any{"id": "a"}, map[string]any{"id": "b"}, map[string]any{"id": "c"}}

			data := map[string]any{
				"foo": append([]any{}, l...),
			}

			changes := []Change{
				{
					Op:    Swap,
					Value: j, // target index
					Path: []PathElement{
						{
							PathType: Direct,
							Name:     "foo",
						},
						{
							PathType: Array,
							Index:    i, // source index
						},
					},
				},
			}

			if err := ApplyChanges(data, changes); err != nil {
				t.Fatal(err)
			}

			ll := data["foo"].([]any)

			if ll[j].(map[string]any)["id"] != l[i].(map[string]any)["id"] {
				t.Fatalf("unexpected: %d - %d - %v - %v", i, j, ll, l)
			}

		}
	}

}

func TestRemoveMapById(t *testing.T) {

	data := map[string]any{
		"foo": []any{map[string]any{"id": "a"}, map[string]any{"id": "b"}, map[string]any{"id": "c"}},
	}

	changes := []Change{
		{
			Op: Remove,
			Path: []PathElement{
				{
					PathType: Direct,
					Name:     "foo",
				},
				{
					PathType:   ById,
					Identifier: "id",
					Value:      "a",
				},
			},
		},
	}

	if err := ApplyChanges(data, changes); err != nil {
		t.Fatal(err)
	}

	l := data["foo"].([]any)

	if len(l) != 2 {
		t.Fatal("expected 2 items")
	}

	if l[0].(map[string]any)["id"] != "b" {
		t.Fatal("wrong")
	}

	if l[1].(map[string]any)["id"] != "c" {
		t.Fatal("wrong")
	}

}

func TestRemoveValueById(t *testing.T) {

	data := map[string]any{
		"foo": []any{map[string]any{"id": "a", "foo": "bar"}, map[string]any{"id": "b", "foo": "bar"}, map[string]any{"id": "c"}},
	}

	changes := []Change{
		{
			Op: Remove,
			Path: []PathElement{
				{
					PathType: Direct,
					Name:     "foo",
				},
				{
					PathType:   ById,
					Identifier: "id",
					Value:      "b",
				},
				{
					PathType: Direct,
					Name:     "foo",
				},
			},
		},
	}

	if err := ApplyChanges(data, changes); err != nil {
		t.Fatal(err)
	}

	l := data["foo"].([]any)

	if len(l) != 3 {
		t.Fatal("expected 3 items")
	}

	if _, ok := l[1].(map[string]any)["foo"]; ok {
		t.Fatal("should be removed")
	}

	if v, ok := l[0].(map[string]any)["foo"]; !ok || v != "bar" {
		t.Fatal("should be there")
	}

}

func TestUpdateById(t *testing.T) {

	data := map[string]any{
		"foo": []any{map[string]any{"id": "a"}, map[string]any{"id": "b"}, map[string]any{"id": "c"}},
	}

	changes := []Change{
		{
			Op:    Update,
			Value: "baz",
			Path: []PathElement{
				{
					PathType: Direct,
					Name:     "foo",
				},
				{
					PathType:   ById,
					Identifier: "id",
					Value:      "c",
				},
				{
					PathType: Direct,
					Name:     "bar",
				},
			},
		},
	}

	if err := ApplyChanges(data, changes); err != nil {
		t.Fatal(err)
	}

	l := data["foo"].([]any)

	if len(l) != 3 {
		t.Fatal("expected 3 items")
	}

	if l[2].(map[string]any)["bar"] != "baz" {
		t.Fatal("wrong")
	}

}

func TestArrayRemove(t *testing.T) {

	testCases := []TestCase{
		{
			-1,
			[]any{"a", "b"},
		},
		{
			0,
			[]any{"b", "c"},
		},
		{
			1,
			[]any{"a", "c"},
		},
		{
			2,
			[]any{"a", "b"},
		},
		{
			3,
			nil,
		},
	}

	for j, testCase := range testCases {

		data := map[string]any{
			"foo": []any{"a", "b", "c"},
		}

		changes := []Change{
			{
				Op: Remove,
				Path: []PathElement{
					{
						PathType: Direct,
						Name:     "foo",
					},
					{
						PathType: Array,
						Index:    testCase.Index,
					},
				},
			},
		}

		if err := ApplyChanges(data, changes); err != nil {

			if testCase.Expected == nil {
				continue
			}

			t.Fatal(err)
		}

		l := data["foo"].([]any)

		if len(l) != len(testCase.Expected) {
			t.Fatalf("expected %d elements, got %d (%v) - test case %d", len(testCase.Expected), len(l), l, j)
		}

		for i, v := range l {
			if testCase.Expected[i] != v {
				t.Fatalf("expected %v, got %v for element %d (%v) - test case %d", testCase.Expected[i], v, i, l, j)
			}
		}

	}

}

func TestArrayInsert(t *testing.T) {

	testCases := []TestCase{
		{
			-1,
			[]any{"a", "b", "c", "one"},
		},
		{
			0,
			[]any{"one", "a", "b", "c"},
		},
		{
			1,
			[]any{"a", "one", "b", "c"},
		},
		{
			2,
			[]any{"a", "b", "one", "c"},
		},
		{
			3,
			[]any{"a", "b", "c", "one"},
		},
		{
			4,
			nil, // we expect an error
		},
	}

	for _, testCase := range testCases {

		data := map[string]any{
			"foo": []any{"a", "b", "c"},
		}

		changes := []Change{
			{
				Op:    Insert,
				Value: "one",
				Path: []PathElement{
					{
						PathType: Direct,
						Name:     "foo",
					},
					{
						PathType: Array,
						Index:    testCase.Index,
					},
				},
			},
		}

		if err := ApplyChanges(data, changes); err != nil {

			if testCase.Expected == nil {
				continue
			}

			t.Fatal(err)
		}

		l := data["foo"].([]any)

		if len(l) != len(testCase.Expected) {
			t.Fatalf("expected %d elements, got %d (%v)", len(testCase.Expected), len(l), l)
		}

		for i, v := range l {
			if testCase.Expected[i] != v {
				t.Fatalf("expected %v, got %v for element %d (%v)", testCase.Expected[i], v, i, l)
			}
		}

	}

}
