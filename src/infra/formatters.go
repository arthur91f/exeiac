package infra

import (
	"bytes"
	"fmt"
	"os"
)

// TODO(half-shell): Change naming; it sucks.
// A formatter is an interface that allows to write to a file a representation of
// a set of key/value pairs.
// Sample formatters differ from the data format they write to the file to.
// NOTE(half-shell): Should this just be a struct with different member functions
// for each format?
//
//	type Formatter struct {
//		Inputs map[string]interface{}
//	}
//
// func (f Formatter) ToJSON() []byte {}
// func (f Formatter) toYAML() []byte {}
// func (f Formatter) toEnvFile() []byte {}
// func (f Formatter) toEnvVars() []string {}
type Formatter interface {
	Write(f *os.File) (n int, err error)
	// NOTE(half-shel): String representation of the inputs.
	// Each element in the slice are considered to be a line on their own.
	String() (input []string)
}

type JsonFormat struct {
	Inputs map[string]interface{}
}

func (i JsonFormat) Write(f *os.File) (n int, err error) {
	return
}

func (i JsonFormat) String() (input []string) {
	return
}

type EnvFormat struct {
	Inputs map[string]interface{}
}

func (i EnvFormat) Write(f *os.File) (n int, err error) {
	buf := new(bytes.Buffer)
	for varName, varVal := range i.Inputs {
		buf.WriteString(fmt.Sprintf("%s=%v\n", varName, varVal))
	}

	n, err = f.Write(buf.Bytes())

	return
}

func (i EnvFormat) String() (input []string) {
	for varName, varVal := range i.Inputs {
		input = append(input, fmt.Sprintf("%s=%v", varName, varVal))
	}

	return
}
