package infra

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	extools "src/exeiac/tools"
	"strings"
)

const ACTION_HELP = "help"
const ACTION_INIT = "init"
const ACTION_LAY = "lay"
const ACTION_OUTPUT = "output"
const ACTION_PLAN = "plan"
const ACTION_REMOVE = "remove"
const ACTION_SHOW_AVAILABLE_ACTIONS = "show_implemented_actions"

type Module struct {
	Name    string
	Path    string
	Actions []string
}

func (m Module) String() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("-\tName: %s\n", m.Name))
	sb.WriteString(fmt.Sprintf("\tPath: %s\n", m.Path))
	sb.WriteString(fmt.Sprintf("\tActions: %v\n", m.Actions))

	return sb.String()
}

// Executes the ACTION_SHOW_AVAILABLE_ACTIONS command on a module to get
// the available actions from the module. Then parses and saves them in the
// *Actions* slice.
// If the *Actions* slice is not empty, bypass the call the the command.
// NOTE(half-shell): Do we have a use to force the call to be triggered again here?
func (module *Module) LoadAvailableActions() error {
	// Actions are already loaded; no need to reprocess it
	if len(module.Actions) > 0 {
		return nil
	}

	path, err := exec.LookPath(module.Path)
	if err != nil {
		return fmt.Errorf("unable to load available actions for module %s: %v", module.Name, err)
	}

	cmd := exec.Cmd{
		Path: path,
		// NOTE(half-shell): We have to manually add something as a first element in args
		// because Cmd **seems** to poorly overwrite it, making any first argument provided disappear
		// e.g. Args:   []string{ACTION_SHOW_AVAILABLE_ACTIONS} won't work
		Args: []string{path, ACTION_SHOW_AVAILABLE_ACTIONS},
	}

	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("unable to load available actions for module %s: %v", module.Name, err)
	}

	for _, action := range bytes.Split(output, []byte("\n")) {
		if len(action) > 0 {
			module.Actions = append(module.Actions, string(action))
		}
	}

	return nil
}

func (m *Module) exec(
	brick *Brick,
	args []string,
	env []string,
	stdout io.Writer,
	stderr io.Writer,
) (
	err error,
) {
	cmd := exec.Cmd{
		Path:   m.Path,
		Args:   append([]string{m.Path}, args...),
		Env:    env,
		Dir:    brick.Path,
		Stdin:  os.Stdin,
		Stdout: stdout,
		Stderr: stderr,
	}

	err = cmd.Run()

	return
}

// Executes a module's action over a brick, the provided CLI arguments and environment
// variables. It takes between 0 and 2 writers; they are used to process the module's
// `stdout` and `stderr`. They'll *usually* match one of infra's writers.
//
// Returns a statusCode returned by the module, and an error if any.
// Note that `err` here is not an error thrown from the external module, but only coming
// from the go execution. Module errors are displayed in `stderr`.
func (m *Module) Exec(
	b *Brick,
	action string,
	args []string,
	env []string,
	writers ...io.Writer,
) (
	statusCode int,
	err error,
) {
	if !extools.ContainsString(m.Actions, action) {
		err = ActionNotImplementedError{Action: action, Module: m}

		return
	}

	if len(env) != 0 {
		env = append(os.Environ(), env...)
	} else {
		env = os.Environ()
	}

	if len(writers) > 1 {
		err = m.exec(b, append([]string{action}, args...), env, writers[0], writers[1])
	} else if len(writers) > 0 {
		err = m.exec(b, append([]string{action}, args...), env, writers[0], os.Stderr)
	} else {
		err = m.exec(b, append([]string{action}, args...), env, os.Stdout, os.Stderr)
	}

	if err != nil {
		if ee, isExitError := err.(*exec.ExitError); isExitError {
			// NOTE(half-shell): We don't consider an exitError an actual error as far as
			// exeiac goes.
			// We return it as a separate value to make that distinction obvious.
			statusCode = ee.ExitCode()
			err = nil
		}
	}

	return
}
