
## How to use:
### Use dropdown menu to select version
<img src="https://kepler-images.s3.us-east-2.amazonaws.com/warrensbox/tgswitch/tgswitch_v1.gif" alt="drawing" style="width: 600px;"/>

1.  You can switch between different versions of terragrunt by typing the command `tgswitch` on your terminal.
2.  Select the version of terragrunt you require by using the up and down arrow.
3.  Hit **Enter** to select the desired version.

The most recently selected versions are presented at the top of the dropdown.

### Supply version on command line
<img src="https://s3.us-east-2.amazonaws.com/kepler-images/warrensbox/tgswitch/tgswitch_v3.gif" alt="drawing" style="width: 600px;"/>

1. You can also supply the desired version as an argument on the command line.
2. For example, `tgswitch 0.37.1` for version 0.37.1 of terragrunt.
3. Hit **Enter** to switch.

### Use environment variables
<img src="https://s3.us-east-2.amazonaws.com/kepler-images/warrensbox/tgswitch/tgswitch_v7.gif" alt="drawing" style="width: 600px;"/>
You can also set the `TG_VERSION` environment variable to your desired terragrunt version. 
For example:   
```bash
export TG_VERSION=0.37.0
tgswitch #will automatically switch to terragrunt version 0.37.0
```

### Use .tgswitch.toml file  (For non-admin AND Apple M1 users with limited privilege on their computers)
Specify a custom binary path for your terragrunt installation

<img src="https://s3.us-east-2.amazonaws.com/kepler-images/warrensbox/tgswitch/tgswitch_v5.gif" alt="drawing" style="width: 600px;"/>      
<img src="https://s3.us-east-2.amazonaws.com/kepler-images/warrensbox/tgswitch/tgswitch_v4.gif" alt="drawing" style="width: 600px;"/>

1. Create a custom binary path. Ex: `mkdir /Users/uveerum/bin` (replace uveerum with your username)
2. Add the path to your PATH. Ex: `export PATH=$PATH:/Users/uveerum/bin` (add this to your bash profile or zsh profile)
3. Pass -b or --bin parameter with your custom path to install terragrunt. Ex: `tgswitch -b /Users/uveerum/bin/terragrunt 0.34.0 `
4. Optionally, you can create a `.tgswitch.toml` file in your terragrunt directory(current directory) OR in your home directory(~/.tgswitch.toml). The toml file in the current directory has a higher precedence than toml file in the home directory
5. Your `.tgswitch.toml` file should look like this:
```ruby
bin = "/usr/local/bin/terragrunt"
version = "0.34.0"
```
4. Run `tgswitch` and it should automatically install the required terragrunt version in the specified binary path

**NOTE** 
1. For linux users that do not have write permission to `/usr/local/bin/`, `tgswitch` will attempt to install terragrunt at `$HOME/bin`. Run `export PATH=$PATH:$HOME/bin` to append bin to PATH  
2. For windows host, `tgswitch` need to be run under `Administrator` mode, and `$HOME/.tgswitch.toml` with `bin` must be defined (with a valid path) as minimum, below is an example for `$HOME/.tgswitch.toml` on windows

```toml
bin = "C:\\Users\\<%USRNAME%>\\bin\\terragrunt.exe"
```
### Use .tgswitchrc file
<img src="https://s3.us-east-2.amazonaws.com/kepler-images/warrensbox/tgswitch/tgswitch_v6.gif" alt="drawing" style="width: 600px;"/>

1. Create a `.tgswitchrc` file containing the desired version
2. For example, `echo "0.33.0" >> .tgswitchrc` for version 0.33.0 of terragrunt
3. Run the command `tgswitch` in the same directory as your `.tgswitchrc`

*Instead of a `.tgswitchrc` file, a `.terragrunt-version` file may be used as well*

### Use terragrunt.hcl file
<img src="https://s3.us-east-2.amazonaws.com/kepler-images/warrensbox/tgswitch/tgswitch_v2.gif" alt="drawing" style="width: 600px;"/>

If a terragrunt.hcl file with the terragrunt constrain is included in the current directory, it should automatically download or switch to that terragrunt version. For example, the following should automatically switch terragrunt to version 0.36.0:     
```ruby
terragrunt_version_constraint = ">= 0.36, < 0.36.1"
...
```

### Get the version from a subdirectory
```bash
tgswitch --chdir terragrunt_dir
tgswitch -c terragrunt_dir
```

**Automatically switch with bash**

Add the following to the end of your `~/.bashrc` file:
(Use either `.tgswitchrc` or `.tgswitch.toml` or `.terragrunt-version`)

```sh
cdtgswitch(){
  builtin cd "$@";
  cdir=$PWD;
  if [ -e "$cdir/.tgswitchrc" ]; then
    tgswitch
  fi
}
alias cd='cdtgswitch'
```

**Automatically switch with zsh**

Add the following to the end of your `~/.zshrc` file:

```sh
load-tgswitch() {
  local tgswitchrc_path=".tgswitchrc"

  if [ -f "$tgswitchrc_path" ]; then
    tgswitch
  fi
}
add-zsh-hook chpwd load-tgswitch
load-tgswitch
```
> NOTE: if you see an error like this: `command not found: add-zsh-hook`, then you might be on an older version of zsh (see below), or you simply need to load `add-zsh-hook` by adding this to your `.zshrc`:
>    ```
>    autoload -U add-zsh-hook
>    ```

*older version of zsh*
```sh
cd(){
  builtin cd "$@";
  cdir=$PWD;
  if [ -e "$cdir/.tgswitchrc" ]; then
    tgswitch
  fi
}
```
**Automatically switch with fish shell**

Add the following to the end of your `~/.config/fish/config.fish` file:

```sh
function switch_terragrunt --on-event fish_postexec
    string match --regex '^cd\s' "$argv" > /dev/null
    set --local is_command_cd $status

    if test $is_command_cd -eq 0 
      if count *.tf > /dev/null

        grep -c "required_version" *.tf > /dev/null
        set --local tf_contains_version $status
        
        if test $tf_contains_version -eq 0      
            command tgswitch
        end
      end
    end
end
```
