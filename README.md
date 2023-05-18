[![Build Status](https://travis-ci.org/warrensbox/tgswitch.svg?branch=master)](https://travis-ci.org/warrensbox/tgswitch)
[![Go Report Card](https://goreportcard.com/badge/github.com/warrensbox/tgswitch)](https://goreportcard.com/report/github.com/warrensbox/tgswitch)
[![CircleCI](https://circleci.com/gh/warrensbox/tgswitch/tree/master.svg?style=shield&circle-token=d74b0de145c45b1d0da97f817363c77350e1a121)](https://circleci.com/gh/warrensbox/tgswitch)

# Terragrunt Switcher 

<img style="text-allign:center" src="https://kepler-images.s3.us-east-2.amazonaws.com/warrensbox/tgswitch/tgswitch-banner.png" alt="drawing"/>


The `tgswitch` command line tool lets you switch between different versions of <a href="https://terragrunt.gruntwork.io/" target="_blank">terragrunt</a>. 
If you do not have a particular version of terragrunt installed, `tgswitch` will download the version you desire.
The installation is minimal and easy. 
Once installed, simply select the version you require from the dropdown and start using terragrunt. 

## Installation

`tgswitch` is available for MacOS and Linux based operating systems.

### Homebrew

Installation for MacOS is the easiest with Homebrew. <a href="https://brew.sh/" target="_blank">If you do not have homebrew installed, click here</a>.

```ruby
brew install warrensbox/tap/tgswitch
```

### Linux

Installation for Linux operation systems.

```sh
curl -L https://raw.githubusercontent.com/warrensbox/tgswitch/release/install.sh | bash
```

### Install from source

Alternatively, you can install the binary from the source <a href="https://github.com/warrensbox/tgswitch/releases" target="_blank">here</a>.

[Having trouble installing](https://tgswitch.warrensbox.com/Troubleshoot/).
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

1. You can also set the `TG_VERSION` environment variable to your desired terragrunt version. 
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

### Install to non-default location

By default `tfswitch` will download the Terraform binary to the user home directory under this path: `/Users/warrenveerasingam/.terraform.versions`

If you want to install the binaries outside of the home directory then you can provide the `-i` or `--install` to install Terraform binaries to a non-standard path. Useful if you want to install versions of Terraform that can be shared with multiple users.

The Terraform binaries will then be placed in the folder `.terraform.versions` under the custom install path e.g. `/opt/terraform/.terraform.versions`

```bash
tfswitch -i /opt/terraform/
```

**NOTE**

* The folder must exists before you run `tfswitch`
* The folder passed in `-i`/`--install` must be created before running `tfswtich`

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
## Issues
Please open  *issues* here:  [New Issue](https://github.com/warrensbox/tgswitch/issues)

## Upcoming Features
N/A
