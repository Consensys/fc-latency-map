package cli

import (
	"github.com/c-bata/go-prompt"
)

type Commander interface {

	// Complete completes the input
	Complete() []prompt.Suggest

	// Execute executes the command
	Execute(in string)
}
