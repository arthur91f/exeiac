package infra

import (
	// "bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	extools "src/exeiac/tools"
	"strings"
	// "strings"
)

const ACTION_HELP = "help"
const ACTION_INIT = "init"
const ACTION_LAY = "lay"
const ACTION_OUTPUT = "output"
const ACTION_PLAN = "plan"
const ACTION_REMOVE = "remove"
const ACTION_DESCRIBE_MODULE = "describe_module_for_exeiac"

type Event struct {
	Type       string                  `json:"type"`
	StatusCode extools.NumbersSequence `json:"status_code"`
	Path       string                  `json:"path" `
}

type Action struct {
	Behaviour      string                  `json:"behaviour"`
	StatusCodeFail extools.NumbersSequence `json:"status_code_fail"`
	Events         map[string]Event        `json:"events"`
}

type Module struct {
	Name    string
	Path    string
	Actions map[string]Action
}

func (m Module) String() string {

	j, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	return string(j)
}

func (m Module) ListActions() (actions []string) {
	for name := range m.Actions {
		actions = append(actions, name)
	}
	return
}

// Executes the describe_module_for_exeiac command on a module to get
// the available actions from the module. Then parses and saves them in the
// *Actions* slice.
// If the *Actions* slice is not empty, bypass the call of the command.
// NOTE(half-shell): Do we have a use to force the call to be triggered again here?
//
//	ANSWER(arthur91f): I don't think so, you shouldn't change your module during a run
func (module *Module) LoadAvailableActions() (err error) {
	// Actions are already loaded; no need to reprocess it
	if len(module.Actions) > 0 {
		return
	}

	path, err := exec.LookPath(module.Path)
	if err != nil {
		return fmt.Errorf("unable to load available actions for module %s: %v", module.Name, err)
	}

	stdout := StoreStdout{}
	cmd := exec.Cmd{
		Path:   path,
		Args:   []string{path, ACTION_DESCRIBE_MODULE},
		Stdout: &stdout,
		Stderr: os.Stderr,
	}

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("unable to load module descriptions for module %s: %v", module.Name, err)
	}

	actions := map[string]Action{}
	err = json.Unmarshal(stdout.Output, &actions)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return fmt.Errorf("unable to load module descriptions for module %s: %v", module.Name, err)
	}
	for actionName, action := range actions {
		if action.Behaviour == "" {
			action.Behaviour = "standard"
		}
		if action.StatusCodeFail == "" {
			action.StatusCodeFail = "1-255"
		} else {
			if !action.StatusCodeFail.IsValid() {
				return fmt.Errorf("Error invalid sequence format: module %s %s in $.%s.status_code_fail: \"%s\"",
					module.Name, ACTION_DESCRIBE_MODULE, actionName, action.StatusCodeFail)
			}
		}
		if action.Events == nil {
			action.Events = map[string]Event{}
		} else {
			for k, v := range action.Events {
				if v.Type == "" {
					return fmt.Errorf("module error: event without type: module: %s, action: %s, event: %s",
						module.Name, actionName, k)
				}
				if v.Type != "status_code" && v.Type != "file" && v.Type != "json" && v.Type != "yaml" {
					return fmt.Errorf("module error: event with bad type (should be 'status_code' or 'file' or json or yaml): module: %s, action: %s, event: %s",
						module.Name, actionName, k)
				}
				if v.StatusCode == "" && v.Path == "" {
					return fmt.Errorf("module error: event without path nor status_code: module: %s, action: %s, event: %s",
						module.Name, actionName, k)
				} else if v.StatusCode != "" {
					if !v.StatusCode.IsValid() {
						return fmt.Errorf("Error invalid sequence format: module %s %s in $.%s.events.%s: \"%s\"",
							module.Name, ACTION_DESCRIBE_MODULE, actionName, k, v.StatusCode)
					}
				}
			}
		}

		module.Actions[actionName] = action
	}

	return
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

// REFACTO: should be a brick method or take the brick path as argument
// REFACTO: should take default result for plan
func (m *Module) getExecutionEvents(
	brickPath string,
	action string,
	statusCode int,
) (
	events map[string]interface{},
	execErr error,
) {
	var boolean bool
	var err error

	events = make(map[string]interface{})

	boolean, err = m.Actions[action].StatusCodeFail.Contains(statusCode)
	if err != nil {
		execErr = fmt.Errorf("%s %s bad status_code_fail: %v",
			m.Name, action, err)
	}
	if boolean {
		execErr = fmt.Errorf("%s %s execution has failed with status code: %d",
			m.Name, action, statusCode)
	}

	for eventName, event := range m.Actions[action].Events {
		if event.Type == "status_code" {
			boolean, err = event.StatusCode.Contains(statusCode)
			if err != nil {
				execErr = fmt.Errorf("%s %s event %s bad status_code: %v",
					m.Name, action, eventName, err)
			}
			if boolean {
				events[eventName] = true
			} else {
				events[eventName] = false
			}
		} else if event.Type == "file" || event.Type == "json" || event.Type == "yaml" {
			var path string
			if strings.HasPrefix(event.Path, "/") {
				path = event.Path
			} else {
				path = brickPath + "/" + event.Path
			}
			events[eventName], err = os.ReadFile(path)
			if err != nil {
				if event.Path == "" { // could be catched before but may be it's simplier to do it it here
					err = fmt.Errorf("%s event %s of type %s doesn't specify a path",
						m.Name, eventName, event.Type)
				}
				if execErr != nil {
					execErr = fmt.Errorf("%v\n%v", execErr, err)
				} else {
					execErr = err
				}
			}
		} else {
			err := fmt.Errorf("module error: event with bad type (should be 'status_code' or 'file'): module: %s, action: %s, event: %s",
				m.Name, action, eventName)
			if execErr != nil {
				execErr = fmt.Errorf("%v\n%v", execErr, err)
			} else {
				execErr = err
			}
		}
	}

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
	confEnv []string,
	writers ...io.Writer,
) (
	events map[string]interface{},
	err error,
) {
	var statusCode int

	if !extools.ContainsString(m.ListActions(), action) {
		err = ActionNotImplementedError{Action: action, Module: m}

		return
	}

	// set envs vars
	env := append(os.Environ(), []string{
		fmt.Sprintf("EXEIAC_BRICK_PATH=%s", b.Path),
		fmt.Sprintf("EXEIAC_BRICK_NAME=%s", b.Name),
		fmt.Sprintf("EXEIAC_ROOM_PATH=%s", b.Room.Path),
		fmt.Sprintf("EXEIAC_ROOM_NAME=%s", b.Room.Name),
		fmt.Sprintf("EXEIAC_MODULE_PATH=%s", m.Path),
		fmt.Sprintf("EXEIAC_MODULE_NAME=%s", m.Name),
	}...)
	if len(confEnv) != 0 {
		env = append(env, confEnv...)
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

	events, err = m.getExecutionEvents(b.Path, action, statusCode)

	return
}
