package api

import (
	"fmt"
)

type Op int

const (
	Insert Op = 0
	Remove    = 1
	Update    = 2
	Swap      = 3
)

type PathType int

const (
	Direct PathType = 0
	ById   PathType = 1
	Array  PathType = 2
)

// Project details

type PathElement struct {
	PathType   PathType `json:"pathType"`
	Name       string   `json:"name"`
	Index      int      `json:"index"`
	Identifier string   `json:"identifier"`
	Value      any      `json:"value"`
}

func (p PathElement) String() string {
	switch p.PathType {
	case Direct:
		return fmt.Sprintf("%s", p.Name)
	case ById:
		return fmt.Sprintf("%s=%v", p.Identifier, p.Value)
	case Array:
		return fmt.Sprintf("%d", p.Index)
	}
	return "unknown"
}

type Change struct {
	Op          Op            `json:"op"`
	Value       any           `json:"value"`
	Path        []PathElement `json:"path"`
	Description string        `json:"description"`
}

func (c Change) String() string {

	var path string
	var op string

	switch c.Op {
	case Insert:
		op = "insert"
	case Remove:
		op = "remove"
	case Update:
		op = "update"
	case Swap:
		op = "swap"
	default:
		op = "unknown"
	}

	for _, pathElement := range c.Path {

		if path != "" {
			path += "/"
		}

		path += pathElement.String()
	}

	return fmt.Sprintf("%s(path=%s value=%v)", op, path, c.Value)
}

func ArrayPath(index int) PathElement {
	return PathElement{
		PathType: Array,
		Index:    index,
	}
}

func DirectPath(name string) PathElement {
	return PathElement{
		PathType: Direct,
		Name:     name,
	}
}

func ByIdPath(identifier string, value any) PathElement {
	return PathElement{
		PathType:   ById,
		Identifier: identifier,
		Value:      value,
	}
}

func prependPath(changes []Change, pathElements []PathElement) []Change {
	newChanges := make([]Change, 0, len(changes))

	for _, change := range changes {
		change.Path = append(pathElements, change.Path...)
		newChanges = append(newChanges, change)
	}

	return newChanges
}

// Performs a diff on two maps
func diffMap(a, b map[string]any) []Change {

	changes := make([]Change, 0)

	// we check if any of a's keys are missing in b
	for ka, va := range a {

		_, ok := b[ka]

		if !ok {
			// this value was removed, we use a remove operation
			changes = append(changes, Change{Op: Remove, Value: va, Path: []PathElement{DirectPath(ka)}})
		}
	}

	// we check if any of b's keys are missing in a
	for kb, vb := range b {

		_, ok := a[kb]

		if !ok {
			// this is a new value, we use an update operation to add it
			changes = append(changes, Change{Op: Update, Value: vb, Path: []PathElement{DirectPath(kb)}})
		}
	}

	// we compare all common keys in a & b
	for ka, va := range a {

		vb, ok := b[ka]

		if !ok {
			continue
		}

		// we diff the changes between a & b, prepend the current path and return them
		changes = append(changes, prependPath(diffAny(va, vb), []PathElement{DirectPath(ka)})...)

	}

	return changes
}

