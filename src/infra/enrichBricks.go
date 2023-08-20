package infra

import (
	"fmt"
	"os"
	extools "src/exeiac/tools"
	"strings"

	"github.com/PaesslerAG/jsonpath"
	"gopkg.in/yaml.v2"
)

type BrickConfYaml struct {
	// The configuration file's format version
	Version string `yaml:"version"`
	// The name or path of the module it uses
	Module string `yaml:"module"`

	// The dependencies
	Dependencies map[string]struct {
		From             string   `yaml:"from"`
		NeededFor        []string `yaml:"needed_for"`
		NotNeededFor     []string `yaml:"not_needed_for"`
		TriggerType      string   `yaml:"trigger_type"`
		TriggeredActions []string `yaml:"triggered_actions"`
	} `yaml:"dependencies"`

	// A slice of different kinds of input needed for this brick
	Inputs []struct {
		Type   string   `yaml:"type"`   // can be file or env_vars
		Format string   `yaml:"format"` // can be json, yaml or env
		Path   string   `yaml:"path"`   // used for type=file let empty else
		Datas  []string `yaml:"datas"`  // reference to Dependencies keys
	} `yaml:"inputs"`
}

func (infra *Infra) EnrichBrick(brick *Brick) error {

	// READ brick.yml FILE ----------------------------------------------------
	in, err := os.ReadFile(brick.ConfigurationFilePath)
	if err != nil {
		return fmt.Errorf("unable to read %s: %v",
			brick.ConfigurationFilePath, err)
	}

	// LOAD brick.yml FILE ----------------------------------------------------
	var c BrickConfYaml
	err = yaml.Unmarshal(in, &c)
	if err != nil {
		return err
	}

	// SET DEFAULTS VALUE FOR brick.yml AND CATCH ERRORS ----------------------
	// -- set version --
	if c.Version == "" {
		c.Version = "1.0.0"
	}

	// -- set module --
	if c.Module == "" {
		return fmt.Errorf("$.module required")
	}

	// -- set dependencies --
	for k, d := range c.Dependencies {
		if d.From == "" {
			return fmt.Errorf("$.dependencies[\"%s\"].from required", k)
		}
		if d.TriggerType == "" {
			d.TriggerType = "classic"
		}
		if d.TriggerType == "weak" &&
			len(d.TriggeredActions) > 0 &&
			d.TriggeredActions[0] == "lay" {
			fmt.Fprintf(os.Stderr, "WARNING: $.dependencies[\"%s\"].trigger_action"+
				"specified and not [] although trigger_type is weak. "+
				"consider that trigger_action is []", k)
			d.TriggeredActions = []string{}
		}
		if d.TriggerType != "weak" && len(d.TriggeredActions) == 0 {
			d.TriggeredActions = []string{"lay"}
		}

		c.Dependencies[k] = d
	}

	// -- set inputs --
	for n, i := range c.Inputs {
		if i.Type == "" && i.Path == "" {
			i.Type = "env_vars"
		}

		if i.Type == "" {
			return fmt.Errorf("$.inputs[\"%d\"].type required", n)
		} else if i.Type != "file" && i.Type != "env_vars" {
			return fmt.Errorf("$.inputs[\"%d\"].type should be \"env_vars\" or \"file\"", n)
		}
		if i.Format == "" {
			i.Format = "env"
		}
		if i.Path == "" && i.Type != "env_vars" {
			return fmt.Errorf("$.inputs[\"%d\"].path required", n)
		}
		if len(i.Datas) == 0 {
			return fmt.Errorf("$.inputs[\"%d\"].datas can't be empty", n)
		}
		c.Inputs[n] = i
	}

	// SET BRICK --------------------------------------------------------------

	// -- set module --
	module, err := infra.GetModule(c.Module, brick)
	if err != nil {
		return err
	}
	brick.Module = module

	// -- define local functions --
	getExceptionNeededFor := func(needed []string, notNeeded []string) (e []string) {
		if infra.Conf.DefaultIsInputNeeded {
			e = append(infra.Conf.ExceptionIsInputNeeded, notNeeded...)
		} else {
			e = append(infra.Conf.ExceptionIsInputNeeded, needed...)
		}
		extools.Deduplicate(e)
		return e
	}
	parseFromField := func(from string) (brickName string, source string, dataKey string, err error) {
		if from == "" {
			err = fmt.Errorf("field from is empty or doesn't exist")

			return
		}
		fields := strings.Split(from, ":")
		if len(fields) != 3 {
			err = fmt.Errorf("field from \"%s\" is not of from: <brick name>:"+
				"<source>:<json path>", from)

			return
		}
		if fields[1] != "output" && fields[1] != "event" {
			err = fmt.Errorf("field from \"%s\" is not valid: "+
				"source part should be \"output\" or \"event\"", from)

			return
		}

		_, err = jsonpath.New(fields[2])
		if err != nil {
			err = fmt.Errorf("field from \"%s\" is not of from: <brick name>:"+
				"<source>:<json path>: json path unvalid: %v", from, err)

			return
		}

		return fields[0], fields[1], fields[2], nil
	}

	// -- set dependencies --
	brick.Dependencies = make(map[string]*Dependency)
	for k, d := range c.Dependencies {
		var fromBrick *Brick

		fromBrickName, fromSource, fromJsonpath, err := parseFromField(d.From)
		if err != nil {
			return fmt.Errorf("dependencies.%s.from unvalid: %v", k, err)
		}

		fromBrick, ok := infra.Bricks[SanitizeBrickName(fromBrickName)]
		if !ok {
			return fmt.Errorf("dependencies.%s.from unvalid: %s doesn't correspond to existing brick name",
				k, fromBrickName)
		}

		exception := getExceptionNeededFor(c.Dependencies[k].NeededFor, c.Dependencies[k].NotNeededFor)

		brick.Dependencies[k] = &Dependency{
			From: From{
				Brick:    fromBrick,
				Source:   fromSource,
				JsonPath: fromJsonpath,
			},
			TriggeredAction:    d.TriggeredActions,
			TriggerType:        d.TriggerType,
			DefaultNeededFor:   infra.Conf.DefaultIsInputNeeded,
			ExceptionNeededFor: exception,
			Value:              nil,
		}
	}

	// -- set inputs --
	for _, i := range c.Inputs {
		for _, d := range i.Datas {
			_, ok := c.Dependencies[d]
			if !ok {
				return fmt.Errorf("%s in inputs doesn't correspond to any dependencies", d)
			}
			inputFormat, ok := SupportedFormats[i.Format]
			if !ok {
				return fmt.Errorf("input format not supported: %s", i.Format)
			}
			brick.Inputs = append(brick.Inputs, Input{
				VarName:    d,
				Dependency: brick.Dependencies[d],
				Type:       i.Type,
				Format:     inputFormat,
				Path:       i.Path,
			})
		}
	}

	// -- set direct previous brick --
	for _, d := range brick.Dependencies {
		if !brick.DirectPrevious.BricksContains(d.From.Brick) {
			brick.DirectPrevious = append(brick.DirectPrevious, d.From.Brick)
		}
	}

	// RETURN NIL WHEN EVERYTHING GOES WELL -----------------------------------
	return nil
}

func (infra *Infra) EnrichBricks() {
	for _, b := range infra.Bricks {
		if b.IsElementary {
			err := infra.EnrichBrick(b)
			if err != nil {
				infra.Bricks[b.Name].EnrichError =
					fmt.Errorf("unable to enrich brick(%s): %v", b.Name, err)
			}

			err = b.Module.LoadAvailableActions()
			if err != nil {
				infra.Bricks[b.Name].EnrichError =
					fmt.Errorf("unable to enrich brick(%s): %v", b.Name, err)
			}
		}
	}
}
