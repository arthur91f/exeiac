package arguments

import "fmt"

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
    new_list := make([]string, (len(list) - 1))
    n := 0
    for i := 0; i < (len(list)); i++ {
        if i != index {
            new_list[n] = list[i]
            n++
        }
    }
    return new_list
}

func consume_opt_and_val(short string, long string, args *[]string) (string, bool) {
    var index int
    var found bool
    var value string

    // check if there is enough args
    if len(*args) < 2 {
        return "", false
    }

    // search long
    if index, found = get_index(long, *args) ; found && (len(*args) > index+1) {
        value = (*args)[index]
        *args = remove_item(index, remove_item(index, *args))
        return value, found
    } else if found && len(*args) < index+2 {
        fmt.Println("Warning: ", long, " should be followed by a value")
        fmt.Println("  ignoring option")
        *args = remove_item(index, *args)
        return "", false
    }

    // search short
    if index, found = get_index(short, *args) ; found && (len(*args) > index+1) {
        value = (*args)[index]
        *args = remove_item(index, remove_item(index, *args))
        return value, found
    } else if found && len(*args) < index+2 {
        fmt.Println("Warning: ", short, " should be followed by a value")
        fmt.Println("  ignoring option")
        *args = remove_item(index, *args)
        return "", false
    }

    // if not found
    return "", false
}

func consume_opt(long string, short string, args *[]string) (bool) {
    var index int
    var found bool

    // check if there is enough args
    if len(*args) < 2 {
        return false
    }

    if index, found = get_index(long, *args) ; found {
        *args = remove_item(index, *args)
    } else if index, found = get_index(short, *args) ; found {
        *args = remove_item(index, *args)
    }
    return found
}


