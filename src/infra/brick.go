package infra

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

type Dependency struct {
	// A reference to the related brick
	Brick *Brick
	// The name the variable is supposed to take
	VarName string
	// The JSON path to access the variable
	JsonPath string
}

type Brick struct {
	// The brick's name. Usually the name of the parent directory
	Name string
	// The absolute path of the brick's directory
	Path string
	// The absolute path of the `brick.yml` file
	ConfigurationFilePath string
	// Wheither or not the brick contains a `brick.yml` file.
	// Meaning it does not contain any other brick.
	IsElementary bool
	// A pointer to a module
	Module *Module
	// Pointer to the bricks it depends on
	Dependencies []Dependency
}

type BrickConfYaml struct {
	// The configuration file's format version
	Version string `yaml:"version"`
	// The name or path of the module it uses
	Module string `yaml:"module"`
	// A slice of different kinds of input needed for this brick
	// It **usually** matches the plain or processed output of another brick
	Input []struct {
		// The type of input this brick is expecting
		// Can match the strings "env_file" or "env_var"
		Type string `yaml:"type"`
		// If the type is a path, it is the path the dependency output should be saved to
		Path string `yaml:"path"`
		Data []struct {
			// The name the variable is expected to have
			Name string `yaml:"name"`
			// The key path the input variable should match
			// It is of the form "<brick_name>:"<json_path>"
			// OR "<brick_path>:"<json_path>"
			// e.g. "super-brick/brick:.object.field
			From string `yaml:"from"`
		}
	} `yaml:"input"`
}

// Reads a the yaml configuration file and parses it.
// Returns the parsed configuration.
// Returns an error because of reading or parsing the file.
func (bcy BrickConfYaml) New(path string) (BrickConfYaml, error) {
	in, err := os.ReadFile(path)
	if err != nil {
		return BrickConfYaml{}, err
	}

	err = yaml.Unmarshal(in, &bcy)
	if err != nil {
		return BrickConfYaml{}, err
	}

	return bcy, nil
}

func (brick Brick) String() string {
	return fmt.Sprintf("Name: %s, Path: %s, ConfigurationFilePath: %s, IsElementary: %t, Module: %v",
		brick.Name,
		brick.Path,
		brick.ConfigurationFilePath,
		brick.IsElementary,
		brick.Module)
}

func (brick *Brick) SetElementary(cfp string) *Brick {
	brick.IsElementary = true
	brick.ConfigurationFilePath = cfp

	return brick
}

// Processes the relevant parts of a brick's configuration and updates the brick itself
// with it.
func (brick *Brick) Enrich(bcy BrickConfYaml, infra *Infra) (*Brick, error) {
	if !brick.IsElementary {
		return brick, errors.New("Cannot enrich a non-elementary brick")
	}

	module, err := GetModule(bcy.Module, &infra.Modules)
	if err != nil {
		return brick, err
	}

	brick.Module = module
	dependencies, err := bcy.getDependencies(infra)
	if err != nil {
		log.Println("An error occured when getting dependencies: ", err)
	}
	brick.Dependencies = dependencies

	return brick, nil
}

func (bcy BrickConfYaml) getDependencies(infra *Infra) ([]Dependency, error) {
	var dependencies []Dependency
	parseFromField := func(from string) (brickName string, dataKey string) {
		fields := strings.Split(from, ":")

		return fields[0], fields[1]
	}

	for _, i := range bcy.Input {
		for _, d := range i.Data {
			brickName, jsonPath := parseFromField(d.From)
			// NOTE(half-shell): We sanitize the brick name here in case they turn
			// out to be brick paths
			brick, _ := GetBrick(sanitizeBrickName(brickName), &infra.Bricks)

			dependencies = append(dependencies, Dependency{
				Brick:    brick,
				VarName:  d.Name,
				JsonPath: jsonPath,
			})
		}
	}

	return dependencies, nil
}
