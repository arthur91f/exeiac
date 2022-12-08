package infra

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	extools "src/exeiac/tools"
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
	if len(m.Actions) > 0 {
		return fmt.Sprintf("name: %s\npath: %s\nactions: %v\n",
			m.Name, m.Path, m.Actions)
	} else {
		return fmt.Sprintf("name: %s\npath: %s\nactions: []",
			m.Name, m.Path)
	}
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
		return err
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
		return err
	}

	for _, action := range bytes.Split(output, []byte("\n")) {
		if len(action) > 0 {
			module.Actions = append(module.Actions, string(action))
		}
	}

	return nil
}

func (m *Module) exec(brick *Brick, args []string, stdout io.Writer, stderr io.Writer) (err error) {
	cmd := exec.Cmd{
		Path:   m.Path,
		Args:   append([]string{m.Path}, args...),
		Env:    []string{},
		Dir:    brick.Path,
		Stdin:  os.Stdin,
		Stdout: stdout,
		Stderr: stderr,
	}

	err = cmd.Run()

	return
}

type ActionNotImplementedError struct {
	Action string
	Module *Module
}

func (err ActionNotImplementedError) Error() string {
	return fmt.Sprintf("Module %s does not implement action %s", err.Module.Name, err.Action)
}

type StdOutput struct {
	Output []byte
}

func (stdout *StdOutput) Write(p []byte) (n int, err error) {
	stdout.Output = append(stdout.Output, p...)

	return len(p), nil
}

type StdError struct {
	Output []byte
}

func (stderr *StdError) Write(p []byte) (n int, err error) {
	stderr.Output = append(stderr.Output, p...)

	return len(p), nil
}

func (m *Module) Exec(b *Brick, action string, args []string, writers ...io.Writer) (exitError *exec.ExitError, err error) {
	if !extools.ContainsString(m.Actions, action) {
		err = ActionNotImplementedError{Action: action, Module: m}

		return
	}

	if len(writers) > 1 {
		err = m.exec(b, append([]string{action}, args...), writers[0], writers[1])
	} else if len(writers) > 0 {
		err = m.exec(b, append([]string{action}, args...), writers[0], os.Stderr)
	} else {
		err = m.exec(b, append([]string{action}, args...), os.Stdout, os.Stderr)
	}

	if err != nil {
		if ee, isExitError := err.(*exec.ExitError); isExitError {
			// NOTE(half-shell): We don't consider an exitError an actual error as far as exeiac goes.
			// We return it as a separate value to make that distinction obvious.
			exitError = ee
			err = nil
		}
	}

	return
}
