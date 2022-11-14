package infra

import (
	"io/fs"
	"log"
	"path/filepath"
	"regexp"
	tools "src/exeiac/tools"
)

var hasDigitPrefixRegexp = regexp.MustCompile(`.*/\d+-\w+$`)
var prefixRegexp = regexp.MustCompile(`\d+-`)

func validateDirName(path string) bool {
	return hasDigitPrefixRegexp.MatchString(path)
}

func sanitizeBrickName(name string) string {
	return prefixRegexp.ReplaceAllString(name, "")
}

// Walks the file system from the provided root, gathers all folders containing a `brick.html` file, and build a Brick struct from it.
// Returns a pointer to a slice of Bricks.
func getRoomsBricks(room tools.NamePathBinding) ([]Brick, error) {
	brickFiles := []Brick{}

	err := filepath.WalkDir(
		room.Path,
		func(path string, d fs.DirEntry, err error) error {
			brickRelPath, err := filepath.Rel(room.Path, path)
			if err != nil {
				log.Fatal(err)
			}

			// A brick can just be described as a sub-path of a room, containing a prefixed folder name with digits, and split with a hypen ("-")
			if d.Type().IsDir() && validateDirName(path) {
				brickName := filepath.Join(room.Name, brickRelPath)
				name := sanitizeBrickName(brickName)

				// Do not duplicate entries
				if len(brickFiles) == 0 || brickFiles[len(brickFiles)-1].Name != name {
					brickFiles = append(brickFiles, Brick{
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
				name := sanitizeBrickName(brickName)

				// Do not duplicate entries
				if brickFiles[len(brickFiles)-1].Name == name {
					brickFiles[len(brickFiles)-1].IsElementary = true
					brickFiles[len(brickFiles)-1].ConfigurationFilePath = path
				}
			}

			return err
		})

	return brickFiles, err
}
