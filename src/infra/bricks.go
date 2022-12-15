package infra

import "fmt"

// A slice of several Brick.
type Bricks []*Brick

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

func (b Bricks) String() string {
	var str string
	if len(b) > 0 {
		for _, brick := range b {
			str = fmt.Sprintf("%s\n- %s", str, brick.Name)
		}
	} else {
		str = " []"
	}
	return str
}
