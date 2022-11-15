package infra

import (
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"regexp"
	extools "src/exeiac/tools"
)

type Infra struct {
	Modules []Module
	Bricks  []Brick
}

type RoomError struct {
	id     string
	path   string
	reason string
	trace  error
}

func (e RoomError) Error() string {
	return fmt.Sprintf("! Error%s:room: %s: %s\n< %s", e.id,
		e.reason, e.path, e.trace.Error())
}

func (i Infra) New(
	rooms []extools.NamePathBinding,
	modules []extools.NamePathBinding) (Infra, error) {

	// create Modules
	for _, m := range modules {
		i.Modules = append(i.Modules, Module{
			Name: m.Name,
			Path: m.Path,
		})
	}

	// create Bricks
	for _, r := range rooms {
		// get all room's bricks
		err := appendBricks(r, &i.Bricks)
		if err != nil {
			fmt.Println("%v\n> Warning63724ff3:infra/CreateInfra:"+
				"can't add bricks of this room: %s", err, r.Path)
		}
	}

	return i, nil
}

var hasDigitPrefixRegexp = regexp.MustCompile(`.*/\d+-\w+$`)
var prefixRegexp = regexp.MustCompile(`\d+-`)

func validateDirName(path string) bool {
	return hasDigitPrefixRegexp.MatchString(path)
}

func sanitizeBrickName(name string) string {
	return prefixRegexp.ReplaceAllString(name, "")
}

// Walks the file system from the provided root, gathers all folders containing a `brick.html` file, and build a Brick struct from it.
func appendBricks(room extools.NamePathBinding, bricks *[]Brick) error {
	err := filepath.WalkDir(
		room.Path,
		func(path string, d fs.DirEntry, err error) error {
			brickRelPath, err := filepath.Rel(room.Path, path)
			if err != nil {
				log.Fatal(err)
			}

			lastBrick := func() Brick {
				if len(*bricks) > 0 {
					return (*bricks)[len(*bricks)-1]
				}

				return Brick{}
			}()

			// A brick can just be described as a sub-path of a room, containing a prefixed folder name with digits, and split with a hypen ("-")
			if d.Type().IsDir() && validateDirName(path) {
				brickName := filepath.Join(room.Name, brickRelPath)
				name := sanitizeBrickName(brickName)

				// Do not duplicate entries
				if len(*bricks) == 0 || lastBrick.Name != name {
					*bricks = append(*bricks, Brick{
						Name:         name,
						Path:         path,
						IsElementary: false,
					})
				}
			}

			// An elementary brick has prefixed folder name, and a brick.yml file.
			// TODO(half-shell): Make the configuration filename more flexible.
			if d.Type().IsRegular() && d.Name() == "brick.yml" {
				lastBrick := &((*bricks)[len(*bricks)-1])
				brickName := filepath.Join(room.Name, filepath.Dir(brickRelPath))
				name := sanitizeBrickName(brickName)

				// Do not duplicate entries
				if lastBrick.Name == name {
					lastBrick.IsElementary = true
					lastBrick.ConfigurationFilePath = path
				}
			}

			return err
		})

	return err
}

func (infra Infra) Display() {
	fmt.Println("Infra:")

	fmt.Println("  Modules:")
	for _, module := range infra.Modules {
		fmt.Println("  - name: " + module.Name)
		fmt.Println("    path: " + module.Path)
		fmt.Printf("    actions: %v\n", module.Actions)
	}

	fmt.Println("  Bricks:")
	for _, brick := range infra.Bricks {
		fmt.Println("  - name: " + brick.Name)
		fmt.Println("    path: " + brick.Path)
		fmt.Printf("    isElementary: %t\n", brick.IsElementary)
		fmt.Println("    confFile: " + brick.ConfigurationFilePath)
	}
}
