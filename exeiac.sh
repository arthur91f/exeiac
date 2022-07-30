source $HOME/.exeiac
    # room_paths_list (list of exeiac repo in the right order)
    # sanitize_function
    # modules_path
room_paths_list="$(sed 's|/ *$||g' <<<"$room_paths_list")"
current_path="$(pwd)"

function init_module {
    source "$exeiac_lib_path/default_module.sh"
}

function dispdebug {
    echo "debug: $1" >&2
}

function install {
    return_code=0
    if [ -z "$room_paths_list" ]; then
        echo "add room_paths_list variable in $HOME/.exeiac"
        return_code=1
    fi
    if [ "$(type -t sanitize_function)" != "function" ]; then
        echo "add a sanitize_function in $HOME/.exeiac"
        return_code=1
    fi
    if [ -z "$modules_path" ]; then
        echo "add a modules_path variable in $HOME/.exeiac"
        return_code=1
    fi
    exit $return_code
}

function get_brick_sanitized_name {
    brick_name="$1"
    sanitize_function "$brick_name"
}

function get_brick_type {
    brick_path="$1"
    if ! grep -q "/[0-9]\+-[^/]*$" <<<"$brick_path" ; then
        echo "not_a_brick"
        return 1
    fi
    if [ -f "$brick_path" ] && [ -x "$brick_path" ] ; then
        echo "elementary_script_brick"
    elif [ -d "$brick_path" ] &&
        ls -1 "$brick_path" | grep -q "^[0-9]\+-[^/]*$" ; then
        echo "super_brick"
    elif [ -d "$brick_path" ]; then
        echo "elementary_directory_brick"
    else
        echo "not_a_brick"
        echo "Warning: As it begins by a number this path should be considered 
    as a brick but is not : $brick_path" >&2
        return 1
    fi
    return 0
}

function get_elementary_bricks_list {
    bricks_path_list="$(for room_path in $room_paths_list ; do
        cd "$room_path"
        find . | grep "/[0-9]\+-[^/]*$" | grep -v '/[^0-9]' | 
        sed "s|^\./|$room_path/|g" ; done)"
    for brick_path in $bricks_path_list ; do
        if [ "$(get_brick_type "$brick_path")" != "super_brick" ]; then
            get_brick_name "$brick_path"
        fi
    done
}

function get_brick_name {
    brick_path="$1"
    for room_path in $room_paths_list ; do
        if grep -q "$room_path" <<<"$brick_path" ; then
            room_name="$(sed 's|^.*/\([^/]*\)$|\1|g' <<<"$room_path")"
            sed "s|$room_path|$room_name|g" <<<"$brick_path"
        fi
    done
}

function get_brick_path {
    brick_name="$1"
    room_name="$(sed 's|^\([^/]*\)/.*$|\1|g' <<<"$brick_name")"
    room_path="$(grep "/$room_name$" <<<"$room_paths_list")"
    if [ -z "$room_path" ]; then # if $1 wasn't a brick name
        return 1
    fi
    parent_room_path="$(sed "s|/$room_name$||g" <<<"$room_path")"
    echo "$parent_room_path/$brick_name"
}

function execute_brick {
    case "$(get_brick_type "$brick_path")" in
        super_brick)
            for sub_brick in $(ls -1 "$brick_path") ; do
                if get_brick_type "$brick_path/$sub_brick" ; then
                    echo "## EXEC BRICK: $(get_brick_name "$brick_path/$sub_brick")"
                    $0 "$brick_path/$sub_brick" "$action" "$@"
                    if [ "$?" != 0 ]; then
                        echo "ERROR: $action fail on $(get_brick_name "$brick_path")"
                        exit 1
                    fi
                    echo "## ---------------- ##"
                fi
            done
        ;;
        elementary_script_brick)
            cd "$(dirname "$brick_path")"
            source "$exeiac_lib_path/default_module.sh"
            source "$brick_path"
            $action "$@"
            cd "$current_path"
        ;;
        elementary_directory_brick)
            cd "$brick_path"
            module_found="false"
            for module in $( ls -1 "$modules_path") ; do
                source "$exeiac_lib_path/default_module.sh"
                source "$modules_path/$module"
                if is_brick_using_this_module ; then
                    module_found="true"
                    break
                fi
            done
            if [ "$module_found" == "false" ]; then
                echo "Error: no known modules can execute that brick: $brick_name" >&2
                exit 5
            fi
            $action "$@"
            cd "$current_path"
        ;;
        not_a_brick)
            echo "Error: not a brick: $brick_name" >&2
            exit 3
        ;;
        *)
            echo "Error: brick with unrecognized type: $brick_name" >&2
            exit 4
        ;;
    esac
}

function list_bricks {
    brick_name="$1"
    if [ -n "$brick_name" ]; then
        elementary_bricks_list="$(get_elementary_bricks_list |
            grep "^$brick_name")"
    else
        elementary_bricks_list="$(get_elementary_bricks_list)"
    fi
    echo "$elementary_bricks_list"
}

