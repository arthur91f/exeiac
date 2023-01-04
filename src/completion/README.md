# Exeiac completion
A small package to provide an auto-completion to different shells (`bash`, `zsh`, `fish`).

It justs references all brick's in a configuration file's "Rooms" section, and return them.

## Bash

## Zsh
To enable zsh's auto-completion, the script needs to be stored under the
`/usr/share/zsh/site-functions/` directory under the name `_exeiac`.
```zsh
$ cp ./scripts/exeiac.zsh /usr/share/zsh/site-functions/_exeiac
```

## Fish
To enable the fish auto-completion, you can just copy the `./scripts/exeiac.fish` file
to `/usr/share/fish/vendor_completions.d/exeiac.fish`.
```fish
$ cp ./scripts/exeiac.fish /usr/share/fish/vendor_completions.d/exeiac.fish
```
You should have completion enabled once
