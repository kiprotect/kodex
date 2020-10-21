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
		return []*GroupByValue{
			&GroupByValue{
				Values: map[string]interface{}{
					field: value,
				},
			},
		}, nil
	}, nil
}
