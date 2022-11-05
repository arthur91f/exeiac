package arguments
import (
  "fmt"
  "os"
)

type Arguments struct {
    Action string
    Bricks_paths []string
    Bricks_specifiers []string
    Conf_file string
    Interactive bool
    Output_specifier string
    Other_options []string
}

var actions_list = []string{
    "plan", "lay", "remove", "output", "init", "validate_code", "help",
    "show_input", "list_elementary_bricks", "cd",
    "get_brick_path", "get_brick_name"}

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

func consume_option_and_value(opt string, args *[]string) (string, bool) {
    var index int
    var found bool
    var value string
    
    index, found = get_index(long, os_args)
    if found && (len(args) > index+1) {
        value = args[index]
        *args = remove_item(index, remove_item(index, *args))
        return value, found
    } else if found && len(args) < index+2 {
        fmt.Println("Error: ", opt, " should be followed by a value")
        fmt.Println("  assuming value is \"\"")
        *args = remove_item(index, *args)
        return "", found, args
    } else {
        return "", found, args
    }
}

func consume_option(opt string, args *[]string) (bool) {
    var index int
    var found bool
    index, found = get_index(long, os_args)
    if found {
        *args = remove_item(index, *args)
    }
    return found
}

func Get_arguments() (Arguments, bool) {
    args := Arguments{
        Action: "",
        Bricks_paths: []string{},
        Bricks_specifiers: []string{"selected"},
        Conf_file: "",
        Interactive: true,
        Output_specifier: "",
        Other_options: []string{},
        }
    os_args := remove_item(0, os.Args)
    pos_args := &os_args
    var value string
    var found bool

    // set action and remove it
    if len(os_args) < 1 {
        fmt.Println("You need to specify at least an action")
        return args, false
    } else if is_string_in_list(os_args[0], actions_list) {
        args.Action = os_args[0]
        os_args = remove_item(0, os_args)
    } else if len(os_args) < 2 {
        fmt.Println("You need to specify at least an action in first or second arg")
        return args, false
    } else if is_string_in_list(os_args[1], actions_list){
        args.Action = os_args[1]
        os_args = remove_item(1, os_args)
    } else {
        fmt.Println("You need to specify an action: \"exeiac help\"")
        return args, false
    }
    
    // set bricks_paths
    if len(os_args) > 0  {
        args.Bricks_paths = []string{os_args[0]}
        os_args = remove_item(0, os_args)
    }

    // set interactive
    if consume_option("--non-interactive", &os_args) {
        args.Interactive = false
    } else if consume_option("-I", &os_args) {
        args.Interactive = false
    }

    // set overload conf file
    value, found = consume_option_and_value("--conf-file", &os_args)
    if found {
        args.Conf_file = value 
    } else {
        args.Conf_file, _ = consume_option_and_value("-c", &os_args)
    }

    // set output specifier

    // return
    args.Other_options = os_args
    return args, true
}

/*
func get_arg_type(arg string) string {
    // return
    // - action
    // - brick_name
    // - brick_path
    // - boolean_option
    // - option_before_value
    // - option_and_value
    // - other_arguments
    if is_string_in_list(arg, actions_list) {
        return "action"
    } else if {
        
    }

}

func get_action(args []string) string {
    
}

func get_selected_brick(args []string) []string {

}
*/
