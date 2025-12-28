package compiler

import (
	"os"
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
	return CompileFile(args.Target, src)
}
