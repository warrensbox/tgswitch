# Terragrunt Switcher 

The `tgswitch` command line tool lets you switch between different versions of [terragrunt](https://www.terragrunt.io/){:target="_blank"}. 
If you do not have a particular version of terragrunt installed, `tgswitch` will download the version you desire.
The installation is minimal and easy. 
Once installed, simply select the version you require from the dropdown and start using terragrunt. 

<hr>

## Installation

`tgswitch` is available for MacOS and Linux based operating systems.

### Homebrew

Installation for MacOS is the easiest with Homebrew. [If you do not have homebrew installed, click here](https://brew.sh/){:target="_blank"}. 


```ruby
brew install warrensbox/tap/tgswitch
```

### Linux

Installation for Linux operation systems.

```sh
curl -L https://raw.githubusercontent.com/warrensbox/tgswitch/release/install.sh | bash
```

### Install from source

Alternatively, you can install the binary from the source [here](https://github.com/warrensbox/tgswitch/releases) 

<hr>

## How to use:
### Use dropdown menu to select version
<img align="center" src="https://s3.us-east-2.amazonaws.com/kepler-images/warrensbox/tgswitch/tgswitch.gif" alt="drawing" style="width: 480px;"/>

1.  You can switch between different versions of terragrunt by typing the command `tgswitch` on your terminal. 
2.  Select the version of terragrunt you require by using the up and down arrow.
3.  Hit **Enter** to select the desired version.

The most recently selected versions are presented at the top of the dropdown.

### Supply version on command line
<img align="center" src="https://s3.us-east-2.amazonaws.com/kepler-images/warrensbox/tgswitch/tgswitch-v4.gif" alt="drawing" style="width: 480px;"/>

1. You can also supply the desired version as an argument on the command line.
2. For example, `tgswitch 0.10.5` for version 0.10.5 of terragrunt.
3. Hit **Enter** to switch.   

### Use custom installation location  (For non-admin - users with limited privilege on their computers)    
You can specify a custom binary path for your terragrunt installation

1. Create a custom binary path. Ex: `mkdir /Users/warrenveerasingam/bin` (replace warrenveerasingam with your username)
2. Add the path to your PATH. Ex: `export PATH=$PATH:/Users/warrenveerasingam/bin` (add this to your bash profile or zsh profile)
3. Pass -b or --bin parameter with your custom path to install terragrunt. Ex: `tgswitch -b /Users/warrenveerasingam/bin/terragrunt 0.14.1 `

### Use .tgswitchrc file
<img align="center" src="https://s3.us-east-2.amazonaws.com/kepler-images/warrensbox/tgswitch/manual-tgswitchrc.gif" alt="drawing" style="width: 480px;"/>

1. Create a `.tgswitchrc` file containing the desired version
2. For example, `echo "0.14.1" >> .tgswitchrc` for version 0.10.5 of terragrunt
3. Run the command `tgswitch` in the same directory as your `.tgswitchrc`

**Automatically switch with bash**

Add the following to the end of your `~/.bashrc` file:
(Use `.tgswitchrc`)

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

<img align="center" src="https://s3.us-east-2.amazonaws.com/kepler-images/warrensbox/tgswitch/auto-tgswitchrc.gif" alt="drawing" style="width: 480px;"/>

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


<hr>

## Issues

Please open  *issues* here: [New Issue](https://github.com/warrensbox/terragrunt-switcher/issues){:target="_blank"}

<hr>

See how to *upgrade*, *uninstall*, *troubleshoot* here:
[Additional Info](additional)