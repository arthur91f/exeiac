#!/bin/bash
source $HOME/.exeiac
    # room_paths_list (list of exeiac repo in the right order)
    # sanitize_function
    # modules_path
room_paths_list="$(sed 's|/ *$||g' <<<"$room_paths_list")"
current_path="$(pwd)"

source "$exeiac_lib_path/general_functions.sh"
source "$exeiac_lib_path/develop_functions.sh"
source "$exeiac_lib_path/brick_arg_functions.sh"
source "$exeiac_lib_path/actions_functions.sh"
source "$exeiac_lib_path/execute_plan_functions.sh"

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
        dispdebug "- $1 - $2 - $3 -"
        exit 1
    fi
else
    echo "Error: bad argument: \"exeiac help\" for help" >&2
    exit 1
fi

declare -a MODULE_OPTS
EXEIAC_OPTS=""
for arg in "$@"; do
    if grep -q '^--exeiac-opts=' <<<"$arg"; then
        exeiac_opt="$(sed 's/^--exeiac-opts=//g' <<<"$arg")"
        EXEIAC_OPTS="$EXEIAC_OPTS,$exeiac_opt"
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
        if get_arg_in_string "plan-dependencies-before" "$EXEIAC_OPTS"; then
            for dependency in $($0 show_dependencies "$brick_path"); do
                if ! $0 "plan" "$dependency"; then
                    execute_brick_authorize="false"
                    break
                fi
            done
            execute_brick_authorize="true"
        elif get_arg_in_string "plan-dependencies-before-recursively" "$EXEIAC_OPTS"; then
            for dependency in $(get_dependencies_recursively); do
                if ! $0 "plan" "$dependency"; then
                    execute_brick_authorize="false"
                    break
                fi
            done
            execute_brick_authorize="true"
        else
            execute_brick_authorize="true"
        fi
        
        if execute_brick ; then
            apply_succeed="true"
        else
            apply_succeed="false"
        fi

        dependents
        if get_arg_in_string "plan-dependents-after" "$EXEIAC_OPTS"; then
            for dependency in $($0 show_dependencies "$brick_path"); do
                if ! $0 "plan" "$dependency"; then
                    execute_brick_authorize="false"
                    break
                fi
            done
        elif get_arg_in_string "plan-dependents-after-recursively" "$EXEIAC_OPTS"; then

        elif get_arg_in_string "apply-dependents-after" "$EXEIAC_OPTS"; then

        elif get_arg_in_string "apply-dependents-after-recursively" "$EXEIAC_OPTS"; then

        fi

    ;;
    esac
fi
cd "$current_path"

