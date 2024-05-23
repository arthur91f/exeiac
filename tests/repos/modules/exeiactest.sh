#!/bin/bash
# This module will display it's arg and it's env
# Then it will search if it finds those files : 
#   - state.exeiactest.json (correspond to the output)
#   - code.exeiactest.json (the state willing)
#   - events.exeiactest.txt (jsonpath list that create event when changed)
#   - special_behaviours.exeiactest.json (to simulate error or particular behaviour)

ACTION="$1"
CURRENT_PATH="$(pwd)"
state_path="$CURRENT_PATH/state.exeiactest.json"
code_path="$CURRENT_PATH/code.exeiactest.json"
special_behaviours_path="$CURRENT_PATH/special_behaviours.exeiactest.json"
test_events_path="$CURRENT_PATH/events.exeiactest.txt"
onchange_events_path="$CURRENT_PATH/.onchange_events"
other_events_path="$CURRENT_PATH/.other_events"

# save the function args
declare -a ALL_ARGS
for arg in "$@"; do
    ALL_ARGS+=("$arg")
done

TF_VAR_brickname="$EXEIAC_BRICK_NAME"

function remove_files {
    option_debug="$(grep -qE '.*--debug-input-file' <<<"${ALL_ARGS[@]}" && echo true || echo false)"
    path_is_brick="$(grep -qE 'tests2/repos(/.*/|/)[0-9]+-[^/]+$' <<<"$CURRENT_PATH" && echo true || echo false)"
    if [ "$option_debug" == "false" ] && [ "$path_is_brick" == "true" ] ; then
        for file in $(find "$CURRENT_PATH" -name 'input.*' | grep -E '/input\.(json|yaml|yml|env)$' | sort) ; do
            if grep -Eq 'tests2/repos/.*/input\.(json|yaml|yml|env)$' <<<"$file" ; then
                rm "$file" || echo "rm '$file' won't work" >&2
            else
                echo "rm '$file' has not been runned" >&2
            fi
        done
    fi
}

function display_env_and_args { #<$ALL_ARGS
    if grep -qE '^(|.* )--debug-env=.*$' <<<"${ALL_ARGS[@]}"; then
        echo "#######"
        echo "# ENV #"
        for arg in "${ALL_ARGS[@]}"; do
            if grep -q  '^--debug-env=.*$' <<<"$arg"; then
                vars="$(sed 's|^--debug-env=||g' <<<"$arg")"
                env | grep -E "$vars"
            fi
        done
    # env | grep -vE '^(SESSION_MANAGER|QT_ACCESSIBILITY|COLORTERM|XDG_CONFIG_DIRS|XDG_SESSION_PATH|GNOME_DESKTOP_SESSION_ID|LANGUAGE|DESKTOP_SESSION|GTK_MODULES|XDG_SEAT|XDG_SESSION_DESKTOP|QT_QPA_PLATFORMTHEME|XDG_SESSION_TYPE|GPG_AGENT_INFO|XAUTHORITY|XDG_GREETER_DATA_DIR|GDM_LANG|HOME|LANG|LS_COLORS|XDG_CURRENT_DESKTOP|VTE_VERSION|XDG_SEAT_PATH|GNOME_TERMINAL_SCREEN|LESSCLOSE|XDG_SESSION_CLASS|TERM|ASDF_DIR|LESSOPEN|GNOME_TERMINAL_SERVICE|DISPLAY|SHLVL|XDG_VTNR|DESKTOP_AUTOSTART_ID|XDG_SESSION_ID|XDG_RUNTIME_DIR|GTK3_MODULES|XDG_DATA_DIRS|PATH|GDMSESSION|DBUS_SESSION_BUS_ADDRESS|OLDPWD|_)=.*$'
    fi

    echo ""

    echo "########"
    echo "# ARGS #"
    for i in "${ALL_ARGS[@]}"; do
        echo "$i"
    done

    echo ""

    for file in $(find "$CURRENT_PATH" -name 'input.*' | grep -E '/input\.(json|yaml|yml|env)$' | sort) ; do
        echo "##############"
        echo "# INPUT FILE #"
        echo "# $(sed "s|^$CURRENT_PATH/|./|g" <<<"$file")"
        if grep -q '\.json$' <<<"$file" ; then
            jq . --sort-keys "$file"
        elif grep -q '\.env$' <<<"$file" ; then
            cat "$file" | sort
        else
            cat "$file"
        fi
        echo "--------------"
    done

}

