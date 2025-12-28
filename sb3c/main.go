package sb3c

import "fmt"

// sb3c entry point with error handling
func Main() error {
	args, err := parseArgs()
	if err != nil {
		return err
	}
	fmt.Println(args.Target)
	return nil
}
