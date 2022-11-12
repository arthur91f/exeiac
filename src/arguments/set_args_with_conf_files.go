package arguments

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	extools "src/exeiac/tools"
	"strings"
)

func get_conf_paths(file_path string) []string {
	// create all possible conf file path and return a list of existing conf file
	err_msg := ":arguments/get_conf_paths:"
	existing_conf_paths := []string{}
	possible_conf_paths := []string{
		"/etc/exeiac.conf",
		"/etc/exeiac.yml",
		"/etc/exeiac.yaml",
	}
	home, err := os.UserHomeDir()
	if err == nil {
		possible_conf_paths = append(possible_conf_paths,
			home+"/.exeiac.conf",
			home+"/.exeiac.yml",
			home+"/.exeiac.yaml")
	}
	if file_path != "" {
		possible_conf_paths = append(possible_conf_paths, file_path)
	}

	for _, path := range possible_conf_paths {
		info, err := os.Stat(path)
		if err == nil {
			if !info.IsDir() {
				existing_conf_paths = append(existing_conf_paths, path)
			} else {
				fmt.Fprintf(os.Stderr, "! Warning636a4d5e%s %s is a directory",
					err_msg, path)
			}
		} else if !os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "! Error00000000:%v\n"+
				"> Warning636a4e1c%s os.Stat(%s)", err, err_msg, path)
		}
	}

	return existing_conf_paths
}

func get_conf(path string) (exeiacConf, error) {
	err_msg := ":arguments/get_conf:"
	var conf exeiacConf

	info, err := os.Stat(path)

	if err != nil {
		return conf, fmt.Errorf("! Error00000000:%w\n"+
			"> Error636a51e1%s os.Stat(%s)", err, err_msg, path)

	} else if info.IsDir() {
		return conf, fmt.Errorf("! Error636a5259%s"+
			" path is a directoryi: %s", err_msg, path)

	} else {
		file, err := ioutil.ReadFile(path)
		if err != nil {
			return conf, fmt.Errorf("! Error00000000:%w\n"+
				"> Error636a52e9%s file can't be read %s", err, err_msg, path)
		}

		err = yaml.Unmarshal(file, &conf)
		if err != nil {
			return conf, fmt.Errorf("! Error00000000:%w\n"+
				"> Error636a5387%s file can't be interpret by yaml: %s",
				err, err_msg, path)
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
			name_test, path_test := extools.AreNamePathBindingEqual(mb, mo)
			if name_test && !path_test {
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
	err_msg := ":arguments/set_args_with_conf_files:"
	args := getDefaultArguments()
	conf_paths_list := get_conf_paths(file_path)

	if len(conf_paths_list) < 1 {
		return args, fmt.Errorf("! Error636a540f%s no conf file found", err_msg)
	}

	for _, conf_file := range conf_paths_list {
		if content, err := get_conf(conf_file); err == nil {
			conf = overload_conf(conf, content)
		} else {
			return args, fmt.Errorf("%w\n> Error636a54c7%s "+
				" conf file unexploitable: %s", err, err_msg, conf_file)
		}
	}

	args.RoomsList = conf.RoomsList
	args.ModulesList = conf.ModulesList
	args.Interactive = !conf.DefaultArgs.NonInteractive
	if conf.DefaultArgs.BricksSpecifiers != "" {
		args.BricksSpecifiers = strings.Split(conf.DefaultArgs.BricksSpecifiers,
			"+")
	}
	if conf.DefaultArgs.OtherOptions != "" {
		args.OtherOptions = strings.Split(conf.DefaultArgs.OtherOptions, " ")
	}
	return args, nil
}
