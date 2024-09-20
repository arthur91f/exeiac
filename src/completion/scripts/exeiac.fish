set -l actions help show get-depends exec smart-exec

# Prevent file auto completion
complete -c exeiac -f

# Auto completion for sub-commands (actions)
complete -c exeiac -n "not __fish_seen_subcommand_from $commands" -a "$actions"

# Auto completion for bricks
complete -c exeiac -n "__fish_seen_subcommand_from $actions" -a "(exeiac -l)"
