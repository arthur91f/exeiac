#/usr/bin/env bash
_exeiac_completions()
{
    local cur prev actions
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"
    actions="init plan lay smart-lay remove help validate_code show clean get-depends"

    # We auto-complete with brick names if we already have a first argument
    if [ "$prev" != "exeiac" ]; then
        COMPREPLY=($(compgen -W "$(exeiac -l)" -- ${cur}))

        return 0
    fi

    # We auto-complete with an action name for the first argument
    COMPREPLY=( $(compgen -W "${actions}" -- ${cur}) )

    return 0
}

complete -F _exeiac_completions exeiac