// Performs a diff on two arrays
func diffArray(a, b []any) []Change {

	changes := make([]Change, 0)

	// we check if both arrays have maps with identifiers

	withIds := true

	for _, va := range a {
		if mapVa, ok := va.(map[string]any); ok {
			if _, ok := mapVa["id"]; !ok {
				withIds = false
				break
			}
		} else {
			withIds = false
			break
		}
	}

	if withIds {
		for _, vb := range b {
			if mapVb, ok := vb.(map[string]any); ok {
				if _, ok := mapVb["id"]; !ok {
					withIds = false
					break
				}
			} else {
				withIds = false
				break
			}
		}
	}

	equal := func(a, b any) bool {
		if withIds {
			mapA := a.(map[string]any)
			mapB := b.(map[string]any)
			return mapA["id"] == mapB["id"]
		}
		return len(diffAny(a, b)) == 0
	}

	// - we find removed objects and add 'Remove' changes for them
	// - we find added objects and add 'Insert' changes for them
	// - after that the two lists should have identical objects
	//   just with different order, we add 'Swap' changes to fix that

	// finding removed objects

	am := make([]any, 0, len(a))

	usedObjects := make(map[int]bool)
	removedObjects := 0

	// we add 'Remove' changes for removed objects
	for i, va := range a {
		found := false
		for j, vb := range b {

			if _, ok := usedObjects[j]; ok {
				continue
			}

			if equal(va, vb) {
				found = true
				usedObjects[j] = true
				am = append(am, va)
				break
			}
		}

		if !found {

			var path []PathElement

			if withIds {
				path = []PathElement{ByIdPath("id", va.(map[string]any)["id"])}
			} else {
				path = []PathElement{ArrayPath(i - removedObjects)}
			}

			// we keep track of how many objects we already removed
			removedObjects += 1

			// this object was removed, we add a corresponding change
			changes = append(changes, Change{
				Op:    Remove,
				Value: va,
				Path:  path,
			})
		}
	}

	usedObjects = make(map[int]bool)

	// we add 'Insert' changes for new objects
	for i, vb := range b {

		found := false

		for j, va := range a {

			if _, ok := usedObjects[j]; ok {
				continue
			}

			if equal(va, vb) {
				found = true
				usedObjects[j] = true
				break
			}
		}

		if !found {

			// this object was added, we add the corresponding change
			changes = append(changes, Change{
				Op:    Insert,
				Value: vb,
				Path:  []PathElement{ArrayPath(i)},
			})

			if i >= len(am) {
				am = append(am, vb)
			} else {
				// we insert the added object to the working list
				am = append(am[:i+1], am[i:]...)
				am[i] = vb
			}

		}
	}

	// now only the order of objects in the working list and b should differ

	for i, vb := range b {

		var j int

		for j = i; j < len(am); j++ {
			if equal(vb, am[j]) {
				break
			}
		}

		if i != j {
			// element is at position i but should be at j, we swap it
			changes = append(changes, Change{
				Op:    Swap,
				Value: j, // target
				Path:  []PathElement{ArrayPath(i)},
			})

			// we swap the items in the working list
			am[i], am[j] = am[j], am[i]
		}
	}

	// now only values should differ
	for i, va := range am {

		var path []PathElement

		if withIds {
			path = []PathElement{ByIdPath("id", va.(map[string]any)["id"])}
		} else {
			path = []PathElement{ArrayPath(i)}
		}

		if elementChanges := diffAny(va, b[i]); len(elementChanges) > 0 {
			changes = append(changes, prependPath(elementChanges, path)...)
		}

	}

	return changes
}

func diffAny(a, b any) []Change {

	taMap, okA := a.(map[string]any)
	tbMap, okB := b.(map[string]any)

	if okA && okB {
		// two maps, we diff them
		return diffMap(taMap, tbMap)
	}

	taArray, okA := a.([]any)
	tbArray, okB := b.([]any)

	if okA && okB {
		// two array, we diff them
		return diffArray(taArray, tbArray)
	}

	if a != b {
		// the value changed, we return an update c hange
		return []Change{
			{
				Op:    Update,
				Value: b,
				Path:  []PathElement{},
			},
		}
	}

	// nothing changed
	return nil
}

// Returns the difference between two data structures as a sequence of changes
func Diff(a, b map[string]any) []Change {
	return diffMap(a, b)
}

