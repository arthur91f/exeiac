#compdef exeiac

# Function to get all of the available bricks
function _exeiac_action_arguments() {
    local bricks=("${(@f)$(exeiac -l)}")
    _describe -t output 'Bricks' bricks
}

function _exeiac {
    local -a actions

    actions=(
        help
        show
        get-depends
        exec
        smart-exec
    )

    _arguments -C \
        "1: :->acts" \
        "*:: :->args" \

        case "$state" in
            (acts)
                _describe -t actions 'actions' actions
                ;;
            (*)
                _exeiac_action_arguments
                ;;
        esac
    }

    _exeiac
