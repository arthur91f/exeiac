#!/bin/bash
# launcher script

#######################
# DECLARE GLOBAL VARS #
#######################
INITIAL_CURRENT_PATH="$(pwd)" # used by soft_exit
BRICK_NAME=""                 # will be set after arguments interpretation
BRICK_PATH=""                 # will be set after arguments interpretation
ACTION=""                     # will be set after arguments interpretation
declare -a OPTS               # will be set after arguments interpretation
ROOMS_LIST=""                 # will be set after sourcing configuration files
MODULES_PATH=""               # will be set after sourcing configuration files 
DEFAULT_MODULE_PATH=""        # will be set after sourcing configuration files 
EXECUTE_SUM_UP_FILE="/tmp/exeiac_execute_sum_up-$(date +%y%m%d-%H%M%S-%N)"
    # fill in during execution and display at the end

##############################
# SOURCE CONFIGURATION FILES #
##############################
configuration_files_list="/usr/lib/exeiac/exeiac.conf
/usr/local/exeiac/exeiac.conf
/opt/exeiac/exeiac.conf
/etc/exeiac.conf
$HOME/.exeiac
$HOME/.exeiac.conf"
for conf_file in $configuration_files_list ; do
    if [ -f "$conf_file" ]; then
        source "$conf_file"
    fi
done

########################
# IMPORT ALL FUNCTIONS #
########################
source "$exeiac_lib_path/general_functions.sh"
source "$exeiac_lib_path/develop_functions.sh"
source "$exeiac_lib_path/brick_arg_functions.sh"
source "$exeiac_lib_path/actions_functions.sh"
source "$exeiac_lib_path/execute_plan_functions.sh"

#######################################
# CHECK CONFIGURATION FILES ARGUMENTS #
#######################################
if [ -z "$room_paths_list" ]; then
    echo "ERROR: room_path_list isn't set in configfile as ~/.exeiac" >&2
else
    ROOMS_LIST="$(sed 's|/ *$||g' <<<"$room_paths_list")"
fi
if [ ! -d "$modules_path" ]; then
    echo "ERROR: modules_path set in configfiles as ~/.exeiac isn't a directory" >&2
fi

#######################################
# SET GLOBAL VARIABLES FROM ARGUMENTS #
#######################################
actions_list=" init validate fmt
 plan apply output destroy 
 show_dependencies show_dependents list_bricks 
 help -h --help debug "

# Identify which parameters is ACTION, BRICK, OPTS
arg_case=null

if grep -q " $1 " <<<"$actions_list" || grep -q " $2 " <<<"$actions_list"; then
    if [ -e "$2" ] ; then
        arg_case="action+brick"
        action="$1"
        brick_path="$(get_absolute_path "$2")"
        brick_name="$(get_brick_name "$brick_path")"
        shift 2
        OPTS=("$@")
    elif [ -e "$1" ] ; then
        arg_case="brick+action"
        action="$2"
        brick_path="$(get_absolute_path "$1")"
        brick_name="$(get_brick_name "$brick_path")"
        shift 2
        OPTS=("$@")
    elif brick_path="$(get_brick_path "$1")"; then
        arg_case="brick+action"
        action="$2"
        brick_name="$1"
        shift 2
        OPTS=("$@")
    elif brick_path="$(get_brick_path "$2")"; then
        arg_case="action+brick"
        action="$1"
        brick_name="$2"
        shift 2
        OPTS=("$@")
    elif grep -q " $1 " <<<"$actions_list"; then
        arg_case="action_only"
        action="$1"
        shift 1
        OPTS=("$@")
    else
        soft_exit 1 "Error: bad argument: \"exeiac help\" for help"
    fi
else
    soft_exit 1 "Error: bad argument: \"exeiac help\" for help"
fi

##################
# BEFORE EXECUTE #
##################
EXECUTE_ALLOWED="true"

if get_arg --boolean=init-before "$@"; then
    echo "init_before not implemented yet"
    # have to check all params to init the before bricks the specifier bricks and the after bricks
fi

if get_arg --boolean=plan-dependencies-recursively-before "$@"; then
    bricks_list="$(get_dependencies_recursively $BRICK_PATH)"
    execute_bricks_list -bricks-paths-list "$bricks_list" -action=plan
elif get_arg --boolean=plan-dependencies-before "$@"; then
    bricks_list="$(get_dependencies $BRICK_PATH)"
    execute_bricks_list -bricks-paths-list "$bricks_list" -action=plan
fi

if get_arg --boolean=plan-dependents-recursively-before "$@"; then
    bricks_list="$(get_dependents_recursively $BRICK_PATH)"
    execute_bricks_list -bricks-paths-list "$bricks_list" -action=plan
elif get_arg --boolean=plan-dependents-before "$@"; then
    bricks_list="$(get_dependents_recursively $BRICK_PATH)"
    execute_bricks_list -bricks-paths-list "$bricks_list" -action=plan
fi


if [ EXECUTE_ALLOWED != "true" ]; then
    cat "$EXECUTE_SUM_UP_FILE" 
    echo "ABORT:$ACTION $BRICK_NAME" >&2
    soft_exit 1 "ABORT:the before execution failed"
fi

###########
# EXECUTE #
###########

if [ "$arg_case" == "action_only" ]; then
    case "$action" in
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
        if get_arg --boolean=all-bricks "$@"; then
            execute_bricks_list \
                -bricks-paths-list "$(list_bricks)"
                -action="$ACTION"
        elif bricks_list="$(get_arg --string=bricks-paths-list "$@")"; then
            execute_bricks_list \
                -bricks-paths-list "$bricks_list"
                -action="$ACTION"
        elif bricks_list="$(get_arg --string=bricks-names-list "$@")"; then
            bricks_list="$(get_bricks_paths_list "$bricks_list")"
            execute_bricks_list \
                -bricks-paths-list "$bricks_list"
                -action="$ACTION"
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
        list_bricks "$BRICK_NAME"
    ;;
    show_dependents)
        show_dependents "$BRICK_NAME"
    ;;
    show_dependents_recursively)
        get_dependents_recursively "$BRICK_NAME"
    ;;
    show_dependencies_recursively)
        get_dependencies_recursively "$BRICK_PATH"
    ;;
    debug)
        cmd_debug
    ;;
    *)
        execute_brick -brick-path=$brick_path -action=$action
    ;;
    esac
fi
#################
# AFTER EXECUTE #
#################


soft_exit 0

