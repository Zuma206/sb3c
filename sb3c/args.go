package sb3c

import (
	"errors"
	"flag"
	"fmt"
)

// All positional arguments and flags used by the compiler
type Args struct {
	Target string
}

// Expected number of positional arguments needed to construct Args
const expectedPositionalArgs = 1

var (
	// Too many positional arguments were passed to the compiler
	UnexpectedPositionalArgs = errors.New("unexpected positional args")
	// Too few positional arguments were passed to the compiler
	MissingPositionalArgs = errors.New("missing positional args")
)

// Parses the command line arguments from os.Args and returns them in a struct
func parseArgs() (*Args, error) {
	flag.Parse()
	args := flag.Args()
	if len(args) != expectedPositionalArgs {
		err := MissingPositionalArgs
		if len(args) > expectedPositionalArgs {
			err = UnexpectedPositionalArgs
		}
		return nil, fmt.Errorf("%w: expected %d got %d", err, expectedPositionalArgs, len(args))
	}
	return &Args{Target: args[0]}, nil
}
