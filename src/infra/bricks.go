package infra

import (
	"fmt"
	"strings"
)

// A slice of several Brick.
type Bricks []*Brick

// A map of bricks, key being the brick's name, and the value being a reference to a `Brick` struct.
// Defined as itw own type mainly for display purposes.
type BricksMap map[string]*Brick

// Allows for sorting over Bricks
func (slice Bricks) Len() int {
	return len(slice)
}

// Allows for sorting over Bricks
func (slice Bricks) Less(i, j int) bool {
	return slice[i].Index < slice[j].Index
}

// Allows for sorting over Bricks
func (slice Bricks) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func (b Bricks) String() string {
	var sb strings.Builder

	for i, brick := range b {
		sb.WriteString(fmt.Sprintf("- \t%d: %s", brick.Index, brick.Name))
		if i < len(b)-1 {
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

// Helper function to check wheither or not a brick was added to a slice of bricks
func (b Bricks) BricksContains(brick *Brick) bool {
	for _, i := range b {
		if i.Index == brick.Index {
			return true
		}
	}
	return false
}

// Remove duplicates in a slice of bricks.
// Return the de-duplicated slice of bricks.
func RemoveDuplicates(bricks Bricks) Bricks {
	allKeys := make(map[int]bool)
	bs := Bricks{}
	for _, b := range bricks {
		if _, ok := allKeys[b.Index]; !ok {
			allKeys[b.Index] = true
			bs = append(bs, b)
		}
	}
	return bs
}

func (b BricksMap) String() string {
	var sb strings.Builder

	for _, brick := range b {
		sb.WriteString(fmt.Sprintf("- %v", brick))
	}

	return sb.String()
}
