package api

import (
	"testing"
)

type TestCase struct {
	Index    int
	Expected []any
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
