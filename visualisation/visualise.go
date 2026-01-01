package visualisation

import (
	"fmt"
	"io"
	"os"
	"reflect"
)

type Visualiser struct {
	file        io.Writer
	indentation int
}

func (visualiser *Visualiser) visualiseSlice(value any) {
	valueof := reflect.ValueOf(value)
	fmt.Fprintln(visualiser.file, "[]"+valueof.Type().Elem().Elem().Name(), "{")
	for i := 0; i < valueof.Len(); i++ {
		fmt.Fprintf(visualiser.file, "[%d]: ", i)
		visualiser.visualise(valueof.Index(i).Interface())
	}
	fmt.Fprintln(visualiser.file, "}")
}

// Visualise a value using reflection
func (visualiser *Visualiser) visualiseWithReflection(value any) {
	typeof := reflect.TypeOf(value)
	switch typeof.Kind() {
	case reflect.Slice:
		visualiser.visualiseSlice(value)
	default:
		fmt.Println(value)
	}
}

// Visualise a data structure using a visualiser
func (visualiser *Visualiser) visualise(value any) {
	visualiser.visualiseWithReflection(value)
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
