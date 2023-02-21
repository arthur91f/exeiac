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

func enrichDatas(bricksToExecute exinfra.Bricks, infra *exinfra.Infra) error {
	// find all bricks that we need to ask output
	var neededBricksForTheirOutputs exinfra.Bricks
	for _, b := range bricksToExecute {
		/* we can assume it's true if it's the bricksToExecute from main
		if b.EnrichError != nil {
			return b.EnrichError
		}*/
		bricks, err := infra.GetCorrespondingBricks(exinfra.Bricks{b}, []string{"selected", "linked_previous"})
		if err != nil {
			return err
		}
		neededBricksForTheirOutputs = append(neededBricksForTheirOutputs, bricks...)
	}
	sort.Sort(neededBricksForTheirOutputs)
	neededBricksForTheirOutputs = exinfra.RemoveDuplicates(neededBricksForTheirOutputs)

	// check we don't have any enrich error on brick we will execute output
	for _, b := range neededBricksForTheirOutputs {
		if b.EnrichError != nil {
			return b.EnrichError
		}

		envs, err := writeEnvFilesAndGetEnvs(b)
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

		b.Output = stdout.Output
	}

	return nil
}

func writeEnvFilesAndGetEnvs(brick *exinfra.Brick) (envs []string, err error) {

	formatters, envFormatter, err := brick.CreateFormatters()
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
	formatters, _, err := brick.CreateFormatters()
	if err != nil {
		return fmt.Errorf("error when searching input files to delete of brick %s: %v", brick.Name, err)
	}

	if len(formatters) > 0 {
		for path := range formatters {
			if _, err := os.Stat(path); err == nil {
				err = os.Remove(path)
				if err != nil {
					return fmt.Errorf("error when deleting %s an input files of brick %s: %v", path, brick.Name, err)
				}
			}
		}
	}

	return err
}