function describe_module_for_exeiac {
    echo "{
    \"init\": {
        \"behaviour\": \"standard\"
    },
    \"help\": {
        \"behaviour\": \"standard\"
    },
    \"plan\": {
        \"behaviour\": \"plan\",
        \"status_code_fail\": \"1,4-255\",
        \"events\": {
            \"exeiac_plan_no_drift\": {
                \"type\": \"status_code\",
                \"status_code\": \"0\"
            },
            \"exeiac_plan_drift\": {
                \"type\": \"status_code\",
                \"status_code\": \"2\"
            },
            \"exeiac_plan_unknown\": {
                \"type\": \"status_code\",
                \"status_code\": \"3\"
            }
        }
    },
    \"lay\": {
        \"behaviour\": \"lay\",
        \"status_code_fail\": \"1,3-255\",
        \"events\": {
            \"nothing_todo\": {
                \"type\": \"status_code\",
                \"status_code\": \"0\"
            },
            \"drift\": {
                \"type\": \"status_code\",
                \"status_code\": \"2\"
            },
            \"events\": {
                \"type\": \"file\",
                \"path\": \".onchange_events\"
            },
            \"other\": {
                \"type\": \"file\",
                \"path\": \".other_events\"
            }
        }
    },
    \"remove\": {
        \"behaviour\": \"remove\"
    },
    \"output\": {
        \"behaviour\": \"output\"
    }
}"
}

function get_json_from_file { #< file
    file="$1"
    if [ -f "$file" ]; then
        jq . --sort-keys "$file"
        return "$?"
    else
        echo "null" | jq .
    fi
}

function help {
    echo "exeiactest ACTION [OPTIONS...]

DESCRIPTION:
    It is an exeiac module for testing exeiac or debugging.

    This module can be use by two different brick : a normal brick that usually
    use an other module (for debugging) or a testing brick for e2e exeiac test.
    - normal brick : in a such case it will only display arguments and env vars.
    - testing brick : it's a brick that do nothing except simulate a brick. 
      Those type of brick contains 3 files
        - state.exeiactest.json (correspond to the output)
        - code.exeiactest.json (the state willing)
        - events.exeiactest.txt (jsonpath list that create event when changed)
        - special_behaviours.exeiactest.json (to simulate errors)

ACTIONS:"
"$0" describe_module_for_exeiac | jq -r 'keys | .[]' | sed 's/^\(.*\)$/    \1/g'

    echo "OPTIONS:
    --debug-env=REGEX you could display the env vars value as
                exeiac plan my/brick -o '--debug-env=EXEIAC.*'
                ./exeiactest.sh --debug-env=SSH_AUTH_SOCK|PWD|USER"
}


function init {
    if which jq >/dev/null ; then
        echo "exeiactest:init: jq installed"
    else
        echo "exeiactest:init: jq not installed"
        echo "  https://stedolan.github.io/jq/download/"
        return 2
    fi
    
    if ! [ -f "$state_path" ] || ! [ -f "$code_path" ]; then
        touch "$state_path"
        touch "$code_path"
    fi
}

function plan {
    echo "########"
    echo "# PLAN #"
    if [ -f "$code_path" ]; then
        state="$(get_json_from_file "$state_path")" ; status_code="$?"
        if [ "$status_code" != "0" ]; then
            echo "Unable to get state :$status_code: $state_path" >&2
            echo "Not able to detect drift"
            status_code=5
        fi
        code="$(get_json_from_file "$code_path")" ; status_code="$?"
        if [ "$status_code" != "0" ]; then
            echo "Unable to get code :$status_code: $state_code" >&2
            echo "Not able to detect drift"
            status_code=5
        fi
        diff <(echo "$state") <(echo "$code")
        case "$?" in
            0)
                echo -e "\nNothing to do"
                status_code=0 ;;
            1)
                echo -e "\nA drift has been detected"
                status_code=2 ;;
            *) 
                echo -e "\nAn unexpected error occurs"
                status_code=4 ;;
        esac
    else
        echo "File not found : $code_path" >&2
        echo "Not able to detect drift"
        status_code=3
    fi

    return "$status_code"
}

