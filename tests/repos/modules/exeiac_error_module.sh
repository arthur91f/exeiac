#!/bin/bash

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

function display_env_and_args { #<$ALL_ARGS
    echo "#######"
    echo "# ENV #"
    env | grep -vE '^(SESSION_MANAGER|QT_ACCESSIBILITY|COLORTERM|XDG_CONFIG_DIRS|XDG_SESSION_PATH|GNOME_DESKTOP_SESSION_ID|LANGUAGE|DESKTOP_SESSION|GTK_MODULES|XDG_SEAT|XDG_SESSION_DESKTOP|QT_QPA_PLATFORMTHEME|XDG_SESSION_TYPE|GPG_AGENT_INFO|XAUTHORITY|XDG_GREETER_DATA_DIR|GDM_LANG|HOME|LANG|LS_COLORS|XDG_CURRENT_DESKTOP|VTE_VERSION|XDG_SEAT_PATH|GNOME_TERMINAL_SCREEN|LESSCLOSE|XDG_SESSION_CLASS|TERM|ASDF_DIR|LESSOPEN|GNOME_TERMINAL_SERVICE|DISPLAY|SHLVL|XDG_VTNR|DESKTOP_AUTOSTART_ID|XDG_SESSION_ID|XDG_RUNTIME_DIR|GTK3_MODULES|XDG_DATA_DIRS|PATH|GDMSESSION|DBUS_SESSION_BUS_ADDRESS|OLDPWD|_)=.*$'

    echo ""

    echo "########"
    echo "# ARGS #"
    for i in "${ALL_ARGS[@]}"; do
        echo "$i"
    done

    echo""
}

echo "ERROR: you called the error module" >&2
display_env_and_args
exit 1
