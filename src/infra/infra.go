package infra

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"regexp"
	extools "src/exeiac/tools"
	"strings"
	"sync"
)

type Infra struct {
	Modules []Module
	Bricks  map[string]*Brick
}

func CreateInfra(rooms []extools.NamePathBinding, modules []extools.NamePathBinding) (Infra, error) {
	i := Infra{
		Bricks: make(map[string]*Brick),
	}

	// create Modules
	for _, m := range modules {
		i.Modules = append(i.Modules, Module{
			Name: m.Name,
			Path: m.Path,
		})
	}

	// Temporary brick storage to have consistent indexing across rooms
	// i.e. having an overall ordering as the index
	bricks := []Brick{}
	for _, r := range rooms {
		b, err := GetBricks(r)
		if err != nil {
			fmt.Printf("%v\n> Warning63724ff3:infra/CreateInfra:"+
				"can't add bricks of this room: %s", err, r.Path)
		}

		bricks = append(bricks, b...)
	}

	for idx := range bricks {
		// We do not want to access &b because it uses the same pointer to copy the new element in
		// meaning we'd reference the same entity over and over
		// c.f. https://stackoverflow.com/questions/20185511/range-references-instead-values
		b := bricks[idx]
		b.Index = idx
		i.Bricks[b.Name] = &b
	}

	return i, nil
}

var hasDigitPrefixRegexp = regexp.MustCompile(`.*/\d+-\w+$`) // TODO(arthur91f): change regex
// .../1-database2-eu_confAdmin is a valid brick dirname
var prefixRegexp = regexp.MustCompile(`\d+-`) // TODO(arthur91f): replace regex
// `^\d+-` or `/\d+-` but we want to change
// OK  rooms/1-database2-eu -> rooms/database2-eu
// NOK rooms/1-database2-eu -> rooms/databaseeu

func validateDirName(path string) bool {
	return hasDigitPrefixRegexp.MatchString(path)
}

func SanitizeBrickName(name string) string {
	return prefixRegexp.ReplaceAllString(name, "")
}

// Walks the file system from the provided root, gathers all folders containing a `brick.html` file, and build a Brick struct from it.
func GetBricks(room extools.NamePathBinding) ([]Brick, error) {
	var bricks []Brick
	var err error

	err = filepath.WalkDir(
		room.Path,
		func(path string, d fs.DirEntry, err error) error {
			brickRelPath, err := filepath.Rel(room.Path, path)
			if err != nil {
				log.Fatal(err)
			}

			lastBrick := func() *Brick {
				if len(bricks) > 0 {
					return &(bricks)[len(bricks)-1]
				}

				return &Brick{}
			}()

			// A brick can just be described as a sub-path of a room, containing a prefixed folder name with digits, and split with a hypen ("-")
			if d.Type().IsDir() && validateDirName(path) {
				brickName := filepath.Join(room.Name, brickRelPath)
				name := SanitizeBrickName(brickName)

				// Do not duplicate entries
				if len(bricks) == 0 || lastBrick.Name != name {
					bricks = append(bricks, Brick{
						Name:         name,
						Path:         path,
						IsElementary: false,
					})
				}
			}

			// An elementary brick has prefixed folder name, and a brick.yml file.
			// TODO(half-shell): Make the configuration filename more flexible.
			if d.Type().IsRegular() && d.Name() == "brick.yml" {
				brickName := filepath.Join(room.Name, filepath.Dir(brickRelPath))
				name := SanitizeBrickName(brickName)

				// Set the last brick as elementary if names match
				// This happens because it means that the parent brick is not a "super-brick"
				// but an elementary brick
				if lastBrick.Name == name {
					lastBrick.SetElementary(path)
				}
			}

			return err
		})

	return bricks, err
}

func (infra Infra) String() string {
	var modulesString string
	var bricksString string

	if len(infra.Modules) > 0 {
		for _, m := range infra.Modules {
			modulesString = fmt.Sprintf("%s%s", modulesString,
				extools.IndentForListItem(m.String()))
		}
		modulesString = fmt.Sprintf("modules:\n%s", modulesString)
	} else {
		modulesString = "modules: []\n"
	}

	if len(infra.Bricks) > 0 {
		for _, b := range infra.Bricks {
			bricksString = fmt.Sprintf("%s%s", bricksString,
				extools.IndentForListItem(b.String()))
		}
		bricksString = fmt.Sprintf("bricks:\n%s", bricksString)
	} else {
		bricksString = "bricks: []\n"
	}

	return fmt.Sprintf("infra:\n%s%s",
		extools.Indent(modulesString),
		extools.Indent(bricksString),
	)
}

