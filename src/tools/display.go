package tools

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const LINE_LENGTH = 80

func AskConfirmation(question string) (confirm bool, err error) {

	reader := bufio.NewReader(os.Stdin)
	var answer string

	fmt.Printf("\033[1m" + question + " \033[0m\033[3m" + "(only yes accepted): \033[0m")
	answer, err = reader.ReadString('\n')

	if err != nil {
		return
	}

	switch answer {
	case "yes\n", "YES\n", "Yes\n":
		confirm = true
	default:
		confirm = false
	}
	return
}

func DisplaySeparator(separatorName string) {
	var s string
	if separatorName == "" {
		s = strings.Repeat("_", LINE_LENGTH)
		fmt.Printf("\033[01;36m" + s + "\033[0m\n")
	} else {
		endLineLength := LINE_LENGTH - len(separatorName) - 4
		if endLineLength > 0 {
			s = strings.Repeat("-", endLineLength)
			fmt.Printf("\033[01;36m-- " + separatorName + " " + s + "\033[0m\n")
		} else {
			lengthToDisplay := LINE_LENGTH - 3
			s = "-- " + separatorName
			s = s[0:lengthToDisplay]
			fmt.Printf("\033[01;36m" + s + "..." + "\033[0m\n")
		}
	}
}
