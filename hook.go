package kodex

type HookDefinition struct {
	Description string `json:"description"`
	Hook        Hook   `json:"-"`
}

type HookDefinitions map[string][]HookDefinition
type Hook func(data interface{}) (interface{}, error)
