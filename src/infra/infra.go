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
	"sync"
)

type Infra struct {
	Modules []Module
	Bricks  BricksMap
}

func CreateInfra(configuration exargs.Configuration) (Infra, error) {
	i := Infra{
		Bricks: make(map[string]*Brick),
	}

	// create Modules
	for name, path := range configuration.Modules {
		i.Modules = append(i.Modules, Module{
			Name: name,
			Path: path,
		})
	}

	// Temporary brick storage to have consistent indexing across rooms
	// i.e. having an overall ordering as the index
	bricks := []Brick{}
	for name, path := range configuration.Rooms {
		b, err := GetBricks(name, path)
		if err != nil {
			fmt.Printf(`Cannot add "%s" (path: %s) room's bricks: %s`, name, path, err)
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
func GetBricks(roomName string, roomPath string) ([]Brick, error) {
	var bricks []Brick
	var err error

	_, err = os.Stat(roomPath)
	if err != nil {
		if os.IsNotExist(err) {
			err = ErrBrickNotFound{brick: roomPath}
			return bricks, err
		}
		return bricks, err
	}
	bricks = []Brick{{
		Name:         roomName,
		Path:         roomPath,
		IsElementary: false,
	}}
	_, err = os.Stat(roomPath + "/brick.yml")
	if err == nil {
		bricks[0].IsElementary = true
	}

	err = filepath.WalkDir(
		roomPath,
		func(path string, d fs.DirEntry, err error) error {
			brickRelPath, err := filepath.Rel(roomPath, path)
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
				brickName := filepath.Join(roomName, brickRelPath)
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

	return bricks, err
}

func (infra Infra) String() string {
	var sb strings.Builder
	var modulesSb strings.Builder
	var bricksSb strings.Builder

	for _, m := range infra.Modules {
		modulesSb.WriteString(fmt.Sprintf("%s", m))
	}

	bricksSb.WriteString(fmt.Sprintf("%v", infra.Bricks))

	sb.WriteString(fmt.Sprintf("Modules: [\n%s]\nBricks:[\n%s]", modulesSb.String(), bricksSb.String()))

	return sb.String()
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
	superBrickPath := brick.Path
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

// Cheks wheither or not the brick that dependends (directly or not) of the given brick are enriched.
// Returns a slice of brick's pointers if they are, or the first error encountered otherwise
func (infra *Infra) GetLinkedNext(bricks *Bricks, brick *Brick) (err error) {
	var wg sync.WaitGroup

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
				*bricks = append(*bricks, infra.Bricks[b.Name])

				// TODO(half-shell): handle cases of cirulare dependencies
				// or misconfiguration. We want to add a check to either limit
				// the number of wait group created, or just a check over the indexes
				// of the bricks handled here. Or prevent adding a brick already added
				// to avoid those cases.
				if len(b.DirectPrevious) > 0 {
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
				var bs Bricks
				infra.GetLinkedNext(&bs, brick)
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
				infra.Bricks[b.Name].EnrichError = err
			}

			err = b.Enrich(conf, infra)
			if err != nil {
				infra.Bricks[b.Name].EnrichError = err
			}
			err = b.Module.LoadAvailableActions()
			if err != nil {
				infra.Bricks[b.Name].EnrichError = err
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
