package tools

import (
	"bufio"
	"fmt"
	"os"
)

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
