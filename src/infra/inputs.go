package infra

import "sort"

type Inputs []Input

// Allows for sorting over Inputs
func (slice Inputs) Len() int {
	return len(slice)
}

// Allows for sorting over Inputs
func (slice Inputs) Less(i, j int) bool {
	return slice[i].VarName < slice[j].VarName
}

// Allows for sorting over Inputs
func (slice Inputs) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func (slice Inputs) SortInputByVarname() {
	sort.Sort(slice)
}

func CreateInputs(i []Input) Inputs {
	return i
}
