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
)

// A slice of several Brick.
type Bricks []*Brick

// Allows for sorting over Bricks
func (slice Bricks) Len() int {
	return len(slice)
}

// Allows for sorting over Bricks
func (slice Bricks) Less(i, j int) bool {
	return slice[i].Index > slice[j].Index
}

// Allows for sorting over Bricks
func (slice Bricks) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

type Infra struct {
	Modules []Module
	Bricks  map[string]*Brick
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

type ErrBrickNotFound struct {
	brick string
}

func (e ErrBrickNotFound) Error() string {
	return fmt.Sprintf("Brick not found: %s", e.brick)
}

func CreateInfra(rooms []extools.NamePathBinding, modules []extools.NamePathBinding) (Infra, error) {
	i := Infra{
		Bricks:  make(map[string]*Brick),
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

func ConvertToName(path string) string {
	return ""
}

func (i Infra) GetBrickIndexWithPath(brickPath string) (int, error) {
	if brick, ok := i.Bricks[ConvertToName(brickPath)]; ok {
		return brick.Index, nil
	}

	return -1, ErrBrickNotFound{brick: brickPath}
}

func (i Infra) GetBrickIndexWithName(brickName string) (int, error) {
	if brick, ok := i.Bricks[brickName]; ok {
		return brick.Index, nil
	}

	return -1, ErrBrickNotFound{brick: brickName}
}

func (i Infra) GetSubBricksIndexes(brickName string) (bricks []Brick) {
	// the infra.Bricks is sorted with super bricks
	// directly before their subbricks
	superBrickPath := i.Bricks[brickName].Path
	for n, b := range i.Bricks {
		if strings.HasPrefix(i.Bricks[n].Path, superBrickPath) {
			bricks = append(bricks, *b)
		}
	}

	return bricks
}

func (i *Infra) GetSubBricks(brick *Brick) []*Brick {
	subBricks := []*Brick{}

	// the infra.Bricks is sorted with super bricks
	// directly before their subbricks
	superBrickPath := i.Bricks[brick.Name].Path
	for _, b := range i.Bricks {
		brickPath := b.Path
		// We ignore the sub-brick if they're the same; strings.HasPrefix returns true in that case
		if brickPath != superBrickPath && strings.HasPrefix(brickPath, superBrickPath) {
			subBricks = append(subBricks, b)
		}
	}

	return subBricks
}

// TODO(half-shell): Can use a generic argument and be merged with
// `GetBrick`
func GetModule(name string, modules *[]Module) (*Module, error) {
	for i, m := range *modules {
		if m.Name == name {
			return &(*modules)[i], nil
		}
	}

	return nil, errors.New("No matching module name")
}

func GetBrick(name string, bricks *[]Brick) (*Brick, error) {
	for i, b := range *bricks {
		if b.Name == name {
			return &(*bricks)[i], nil
		}
	}

	return nil, errors.New("No matching module name")
}
