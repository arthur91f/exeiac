package tools

import (
	"fmt"
	"strings"
)

func Indent(i string) (o string) {
	for _, line := range strings.Split(i, "\n") {
		o = fmt.Sprintf("%s  %s\n", o, line)
	}
	return
}

func IndentForListItem(i string) (o string) {
	lines := strings.Split(i, "\n")
	o = fmt.Sprintf("- %s\n", lines[0])
	for index := 1; index < len(lines); index++ {
		o = fmt.Sprintf("%s  %s\n", o, lines[index])
	}
	return
}

func StringListOfString(i []string) (o string) {
	if len(i) > 0 {
		for _, line := range i {
			o = fmt.Sprintf("%s\n- %s", o, line)
		}
	} else {
		o = " []"
	}
	return
}
