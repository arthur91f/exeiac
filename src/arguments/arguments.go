package arguments
import (
  "fmt"
  "os"
  "strings"
  "path/filepath"
)

func get_brick_full_path(path string) (string, error) {
    info, err := os.Stat(path)
    if err != nil {
        return "", fmt.Errorf("Error:arguments/get_brick_full_path: %s %w\n",
            path, err)
    } else if !info.IsDir() {
        return "", fmt.Errorf("Error:arguments/get_brick_full_path: " +
            "%s not a directory\n", path)
    }

    if _, err1 := os.Stat(path + "/brick.yml") ; err1 == nil {
        return filepath.Abs(path)
    } else if _, err2 := os.Stat(path + "/brick.yaml") ; err2 == nil {
        return filepath.Abs(path)
    } else if os.IsNotExist(err1) {
        return "", fmt.Errorf("Error:arguments/get_brick_full_path: " +
            "%s/brick.yml don't exist\n", path)
    } else if os.IsNotExist(err2) {
        return "", fmt.Errorf("Error:arguments/get_brick_full_path: " +
            "%s/brick.yaml don't exist\n", path)
    } else {
        return "", fmt.Errorf("Error:arguments/get_brick_full_path: " +
            "%s: {%s: %w} {%s:%v}", path, "brick.yml", err1, "brick.yaml", err2)
    }
}

func GetArguments() (Arguments, error) {
    os_args := remove_item(0, os.Args)
    value := ""
    found := false

    // set overload conf file and read it
    value, _ = consume_opt_and_val("--conf-file", "-c", &os_args)
    args, err := set_args_with_conf_files(value)
    if err != nil {
        return args, fmt.Errorf("Error:arguments/GetArguments:" +
            "get configuration from files error: %w\n", err)
    }

    // set action and remove it
    if len(os_args) < 1 {
        return args, fmt.Errorf("Error:arguments/GetArguments: " + 
            "You need to specify at least an action\n")
    } else if is_string_in_list(os_args[0], actions_list) {
        args.Action = os_args[0]
        os_args = remove_item(0, os_args)
    } else if len(os_args) < 2 {
        return args, fmt.Errorf("Error:arguments/GetArguments: " +
            "You need to specify at least an action in first or second arg\n")
    } else if is_string_in_list(os_args[1], actions_list){
        args.Action = os_args[1]
        os_args = remove_item(1, os_args)
    } else {
        return args, fmt.Errorf("Error:arguments/GetArguments: " +
            "You need to specify an action: \"exeiac help\"\n")
    }
    
    // set bricks_paths
    // TODO: supports --bricks-names|-n
    if value, found = consume_opt_and_val("--bricks-paths", "-p", &os_args) ;
        found {
        args.BricksPaths = strings.Split(value, ",")
        for i, p := range args.BricksPaths {
            args.BricksPaths[i], err = get_brick_full_path(p)
            if err != nil {
                return args, fmt.Errorf("SorryError:arguments/GetArguments: " +
                    "%s not vaid brick\n%w\n", p, err)
            }
        }
    } else if _, found = consume_opt_and_val("--bricks-names", "-n", &os_args) ;
        found {
        return args, fmt.Errorf("SorryError:arguments/GetArguments: " +
            "exeiac do not support --bricks-names|-n options\n")
    } else {
        // TODO: check if it's a brick_name or a brick_path
        if len(os_args) > 0  {
            args.BricksPaths = []string{os_args[0]}
            os_args = remove_item(0, os_args)
        }
    }

    // set interactive
    if consume_opt("--non-interactive", "-I", &os_args) {
        args.Interactive = false
    }
    
    // set output specifier
    value, found = consume_opt_and_val("--output-specifier", "-S", &os_args)
    if found {
        args.OutputSpecifier = value
    }
    
    // set brick specifier
    value, found = consume_opt_and_val("--bricks-specifiers", "-s", &os_args)
    if found {
        args.BricksSpecifiers = strings.Split(value, "+")
    }

    // return
    if len(os_args) > 0 {
        args.OtherOptions = os_args
    }

    return args, nil
}

