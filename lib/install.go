package lib

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"regexp"
	"runtime"

	"github.com/warrensbox/terragrunt-switcher/modal"
)

const (
	gruntURL       = "https://github.com/gruntwork-io/terragrunt/releases/download/"
	installFile    = "terragrunt"
	installVersion = "terragrunt_"
	binLocation    = "/usr/local/bin/terragrunt"
	installPath    = "/.terragrunt.versions/"
	recentFile     = "RECENT"
)

var (
	installLocation  = "/tmp"
	installedBinPath = "/tmp"
)

func init() {
	/* get current user */
	usr, errCurr := user.Current()
	if errCurr != nil {
		log.Fatal(errCurr)
	}

	/* set installation location */
	installLocation = usr.HomeDir + installPath

	/* set default binary path for terragrunt */
	installedBinPath = binLocation

	/* find terragrunt binary location if terragrunt is already installed*/
	cmd := NewCommand("terragrunt")
	next := cmd.Find()

	/* overrride installation default binary path if terragrunt is already installed */
	/* find the last bin path */
	for path := next(); len(path) > 0; path = next() {
		installedBinPath = path
	}
	/* Create local installation directory if it does not exist */
	CreateDirIfNotExist(installLocation)
}

//Install : Install the provided version in the argument
func Install(url string, appversion string, assests []modal.Repo) string {

	/* check if selected version already downloaded */
	fileExist := CheckFileExist(installLocation + installVersion + appversion)

	//fmt.Println(fileExist)

	/* if selected version already exist, */
	if fileExist {

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

	/* remove current symlink if exist*/
	symlinkExist := CheckSymlink(installedBinPath)

	if symlinkExist {
		RemoveSymlink(installedBinPath)
	}

	/* if selected version already exist, */
	/* proceed to download it from the terragrunt release page */
	//url := gruntURL + "v" + tgversion + "/" + "terragrunt" + "_" + goos + "_" + goarch

	goarch := runtime.GOARCH
	goos := runtime.GOOS
	urlDownload := ""

	//	fmt.Println(assests)

	for _, v := range assests {

		if v.TagName == "v"+appversion {
			if len(v.Assets) > 0 {
				for _, b := range v.Assets {
					fmt.Println(b)
					matchedOS, _ := regexp.MatchString(goos, b.BrowserDownloadURL)
					matchedARCH, _ := regexp.MatchString(goarch, b.BrowserDownloadURL)
					if matchedOS && matchedARCH {
						urlDownload = b.BrowserDownloadURL
						break
					}
				}
			}
			break
		}
	}

	fmt.Println(urlDownload)
	fmt.Println(urlDownload)

	fileInstalled, _ := DownloadFromURL(installLocation, urlDownload)

	/* rename file to terragrunt version name - terragrunt_x.x.x */
	RenameFile(fileInstalled, installLocation+installVersion+appversion)

	err := os.Chmod(installLocation+installVersion+appversion, 0755)
	if err != nil {
		log.Println(err)
	}

	/* set symlink to desired version */
	CreateSymlink(installLocation+installVersion+appversion, installedBinPath)
	fmt.Printf("Switched terragrunt to version %q \n", appversion)
	return installLocation
}

// AddRecent : add to recent file
func AddRecent(requestedVersion string, installLocation string) {

	semverRegex := regexp.MustCompile(`\d+(\.\d+){2}\z`)

	fileExist := CheckFileExist(installLocation + recentFile)
	if fileExist {
		lines, errRead := ReadLines(installLocation + recentFile)

		if errRead != nil {
			fmt.Printf("Error: %s\n", errRead)
			return
		}

		for _, line := range lines {
			if !semverRegex.MatchString(line) {
				RemoveFiles(installLocation + recentFile)
				CreateRecentFile(requestedVersion)
				return
			}
		}

		versionExist := VersionExist(requestedVersion, lines)

		if !versionExist {
			if len(lines) >= 3 {
				_, lines = lines[len(lines)-1], lines[:len(lines)-1]

				lines = append([]string{requestedVersion}, lines...)
				WriteLines(lines, installLocation+recentFile)
			} else {
				lines = append([]string{requestedVersion}, lines...)
				WriteLines(lines, installLocation+recentFile)
			}
		}

	} else {
		CreateRecentFile(requestedVersion)
	}
}

// GetRecentVersions : get recent version from file
func GetRecentVersions() ([]string, error) {

	fileExist := CheckFileExist(installLocation + recentFile)
	if fileExist {
		semverRegex := regexp.MustCompile(`\A\d+(\.\d+){2}\z`)

		lines, errRead := ReadLines(installLocation + recentFile)

		if errRead != nil {
			fmt.Printf("Error: %s\n", errRead)
			return nil, errRead
		}

		for _, line := range lines {
			if !semverRegex.MatchString(line) {
				RemoveFiles(installLocation + recentFile)
				return nil, errRead
			}
		}
		return lines, nil
	}
	return nil, nil
}

//CreateRecentFile : create a recent file
func CreateRecentFile(requestedVersion string) {
	WriteLines([]string{requestedVersion}, installLocation+recentFile)
}
