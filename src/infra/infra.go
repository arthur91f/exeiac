package infra

import (
	"fmt"
	extools "src/exeiac/tools"
)

type Infra struct {
	Modules []Module
	Bricks  []Brick
	Rooms   []extools.NamePathBinding
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

func CreateInfra(
	roomsList []extools.NamePathBinding,
	modulesList []extools.NamePathBinding) (
	infra Infra, err error) {

	// create Modules
	for _, m := range modulesList {
		infra.Modules = append(infra.Modules, Module{
			Name: m.Name,
			Path: m.Path,
		})
	}

	// create Bricks
	for _, r := range roomsList {
		// get all room's bricks
		newBricks, err := getRoomsBricks(r)
		if err != nil {
			fmt.Println("%v\n> Warning63724ff3:infra/CreateInfra:"+
				"can't add bricks of this room: %s", err, r.Path)
		}
		infra.Bricks = append(infra.Bricks, newBricks...)
	}

	// add Rooms
	infra.Rooms = roomsList

	return infra, nil
}
