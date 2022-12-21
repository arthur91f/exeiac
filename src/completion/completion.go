package completion

import (
	"fmt"
	"os"
	exargs "src/exeiac/arguments"
	"src/exeiac/infra"
)

func ListBricks(configuration exargs.Configuration) {
	var bricks infra.Bricks
	for roomName, roomPath := range configuration.Rooms {
		bs, err := infra.GetBricks(roomName, roomPath)
		if err == nil {
			bricks = append(bricks, bs...)
		}
	}

	for _, b := range bricks {
		fmt.Fprintln(os.Stdout, b.Name)
	}
}
