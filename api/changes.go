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
func diffMap(a, b map[string]any, options DiffOptions) []Change {

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
		changes = append(changes, prependPath(diffAny(va, vb, options), []PathElement{DirectPath(ka)})...)

	}

	return changes
}

// Performs a diff on two arrays
func diffArray(a, b []any, options DiffOptions) []Change {

	changes := make([]Change, 0)

	// we check if both arrays have maps with identifiers

	withIds := false

	var identifier string

identifiers:
	for _, identifier = range options.Identifiers {

		for _, va := range a {
			if mapVa, ok := va.(map[string]any); ok {
				if _, ok := mapVa[identifier]; !ok {
					continue identifiers
				}
			} else {
				continue identifiers
			}
		}

		for _, vb := range b {
			if mapVb, ok := vb.(map[string]any); ok {
				if _, ok := mapVb[identifier]; !ok {
					continue identifiers
				}
			} else {
				continue identifiers
			}
		}

		// we found a working identifier!
		withIds = true
		break

	}

	fmt.Println(withIds)

	similarity := func(a, b any) int {

		mapA, okA := a.(map[string]any)
		mapB, okB := b.(map[string]any)

		if withIds {
			if mapA[identifier] == mapB[identifier] {
				return 0
			}
			return 1
		}

		d := diffAny(a, b, options)

		if okA && okB {

			innerChanges := 0

			for _, change := range d {
				if len(change.Path) == 1 {
					innerChanges += 1
				}
			}

			if innerChanges == 0 && len(d) > 0 {
				return 1
			}

			return innerChanges
		}

		return len(d)
	}

	// - we find removed objects and add 'Remove' changes for them
	// - we find added objects and add 'Insert' changes for them
	// - after that the two lists should have identical objects
	//   just with different order, we add 'Swap' changes to fix that

	// finding removed objects

	am := make([]any, 0, len(a))

	bToAm := make(map[int]int)
	removedObjects := 0

	// we add 'Remove' changes for removed objects
	for i, va := range a {

		bestSimilarity := 0
		bestJ := -1

		for j, vb := range b {

			if _, ok := bToAm[j]; ok {
				continue
			}

			s := similarity(va, vb)

			if bestJ == -1 || s < bestSimilarity {
				bestJ = j
				bestSimilarity = s
			}
		}

		fmt.Println(bestJ, bestSimilarity, va, b[bestJ])

		if bestJ != -1 && bestSimilarity == 0 {
			bToAm[bestJ] = len(am)
			am = append(am, va)
		} else {

			var path []PathElement

			if withIds {
				path = []PathElement{ByIdPath(identifier, va.(map[string]any)[identifier])}
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

	// we add 'Insert' changes for new objects
	for i, vb := range b {

		if _, ok := bToAm[i]; !ok {

			if i >= len(am) {
				i = len(am) - 1
			}

			// this object was added, we add the corresponding change
			changes = append(changes, Change{
				Op:    Insert,
				Value: vb,
				Path:  []PathElement{ArrayPath(i)},
			})

			// we insert the added object to the working list
			am = append(am[:i+1], am[i:]...)
			am[i] = vb
			bToAm[i] = i
			for bi, j := range bToAm {
				// if j >= i it will increase by one
				if j >= i {
					bToAm[bi] = j + 1
				}
			}

		}
	}

	// now only the order of objects in the working list and b should differ

	amToB := make(map[int]int)

	for i, j := range bToAm {
		amToB[j] = i
	}

	for i, _ := range b {

		j := bToAm[i]

		if i != j {

			// element is at position i but should be at j, we swap it
			changes = append(changes, Change{
				Op:    Swap,
				Value: j, // target
				Path:  []PathElement{ArrayPath(i)},
			})

			// we swap the items in the working list
			am[i], am[j] = am[j], am[i]

			// the value that was at j > i will now be at i
			bToAm[amToB[i]] = j

		}
	}

	// now only values should differ
	for i, va := range am {

		var path []PathElement

		if withIds {
			path = []PathElement{ByIdPath(identifier, va.(map[string]any)[identifier])}
		} else {
			path = []PathElement{ArrayPath(i)}
		}

		if elementChanges := diffAny(va, b[i], options); len(elementChanges) > 0 {
			changes = append(changes, prependPath(elementChanges, path)...)
		}

	}

	return changes
}

func diffAny(a, b any, options DiffOptions) []Change {

	taMap, okA := a.(map[string]any)
	tbMap, okB := b.(map[string]any)

	if okA && okB {
		// two maps, we diff them
		return diffMap(taMap, tbMap, options)
	}

	taArray, okA := a.([]any)
	tbArray, okB := b.([]any)

	if okA && okB {
		// two array, we diff them
		return diffArray(taArray, tbArray, options)
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

type DiffOptions struct {
	Identifiers []string
}

func DiffWithOptions(a, b map[string]any, options DiffOptions) []Change {
	return diffMap(a, b, options)
}

// Returns the difference between two data structures as a sequence of changes
func Diff(a, b map[string]any) []Change {
	return diffMap(a, b, DiffOptions{Identifiers: []string{"id"}})
}

// Applies a sequence of changes to an object
func ApplyChanges(object map[string]any, changes []Change) error {

	for _, change := range changes {

		makeErr := func(err error) error {
			return fmt.Errorf("%v: %v", change, err)
		}

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
					return makeErr(fmt.Errorf("direct: expected a map for path %v, got %T", pathElement, obj))
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
					obj = arrayObj[pathElement.Index]
				}
			// we expect and array consisting of map elements that we can descend into...
			case ById:
				if arrayObj, ok := obj.([]any); !ok {
					return fmt.Errorf("expected a array")
				} else {
					found := false
					for _, arrayItem := range arrayObj {
						if mapObj, ok := arrayItem.(map[string]any); !ok {
							return fmt.Errorf("byID: expected a map")
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
