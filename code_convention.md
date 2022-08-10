# Code convention

## FUNCTIONS WRITING

We use only that notation `function my_func { cmd; }` not this one `my_func() { cmd; }`
On the declaration function line we put commentary that specify arguments
On the second line we put commentaries that specifies output beginning by ->.
Of course everything is string, so we don't specify type. We specify sense of the string.

### Naming convention

We try to work only with brick_path. We display brick_name for prettinessi.
To avoid confusing: brick_path are always absolute and directory doesn't finsh with /
- function get_BLA_brick: return a brick_path
- function show_BLA_brick: return a brick_name

### Arguments documentation notation

- #< indicates arguments documentation
- positionnal arguments use _ as bash vars convention
- named arguments begins by a - and use - to be separate as current command line option. In exeiac they can get by calling get_arg.
- optionnal arg are between square brackets [] as for man documentation
- (brick_name|brick_path) indicates that you have to specify a brick but you can pass the brick name or its path.
- =(init|plan|apply) indicates short list of authorise values
  - `#< =(init|plan|apply)` is a positionnl argument that can only take init plan or apply
  - `#< -action=(init|plan|apply)` is a named arguments that can only take those 3 values

### Output documentation

- #> indicates output documenation
- ? indicates that the function is really made to be used as a boolean or that you interprete the return code.
- ~ indicates that the standard output isn't exploitable and is made to be displayed. If the function run a terraform apply or is interactive. So you shouldn't redirect the output in variable `var="$(my_fx)"` or in a file `my_fx > /tmp/myfile`
- #>2 indicates error outputs documentation (if it's only an error message don't document it.
- #>myfile indicates that an output is written in myfile
- then discribe what will be displayed in the standard output. Note that as you can't display a list so the list word indicates that it is a multiline output

So with a `egrep -r '(^function |#>|#<)'` you should easily see all options their args and outputs

### Examples

```
function is_value_in_list { #< value list [-mode=(ignore_case|strict)]
    #> ? 
    value="$1"
    list="$2"
    if mode="$(get_arg --string=mode "$@")"; then
        mode=normal
    fi
    
    case "$mode" in
    ignore_case)
        grep -qi "$value" <<<"$list"
        return $?
        ;;
    strict)
        grep -q "^$value$" <<<"$list"
        return $?
        ;;
    normal)
        grep -q "$value" <<<"$list"
        return $?
        ;;
    esac
}

fucntion display_help {
    #> ~
    cat help.txt
}
```

## CONDITION WRITING

To be easily readable by anyone that doesn't knwo well bash, each time you want to test something you have to use if or case settement.
- **Forbiden:**
    ```
    is_brick_using_this_module && module="$current_module" || echo "module not found"
    ```
- **To do instead**
    ```
    if is_brick_using_this_module ; then
    	module="$current_module"
    else
    	echo "module not found"
    fi
    ```
##


