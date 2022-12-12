package actions

import (
	"os"
	"sort"
	exargs "src/exeiac/arguments"
	exinfra "src/exeiac/infra"
)

var BehaviourMap = map[string]func(*exinfra.Infra, *exargs.Arguments, exinfra.Bricks) (int, error){
	"clean":         Clean,
	"help":          Help,
	"init":          Default,
	"lay":           Lay,
	"plan":          Plan,
	"remove":        Remove,
	"show":          Show,
	"validate_code": Default,
	"debug_args":    DebugArgs,
	"debug_infra":   DebugInfra,
	"default":       Default,
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

		formatters, err := b.CreateFormatters()
		if err != nil {
			return err
		}

		if len(formatters) > 0 {
			for path, formatter := range formatters {
				f, err := os.Create(path)
				if err != nil {
					return err
				}

				data, err := formatter.Format()
				_, err = f.Write(data)
				if err != nil {
					return err
				}
			}
		}

		stdout := exinfra.StoreStdout{}
		exitError, err := b.Module.Exec(b, "output", []string{}, &stdout)
		if err != nil {
			return err
		}

		if exitError != nil {
			if exitError.ExitCode() != 0 {
				return exitError
			}
		}
		b.Output = stdout.Output
	}

	return nil
}
