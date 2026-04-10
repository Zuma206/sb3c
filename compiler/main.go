package compiler

import (
	"os"
	"path/filepath"

	"github.com/zuma206/sb3c/visualisation"
)

type Inputs struct {
	args *Args
	src  string
	dir  string
}

func getInputs() (*Inputs, error) {
	args, err := parseArgs()
	if err != nil {
		return nil, err
	}
	src, err := os.ReadFile(args.Target)
	if err != nil {
		return nil, err
	}
	targetPath, err := filepath.Abs(args.Target)
	if err != nil {
		return nil, err
	}
	return &Inputs{
		args: args,
		src:  string(src),
		dir:  filepath.Dir(targetPath),
	}, nil
}

func Main() error {
	inputs, err := getInputs()
	if err != nil {
		return err
	}
	result := Compile(inputs.src, os.DirFS(inputs.dir))
	if result.Tokens != nil && inputs.args.Tokens {
		visualisation.Visualise(result.Tokens)
	}
	if result.Program != nil && inputs.args.Syntax {
		visualisation.Visualise(result.Program)
	}
	if result.Err != nil {
		return result.Err
	}
	return os.WriteFile(inputs.args.Outfile, result.Output, 0644)
}
