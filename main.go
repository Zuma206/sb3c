package main

import (
	"fmt"
	"os"

	"github.com/zuma206/sb3c/compiler"
)

const (
	colorRed  = "\033[0;31m"
	colorNone = "\033[0m"
)

// Main compiled binary entry point, calls the compiler and prints fatal errors to stderr
func main() {
	if err := compiler.Main(); err != nil {
		fmt.Fprintln(os.Stderr, colorRed+"encountered fatal error(s)\n"+err.Error()+colorNone)
		os.Exit(1)
	}
}
