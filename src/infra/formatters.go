package infra

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type InputFormat string

const (
	Json = InputFormat("json")
	Yaml = InputFormat("yaml")
	Env  = InputFormat("env")
)

var SupportedFormats = map[string]InputFormat{
	"json": Json,
	"yaml": Yaml,
	"env":  Env,
}

// A formatter is an interface that allows to produce an output in a specific format.
type Formatter interface {
	Format() (input []byte, err error)
}

type JsonFormat map[string]interface{}
type EnvFormat map[string]interface{}

func (i JsonFormat) Format() (input []byte, err error) {
	input, err = json.MarshalIndent(i, "", "\t")

	return
}

func (i EnvFormat) Format() (input []byte, err error) {
	buf := new(bytes.Buffer)
	for varName, varVal := range i {
		buf.WriteString(fmt.Sprintf("%s=\"%v\"\n", varName, varVal))
	}

	input = buf.Bytes()

	return
}
