package tools

import (
	"fmt"
	"strings"
)

func IndentIfMultiline(i string) (o string) {
	lines := strings.Split(i, "\n")
	if len(strings.Split(i, "\n")) > 1 {
		o = fmt.Sprintf("%s\n%s", lines[0],
			Indent(strings.Join(lines[1:], "\n"), "        "))
	} else {
		o = i
	}
	return
}

func Indent(i string, indent string) (o string) {
	lines := strings.Split(i, "\n")
	last_line := len(lines) - 1

	for index, line := range lines {
		if index != last_line {
			o = fmt.Sprintf("%s%s%s\n", o, indent, line)
		} else {
			o = fmt.Sprintf("%s%s%s", o, indent, line)
		}
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
