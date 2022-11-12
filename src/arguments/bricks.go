package arguments

import (
	"io/fs"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
)

type Brick struct {
	// The brick's name. Usually the name of the parent directory
	Name string
	// The absolute path of the brick's directory
	Path string
	// The absolute path of the `brick.yml` file
	ConfigurationFilePath string
	// The content of the `brick.yml` to be parsed later on
	ConfigurationFileContent []byte
}

// Walks the file system from the provided root, gathers all folders containing a `brick.html` file, and build a Brick struct from it.
// Returns a pointer to a slice of Bricks.
func GetBricksAndContent(root string) (*[]Brick, error) {
	brickFiles := []Brick{}

	err := filepath.WalkDir(
		root,
		func(path string, d fs.DirEntry, err error) error {
			if d.Type().IsRegular() && d.Name() == "brick.yml" {
				content, err := ioutil.ReadFile(path)
				if err != nil {
					log.Fatal(err)
				}

				parentPath := filepath.Dir(path)
				dirs := strings.Split(parentPath, "/")

				brickFiles = append(brickFiles, Brick{
					Name:                     dirs[len(dirs)-1],
					Path:                     parentPath,
					ConfigurationFilePath:    path,
					ConfigurationFileContent: content,
				})
			}

			return err
		})

	if err != nil {
		log.Fatal(err)
	}

	return &brickFiles, nil
}

// Walks the file system from the provided root, gathers all folders containing a `brick.html` file and returns a slice of their absolute path.
func GetBricksPaths(root string) ([]string, error) {
	brickFiles := []string{}

	err := filepath.WalkDir(
		root,
		func(path string, d fs.DirEntry, err error) error {
			if d.Type().IsRegular() && d.Name() == "brick.yml" {
				brickFiles = append(brickFiles, path)
			}

			return err
		})

	if err != nil {
		log.Fatal(err)
	}

	return brickFiles, nil
}
