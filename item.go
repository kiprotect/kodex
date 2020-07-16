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

package kiprotect

import (
	"encoding/json"
	"fmt"
)

type Item struct {
	d map[string]interface{}
}

func MakeItem(d map[string]interface{}) *Item {
	item := new(Item)
	item.d = d
	return item
}

func (f *Item) Keys() []string {
	keys := make([]string, 0)
	for key := range f.d {
		keys = append(keys, key)
	}
	return keys
}

func (f *Item) Values() []interface{} {
	values := make([]interface{}, 0)
	for _, value := range f.d {
		values = append(values, value)
	}
	return values
}

func (f *Item) Delete(key string) {
	delete(f.d, key)
}

func (f *Item) Get(key string) (interface{}, bool) {
	v, ok := f.d[key]
	return v, ok
}

func (f *Item) All() map[string]interface{} {
	return f.d
}

func (f *Item) Set(key string, value interface{}) {
	f.d[key] = value
}

func (f *Item) Serialize(format string) ([]byte, error) {
	switch format {
	case "json":
		return f.SerializeJSON()
	default:
		return nil, fmt.Errorf("Unknown format: %s", format)
	}
}

func (f *Item) SerializeJSON() ([]byte, error) {
	return json.Marshal(f.d)
}

// Validates the format of a list of items
type IsItem struct{}

func (i IsItem) Validate(value interface{}, values map[string]interface{}) (interface{}, error) {
	strMap, ok := value.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("expected a string map")
	}
	return MakeItem(strMap), nil
}

type IsItems struct{}

func (i IsItems) Validate(value interface{}, values map[string]interface{}) (interface{}, error) {
	list, ok := value.([]interface{})
	if !ok {
		return nil, fmt.Errorf("expected a list")
	}
	itemList := make([]*Item, len(list))
	for i, listValue := range list {
		item, ok := listValue.(*Item)
		if !ok {
			return nil, fmt.Errorf("not an item")
		}
		itemList[i] = item
	}
	return itemList, nil
}
