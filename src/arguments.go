package arguments
import (
  "fmt"
  "os"
  "strings"
)

type Arguments struct {
    action string
    bricks_paths []string
    bricks_specifiers []string
    interactive bool
    output_specifier string
    other_options []string
}

actions_list = [
    "plan", "lay", "remove", "output", "init", "validate_code", "help",
    "show_input", "list_elementary_bricks", "cd",
    "get_brick_path", "get_brick_name"
]
func is_string_in_list(str string, list []string) bool {
    for _, value := range list {
        if value == str {
            return true, index
        }
    }
    return false, -1
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

func Get_arguments() (Arguments, bool) {
    var args := Arguments{
        action: "",
        bricks_paths: [],
        bricks_specifiers: ["selected"]
        interactive: true
        output_specifier: ""
        other_options: []
    }
    var os_args := remove_item(0, os.Args)

    // set action and remove it
    if is_string_in_list(os_args[0], actions_list) {
        args.action = os_args[0]
        os_args = remove_item(0, os.Args)
    } else if is_string_in_list(os_args[1], actions_list){
        args.action = os_args[1]
        os_args = remove_item(1, os.Args)
    } else {
        fmt.Println("You need to specify an action: \"exeiac help\"")
        return args, false
    }
    
    // set bricks_paths
    if is_path(os_args[0]) {
        args.bricks_paths = [get_full_path(os_args[0])]
    } else {
        fmt.Println("Sorry for the moment the only brick outputs is: \"exeiac help\"")
        return args, false
    }

    // return
    args.other_options = os_args
    return args, true
}

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

