package sb3c

import (
	"errors"
	"flag"
	"fmt"
)

type Args struct {
	Target string
}

const expectedPositionalArgs = 1

var (
	UnexpectedPositionalArgs = errors.New("unexpected positional args")
	MissingPositionalArgs    = errors.New("missing positional args")
)

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
