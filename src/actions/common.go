package actions

import (
	"fmt"
	"os"
	"sort"
	"strings"

	exargs "src/exeiac/arguments"
	exinfra "src/exeiac/infra"
	extools "src/exeiac/tools"

	"github.com/fatih/color"
)

var BehaviourMap = map[string]func(*exinfra.Infra, *exargs.Configuration, exinfra.Bricks) (int, error){
	"clean":         Clean,
	"help":          Help,
	"init":          PassthroughAction,
	"lay":           Lay,
	"smart-lay":     SmartLay,
	"plan":          Plan,
	"remove":        Remove,
	"show":          Show,
	"validate_code": PassthroughAction,
	"debug_args":    DebugArgs,
	"debug_infra":   DebugInfra,
	"default":       PassthroughAction,
	"get-depends":   GetDepends,
}

const TAG_OK = "OK"
const TAG_NO_CHANGE = "NO_CHANGE"
const TAG_DONE = "DONE"
const TAG_ERROR = "ERR"
const TAG_SKIP = "SKIP"
const TAG_DRIFT = "DRIFT"
const TAG_MAY_DRIFT = "DRIFT?"

type ExecSummary []ExecReport

type ExecReport struct {
	Brick  *exinfra.Brick
	Status string // "" red"ERR" blue"SKIP" green"OK" cyan"DONE" cyan"DRIFT"
	Error  error  //
}

func (es ExecSummary) Display() {
	var sb strings.Builder

	sb.WriteString(color.New(color.Bold).Sprint("Summary:\n"))
	for _, report := range es {
		var str string
		if report.Error != nil {
			str = fmt.Sprintf("%s : %s", report.Brick.Name,
				extools.IndentIfMultiline(report.Error.Error()))
		} else {
			str = report.Brick.Name
		}
		switch report.Status {
		case TAG_ERROR:
			sb.WriteString(color.RedString("ERR     "))
		case TAG_SKIP:
			sb.WriteString(color.BlueString("SKIP    "))
		case TAG_OK:
			sb.WriteString(color.GreenString("OK      "))
		case TAG_NO_CHANGE:
			sb.WriteString(color.GreenString("OK      "))
		case TAG_DONE:
			sb.WriteString(color.CyanString("DONE    "))
		case TAG_DRIFT:
			sb.WriteString(color.CyanString("DRIFT   "))
		case TAG_MAY_DRIFT:
			sb.WriteString(color.CyanString("?DRIFT? "))
		case "":
			sb.WriteString(color.RedString("NO FLAG  "))
		default:
			sb.WriteString(color.YellowString(report.Status))
		}
		sb.WriteString(fmt.Sprintf("%s\n", str))
	}

	fmt.Print(sb.String())
}

func (es ExecSummary) String() string {
	var sb strings.Builder

	sb.WriteString("Summary:\n")
	for _, report := range es {
		if report.Error != nil {
			sb.WriteString("Failed")
		} else {
			sb.WriteString("Succeess")
		}
		sb.WriteString(fmt.Sprintf(" %s", report.Brick.Name))
		sb.WriteString("\n")
	}

	return sb.String()
}

func getBricksToOutput(
	bricksToExecute exinfra.Bricks,
	infra *exinfra.Infra,
	action string,
) (
	bricksToOutput exinfra.Bricks,
	err error,
) {
	// 1. search needed inputs corresponding bricks
	var neededBricksForOutput exinfra.Bricks
	for _, b := range bricksToExecute {
		if b.EnrichError != nil {
			err = fmt.Errorf("%s can't be enriched: %v", b.Name, b.EnrichError)

			return
		}
		for _, i := range b.Inputs {
			if i.IsInputNeeded(action) {
				neededBricksForOutput = append(neededBricksForOutput, i.Brick)
			}
		}
	}
	neededBricksForOutput = exinfra.RemoveDuplicates(neededBricksForOutput)
	bricksToOutput = neededBricksForOutput

	// 2. add recursively all bricks that are needed for their output
	for _, b := range neededBricksForOutput {
		var bricksToAdd exinfra.Bricks
		bricksToAdd, err = infra.GetLinkedPreviousFor(b, "output")
		if err != nil {

			return
		}
		bricksToOutput = append(bricksToOutput, bricksToAdd...)
	}

	// 3. sort
	bricksToOutput = exinfra.RemoveDuplicates(bricksToOutput)
	sort.Sort(bricksToOutput)

	// 4. return
	return
}

func enrichOutputs(bricksToOutput exinfra.Bricks) error {
	// 1. check any enrich error
	for _, b := range bricksToOutput {
		if b.EnrichError != nil {

			return b.EnrichError
		}
	}

	// 2. execute all output and register output in brick
	for _, b := range bricksToOutput {
		// 2.1. create input
		envs, err := writeEnvFilesAndGetEnvs(b, "output")

		// 2.2. execute output
		if err != nil {
			return err
		}

		stdout := exinfra.StoreStdout{}
		statusCode, err := b.Module.Exec(b, "output", []string{}, envs, &stdout)
		if err != nil {
			return err
		}

		if statusCode != 0 {
			return fmt.Errorf("unable to get output of %s", b.Name)
		}

		// 2.3. set brick.Outputs
		b.Output = stdout.Output

	}

	return nil
}

func enrichOutputsBeforeExec(
	brick *exinfra.Brick,
	infra *exinfra.Infra,
	action string,
) error {
	bricksWhomOutputIsNeeded, err := getBricksToOutput(exinfra.Bricks{brick}, infra, action)
	if err != nil {
		return err
	}

	var bricksToOutput exinfra.Bricks
	for _, b := range bricksWhomOutputIsNeeded {
		if b.Output == nil {
			bricksToOutput = append(bricksToOutput, b)
		}
	}
	bricksToOutput = append(bricksWhomOutputIsNeeded, brick) // we can add it at the end because its needed output are before

	err = enrichOutputs(bricksToOutput)
	if err != nil {
		return err
	}

	return nil
}

func writeEnvFilesAndGetEnvs(brick *exinfra.Brick, action string) (envs []string, err error) {

	formatters, envFormatter, err := brick.CreateFormatters(action)
	if err != nil {
		return
	}

	if len(formatters) > 0 {
		for path, formatter := range formatters {
			var f *os.File
			f, err = os.Create(path)
			if err != nil {
				return
			}

			var data []byte
			data, err = formatter.Format()
			_, err = f.Write(data)
			if err != nil {
				return
			}
		}
	}

	envs = envFormatter.Environ()
	return
}

func cleanEnvFiles(brick *exinfra.Brick) error {

	var pathsToClean []string
	for _, i := range brick.Inputs {
		if i.Type == "file" {
			pathsToClean = append(pathsToClean, i.Path)
		}
	}
	pathsToClean = extools.Deduplicate(pathsToClean)

	for _, path := range pathsToClean {
		if _, err := os.Stat(path); err == nil {
			err = os.Remove(path)
			if err != nil {
				return fmt.Errorf("error when deleting %s an input files of brick %s: %v", path, brick.Name, err)
			}
		}
	}

	return nil
}
