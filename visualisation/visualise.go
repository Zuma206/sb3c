package visualisation

import (
	"fmt"
	"io"
	"os"
)

type Visualiser struct {
	file        io.Writer
	indentation int
}

// Visualise a data structure using a visualiser
func (visualiser *Visualiser) visualise(value any) {
	fmt.Println(value)
}

// Visualise a data structure on the given file
func Fvisualise(file io.Writer, value any) {
	visualiser := Visualiser{file: file, indentation: 0}
	visualiser.visualise(value)
}

// Visualise a data structure
func Visualise(value any) {
	Fvisualise(os.Stdout, value)
}
