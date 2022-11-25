package arguments

import (
	"fmt"
	"os"
	"path/filepath"
	extools "src/exeiac/tools"
	"strings"
)

func (a Arguments) String() string {
	var modulesString string
	var roomsString string

	if len(a.Modules) > 0 {
		for _, m := range a.Modules {
			modulesString = fmt.Sprintf("%s%s", modulesString,
				extools.IndentForListItem(m.String()))
		}
		modulesString = fmt.Sprintf("modules:\n%s", modulesString)
	} else {
		modulesString = "modules: []\n"
	}

	if len(a.Rooms) > 0 {
		for _, r := range a.Rooms {
			roomsString = fmt.Sprintf("%s%s", roomsString,
				extools.IndentForListItem(r.String()))
		}
		roomsString = fmt.Sprintf("rooms:\n%s", roomsString)
	} else {
		roomsString = "rooms: []\n"
	}

	return fmt.Sprintf("Arguments:\n%s", extools.Indent(
		"bricksNames:"+extools.StringListOfString(a.BricksNames)+"\n"+
			"bricksSpecifiers:"+extools.StringListOfString(a.BricksSpecifiers)+
			"\n"+fmt.Sprintf("interactive: %t", a.Interactive)+"\n"+
			fmt.Sprintf("outputSpecifier: %s", a.OutputSpecifier)+"\n"+
			"otherOptions:"+extools.StringListOfString(a.OtherOptions)+"\n"+
			modulesString+roomsString))
}
func splitBricksPathsAndNames(bricks []string) (
	paths []string, names []string, err error) {
	var absPath string
	for _, brick := range bricks {
		info, statErr := os.Stat(brick)
		if statErr == nil {
			if info.IsDir() {
				absPath, err = filepath.Abs(brick)
				if err == nil {
					paths = append(paths, absPath)
				} else {
					return
				}
			} else {
				err = ErrBadArg{Reason: "Not a brick : path is not a directory",
					Value: brick}
				return
			}
		} else {
			names = append(names, brick)
		}
	}
	return
}

func getRegexList(rooms []extools.NamePathBinding) (
	replaceList []extools.ReplaceOperation, err error) {
	var replacement extools.ReplaceOperation
	for _, room := range rooms {
		replacement, err = extools.CreateReplaceOperation("^"+room.Path, room.Name)
		if err != nil {
			return
		}
		replaceList = append(replaceList, replacement)
	}
	replacement, err = extools.CreateReplaceOperation(`/\d+-`, "/")
	if err != nil {
		err = fmt.Errorf("INTERNAL ERROR: 63809c7f:%v", err)
		return
	}
	replaceList = append(replaceList, replacement)
	return
}

func convertToGetOnlyBrickName(
	bricks []string, roomsPaths []extools.NamePathBinding) (
	names []string, err error) {

	var absPath string
	var name string
	var replaceList []extools.ReplaceOperation

	for _, brick := range bricks {
		info, statErr := os.Stat(brick)
		if statErr != nil { // if it's not a path we assume it's already a name
			names = append(names, brick)
		} else {
			if !info.IsDir() {
				err = ErrBadArg{Reason: "Not a brick: path is not a directory",
					Value: brick}
				return
			}
			absPath, err = filepath.Abs(brick)
			if err != nil {
				return
			}
			if len(replaceList) == 0 {
				replaceList, err = getRegexList(roomsPaths)
				if err != nil {
					return
				}
			}
			name = absPath
			for _, r := range replaceList {
				name = r.Replace(name)
			}
			names = append(names, name)
		}
	}
	return
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
		return args, fmt.Errorf("%w\n> Error636a4927%s "+
			" unable to get configuration from files", err, err_msg)
	}

	// set action and remove it
	if len(os_args) < 1 {
		return args, fmt.Errorf("! Error636a49b4%s "+
			"You need to specify at least an action", err_msg)
	} else if is_string_in_list(os_args[0], actions_list) {
		args.Action = os_args[0]
		os_args = remove_item(0, os_args)
	} else if len(os_args) < 2 {
		return args, fmt.Errorf("! Error636a4ac7%s You need to specify "+
			"at least an action in first or second arg", err_msg)
	} else if is_string_in_list(os_args[1], actions_list) {
		args.Action = os_args[1]
		os_args = remove_item(1, os_args)
	} else {
		return args, fmt.Errorf("! Error636a4b37%s "+
			"You need to specify an action: \"exeiac help\"", err_msg)
	}

	// set bricks (need args.Rooms)
	var bricks []string
	if len(os_args) > 0 {
		if !strings.HasPrefix(os_args[0], "-") {
			if strings.Contains(os_args[0], ",") {
				bricks = strings.Split(os_args[0], ",")
			} else if strings.Contains(os_args[0], "\n") {
				bricks = strings.Split(os_args[0], "\n")
			} else {
				bricks = []string{os_args[0]}
			}
			os_args = remove_item(0, os_args)
		}
	}
	fmt.Println(bricks)
	args.BricksNames, err = convertToGetOnlyBrickName(bricks, args.Rooms)
	if err != nil {
		return args, err
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
