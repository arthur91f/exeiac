package completion

import (
	"fmt"
	"os"
	exargs "src/exeiac/arguments"
	"src/exeiac/infra"
)

func ListBricks(configuration exargs.Configuration) {
	var bricks infra.Bricks
	for _, room := range configuration.Rooms {
		bs, err := infra.DiscoverRoomsBricks(room.Name, room.Path)
		if err == nil {
			bricks = append(bricks, bs...)
		}
	}

	for _, b := range bricks {
		fmt.Fprintln(os.Stdout, b.Name)
	}
}
