package infra

import "fmt"

type ActionNotImplementedError struct {
	Action string
	Module *Module
}

func (err ActionNotImplementedError) Error() string {
	return fmt.Sprintf("Module %s does not implement action %s", err.Module.Name, err.Action)
}
