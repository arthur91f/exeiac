package infra

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	exargs "src/exeiac/arguments"
	extools "src/exeiac/tools"
	"strings"
)

const BRICK_FILE_NAME = "brick.yml"

type Infra struct {
	Modules []Module
	Bricks  BricksMap
	Conf    struct {
		ExceptionIsInputNeeded []string
		DefaultIsInputNeeded   bool
	}
}

func CreateInfra(configuration exargs.Configuration) (Infra, error) {
	i := Infra{
		Bricks: make(map[string]*Brick),
	}

	i.Conf.ExceptionIsInputNeeded = configuration.ExceptionIsInputNeeded
	i.Conf.DefaultIsInputNeeded = configuration.DefaultIsInputNeeded

	// create Modules
	for name, path := range configuration.Modules {
		i.Modules = append(i.Modules, Module{
			Name:    name,
			Path:    path,
			Actions: map[string]Action{},
		})
	}

	// Temporary brick storage to have consistent indexing across rooms
	// i.e. having an overall ordering as the index
	bricks := Bricks{}
	for _, room := range configuration.Rooms {
		b, err := GetBricks(room.Name, room.Path)
		if err != nil {
			fmt.Printf("Cannot add \"%s\" (path: %s) room's bricks: %v\n", room.Name, room.Path, err)
		}

		bricks = append(bricks, b...)
	}

	for idx := range bricks {
		b := bricks[idx]
		b.Index = idx
		i.Bricks[b.Name] = b
	}

	return i, nil
}

var hasDigitPrefixRegexp = regexp.MustCompile(`.*/\d+-[^/]+$`)
var prefixRegexp = regexp.MustCompile(`/\d+-`) // rooms/1-database2-eu -> rooms/database2-eu

func validateDirName(path string) bool {
	return hasDigitPrefixRegexp.MatchString(path)
}

func SanitizeBrickName(name string) string {
	return prefixRegexp.ReplaceAllString(name, "/")
}

