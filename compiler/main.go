package compiler

import (
	"os"

	"github.com/zuma206/sb3c/visualisation"
)

// sb3c entry point with error handling
func Main() error {
	args, err := parseArgs()
	if err != nil {
		return err
	}
	src, err := os.ReadFile(args.Target)
	if err != nil {
		return err
	}
	result, err := Compile(string(src))
	if err != nil {
		return err
	}
	return output(result, args)
}

func output(result *CompileResult, args *Args) error {
	if args.Tokens {
		visualisation.Visualise(result.Tokens)
	}
	if args.Syntax {
		visualisation.Visualise(result.Program)
	}
	return os.WriteFile(args.Outfile, result.Output, 0644)
}