// Applies a sequence of changes to an object
func ApplyChanges(object map[string]any, changes []Change) error {
	for _, change := range changes {

		if len(change.Path) < 1 {
			return fmt.Errorf("invalid change path")
		}

		var obj, previousObj any

		obj = object

		for j, pathElement := range change.Path {

			// we do not follow the last element
			if j == len(change.Path)-1 {
				break
			}

			previousObj = obj

			switch pathElement.PathType {

			// we expect a map with a given key that we can descend into...
			case Direct:
				if mapObj, ok := obj.(map[string]any); !ok {
					return fmt.Errorf("expected a map")
				} else if subObj, ok := mapObj[pathElement.Name]; !ok {
					return fmt.Errorf("unknown key: %s", pathElement.Name)
				} else {
					obj = subObj
				}
			// we expect an array with a specific index that we can descend into...
			case Array:
				if arrayObj, ok := obj.([]any); !ok {
					return fmt.Errorf("expected a array")
				} else if len(arrayObj) <= pathElement.Index {
					return fmt.Errorf("array element out of bounds")
				} else {
					obj = arrayObj
				}
			// we expect and array consisting of map elements that we can descend into...
			case ById:
				if arrayObj, ok := obj.([]any); !ok {
					return fmt.Errorf("expected a array")
				} else {
					found := false
					for _, arrayItem := range arrayObj {
						if mapObj, ok := arrayItem.(map[string]any); !ok {
							return fmt.Errorf("expected a map")
						} else if idValue, ok := mapObj[pathElement.Identifier]; !ok {
							return fmt.Errorf("identifier missing")
						} else if idValue == pathElement.Value {
							found = true
							obj = mapObj
							break
						}
					}
					if !found {
						return fmt.Errorf("object with id '%s=%v' not found", pathElement.Identifier, pathElement.Value)
					}
				}
			}

		}

		lastPathElement := change.Path[len(change.Path)-1]

		switch change.Op {
		// we expect an array element and a 'Array' path element
		case Swap:

			targetIndex, ok := change.Value.(int)

			if !ok {
				return fmt.Errorf("expected a target index")
			}

			if lastPathElement.PathType != Array {
				return fmt.Errorf("expected an array path element")
			}

			if arrayObj, ok := obj.([]any); !ok {
				return fmt.Errorf("expected an array for insertion, got %T", obj)
			} else if lastPathElement.Index >= len(arrayObj) || lastPathElement.Index < 0 {
				return fmt.Errorf("move: source out of bounds")
			} else {

				if targetIndex < 0 || targetIndex > len(arrayObj)-1 {
					return fmt.Errorf("move: target out of bounds")
				}

				if len(change.Path) < 2 {
					return fmt.Errorf("expected at least 2 path elements")
				}

				beforeLastPathElement := change.Path[len(change.Path)-2]

				if beforeLastPathElement.PathType != Direct {
					return fmt.Errorf("expected a direct path")
				}

				index := lastPathElement.Index

				arrayObj[index], arrayObj[targetIndex] = arrayObj[targetIndex], arrayObj[index]

			}
		// we expect an array element and a 'Array' path element
		case Insert:

			if lastPathElement.PathType != Array {
				return fmt.Errorf("expected an array path element")
			}

			if arrayObj, ok := obj.([]any); !ok {
				return fmt.Errorf("expected an array for insertion, got %T", obj)
			} else if lastPathElement.Index > len(arrayObj) || lastPathElement.Index < -1 {
				return fmt.Errorf("insert: out of bounds")
			} else {

				// the parent object should be a map
				previousMapObj, ok := previousObj.(map[string]any)

				if !ok {
					return fmt.Errorf("expected a map, got %T", previousObj)
				}

				if len(change.Path) < 2 {
					return fmt.Errorf("expected at least 2 path elements")
				}

				beforeLastPathElement := change.Path[len(change.Path)-2]

				if beforeLastPathElement.PathType != Direct {
					return fmt.Errorf("expected a direct path")
				}

				index := lastPathElement.Index

				if index == -1 || index == len(arrayObj) {
					arrayObj = append(arrayObj, change.Value)
				} else {
					arrayObj = append(arrayObj[:index+1], arrayObj[index:]...)
					arrayObj[index] = change.Value
				}

				previousMapObj[beforeLastPathElement.Name] = arrayObj

			}
		// we expect either a map and a 'Direct' path element, or an array
		// and a 'ById' or 'Array' path element. The resulting element
		// will be removed from the array
		case Remove:
			if mapObj, ok := obj.(map[string]any); ok {

				if lastPathElement.PathType != Direct {
					return fmt.Errorf("expected a direct path element")
				}

				// we remove the key from the map
				delete(mapObj, lastPathElement.Name)

			} else if arrayObj, ok := obj.([]any); ok {

				previousMapObj, ok := previousObj.(map[string]any)

				if !ok {
					return fmt.Errorf("expected a map with a array")
				}

				var index int

				if lastPathElement.PathType == Array {
					if lastPathElement.Index >= len(arrayObj) || lastPathElement.Index < -1 {
						return fmt.Errorf("remove: out of bounds")
					}
					index = lastPathElement.Index
				} else if lastPathElement.PathType == ById {
					found := false
					for i, elem := range arrayObj {
						if mapElem, ok := elem.(map[string]any); !ok {
							return fmt.Errorf("expected a map element")
						} else if value, ok := mapElem[lastPathElement.Identifier]; !ok {
							return fmt.Errorf("identifier missing")
						} else if value == lastPathElement.Value {
							index = i
							found = true
							break
						}
					}
					if !found {
						return fmt.Errorf("element not found")
					}
				}

				if len(change.Path) < 2 {
					return fmt.Errorf("expected at least 2 path elements")
				}

				beforeLastPathElement := change.Path[len(change.Path)-2]

				if beforeLastPathElement.PathType != Direct {
					return fmt.Errorf("expected a direct path")
				}

				if index == -1 {
					index = len(arrayObj)
				}

				if index == len(arrayObj) {
					arrayObj = arrayObj[:index-1]
				} else {
					arrayObj = append(arrayObj[:index], arrayObj[index+1:]...)
				}

				// we remove the item from the array
				previousMapObj[beforeLastPathElement.Name] = arrayObj

			} else {
				return fmt.Errorf("expected a map or array")
			}
		// we expect a map element and a 'Direct' path element
		case Update:

			if lastPathElement.PathType != Direct {
				return fmt.Errorf("expected a direct path")
			}

			if mapObj, ok := obj.(map[string]any); !ok {
				return fmt.Errorf("object is not a map")
			} else {
				mapObj[lastPathElement.Name] = change.Value
			}
		}
	}
	return nil
}
