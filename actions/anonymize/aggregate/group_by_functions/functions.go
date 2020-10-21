package groupByFunctions

import (
	"github.com/kiprotect/kodex"
)

type GroupByFunction func(item *kodex.Item) []map[string]interface{}
type GroupByFunctionMaker func(map[string]interface{}) (GroupByFunction, error)

var Functions = map[string]GroupByFunctionMaker{
	"time-window": MakeTimeWindowFunction,
}
