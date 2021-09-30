package lib

import (
        "fmt"
        "log"
        "os"
        "os/user"
        "regexp"
        "runtime"

        "github.com/warrensbox/tgswitch/modal"
)

const (
        gruntURL       = "https://github.com/gruntwork-io/terragrunt/releases/download/"
        installFile    = "terragrunt"
        installVersion = "terragrunt_"
        installPath    = "/.terragrunt.versions/"
        recentFile     = "RECENT"
)

var (
        installLocation = "/tmp"
)

// initialize : removes existing symlink to terragrunt binary
func initialize() {

        /* initilize default binary path for terraform */
        /* assumes that terraform is installed here */
        /* we will find the terraform path instalation later and replace this variable with the correct installed bin path */
        installedBinPath := "/usr/local/bin/terragrunt"

        /* find terragrunt binary location if terragrunt is already installed*/
        cmd := NewCommand("terragrunt")
        next := cmd.Find()

        /* overrride installation default binary path if terragrunt is already installed */
        /* find the last bin path */
        for path := next(); len(path) > 0; path = next() {
                installedBinPath = path
        }

        /* remove current symlink if exist*/
        symlinkExist := CheckSymlink(installedBinPath)

        if symlinkExist {
                RemoveSymlink(installedBinPath)
        }
}

// GetInstallLocation : get location where the terraform binary will be installed,
// will create a directory in the home location if it does not exist
func GetInstallLocation() string {
        /* get current user */
        usr, errCurr := user.Current()
        if errCurr != nil {
                log.Fatal(errCurr)
        }
        /* set installation location */
        installLocation = usr.HomeDir + installPath
        /* Create local installation directory if it does not exist */
        CreateDirIfNotExist(installLocation)
        return installLocation
}

//Install : Install the provided version in the argument
func Install(url string, appversion string, assests []modal.Repo, installedBinPath string) string {

        initialize()
        installLocation = GetInstallLocation() //get installation location -  this is where we will put our terraform binary file

        /* If user provided bin path use user one instead of default */
        // if userBinPath != nil {
        //      installedBinPath = *userBinPath
        // }

        pathDir := Path(installedBinPath)     //get path directory from binary path
        binDirExist := CheckDirExist(pathDir) //check bin path exist

        if !binDirExist {
                fmt.Printf("Binary path does not exist: %s\n", pathDir)
                fmt.Printf("Please create binary path: %s for terragrunt installation\n", pathDir)
                os.Exit(1)
        }

        /* check if selected version already downloaded */
        fileExist := CheckFileExist(installLocation + installVersion + appversion)
        if fileExist {
                installLocation := ChangeSymlink(installedBinPath, appversion)
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

        for _, v := range assests {

                if v.TagName == "v"+appversion {
                        if len(v.Assets) > 0 {
                                for _, b := range v.Assets {

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

        installLocation = GetInstallLocation()

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

        installLocation = GetInstallLocation()

        fileExist := CheckFileExist(installLocation + recentFile)
        if fileExist {
                semverRegex := regexp.MustCompile(`\A\d+(\.\d+){2}\z`)

                lines, errRead := ReadLines(installLocation + recentFile)
                outputRecent := []string{}

                if errRead != nil {
                        fmt.Printf("Error: %s\n", errRead)
                        return nil, errRead
                }

                for _, line := range lines {
                        if !semverRegex.MatchString(line) {
                                RemoveFiles(installLocation + recentFile)
                                return nil, errRead
                        }

                        /*      output can be confusing since it displays the 3 most recent used terraform version
                        append the string *recent to the output to make it more user friendly
                        */
                        outputRecent = append(outputRecent, fmt.Sprintf("%s *recent", line))
                }
                return outputRecent, nil
        }
        return nil, nil
}

//CreateRecentFile : create a recent file
func CreateRecentFile(requestedVersion string) {

        installLocation = GetInstallLocation()
        WriteLines([]string{requestedVersion}, installLocation+recentFile)
}

// ValidVersionFormat : returns valid version format
/* For example: 0.1.2 = valid
// For example: 0.1.2-beta1 = valid
// For example: 0.1.2-alpha = valid
// For example: a.1.2 = invalid
// For example: 0.1. 2 = invalid
*/
func ValidVersionFormat(version string) bool {

        // Getting versions from body; should return match /X.X.X-@/ where X is a number,@ is a word character between a-z or A-Z
        // Follow https://semver.org/spec/v1.0.0-beta.html
        // Check regular expression at https://rubular.com/r/ju3PxbaSBALpJB
        semverRegex := regexp.MustCompile(`^(\d+\.\d+\.\d+)(-[a-zA-z]+\d*)?$`)

        return semverRegex.MatchString(version)
}
