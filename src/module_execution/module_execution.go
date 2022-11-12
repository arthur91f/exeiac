package module_execution

import (
	"log"
	"os"
	exec "os/exec"
	exeiac "src/exeiac/arguments"
)

// Executes a provided module with the provided arguments and redirects stdin to the modules' one.
// It returns a pointer to the Cmd, and an error if any.
// Considering exeiac relies on a module's exit codes to provide additionnal information, this function can be used the following way to properly handle those:
//
//	cmd, err := module_execution.ExecInteractive(module, []string{"-h"})
//	exitError, hasExitError := err.(*exec.ExitError)
//
//	if hasExitError {
//		exitCode := exitError.ExitCode()
//
//		switch exitCode {
//		default:
//			log.Println("This is another kind of exit code")
//		case 1:
//			log.Fatal("There's been an actual error")
//		case 2:
//			log.Println("Careful now, it seems a drift happened")
//		}
//	} else if err != nil {
//		log.Fatal(err)
//	}
func ExecInteractive(module *exeiac.NamePathMapping, args []string) (*exec.Cmd, error) {
	path, err := exec.LookPath(module.Path)

	// We want to panic if the module has an invalid path
	if err != nil {
		log.Fatal(err)
	}

	command := &exec.Cmd{
		Path:   path,
		Args:   args,
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	// We can wait for it to run since this'll most likely be run synchronously
	err = command.Run()

	return command, err
}
