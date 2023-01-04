package actions

import (
	"errors"
	"reflect"
	"testing"

	"src/exeiac/arguments"
	exinfra "src/exeiac/infra"
)

var infra *exinfra.Infra
var conf *arguments.Configuration
var bricksToExecute exinfra.Bricks

func TestHelpWhenNoBricksToExecute(t *testing.T) {
	bricksToExecute = exinfra.Bricks{}
	expectedStatusCode := 3

	statusCode, err := Help(infra, conf, bricksToExecute)

	if statusCode != expectedStatusCode {
		t.Fatalf(`statusCode should be %v; is "%v"`, expectedStatusCode, statusCode)
	}

	if errors.Is(err, exinfra.ErrBadArg{}) {
		t.Fatalf("Errors should be of type ErrBadArg, is %v", reflect.TypeOf(err))
	}
}
