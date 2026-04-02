package compiler

import (
	"errors"
	"flag"
	"fmt"
)

// All positional arguments and flags used by the compiler
type Args struct {
	Tokens  bool
	Syntax  bool
	Outfile string
	Target  string
}

// Expected number of positional arguments needed to construct Args
const expectedPositionalArgs = 1

var (
	// Too many positional arguments were passed to the compiler
	UnexpectedPositionalArgs = errors.New("unexpected positional args")
	// Too few positional arguments were passed to the compiler
	MissingPositionalArgs = errors.New("missing positional args")
)

func newArgs() *Args {
	args := &Args{}
	flag.StringVar(&args.Outfile, "o", "project.sb3", "sb3 file to output to")
	flag.BoolVar(&args.Tokens, "t", false, "print lex tokens to stdout")
	flag.BoolVar(&args.Syntax, "s", false, "print syntax tree to stdout")
	return args
}

// Parses the command line arguments from os.Args and returns them in a struct
func parseArgs() (*Args, error) {
	args := newArgs()
	flag.Parse()
	var err error
	args.Target, err = getTarget()
	return args, err
}

func getTarget() (string, error) {
	args := flag.Args()
	if len(args) != expectedPositionalArgs {
		parentErr := MissingPositionalArgs
		if len(args) > expectedPositionalArgs {
			parentErr = UnexpectedPositionalArgs
		}
		childErr := fmt.Errorf("expected %d got %d", expectedPositionalArgs, len(args))
		return "", errors.Join(parentErr, childErr)
	}
	return args[0], nil
}
