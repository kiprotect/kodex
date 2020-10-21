package groupByFunctions

import (
	"github.com/kiprotect/kodex"
)

type GroupByValue struct {
	Values     map[string]interface{}
	Expiration int64
}

type GroupByFunction func(item *kodex.Item) ([]*GroupByValue, error)
type GroupByFunctionMaker func(map[string]interface{}) (GroupByFunction, error)

var Functions = map[string]GroupByFunctionMaker{
	"time-window": MakeTimeWindowFunction,
}