// Walks the file system from the provided root, gathers all folders containing a `brick.html` file, and build a Brick struct from it.
func GetBricks(roomName string, roomPath string) (bricks Bricks, err error) {
	_, err = os.Stat(roomPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			err = ErrBrickNotFound{brick: roomPath}

			return
		}

		return
	}

	bricks = Bricks{{
		Name:         roomName,
		Path:         roomPath,
		IsElementary: false,
	}}
	bricks[0].Room = bricks[0]

	confFilePath, err := GetConfFilePath(roomPath)
	if err != nil {
		bricks[0].SetElementary(confFilePath)
	}

	err = filepath.WalkDir(
		roomPath,
		func(path string, d fs.DirEntry, err error) error {
			brickRelPath, err := filepath.Rel(roomPath, path)
			if err != nil {
				log.Fatalf("Fatal: can't find relative path %s from room: %s\n  -> %v",
					path, roomPath, err)
			}

			lastBrick := func() *Brick {
				if len(bricks) > 0 {
					return (bricks)[len(bricks)-1]
				}

				return &Brick{}
			}()

			// A brick can just be described as a sub-path of a room, containing a prefixed folder name with digits, and split with a hypen ("-")
			if d.Type().IsDir() && validateDirName(path) {
				brickName := filepath.Join(roomName, brickRelPath)
				name := SanitizeBrickName(brickName)

				// Do not duplicate entries
				if len(bricks) == 0 || lastBrick.Name != name {
					bricks = append(bricks, &Brick{
						Name:         name,
						Path:         path,
						IsElementary: false,
						Room:         bricks[0],
					})
				}
			}

			// An elementary brick has prefixed folder name, and a brick.yml file.
			// TODO(half-shell): Make the configuration filename more flexible.
			if d.Type().IsRegular() && d.Name() == BRICK_FILE_NAME {
				brickName := filepath.Join(roomName, filepath.Dir(brickRelPath))
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

	return
}

func (infra Infra) String() string {
	var sb strings.Builder
	var modulesSb strings.Builder
	var bricksSb strings.Builder

	for _, m := range infra.Modules {
		modulesSb.WriteString(fmt.Sprintf("%s,\n", m))
	}

	bricksSb.WriteString(fmt.Sprintf("%v", infra.Bricks))

	sb.WriteString(fmt.Sprintf("Modules: [\n%s]\nBricks:[\n%s]", modulesSb.String(), bricksSb.String()))

	return sb.String()
}

func (infra *Infra) GetModule(name string, b *Brick) (*Module, error) {
	for i, m := range infra.Modules {
		if m.Name == name {
			return &(infra.Modules[i]), nil
		}
	}
	if strings.HasPrefix(name, "./") {
		infra.Modules = append(infra.Modules, Module{
			Name:    b.Name,
			Path:    b.Path + "/" + strings.TrimPrefix(name, "./"),
			Actions: map[string]Action{},
		})
		return &(infra.Modules[len(infra.Modules)-1:][0]), nil
	}

	return nil, errors.New("No matching module name")
}

// Resolve a brick to a slice of elementary bricks.
// If the provided brick is an elementary one, we just return a slice
// of length 1 with that brick.
func (i *Infra) GetSubBricks(brick *Brick) (subBricks Bricks, err error) {
	// the infra.Bricks is sorted with super bricks
	// directly before their subbricks
	superBrickPath := brick.Path
	for _, b := range i.Bricks {
		// We ignore the sub-brick if they're the same by checking if they have a common prefix
		if b.IsElementary && strings.HasPrefix(b.Path, superBrickPath) {
			if b.EnrichError != nil {
				err = b.EnrichError
				return
			}

			subBricks = append(subBricks, i.Bricks[b.Name])
		}
	}

	return
}

// Return only elementary bricks (although Brick.DirectPrevious can contains super brick)
// return error when a elementary bricks b to return has b.EnrichError != nil
func (infra *Infra) GetDirectPrevious(brick *Brick) (results Bricks, err error) {
	for _, dp := range brick.DirectPrevious {
		err = dp.EnrichError
		if dp.EnrichError != nil {
			return
		}

		if dp.IsElementary {
			results = append(results, dp)
		} else {
			var elementaryBricks Bricks
			elementaryBricks, err = infra.GetSubBricks(dp)
			if err != nil {
				return
			}
			results = append(results, elementaryBricks...)
		}
	}

	return
}

// Return only elementary bricks,
// return error when a elementary bricks b to return has b.EnrichError != nil
func (infra *Infra) GetLinkedPrevious(brick *Brick) (results Bricks, err error) {
	var directPrevious Bricks
	added, err := infra.GetDirectPrevious(brick)
	if err != nil {
		return
	}
	results = added
	for {
		toAdd := Bricks{}
		for _, b := range added {
			directPrevious, err = infra.GetDirectPrevious(b)
			for _, dp := range directPrevious {
				if !results.BricksContains(dp) {
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

func (infra *Infra) GetDirectPreviousFor(
	brick *Brick,
	action string,
) (
	results Bricks, err error,
) {
	var bricksFromInputs Bricks
	for _, i := range brick.Inputs {
		if i.IsInputNeeded(action) {
			bricksFromInputs = append(bricksFromInputs, i.Brick)
		}
	}
	bricksFromInputs = RemoveDuplicates(bricksFromInputs)

	for _, dp := range bricksFromInputs {
		err = dp.EnrichError
		if dp.EnrichError != nil {
			return
		}

		if dp.IsElementary {
			results = append(results, dp)
		} else {
			var elementaryBricks Bricks
			elementaryBricks, err = infra.GetSubBricks(dp)
			if err != nil {
				return
			}
			results = append(results, elementaryBricks...)
		}
	}

	return
}

func (infra *Infra) GetLinkedPreviousFor(
	brick *Brick,
	action string,
) (
	results Bricks, err error,
) {
	var directPrevious Bricks
	added, err := infra.GetDirectPreviousFor(brick, action)
	if err != nil {
		return
	}
	results = added
	for {
		toAdd := Bricks{}
		for _, b := range added {
			directPrevious, err = infra.GetDirectPreviousFor(b, action)
			for _, dp := range directPrevious {
				if !results.BricksContains(dp) {
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

// Checks wheither or not a the brick that dependends directly of the given brick are enriched.
// Returns a slice of brick's pointers if they are, or the first error encountered otherwise
func (infra *Infra) GetDirectNext(brick *Brick) (results Bricks, err error) {
	// browse all infra.Bricks and return bricks that have brick in their directPrevious bricks
	for _, b := range infra.Bricks {
		err = b.EnrichError
		if err != nil {
			return
		}

		var directPrevious Bricks
		directPrevious, err = infra.GetDirectPrevious(b)
		if err != nil {
			return
		}
		for _, dp := range directPrevious {
			if dp.Name == brick.Name {
				results = append(results, infra.Bricks[b.Name])
			}
		}
	}
	return
}

func (infra *Infra) GetLinkedNext(brick *Brick) (results Bricks, err error) {
	var directNext Bricks
	added, err := infra.GetDirectNext(brick)
	if err != nil {
		return
	}
	results = added
	for {
		toAdd := Bricks{}
		for _, b := range added {
			directNext, err = infra.GetDirectNext(b)
			for _, dp := range directNext {
				if !results.BricksContains(dp) {
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

func (infra *Infra) GetBricksFromNames(names []string) (bricks Bricks, err error) {
	for _, name := range names {
		b, exist := infra.Bricks[name]
		if !exist {
			err = ErrBrickNotFound{brick: name}
			return
		}
		bricks = append(bricks, b)
	}
	return
}

func (infra *Infra) GetCorrespondingBricks(
	bricks Bricks,
	specifiers []string,
) (
	correspondingBricks Bricks,
	err error,
) {
	var elementaryBricks Bricks
	for _, brick := range bricks {
		if !brick.IsElementary {
			var sb Bricks
			sb, err = infra.GetSubBricks(brick)
			if err != nil {
				return
			}
			elementaryBricks = append(elementaryBricks, sb...)
		} else {
			elementaryBricks = append(elementaryBricks, brick)
		}
	}

	var bricksToAdd Bricks
	for _, specifier := range specifiers {
		bricksToAdd = Bricks{}
		switch specifier {
		case "linked_previous", "all_previous", "lp", "ap":
			for _, brick := range elementaryBricks {
				bs, _ := infra.GetLinkedPrevious(brick)
				bricksToAdd = append(bricksToAdd, bs...)
			}
		case "direct_previous", "dp":
			for _, brick := range elementaryBricks {
				bs, _ := infra.GetDirectPrevious(brick)
				bricksToAdd = append(bricksToAdd, bs...)
			}
		case "selected", "s":
			for _, brick := range elementaryBricks {
				bs, _ := infra.GetSubBricks(brick)
				bricksToAdd = append(bricksToAdd, bs...)
			}
		case "direct_next", "dn":
			for _, brick := range elementaryBricks {
				bs, _ := infra.GetDirectNext(brick)
				bricksToAdd = append(bricksToAdd, bs...)
			}
		case "linked_next", "all_next", "ln", "an":
			for _, brick := range elementaryBricks {
				bs, _ := infra.GetLinkedNext(brick)
				bricksToAdd = append(bricksToAdd, bs...)
			}
		default:
			err = fmt.Errorf("Error: Brick's specifier doesn't exist: %s", specifier)
			return
		}

		correspondingBricks = append(correspondingBricks, bricksToAdd...)
	}

	sort.Sort(correspondingBricks)
	correspondingBricks = RemoveDuplicates(correspondingBricks)

	// check if some correspondingBricks are in error
	for _, b := range correspondingBricks {
		if b.EnrichError != nil {
			err = fmt.Errorf(
				"Error: a specified brick couldn't beeing enriched:%s: %v",
				b.Name, b.EnrichError)
			return
		}
	}

	// for direct next or linked next: check if some elementaryBricks with an higher
	// index have an enrich error

	return
}

func (infra *Infra) EnrichBricks() {
	for _, b := range infra.Bricks {
		if b.IsElementary {
			conf, err := BrickConfYaml{}.New(b.ConfigurationFilePath)
			if err != nil {
				infra.Bricks[b.Name].EnrichError =
					fmt.Errorf("unable to enrich brick(%s): %v", b.Name, err)
			}

			err = b.Enrich(conf, infra)
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

func (infra *Infra) ValidateConfiguration(configuration *exargs.Configuration) (err error) {
	// validate brick names
	for _, brickName := range configuration.BricksNames {
		if _, ok := infra.Bricks[brickName]; !ok {
			return ErrBadArg{Reason: "Brick doesn't exist:", Value: brickName}
		}
	}

	// validate BricksSpecifiers
	for _, specifier := range configuration.BricksSpecifiers {
		if !extools.ContainsString(exargs.AvailableBricksSpecifiers[:], specifier) {
			return ErrBadArg{Reason: "Brick's specifier doesn't exist:", Value: specifier}
		}
	}
	return nil
}

func GetConfFilePath(path string) (string, error) {
	confFilePath := filepath.Join(path + BRICK_FILE_NAME)
	_, err := os.Stat(confFilePath)

	if err != nil {
		return confFilePath, nil
	}

	return confFilePath, fmt.Errorf("No configuration file was found in %s", confFilePath)
}

func GetBricksThatCallthisOutput(linkedBricks Bricks, brick *Brick, jsonpath string) (depends Bricks) {
	for _, b := range linkedBricks {
		for _, i := range b.Inputs {
			if i.Brick == brick {
				if extools.AreJsonPathsLinked(jsonpath, i.JsonPath) {
					depends = append(depends, b)
					break
				}
			}
		}
	}
	return
}

func (infra *Infra) GetBricksThatCallthisOutput(brick *Brick, jsonpath string) (depends Bricks, err error) {

	var linkedBricks Bricks

	linkedBricks, err = infra.GetDirectNext(brick)
	if err != nil {
		return
	}

	depends = GetBricksThatCallthisOutput(linkedBricks, brick, jsonpath)
	return
}
