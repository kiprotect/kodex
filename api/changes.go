package api

import (
	"fmt"
)

type Op int

const (
	Insert Op = 0
	Remove    = 1
	Update    = 2
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

type Change struct {
	Op          Op            `json:"op"`
	Value       any           `json:"value"`
	Path        []PathElement `json:"path"`
	Description string        `json:"description"`
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

// Applies a sequence of changes to an object
func ApplyChanges(object map[string]any, changes []Change) error {
	for _, change := range changes {

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
					return fmt.Errorf("out of bounds")
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
		case Insert:

			if lastPathElement.PathType != Array {
				return fmt.Errorf("expected an array path element")
			}

			if arrayObj, ok := obj.([]any); !ok {
				return fmt.Errorf("expected an array for insertion, got %T", obj)
			} else if lastPathElement.Index > len(arrayObj) || lastPathElement.Index < -1 {
				return fmt.Errorf("out of bounds")
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
						return fmt.Errorf("out of bounds")
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
