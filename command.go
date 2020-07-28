package kiprotect

import (
	"github.com/urfave/cli"
)

type CommandsDefinition struct {
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Maker       CommandsMaker `json:"-"`
}

type CommandsMaker func(controller Controller) ([]cli.Command, error)
type CommandsDefinitions map[string]CommandsDefinition
