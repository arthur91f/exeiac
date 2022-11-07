package arguments

import (
    "os"
    "fmt"
    "io/ioutil"
    "gopkg.in/yaml.v2"
    "strings"
)

func get_conf_paths(file_path string) ([]string) {
    // create all possible conf file path and return a list of existing conf file
    existing_conf_paths := []string{}
    possible_conf_paths := []string{
        "/etc/exeiac.conf",
        "/etc/exeiac.yml",
        "/etc/exeiac.yaml",
    }
    home, err := os.UserHomeDir()
    if err != nil {
        possible_conf_paths = append(possible_conf_paths,
            home + "/.exeiac.conf",
            home + "/.exeiac.yml",
            home + "/.exeiac.yaml",)
    }
    if file_path != "" {
        possible_conf_paths = append(possible_conf_paths, file_path)
    }

    for _, path := range possible_conf_paths {
        if info, err := os.Stat(path) ; err != nil && !info.IsDir() {
            existing_conf_paths = append(existing_conf_paths, path)
        }
    }

    return existing_conf_paths
}

func get_conf(path string) (exeiacConf, error) {
    var conf exeiacConf
    
    info, err := os.Stat(path)
    
    if err == nil {
        return conf, fmt.Errorf("Error:arguments/get_conf: %w\n", err)
    
    } else if info.IsDir() {
        return conf, fmt.Errorf("Error:arguments/get_conf:" + 
            " %s is a directory\n", path)

    } else {    
        file, err := ioutil.ReadFile(path)
        if err != nil {
            return conf, fmt.Errorf("Error:arguments/get_conf:" +     
                " %s can't be read: %w\n", path, err)
        }

        err = yaml.Unmarshal(file, &conf)
        if err != nil {
            return conf, fmt.Errorf("Error:arguments/get_conf:" +      
                " %s can't be interpret by yaml: %w\n", path, err)
        }
        
        return conf, nil
    }
}

func overload_conf(base exeiacConf, overload exeiacConf) exeiacConf {
    // overload RoomsList by replacing
    if len(overload.RoomsList) > 0 {
        base.RoomsList = overload.RoomsList
    }
    
    // overload ModulesList by adding new item
    for _, mo := range overload.ModulesList {
        for i, mb := range base.ModulesList {
            // search same name module to overload or append
            name_test, path_test := are_NamePathMappings_equal(mb, mo)
            if name_test && !path_test{
                base.ModulesList[i].Path = mo.Path
                break
            } else if !name_test {
                base.ModulesList = append(base.ModulesList, mo)
                break
            }
        }
    }

    // overload DefaultArgs.NonInteractive by replacing if it's not default
    if overload.DefaultArgs.NonInteractive {
        base.DefaultArgs.NonInteractive = true
    }

    // overload DefaultArgs.BricksSpecifiers by replacing
    if overload.DefaultArgs.BricksSpecifiers != "" {
        base.DefaultArgs.BricksSpecifiers =
        overload.DefaultArgs.BricksSpecifiers
    }

    // overload DefaultArgs.OtherOptions by replacing
    if overload.DefaultArgs.OtherOptions != "" {
        base.DefaultArgs.OtherOptions = 
        overload.DefaultArgs.OtherOptions
    }

    return base
}

func set_args_with_conf_files(file_path string) (Arguments, error) {
    var conf exeiacConf
    args := getDefaultArguments()
    conf_paths_list := get_conf_paths(file_path)

    if len(conf_paths_list) < 1 {
        return args, fmt.Errorf("Error:arguments/set_args_with_conf_files:" +
            " no conf file found\n")
    }

    for _, conf_file := range conf_paths_list {
        if content, err := get_conf(conf_file); err != nil {
            overload_conf(conf, content)
        } else {
            return args, fmt.Errorf("Error:arguments/set_args_with_conf_files:" +
                " conf file unexploitable: %s\n" +
                "Error:arguments/get_conf: %w",
                conf_file, err)
        }
    }

    args.RoomsList = conf.RoomsList
    args.ModulesList = conf.ModulesList
    args.Interactive = !conf.DefaultArgs.NonInteractive
    args.BricksSpecifiers = strings.Split(conf.DefaultArgs.BricksSpecifiers, "+")
    args.OtherOptions = strings.Split(conf.DefaultArgs.OtherOptions, " ")

    return args, nil
}