func GetModule(name string, modules *[]Module) (*Module, error) {
	for i, m := range *modules {
		if m.Name == name {
			return &(*modules)[i], nil
		}
	}

	return nil, errors.New("No matching module name")
}

// Resolve a brick to a slice of elementary bricks.
// If the provided brick is an elementary one, we just return a slice
// of length 1 with that brick.
func (i *Infra) GetSubBricks(brick *Brick) (subBricks Bricks, err error) {
	// the infra.Bricks is sorted with super bricks
	// directly before their subbricks
	superBrickPath := i.Bricks[brick.Name].Path
	for _, b := range i.Bricks {
		brickPath := b.Path
		// We ignore the sub-brick if they're the same; strings.HasPrefix returns true in that case
		if b.IsElementary && strings.HasPrefix(brickPath, superBrickPath) {
			if b.EnrichError != nil {
				err = b.EnrichError
				return
			}

			subBricks = append(subBricks, i.Bricks[b.Name])
		}
	}

	return
}

// Checks wheither or not a brick's dependencies are enriched.
// Returns a brick's dependency if they are, or an error otherwise.
func (infra *Infra) GetDirectPrevious(brick *Brick) (results Bricks, err error) {
	deps := brick.Dependencies
	for _, d := range deps {
		err = infra.Bricks[d.BrickName].EnrichError
		if infra.Bricks[d.BrickName].EnrichError != nil {
			return
		}

		var elementaryBricks Bricks
		elementaryBricks, err = infra.GetSubBricks(infra.Bricks[d.BrickName])
		if err != nil {
			return
		}

		results = append(results, elementaryBricks...)
	}

	return
}

func (infra *Infra) GetLinkedPrevious(brick *Brick) (results Bricks, err error) {
	var directPrevious Bricks
	added, err := infra.GetDirectPrevious(brick)

	results = added
	for {
		toAdd := Bricks{}
		for _, b := range added {
			directPrevious, err = infra.GetDirectPrevious(b)
			for _, dp := range directPrevious {
				if !toAdd.BricksContains(dp) {
					toAdd = append(toAdd, dp)
				}
			}
		}
		if len(toAdd) == 0 {
			break
		}
		results = append(results, toAdd...)
		added = toAdd
	}
	return
}

// Checks wheither or not a brick's dependents are enriched.
// Returns a slice of brick's pointers if they are, or the first error encountered otherwise
func (infra *Infra) GetDirectNext(brick *Brick) (results Bricks, err error) {
	// browse all infra.Bricks and return bricks that have brickName in dependencies
	for _, b := range infra.Bricks {
		for _, d := range b.Dependencies {
			err = infra.Bricks[b.Name].EnrichError
			if err != nil {
				return
			}

			if d.BrickName == brick.Name {
				results = append(results, infra.Bricks[b.Name])
			}
		}
	}

	return
}

// Cheks wheither or not a brick's dependents (and their dependents and so on) are enriched.
// Returns a slice of brick's pointers if they are, or the first error encountered otherwise
func (infra *Infra) GetLinkedNext(bricks *Bricks, brick *Brick) (err error) {
	var wg sync.WaitGroup

	for _, b := range infra.Bricks {
		for _, d := range b.Dependencies {
			err = infra.Bricks[b.Name].EnrichError
			if err != nil {
				return
			}

			if d.BrickName == brick.Name {
				*bricks = append(*bricks, infra.Bricks[b.Name])

				// TODO(half-shell): handle cases of cirulare dependencies
				// or misconfiguration. We want to add a check to either limit
				// the number of wait group created, or just a check over the indexes
				// of the bricks handled here. Or prevent adding a brick already added
				// to avoid those cases.
				if len(infra.Bricks[b.Name].Dependencies) > 0 {
					wg.Add(1)

					go func(bricks *Bricks, brick *Brick) {
						defer wg.Done()
						infra.GetLinkedNext(bricks, brick)
					}(bricks, b)
				}
			}
		}
	}

	wg.Wait()

	return
}