function lay {
    echo "#######"
    echo "# LAY #"
    echo -n '' > "$onchange_events_path"

    if [ -f "$test_events_path" ]; then
        while read jsonpath ; do
            value="$(jq . "$jsonpath" "$code_path")"
            if [ "$value" != "$(jq . "$jsonpath" "$state_path")" ]; then
                echo "$(sed 's/./_/g' <<<"$jsonpath")=$value" >> "$onchange_events_path"
            fi
        done <"$test_events_path"
    fi

    if [ -f "$code_path" ]; then
        if [ -f "$state_path" ]; then
            state="$(get_json_from_file "$state_path")" ; status_code="$?"
            if [ "$status_code" != "0" ]; then
                echo "Unable to get state :$status_code: $state_path" >&2
                echo "Not able to detect drift"
                status_code=5
            fi
            code="$(get_json_from_file "$code_path")" ; status_code="$?"
            if [ "$status_code" != "0" ]; then
                echo "Unable to get code :$status_code: $state_code" >&2
                echo "Not able to detect drift"
                status_code=5
            fi
            diff <(echo "$state") <(echo "$code")
                case "$?" in
                    0)
                        echo -e "\nNothing to do"
                        status_code=0 ;;
                    1)
                        echo -e "\nA drift has been detected"
                        status_code=2 ;;
                    *) 
                        echo -e "\nAn unexpected error occurs"
                        status_code=4 ;;
                esac
                if [ "$status_code" == 2 ]; then
                    cp "$code_path" "$state_path"
                fi
        else
            echo "The brick is not layed"
            cp "$code_path" "$state_path"
            status_code=2
        fi
    else
        if ! [ -f "$state_path" ] ; then
            echo "File not found : $state_path" >&2
        fi
        if ! [ -f "$code_path" ] ; then
            echo "File not found : $code_path" >&2
        fi
        echo "Not able to lay"
        status_code=3
    fi

    return "$status_code"
}

function remove {
    echo "##########"
    echo "# REMOVE #"
    if [ -f "$state_path" ]; then
        diff <(jq . --sort-keys "$state_path") <(echo "{}")
        echo "{}" > "$state_path"
    else
        echo "Nothing to do"
    fi
}

function output {
    jq . "$state_path" ; err="$?"
    if [ "$err" != 0 ]; then
        echo "jq command failed with status code $err" >&2
        echo "{}"
    fi
    return "$err"
}

function clean {
    remove_files
    if [ -f "$CURRENT_PATH/.onchange_events" ]; then
        rm "$CURRENT_PATH/.onchange_events"
    fi
    if [ -f "$CURRENT_PATH/.other_events" ]; then
        rm "$CURRENT_PATH/.other_events"
    fi
}

function default_exec {
    if grep -qE 'describe_module_for_exeiac|output|help' <<<"$ACTION" ; then
        $ACTION
        return "$?"
    elif grep -q "^$ACTION$" <(describe_module_for_exeiac | jq -r 'keys | .[]') ; then
        display_env_and_args
        $ACTION
        return "$?"
    else
        echo "action not implemented: $ACTION"
        return 21
    fi
}

# function special_behaviours_exec {
#     stdout="$(jq -r ".$ACTION.stdout" "$special_behaviours_path")"
#     stderr="$(jq -r ".$ACTION.stderr" "$special_behaviours_path")"
#     final_state="$(jq ".$ACTION.stderr" "$special_behaviours_path")"
#     events="$(jq ".$ACTION.events" "$special_behaviours_path")"
    
#     if [ "$stdout" != null ]; then
#         echo -e "$stdout"
#     fi
#     if [ "$stdout" != null ]; then
#         echo -e "$stderr" >&2
#     fi
#     if [  ]
#     jq --sort-keys <<<"$final_state" > "$state_path"
#     jq --sort-keys <<<"events" > "$other_events_path"
# }

# {
#     "lay": {
#         "stdout": "TEST: begin lay",
#         "stderr": "TEST: an error is emulated",
#         "final_state": {},
#         "return_code": 1,
#         "events": {},
#         "run": [ "stdout", "stderr", "default", "final_state", "return_code", "events" ]
#     }
# }

if [ -f "$special_behaviours_path" ] && 
    [ "$(jq ".$ACTION" "$special_behaviours_path")" != "null" ]; then
        return_code="$(jq -r ".$ACTION.return_code" "$special_behaviours_path")"
        run="$(jq -r ".$ACTION.run[]" "$special_behaviours_path")"
        stdout="$(jq -r ".$ACTION.stdout" "$special_behaviours_path")"
        stderr="$(jq -r ".$ACTION.stderr" "$special_behaviours_path")"
        final_state="$(jq ".$ACTION.stderr" "$special_behaviours_path")"
        events="$(jq ".$ACTION.events" "$special_behaviours_path")"
        
        for action in $run ; do

            case "$action" in
            stdout) echo -e "$stdout" ;;
            stderr) echo -e "$stderr" >&2 ;;
            final_state) jq --sort-keys <<<"$final_state" > "$state_path" ;;
            events) jq --sort-keys <<<"events" > "$other_events_path" ;;
            default) default_exec ; exit_code="$?" ;;
            return_code) exit_code="$return_code" ;;
            *)
                echo "ERR:special_behaviour bad action in .$ACTION.run it should be a list of " >&2
                echo "[ "stdout", "stderr", "default", "final_state", "return_code", "events" ]" >&2
                exit_code=201
                ;;
            esac

        done

        remove_files
        exit "$exit_code"
else
    default_exec ; return_code="$?"
    remove_files
    exit "$return_code"
fi
