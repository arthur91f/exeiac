#!/bin/bash
# launcher script

#######################
# DECLARE GLOBAL VARS #
#######################
actions_list=" pass init validate fmt 
 public_domain_name   = local.public_domain_name
 plan apply output destroy 
 show_dependencies show_dependents list_bricks 
 show_dependencies_recursively show_dependents_recursively 
 help -h --help debug "
configuration_files_list="/usr/lib/exeiac/exeiac.conf
/usr/local/exeiac/exeiac.conf
/opt/exeiac/exeiac.conf
/etc/exeiac.conf
$HOME/.exeiac
$HOME/.exeiac.conf"

INITIAL_CURRENT_PATH="$(pwd)" # used by soft_exit
ROOMS_LIST=""                 # will be set after sourcing configuration files
MODULES_PATH=""               # will be set after sourcing configuration files 
DEFAULT_MODULE_PATH=""        # will be set after sourcing configuration files 
EXECUTE_SUM_UP_FILE=""        # fill in during execution and display at the end
                              # can be set by option or get a default value
selected_bricks=""            # will be set after arguments interpretation
execute_plan=""               # will be set after arguments interpretation
declare -a OPTS               # will be set after arguments interpretation
brick_name=""                 # will be set after arguments interpretation
brick_path=""                 # will be set after arguments interpretation
action=""                     # will be set after arguments interpretation

##############################
# SOURCE CONFIGURATION FILES #
##############################
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
source "$exeiac_lib_path/bricks_functions.sh"
source "$exeiac_lib_path/dependencies_functions.sh"
source "$exeiac_lib_path/single_called_functions.sh"

#######################################
# CHECK CONFIGURATION FILES ARGUMENTS #
#######################################
if [ -z "$room_paths_list" ]; then
    echo "ERROR:exeiac: room_path_list isn't set in configfile as ~/.exeiac" >&2
else
    ROOMS_LIST="$(sed 's|/ *$||g' <<<"$room_paths_list")"
fi
if [ ! -d "$modules_path" ]; then
    echo "ERROR:exeiac: modules_path set in configfiles as ~/.exeiac isn't a directory" >&2
fi

################################
# SET VARIABLES FROM ARGUMENTS #
################################
if EXECUTE_SUM_UP_FILE="$(get_arg --string=execute-plan-file "$@")"; then
    touch "$EXECUTE_SUM_UP_FILE"
    if [ ! -w "$EXECUTE_SUM_UP_FILE" ]; then
        soft_exit 1 "ERROR:exeiac:execute-sum-up-file not writable:$EXECUTE_SUM_UP_FILE"
    fi
else
    EXECUTE_SUM_UP_FILE="/tmp/exeiac_execute_sum_up-$(date +%y%m%d-%H%M%S-%N)"
fi

# Identify which parameters is ACTION, BRICK, OPTS
if grep -q " $1 " <<<"$actions_list" || grep -q " $2 " <<<"$actions_list"; then
    if [ -e "$2" ] ; then
        action="$1"
        brick_path="$(get_absolute_path "$2")"
        shift 2
        OPTS=("$@")
    elif [ -e "$1" ] ; then
        action="$2"
        brick_path="$(get_absolute_path "$1")"
        shift 2
        OPTS=("$@")
    elif brick_path="$(get_brick_path "$1")"; then
        action="$2"
        shift 2
        OPTS=("$@")
    elif brick_path="$(get_brick_path "$2")"; then
        action="$1"
        shift 2
        OPTS=("$@")
    elif grep -q " $1 " <<<"$actions_list"; then
        action="$1"
        shift 1
        OPTS=("$@")
    else
        soft_exit 1 "Error: bad argument: \"exeiac help\" for help"
    fi
else
    soft_exit 1 "Error: bad argument: \"exeiac help\" for help"
fi

if [ -n "$brick_path" ]; then
    if ! selected_bricks="$(get_selected_bricks \
        -brick-path="$brick_path" "${OPTS[@]}")"; then
        soft_exit 1 "ERROR:exeiac:get_selected_bricks"
    fi
else
    if ! selected_bricks="$(get_selected_bricks "${OPTS[@]}")" ; then
        soft_exit 1 "ERROR:exeiac:get_selected_bricks"
    fi
fi

if specifiers_list="$(get_arg --string=bricks-specifiers "$@")"; then
    execute_plan="$(get_specified_bricks \
        -selected-bricks "$selected_bricks" \
        -bricks-specifiers "$specifiers_list")"
    if [ "$?" != 0 ]; then
        soft_exit 1 "ERROR:exeiac:get_specified_bricks"
    fi
else
    execute_plan="$selected_bricks"
fi

################
# EXECUTE PLAN #
################
if [ -n "$execute_plan" ]; then
    execute_bricks_list -action=$action -bricks-list "$execute_plan"
else
    case "$action" in
    pass)
        true
        return_code=0
        ;;
    help|--help|-h)
        display_help
        return_code=$?
        ;;
    list_bricks)
        get_all_bricks_names
        return_code=$?
        ;;
    debug)
        cmd_debug "$@" -disbale-debug-classic-output
        return_code=$?
        ;;
    *)
        soft_exit 1 "ERROR:exeiac:unrecognized_action_for_null_plan:$action"
        ;;
    esac
fi
soft_exit $return_code

