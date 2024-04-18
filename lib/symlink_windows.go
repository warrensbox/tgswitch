package lib

import (
	"io"
	"log"
	"os"
)

func Copy(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

//CreateSymlink : create symlink
func CreateSymlink(cwd string, dir string) {

	err := Copy(cwd, dir)
	if err != nil {
		log.Fatalf(`
		Unable to create new symlink.
		Maybe symlink already exist. Try removing existing symlink manually.
		Try running "unlink %s" to remove existing symlink.
		If error persist, you may not have the permission to create a symlink at %s.
		Error: %s
		`, dir, dir, err)
		os.Exit(1)
	}
}

//RemoveSymlink : remove symlink
func RemoveSymlink(symlinkPath string) {

	_, err := os.Lstat(symlinkPath)
	if err != nil {
		log.Fatalf(`
		Unable to stat symlink.
		Maybe symlink already exist. Try removing existing symlink manually.
		Try running "unlink %s" to remove existing symlink.
		If error persist, you may not have the permission to create a symlink at %s.
		Error: %s
		`, symlinkPath, symlinkPath, err)
		os.Exit(1)
	} else {
		errRemove := os.Remove(symlinkPath)

		if errRemove != nil {
			log.Fatalf(`
			Unable to remove symlink.
			Maybe symlink already exist. Try removing existing symlink manually.
			Try running "unlink %s" to remove existing symlink.
			If error persist, you may not have the permission to create a symlink at %s.
			Error: %s
			`, symlinkPath, symlinkPath, errRemove)
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
func ChangeSymlink(binVersionPath string, binPath string) {

	//installLocation = GetInstallLocation() //get installation location -  this is where we will put our terraform binary file
	binPath = InstallableBinLocation(binPath)

	/* remove current symlink if exist*/
	symlinkExist := CheckSymlink(binPath)
	if symlinkExist {
		RemoveSymlink(binPath)
	}

	/* set symlink to desired version */
	CreateSymlink(binVersionPath, binPath)

}
