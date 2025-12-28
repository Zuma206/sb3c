package main

import (
	"fmt"
	"os"

	"github.com/zuma206/sb3c/sb3c"
)

// Main compiled binary entry point, calls sb3c and prints fatal errors to stderr
func main() {
	if err := sb3c.Main(); err != nil {
		fmt.Fprintf(os.Stderr, "fatal error: %s\n", err.Error())
		os.Exit(1)
	}
}
