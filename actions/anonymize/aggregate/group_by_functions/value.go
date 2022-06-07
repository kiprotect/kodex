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

package groupByFunctions

import (
	"github.com/kiprotect/go-helpers/errors"
	"github.com/kiprotect/kodex"
)

func MakeValueFunction(config map[string]interface{}) (GroupByFunction, error) {

	field := config["field"].(string)

	return func(item *kodex.Item) ([]*GroupByValue, error) {
		value, ok := item.Get(field)
		if !ok {
			return nil, errors.MakeExternalError("group-by value not defined",
				"VALUE-NOT-DEFINED",
				field,
				nil)
		}
		if isList := config["is-list"].(bool); isList {
			listValue, ok := value.([]interface{})
			if !ok {
				return nil, errors.MakeExternalError("expected a list value",
					"VALUE-EXPECTED-LIST",
					value,
					nil)
			}
			i := int(config["index"].(int64))
			if i >= len(listValue) {
				// the value is undefined, we return nothing
				return nil, nil
			} else {
				if mapValue, ok := listValue[i].(map[string]interface{}); ok {
					// this is a map value, we return it directly
					return []*GroupByValue{
						&GroupByValue{
							Values:     mapValue,
							Expiration: 0,
						},
					}, nil
				}
				value = listValue[i]
			}
		}
		return []*GroupByValue{
			&GroupByValue{
				Values: map[string]interface{}{
					field: value,
				},
				Expiration: 0,
			},
		}, nil
	}, nil
}
