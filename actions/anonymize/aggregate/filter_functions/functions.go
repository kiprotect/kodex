package filterFunctions

import (
	"github.com/kiprotect/kodex"
)

type FilterFunction func(item *kodex.Item) (bool, error)
type FilterFunctionMaker func(map[string]interface{}) (FilterFunction, error)

var Functions = map[string]FilterFunctionMaker{}
