#!/bin/bash
ACTION="$1"
BRICK_PATH="$2"
ALL_ARGS="$@"
CURRENT_PATH="$(pwd)"

function show_implemented_actions {
    grep "^function " $0 |
        sed 's|^function \([^ ]*\) .*$|\1|g' |
        grep -v "internal"
    exit 0
}

function init {
    echo "execute init: install tools and deps"
}

function internal_interactive {
    echo "execute $ACTION: for test you can choose what's happen"
    read -p "choose: ok|drift|fail: " answer
    if [ "$answer" == "ok" ]; then
        echo "no drift"
        return 0
    elif [ "$answer" == "drift" ]; then
        echo "some diff have been found"
        return 1
    else
        echo "fail"
        return $(( $RANDOM % 255 + 2 ))
    fi
}

function plan {
    # for the test the result will depends of the presence of word in path
    if grep -q "drift$" <<<"$BRICK_PATH"; then
        echo "execute $ACTION: some diff have been found, need to be layed"
        return 1
    elif grep -q "fail$" <<<"$BRICK_PATH"; then
        echo "execute $ACTION: the $ACTION have failed"
        echo "remove the ending \"fail\" of the brick name" >&2
        return $(( $RANDOM % 255 + 2 ))
    elif grep -q ".*--non-interactive" <<<"$ALL_ARGS"; then
        echo "execute $ACTION: ok no drift"
        return 0
    elif grep -q "interactive" <<<"$BRICK_PATH"; then
        internal_interactive
        return $?
    else
        echo "execute plan: ok no drift"
        return 0
    fi
}

function lay {
    # for the test the result will depends of the presence of word in path
    if grep -q "drift$" <<<"$BRICK_PATH"; then
        echo "execute $ACTION: some diff have been found, need to be layed"
        echo "  - null_resource: test"
        echo "Do you want to lay (only \"yes\" accepted):"
        read answer
        if [ "$answer" == "yes" ]; then
            echo "lay ok"
            return 1
        else
            return 2
        fi
    elif grep -q "fail$" <<<"$BRICK_PATH"; then
        echo "execute $ACTION: the $ACTION have failed"
        echo "remove the ending \"fail\" of the brick name" >&2
        return $(( $RANDOM % 255 + 2 ))
    elif grep -q ".* --non-interactive" <<<"$ALL_ARGS"; then
        echo "execute $ACTION: ok no drift"
        return 0
    elif grep -q "interactive" <<<"$BRICK_PATH"; then
        internal_interactive
        return $?
    else
        echo "execute plan: ok no drift"
        return 0
    fi
}

function  remove {
    lay
}

cd "$BRICK_PATH"
if [ "$?" != 0 ]; then
    echo "ERROR: \`cd $BRICK_PATH\` failed" >&2
    exit 3
fi

case "$ACTION" in
    init)
        init
        ;;
    plan)
        plan
        ;;
    lay)
        lay
        ;;
    remove)
        remove
        ;;
    *)
        echo "Specified actions doesn't exist"
        ;;
esac
EXIT_CODE="$?"

cd "$CURRENT_PATH"
if [ "$?" != 0 ]; then
    echo "ERROR: \`cd $CURRENT_PATH\` failed" >&2
    exit 3
fi

exit "$EXIT_CODE"
