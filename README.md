# What is this?

this is a way to recreate a set of terminals and commands with tmux 

![example](./demo.gif)
made with [vhs](https://github.com/charmbracelet/vhs/)

## Why?
I just find anoying manually setting up local servers every day

## Available Commands


### create
```bash
phoemux create <alias>
```
create a config file that I like to call ash with the default values pointing to the current path
and open the "ash" with the $EDITOR env variable or vi as a default

### delete
```bash
phoemux delete <alias>
```
delete a config file (ash)

### edit
```bash
phoemux edit <alias>
```
open the "ash" with the $EDITOR env variable or vi as a default

### list
```bash
phoemux list
```
just list all the configs files (or ashes) created

### execute
```bash
phoemux <alias>
```
set up tmux session following the config file or ash related to that alias
