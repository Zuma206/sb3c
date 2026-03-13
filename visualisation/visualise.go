package visualisation

import (
	"fmt"
	"io"
	"iter"
	"os"
	"reflect"

	"github.com/zuma206/sb3c/lexer"
)

type Visualiser struct {
	file        io.Writer
	indentation int
}

func (visualiser *Visualiser) indent(f func()) {
	visualiser.indentation++
	f()
	visualiser.indentation--
}

func (visualiser *Visualiser) print(args ...any) {
	for range visualiser.indentation {
		fmt.Fprint(visualiser.file, "  ")
	}
	fmt.Fprint(visualiser.file, args...)
}

func (visualiser *Visualiser) visualiseSlice(value any) {
	valueof := reflect.ValueOf(value)
	fmt.Fprintln(visualiser.file, "[]"+valueof.Type().Elem().Elem().Name(), "{")
	visualiser.indent(func() {
		for i := 0; i < valueof.Len(); i++ {
			visualiser.print("[", i, "]: ")
			visualiser.visualise(valueof.Index(i).Interface())
		}
	})
	fmt.Fprintln(visualiser.file, "}")
}

func (visualiser *Visualiser) visualiseStruct(value any) {
	valueof := reflect.ValueOf(value)
	typeof := valueof.Type()
	fmt.Fprint(visualiser.file, typeof.Name(), " {\n")
	visualiser.indent(func() {
		for i := range valueof.NumField() {
			field := typeof.Field(i)
			if field.IsExported() {
				visualiser.print(field.Name, ": ")
				visualiser.visualise(valueof.Field(i).Interface())
			}
		}
	})
	visualiser.print("}\n")
}

func (visualiser *Visualiser) visualisePointer(value any) {
	valueof := reflect.ValueOf(value)
	if valueof.IsNil() {
		fmt.Fprintln(visualiser.file, "nil")
	} else {
		visualiser.visualiseWithReflection(valueof.Elem().Interface())
	}
}

// Visualise a value using reflection
func (visualiser *Visualiser) visualiseWithReflection(value any) bool {
	typeof := reflect.TypeOf(value)
	switch typeof.Kind() {
	case reflect.Slice:
		visualiser.visualiseSlice(value)
	case reflect.Pointer:
		visualiser.visualisePointer(value)
	case reflect.Struct:
		visualiser.visualiseStruct(value)
	default:
		return false
	}
	return true
}

type IterAny interface {
	IterAny() iter.Seq[any]
}

func (visualiser *Visualiser) visualiseIterAny(iterAny IterAny) {
	fmt.Fprint(visualiser.file, reflect.TypeOf(iterAny).Elem().Name(), " {\n")
	visualiser.indent(func() {
		for i := range iterAny.IterAny() {
			visualiser.print()
			visualiser.visualise(i)
		}
	})
	visualiser.print("}\n")
}

func (visualiser *Visualiser) visualiseSpecialCase(value any) bool {
	if token, ok := value.(*lexer.Token); ok && token != nil {
		fmt.Fprintf(visualiser.file, "%s(%q, %d:%d)\n",
			token.Type.Name, token.Src, token.Pos.LineNumber, token.Pos.LineOffset)
	} else if error, ok := value.(*lexer.Section); ok {
		fmt.Fprintf(visualiser.file, "%q, %d:%d\n",
			error.Src, error.Pos.LineNumber, error.Pos.LineOffset)
	} else if iterAny, ok := value.(IterAny); ok {
		visualiser.visualiseIterAny(iterAny)
	} else {
		return false
	}
	return true
}

// Visualise a data structure using a visualiser
func (visualiser *Visualiser) visualise(value any) {
	if !(visualiser.visualiseSpecialCase(value) ||
		visualiser.visualiseWithReflection(value)) {
		fmt.Fprintln(visualiser.file, value)
	}
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
