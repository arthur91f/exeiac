#compdef exeiac

# Function to get all of the available bricks
function _exeiac_action_arguments() {
    local bricks=("${(@f)$(exeiac -l)}")
    _describe -t output 'Bricks' bricks
}

function _exeiac {
    local -a actions

    actions=(
        init
        plan
        lay
        remove
        help
        validate_code
        show
        clean
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
