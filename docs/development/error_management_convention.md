# Error management convention

## What we want

### Display convention

X SeverityID:package/function: message[: variables]

- X: possible values:
  - "!" when you catch error
  - ">" for errors or warning that are consequence of a previous error
  - "-" when you write a new line for the same error
- Severity: possible values:
  - Warning: when you want to print something duringthe execution and the error
    is not blocking
  - Error: when your function return an error type (means the runs has failed)
    or when you encounter an error that make the execution exit
- ID: uniq ID for each error or warning that permits to identify easily message 
  in the code with a grep -r. 00000000 when it's an error from an imported 
  package. timestamp of code line writing in hexadecimal for exeiac code.
  ```bash
  printf '%x\n' $(date +%s)
  ```
- package: the package name
- function: the function name
- message: a standard message in english that let understand what have failed
- variables: optionnal value of variables to know for example on which file the
  execution has failed

### Example

```
! Error00000000: yaml: line 2: mapping values are not allowed in this context
> Warning636a5bbf:arguments/read_conf_file: unable to read file as a yaml
- Warning636a5bbf:arguments/read_conf_file: so conf file will be ignored: /etc/exeiac.conf
! Error636a540f:arguments/set_args_with_conf_files: no conf file found
> Error636a4927:arguments/GetArguments: unable to get configuration from files
> Error636a4c9e:main/main: unable to get arguments
```
We can read that output as follow:

First exeiac warned because it has ignored a configuration file 
(/etc/exeiac.conf) because it was unable to read it due to an error from
an imported package (here yaml)

Then (eventually completely independently) main function quit due to an error 
returned by GetArguments that is due to an error catched by 
set_args_with_conf_files() from package arguments.

## How to

### Catch an exeiac error
```go
err_msg := ":arguments/my_function:"
if len(conf_paths_list) < 1 {
    return conf, fmt.Errorf("! Error636a540f%s no conf file found", err_msg)
}
```

### Catch an imported package error
```go
info, err := os.Stat(path)
if err != nil {
    return conf, fmt.Errorf("! Error00000000: %w\n" +
        "> Error636a51e1%s os.Stat(%s)", err, err_msg, path)
}
```

### Follow an error
```go
content, err := get_conf(conf_file)
if err == nil {
    return conf, fmt.Errorf("%w\n> Error636a54c7%s " +
        " conf file unexploitable: %s", err, err_msg, conf_file)
} else {
    retrun conf, nil
}
```

### Print an error and exit
Try as much as possible to return error recusively to main() function from 
main package and let it print and quit.
```go
args, err := exeiacArgs.GetArguments()
if err != nil {
     fmt.Printf("%v\n> Error636a4c9e:main/main: unable to get arguments\n",
        err)
    os.Exit(1)
}
```

