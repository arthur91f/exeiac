#!/bin/bash
current_path="$(pwd)"
script_path="$0"
script_dir="$(cd "$(dirname "$0")"; pwd)"
git_dir="$(cd "$(dirname "$0")" ; cd .. ; pwd)"

source "$git_dir/general_functions.sh" 
source "$git_dir/develop_functions.sh"
tests_dir="$script_dir/units-tests"
PATHS_TO_TESTS="$(ls -1 "$tests_dir")"
rooms_path="$script_dir/git-repo"
return_code=0

exeiac_path="$git_dir/exeiac.sh"

function exeiac {
    $exeiac_path "$@"
}

for file in $PATHS_TO_TESTS; do
    
    echo "-- $file ------------"
    test_commands="$(display_line_before_match \
        "$(cat "$tests_dir/$file")" "## COMMANDS RESULTS ##")"
    
    results_expected="$(display_line_after_match \
        "$(cat "$tests_dir/$file")" "## COMMANDS RESULTS ##" |
        sed '/^## COMMANDS RESULTS ##$/d'| 
        sed "s|\${ROOMS_PATH}|$rooms_path|g")"

    results="$(eval "$test_commands" 2>&1)"
    if [ "$?" != 0 ]; then
        echo "## ERROR:test command failed: $file"
        return_code=1
    fi

    diff --color <(echo "$results_expected") <(echo "$results")
    if [ "$?" != 0 ]; then
        echo "## ERROR:test outputs are different from expected: $file"
        return_code=1
    else
        echo "  test passed !"
    fi
done

exit "$return_code"

