
Problem:
```sh
install: can't change permissions of /usr/local/bin: Operation not permitted
```

```sh
"Unable to remove symlink. You must have SUDO privileges"
```

```sh
"Unable to create symlink. You must have SUDO privileges"
```

```sh
install: cannot create regular file '/usr/local/bin/tgswitch': Permission denied
```

Solution: You probably need to have privileges to install *tgswitch* at /usr/local/bin.

Try the following:

```sh
wget https://raw.githubusercontent.com/warrensbox/tgswitch/release/install.sh  #Get the installer on to your machine:

chmod 755 install.sh #Make installer executable

./install.sh -b $HOME/.bin #Install tgswitch in a location you have permission:

$HOME/.bin/tgswitch #test

export PATH=$PATH:$HOME/.bin #Export your .bin into your path

#You should probably add step 4 in your `.bash_profile` in your $HOME directory.

#Next, try:
`tgswitch -b $HOME/.bin/terragrunt 0.38.0` 

#or simply 

`tgswitch -b $HOME/.bin/terragrunt`


```

See the custom directory option `-b`:    
<img src="https://s3.us-east-2.amazonaws.com/kepler-images/warrensbox/tgswitch/tgswitch_v5.gif" alt="drawing" style="width: 670px;"/>    

