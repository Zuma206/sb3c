package main

import (
	"fmt"
	"os"

	"github.com/zuma206/sb3c/compiler"
)

// Main compiled binary entry point, calls the compiler and prints fatal errors to stderr
func main() {
	if err := compiler.Main(); err != nil {
		fmt.Fprintf(os.Stderr, "fatal error: %s\n", err.Error())
		os.Exit(1)
	}
}
