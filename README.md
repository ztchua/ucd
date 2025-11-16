[![go-tests](https://github.com/ztcjoe93/ucd/actions/workflows/test.yml/badge.svg?branch=main&event=workflow_dispatch)](https://github.com/ztcjoe93/ucd/actions/workflows/test.yml)

# Unnecessary chdir (UCD)
A wrapper for the common `cd` shell utility, with totally unnecessary features.

## Setup

### Compiling ucd binary

If you have golang installed in your environment, you can build the golang binary and shift it into your usr/bin directory.  
```shell
go build . && sudo chmod +x ucd && sudo mv ucd /usr/local/bin/ucd
```

Otherwise, you can download the binary and shift it into your usr/bin directory.

### Redirecting stdout to builtin shell

Append the following to your specific shell [runcom](https://en.wikipedia.org/wiki/RUNCOM) file to forward stdout from `ucd` to shell's builtin `cd` command.

Example for .zshrc  
```shell
function cd() { builtin cd $(ucd $@) }
```

## Usage

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| -h | - | - | display help |
| -v | - | - | display ucd version | 
| -a | string |  | alias for stashed path, used in conjunction with -s |
| -c | bool | false | clear history list |
| -cs | bool | false | clear stash list |
| -d | int | 0 | swap directory at -d parent directories |
| -l | - | - | display Most Recently Used (MRU) list of paths chdir-ed into |
| -ls | - | - | display list of stashed cd commands |
| -ma | int | 0 | modify alias of indicated # from the stash list |
| -p | int | 0 | chdir to the indicated # from MRU list |
| -ps | int | 0 | chdir to the indicated # from stash list |
| -pa | string |  | chdir to path with matching alias from stash list |
| -n | int | 1 | no. of times to execute chdir |
| -s | bool | false | stash cd path into a separate list |

## Configuration

On `ucd`'s first run, a `ucd.conf` JSON file is generated at `$HOME/.config/ucd`, where k-v pairs can be passed in to tweak certain features.  

| Parameter | Type | Default | Description |
| --- | --- | --- | --- |
| MaxMRUDisplay | int | -1 | Limits the total number of paths displayed when using `-l` or `-ls`. Set this to `-1` to show all records. |
| FileFallbackBehavior | bool | true | Toggle to set if the default behavior of cd-ing to a file is to fallback to its parent directory |

## Testing

To run all test suites, run the following in the repo root directory: 
```shell
go test -v ./...
```

## License

`ucd` is released under the [MIT](LICENSE.md) license.
