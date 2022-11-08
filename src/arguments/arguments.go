package arguments
import (
  "fmt"
  "os"
  "strings"
  "path/filepath"
)

func get_brick_full_path(path string) (string, error) {
    err_msg := ":arguments/get_brick_full_path:"
    info, err := os.Stat(path)
    if err != nil {
        return "", fmt.Errorf("! Error00000000:%w\n" + 
            "> Error636a3f09%s os.Stat(%s)", 
            err, err_msg, path)
    } else if !info.IsDir() {
        return "", fmt.Errorf("! Error636a40ea%s " +
            "path is not a directory: path=%s", err_msg, path)
    }

    if _, err1 := os.Stat(path + "/brick.yml") ; err1 == nil {
        return filepath.Abs(path)
    } else if _, err2 := os.Stat(path + "/brick.yaml") ; err2 == nil {
        return filepath.Abs(path)
    } else if os.IsNotExist(err1) {
        return "", fmt.Errorf("! Error636a4412%s %s/brick.yml don't exist", 
            err_msg, path)
    } else if os.IsNotExist(err2) {
        return "", fmt.Errorf("! Error636a44bd%s %s/brick.yaml don't exist",
            err_msg, path)
    } else {
        return "", fmt.Errorf("! Error00000000:%w\n" + "& Error00000000:%w\n" + 
            "> Error636a4711%s os.Stat(%s/brick.yml) && Os.Stat(%s/brick.yaml)",
            err1, err2, err_msg, path, path)
    }
}

func GetArguments() (Arguments, error) {
    err_msg := ":arguments/GetArguments:"
    os_args := remove_item(0, os.Args)
    value := ""
    found := false

    // set overload conf file and read it
    value, _ = consume_opt_and_val("--conf-file", "-c", &os_args)
    args, err := set_args_with_conf_files(value)
    if err != nil {
        return args, fmt.Errorf("%w\n> Error636a4927%s " +
            " unable to get configuration from files", err, err_msg)
    }

    // set action and remove it
    if len(os_args) < 1 {
        return args, fmt.Errorf("! Error636a49b4%s " + 
            "You need to specify at least an action", err_msg)
    } else if is_string_in_list(os_args[0], actions_list) {
        args.Action = os_args[0]
        os_args = remove_item(0, os_args)
    } else if len(os_args) < 2 {
        return args, fmt.Errorf("! Error636a4ac7%s You need to specify " +
            "at least an action in first or second arg", err_msg)
    } else if is_string_in_list(os_args[1], actions_list){
        args.Action = os_args[1]
        os_args = remove_item(1, os_args)
    } else {
        return args, fmt.Errorf("! Error636a4b37%s " +
            "You need to specify an action: \"exeiac help\"", err_msg)
    }
    
    // set bricks_paths
    // TODO: supports --bricks-names|-n
    if value, found = consume_opt_and_val("--bricks-paths", "-p", &os_args) ;
        found {
        args.BricksPaths = strings.Split(value, ",")
        for i, p := range args.BricksPaths {
            args.BricksPaths[i], err = get_brick_full_path(p)
            if err != nil {
                return args, fmt.Errorf("%w\n> Error636a4b4d%s " +
                    "%s is not valid brick", err, err_msg, p)
            }
        }
    } else if _, found = consume_opt_and_val("--bricks-names", "-n", &os_args) ;
        found {
        return args, fmt.Errorf("! SorryError636a4c00%s for the moment " +
            "exeiac do not support options --bricks-names|-n", err_msg)
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

