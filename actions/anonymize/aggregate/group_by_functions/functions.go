package groupByFunctions

import (
	"github.com/kiprotect/kodex"
)

type GroupByFunction func(item *kodex.Item) ([]map[string]interface{}, error)
type GroupByFunctionMaker func(map[string]interface{}) (GroupByFunction, error)

var Functions = map[string]GroupByFunctionMaker{
	"time-window": MakeTimeWindowFunction,
}
