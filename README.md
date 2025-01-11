# gsp
A simple tool for quickly selecting projects.

<img alt="Demo" src="examples/demo.gif" width="600" />

## Installation
To download the official binary, please visit the [releases page](https://github.com/gabefiori/gsp/releases). 

You will also need to have one of the supported fuzzy finders installed: [fzf](https://github.com/junegunn/fzf), [fzy](https://github.com/jhawthorn/fzy), or [skim](https://github.com/skim-rs/skim).

Once the installation is complete, you can use the `gsp` command along with other commands in your shell.
### Examples with `cd`:

<details>
<summary>Bash</summary>

> Add to your `.bashrc` file:
>
> ```sh
> alias sp='gsp_dir=$(gsp) && [ -n "$gsp_dir" ] && cd "$gsp_dir"'
> ```

</details>

<details>
<summary>Zsh</summary>

> Add to your `.zshrc` file:
>
> ```sh
> alias sp='gsp_dir=$(gsp) && [ -n "$gsp_dir" ] && cd "$gsp_dir"'
> ```

</details>

<details>
<summary>Fish</summary>

> Add to your `config.fish` file or create a new file inside the fish's `functions` folder:
>
> ```fish
> function sp
>     set dir (gsp)
>
>     if test -n "$dir"
>         cd "$dir"
>     end
> end
> ```

</details>

### Using with tmux
You can utilize this [script](/scripts/gsp-tmux.sh), which enables you to easily attach to or switch between Tmux sessions using the `gsp` command for selection.

<details>
<summary>Install</summary>

>```sh
>sudo wget -O /usr/local/bin/tms https://raw.githubusercontent.com/gabefiori/gsp/refs/heads/main/scripts/gsp-tmux.sh
>sudo chmod +x /usr/local/bin/tms
>```

</details>

## Configuration
Create a configuration file at `~/.config/gsp/config`:

```sh
# Specifies the tool used for displaying projects. 
# Available options are 'fzf', 'fzy' and 'sk'.
selector = fzf

# Specifies the order in which the entries are displayed.
# Available options are 'asc', 'desc' and 'nosort'.
sort = asc

# When set to 'true', the output will only display unique projects.
# Optional. Defaults to 'false'.
unique = false

# Determines whether the output should be expanded to show additional details. 
# Optional. Defaults to 'true'.
expand-output = true

# Sources are defined with <depth>:<path>.
# Depth must be an unsigned 8-bit integer.
source = 1:~/your/path
source = 3:/home/you/your_other/path
```

## CLI options
```sh
--config file, -c file        Load configuration from the specified file (default: "~/.config/gsp/config")
--list, -l                    Print entries to stdout (default: false)
--measure, -m                 Measure performance (time taken and number of entries processed) (default: false)
--selector value, --sl value  Selector for displaying entries (available options: 'fzf', 'fzy', 'sk')
--sort value, -s value        Specify the sort order for displaying entries (available options: 'asc', 'desc', 'nosort') (default: "nosort")
--unique, -u                  Display only unique entries (default: false)
--expand-output, --eo         Expand selection output (default: true)
--help, -h                    show help
--version, -v                 print the version
```
