package arguments

import (
	"fmt"
	"os"
)

func is_string_in_list(str string, list []string) bool {
	for _, value := range list {
		if value == str {
			return true
		}
	}
	return false
}

func get_index(str string, list []string) (int, bool) {
	for index, value := range list {
		if value == str {
			return index, true
		}
	}
	return -1, false
}

func remove_item(index int, list []string) []string {
	return append(list[:index], list[index+1:]...)
}

func consume_opt_and_val(long string, short string, args *[]string) (string, bool) {
	err_msg := ":arguments/consume_opt_and_val:"
	var index int
	var found bool
	var value string

	// check if there is enough args
	if len(*args) < 2 {
		return "", false
	}

	// search long
	if index, found = get_index(long, *args); found && (len(*args) > index+1) {
		value = (*args)[index+1]
		*args = remove_item(index, remove_item(index, *args))
		return value, found
	} else if found && len(*args) < index+2 {
		fmt.Fprintf(os.Stderr, "! Warning636a575d%s "+
			"option should be followed by a value: %s\n"+
			"> Warning636a57eb%s option ignored: %s\n",
			err_msg, long, err_msg, long)
		*args = remove_item(index, *args)
		return "", false
	}

	// search short
	if index, found = get_index(short, *args); found && (len(*args) > index+1) {
		value = (*args)[index+1]
		*args = remove_item(index, remove_item(index, *args))
		return value, found
	} else if found && len(*args) < index+2 {
		fmt.Fprintf(os.Stderr, "! Warning636a58b2%s "+
			"option should be followed by a value: %s\n"+
			"> Warning636a58cd%s option ignored: %s\n",
			err_msg, short, err_msg, short)
		*args = remove_item(index, *args)
		return "", false
	}

	// if not found
	return "", false
}

func consume_opt(long string, short string, args *[]string) bool {
	var index int
	var found bool

	// check if there is enough args
	if len(*args) < 2 {
		return false
	}

	if index, found = get_index(long, *args); found {
		*args = remove_item(index, *args)
	} else if index, found = get_index(short, *args); found {
		*args = remove_item(index, *args)
	}
	return found
}