function display_help {
    echo "exeiac (BRICK_PATH|BRICK_NAME) (install|init|plan|apply|output|destroy|validate|fmt|help|show_dependencies|show_dependents|list_bricks) [\"--exeiac-opts=OPT1,OPT2='VAL 1'\"] [MODUE_ARGS]"
    echo "exeiac (list_bricks|help)"
    echo "exeiac debug FUNCTION"
    echo ""
    echo "ACTIONS:"
    echo "install: installs or check that all tools are installs to execute 
    the brick"
    echo "init: get some dependencies, typically download terraform modules
    or ansible deps"
    echo "plan: a dry run to check what we want to apply"
    echo "apply: run the IaC with the right tools"
    echo "output: display some outputs that can be used by other bricks"
    echo "destroy: revert the apply"
    echo "validate: validate that the syntaxe is ok"
    echo "fmt: rewrite file to pass linter"
    echo "help: display help for this brick or general help"
    echo "show_dependencies: get brick_name whom output is needed"
    echo "show_dependents: get all brick_name that call the brick's output"
    echo "list_bricks: list all elementary brick of the super brick or all
    elementary brick"
    echo "debug: run a function of this file (use for
    unitary test"
    echo ""
    echo "EXEIAC_OPTS"
    echo "all-room|all-bricks: make the action on all bricks"
    echo "plan-dependencies-before"
    echo "plan-dependents-after"
    echo "apply-dependents-after"
    echo "plan-dependencies-before-recursively"
    echo "plan-dependents-after-recursively"
    echo "apply-dependents-after-recursively"
    echo "non-interactive"
    echo ""
    echo "MODULE_ARGS:"
    echo "Very useful for filter a brick output. But despite it is possible to send ARGS to module, remember that exeIaC aim is to manage the apply of the whole IaC, not for debugging and send ARGS to super_brick is dangerous as it will send useless ARGS to some elementary_bricks"
}

function display_line_after_match {
    text="$1"
    regex="$2"
    match_line="$(echo "$text" | grep -n "$regex" | cut -d: -f1)"
    if [ -z "$match_line" ]; then
        return 1
    else
        echo "$text" | sed -n $match_line',$p'
    fi
}

function show_dependents {
    bricks_to_check="$(
        display_line_after_match "$(get_elementary_bricks_list)" "$brick_name")"
    for brick in $bricks_to_check ; do
        if grep -q "$brick_name" <<<"$($0 "$brick" show_dependencies)"; then
            echo "$brick"
        fi
    done
}

function get_absolute_path {
    path="$1"
    if [ -e "$path" ]; then
        echo "$(cd "$(dirname "$path")"; pwd)/$(basename "$path")"
    else
        echo "Error path not exist: $path"
        return 1
    fi
}

function cmd_debug {
     echo "arg_case: $arg_case"
     echo "action: $action"
     echo "brick_path: $brick_path"
     echo "brick_name: $brick_name"
     get_brick_type "$brick_path" 
}

# Identify which parameters is ACTION, BRICK, EXEIAC_OPTS or MODULE_OPTS
arg_case=null
actions_list=" install init 
 plan apply output destroy 
 validate fmt 
 show_dependencies show_dependents list_bricks 
 help -h --help debug "

if grep -q " $1 " <<<"$actions_list" || grep -q " $2 " <<<"$actions_list"; then
    if [ -e "$2" ] ; then
        arg_case="action+brick"
        action="$1"
        brick_path="$(get_absolute_path "$2")"
        brick_name="$(get_brick_name "$brick_path")"
        shift 2
    elif [ -e "$1" ] ; then
        arg_case="brick+action"
        action="$2"
        brick_path="$(get_absolute_path "$1")"
        brick_name="$(get_brick_name "$brick_path")"
        shift 2
    elif brick_path="$(get_brick_path "$1")"; then
        arg_case="brick+action"
        action="$2"
        brick_name="$1"
        shift 2
    elif brick_path="$(get_brick_path "$2")"; then
        arg_case="action+brick"
        action="$1"
        brick_name="$2"
        shift 2
    elif grep -q " $1 " <<<"$actions_list"; then
        arg_case="action_only"
        action="$1"
        shift 1
    else
        echo "Error: bad argument: \"exeiac help\" for help" >&2
        exit 1
    fi
else
    echo "Error: bad argument: \"exeiac help\" for help" >&2
    exit 1
fi

declare -a MODULE_OPTS
declare -a EXEIAC_OPTS
for arg in "$@"; do
    if [ grep -q '^--exeiac-opts=' <<<"$arg" ]; then
        exeiac_opt="$(sed 's/^--exeiac-opts=//g' <<<"$arg")"
        EXEIAC_OPTS[${#EXEIAC_OPTS[@]}]="$exeiac_opt"
    fi
    MODULE_OPTS[${#MODULE_OPTS[@]}]="$arg"
done

if [ "$arg_case" == "action_only" ]; then
    case "$action" in
    install)
        install
    ;;
    list_bricks)
        list_bricks
    ;;
    help|-h|--help)
        display_help
    ;;
    debug)
        cmd_debug        
    ;;
    *)
        if grep -q "all-room" <<<"$EXEIAC_OPTS"; then
            for room_path in $room_paths_list ; do
                $0 "$action" "$room_path" "$@"
            done
        else
            echo "Error: bad argument" >&2
            display_help >&2
            exit 1
        fi
    ;;
    esac
else
    case "$action" in
    list_bricks)
        list_bricks "$brick_name"
    ;;
    show_dependents)
        show_dependents
    ;;
    debug)
        cmd_debug
    ;;
    *)
        execute_brick
    ;;
    esac
fi
cd "$current_path"

