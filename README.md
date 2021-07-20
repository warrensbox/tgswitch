[![Build Status](https://travis-ci.org/warrensbox/tgswitch.svg?branch=master)](https://travis-ci.org/warrensbox/tgswitch)
[![Go Report Card](https://goreportcard.com/badge/github.com/warrensbox/tgswitch)](https://goreportcard.com/report/github.com/warrensbox/tgswitch)
[![CircleCI](https://circleci.com/gh/warrensbox/tgswitch/tree/master.svg?style=shield&circle-token=d74b0de145c45b1d0da97f817363c77350e1a121)](https://circleci.com/gh/warrensbox/tgswitch)

# Terragrunt Switcher 

<img style="text-allign:center" src="https://s3.us-east-2.amazonaws.com/kepler-images/warrensbox/tgswitch/smallerlogo.png" alt="drawing" width="110" height="140"/>


The `tgswitch` command line tool lets you switch between different versions of [terragrunt](https://www.terragrunt.io/). 
If you do not have a particular version of terragrunt installed, `tgswitch` will download the version you desire.
The installation is minimal and easy. 
Once installed, simply select the version you require from the dropdown and start using terragrunt. 


See installation guide here: [tgswitch installation](https://warrensbox.github.io/tgswitch/)

## Installation

`tgswitch` is available for MacOS and Linux based operating systems.

### Homebrew

Installation for MacOS is the easiest with Homebrew. [If you do not have homebrew installed, click here](https://brew.sh/). 


```ruby
brew install warrensbox/tap/tgswitch
```

### Linux

Installation for other linux operation systems.

```sh
curl -L https://raw.githubusercontent.com/warrensbox/tgswitch/release/install.sh | bash
```

### Install from source

Alternatively, you can install the binary from source [here](https://github.com/warrensbox/tgswitch/releases) 

## How to use:
### Use dropdown menu to select version
<img src="https://s3.us-east-2.amazonaws.com/kepler-images/warrensbox/tgswitch/tgswitch.gif" alt="drawing" style="width: 180px;"/>

1.  You can switch between different versions of terragrunt by typing the command `tgswitch` on your terminal. 
2.  Select the version of terragrunt you require by using the up and down arrow.
3.  Hit **Enter** to select the desired version.

The most recently selected versions are presented at the top of the dropdown.

### Supply version on command line
<img src="https://s3.us-east-2.amazonaws.com/kepler-images/warrensbox/tgswitch/tgswitch-v4.gif" alt="drawing" style="width: 170px;"/>

1. You can also supply the desired version as an argument on the command line.
2. For example, `tgswitch 0.10.7` for version 0.10.7 of terragrunt.
3. Hit **Enter** to switch version.

### Use custom installation location  (For non-admin - users with limited privilege on their computers)    
You can specify a custom binary path for your terragrunt installation

1. Create a custom binary path. Ex: `mkdir /Users/warrenveerasingam/bin` (replace warrenveerasingam with your username)
2. Add the path to your PATH. Ex: `export PATH=$PATH:/Users/warrenveerasingam/bin` (add this to your bash profile or zsh profile)
3. Pass -b or --bin parameter with your custom path to install terragrunt. Ex: `tgswitch -b /Users/warrenveerasingam/bin/terragrunt 0.14.1 `

### Use .tgswitchrc file
<img src="https://s3.us-east-2.amazonaws.com/kepler-images/warrensbox/tgswitch/manual-tgswitchrc.gif" alt="drawing" style="width: 170px;"/>

1. Create a `.tgswitchrc` file containing the desired version
2. For example, `echo "0.14.1" >> .tgswitchrc` for version 0.14.1 of terragrunt
3. Run the command `tgswitch` in the same directory as your `.tgswitchrc`

#### *Instead of a `.tgswitchrc` file, a `.terragrunt-version` file may be used for compatibility with [`tgenv`](https://github.com/cunymatthieu/tgenv#terragrunt-version) and other tools which use it*

**Automatically switch with bash**

Add the following to the end of your `~/.bashrc` file:
(Use either `.tgswitchrc` or `.terragrunt-version`)

```
cdtgswitch(){
  builtin cd "$@";
  cdir=$PWD;
  if [ -f "$cdir/.tgswitchrc" ]; then
    tgswitch
  fi
}
alias cd='cdtgswitch'
```

<img src="https://s3.us-east-2.amazonaws.com/kepler-images/warrensbox/tgswitch/auto-tgswitchrc.gif" alt="drawing" style="width: 170px;"/>   

**Automatically switch with zsh**

Add the following to the end of your `~/.zshrc` file:

```
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
```
cd(){
  builtin cd "$@";
  cdir=$PWD;
  if [ -f "$cdir/.tgswitchrc" ]; then
    tgswitch
  fi
}
```

**Automatically switch with fish**

Add the following to your `~/.config/fish/config.fish` file:

```
function cdtgswitch
  builtin cd "$argv"
  set cdir $PWD
  if test -f "$cdir/.tgswitchrc"
    tgswitch
  end
end
alias cd='cdtgswitch'
```

## Additional Info

See how to *upgrade*, *uninstall*, *troubleshoot* here:[More info](https://warrensbox.github.io/tgswitch/additional)


## Issues

Please open  *issues* here:  [New Issue](https://github.com/warrensbox/tgswitch/issues)
