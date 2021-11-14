package lib

import (
	"fmt"
	"log"
	"os"
)

//CreateSymlink : create symlink
//CreateSymlink : create symlink
func CreateSymlink(cwd string, dir string) {

	err := os.Symlink(cwd, dir)
	if err != nil {
		log.Fatalf(`
		Unable to create new symlink.
		Maybe symlink already exist. Try removing existing symlink manually.
		Try running "unlink" to remove existing symlink.
		If error persist, you may not have the permission to create a symlink at %s.
		Error: %s
		`, dir, err)
		os.Exit(1)
	}
}

//RemoveSymlink : remove symlink
func RemoveSymlink(symlinkPath string) {

	_, err := os.Lstat(symlinkPath)
	if err != nil {
		log.Fatalf(`
		Unable to remove symlink.
		Maybe symlink already exist. Try removing existing symlink manually.
		Try running "unlink" to remove existing symlink.
		If error persist, you may not have the permission to create a symlink at %s.
		Error: %s
		`, symlinkPath, err)
		os.Exit(1)
	} else {
		errRemove := os.Remove(symlinkPath)
		if errRemove != nil {
			log.Fatalf(`
			Unable to remove symlink.
			Maybe symlink already exist. Try removing existing symlink manually.
			Try running "unlink" to remove existing symlink.
			If error persist, you may not have the permission to create a symlink at %s.
			Error: %s
			`, symlinkPath, errRemove)
			os.Exit(1)
		}
	}
}

// CheckSymlink : check file is symlink
func CheckSymlink(symlinkPath string) bool {

	fi, err := os.Lstat(symlinkPath)
	if err != nil {
		return false
	}

	if fi.Mode()&os.ModeSymlink != 0 {
		return true
	}

	return false
}

// ChangeSymlink : move symlink to existing binary
func ChangeSymlink(installedBinPath string, appversion string) string {

	installLocation = GetInstallLocation() //get installation location -  this is where we will put our terraform binary file

	/* remove current symlink if exist*/
	symlinkExist := CheckSymlink(installedBinPath)
	if symlinkExist {
		RemoveSymlink(installedBinPath)
	}
	/* set symlink to desired version */
	CreateSymlink(installLocation+installVersion+appversion, installedBinPath)
	fmt.Printf("Switched terragrunt to version %q \n", appversion)
	return installLocation

}
