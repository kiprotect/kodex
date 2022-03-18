package fixtures

import (
	"github.com/kiprotect/kodex/api"
)

type Definitions struct {
	Definitions api.Definitions
}

func (c Definitions) Setup(fixtures map[string]interface{}) (interface{}, error) {
	return c.Definitions, nil
}

func (c Definitions) Teardown(fixture interface{}) error {
	return nil
}
