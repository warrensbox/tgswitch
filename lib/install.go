package lib

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

const (
	gruntURL       = "https://github.com/gruntwork-io/terragrunt/releases/download/"
	installFile    = "terragrunt"
	installVersion = "terragrunt_"
	installFolder  = ".terragrunt.versions"
	recentFile     = "RECENT"
)

var (
	installLocation = "/tmp"
)

// initialize : removes existing symlink to terragrunt binary
func initialize() {

	/* initilize default binary path for terragrunt */
	/* assumes that terragrunt is installed here */
	/* we will find the terragrunt path instalation later and replace this variable with the correct installed bin path */
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

// GetInstallLocation : get location where the terragrunt binary will be installed,
// will create the installFolder if it does not exist
func GetInstallLocation(installPath string) string {
	/* set installation location */
	installLocation = filepath.Join(installPath, installFolder)

	/* Create local installation directory if it does not exist */
	CreateDirIfNotExist(installLocation)

	return installLocation

}

// AddRecent : add to recent file
func AddRecent(requestedVersion string, installPath string) {

	installLocation = GetInstallLocation(installPath)

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
				CreateRecentFile(requestedVersion, installPath)
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
		CreateRecentFile(requestedVersion, installPath)
	}
}

// GetRecentVersions : get recent version from file
func GetRecentVersions(installPath string) ([]string, error) {

	installLocation = GetInstallLocation(installPath)

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

			/*      output can be confusing since it displays the 3 most recent used terragrunt version
			append the string *recent to the output to make it more user friendly
			*/
			outputRecent = append(outputRecent, fmt.Sprintf("%s *recent", line))
		}
		return outputRecent, nil
	}
	return nil, nil
}

// CreateRecentFile : create a recent file
func CreateRecentFile(requestedVersion string, installPath string) {

	installLocation = GetInstallLocation(installPath)
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

	if !semverRegex.MatchString(version) {
		fmt.Println("Invalid terragrunt version format. Format should be #.#.# or #.#.#-@# where # is numbers and @ is word characters. For example, 0.11.7 and 0.11.9-beta1 are valid versions")
	}

	return semverRegex.MatchString(version)
}

// Install : Install the provided version in the argument
func Install(tgversion string, usrBinPath string, installPath string, mirrorURL string) string {
	/* Check to see if user has permission to the default bin location which is  "/usr/local/bin/terragrunt"
	 * If user does not have permission to default bin location, proceed to create $HOME/bin and install the tgswitch there
	 * Inform user that they dont have permission to default location, therefore tgswitch was installed in $HOME/bin
	 * Tell users to add $HOME/bin to their path
	 */
	binPath := InstallableBinLocation(usrBinPath)

	initialize()                                      //initialize path
	installLocation = GetInstallLocation(installPath) //get installation location -  this is where we will put our terragrunt binary file

	goarch := runtime.GOARCH
	goos := runtime.GOOS

	/* check if selected version already downloaded */
	installFileVersionPath := ConvertExecutableExt(filepath.Join(installLocation, installVersion+tgversion))
	fileExist := CheckFileExist(installLocation + installVersion + tgversion)

	/* if selected version already exist, */
	if fileExist {

		/* remove current symlink if exist*/
		symlinkExist := CheckSymlink(binPath)

		if symlinkExist {
			RemoveSymlink(binPath)
		}

		/* set symlink to desired version */
		CreateSymlink(installFileVersionPath, binPath)
		fmt.Printf("Switched terragrunt to version %q \n", tgversion)
		AddRecent(tgversion, installPath) //add to recent file for faster lookup
		os.Exit(0)
	}

	//if does not have slash - append slash
	hasSlash := strings.HasSuffix(mirrorURL, "/")
	if !hasSlash {
		mirrorURL = fmt.Sprintf("%s/", mirrorURL)
	}

	/* if selected version already exist, */
	/* proceed to download it from the hashicorp release page */
	url := mirrorURL + "v" + tgversion + "/" + "terragrunt" + "_" + goos + "_" + goarch

	downloadedFile, errDownload := DownloadFromURL(installLocation, url)

	/* If unable to download file from url, exit(1) immediately */
	if errDownload != nil {
		fmt.Println(errDownload)
		os.Exit(1)
	}

	/* rename unzipped file to terragrunt version name - terraform_x.x.x */
	RenameFile(downloadedFile, installFileVersionPath)

	err := os.Chmod(installFileVersionPath, 0755)
	if err != nil {
		log.Println(err)
	}
	/* remove current symlink if exist*/
	symlinkExist := CheckSymlink(binPath)

	if symlinkExist {
		RemoveSymlink(binPath)
	}

	/* set symlink to desired version */
	CreateSymlink(installFileVersionPath, binPath)
	fmt.Printf("Switched terragrunt to version %q \n", tgversion)
	AddRecent(tgversion, installPath) //add to recent file for faster lookup
	os.Exit(0)
	return ""
}

// InstallableBinLocation : Checks if terragrunt is installable in the location provided by the user.
// If not, create $HOME/bin. Ask users to add  $HOME/bin to $PATH
// Return $HOME/bin as install location
func InstallableBinLocation(userBinPath string) string {

	usr, errCurr := user.Current()
	if errCurr != nil {
		log.Fatal(errCurr)
	}

	binDir := Path(userBinPath)           //get path directory from binary path
	binPathExist := CheckDirExist(binDir) //the default is /usr/local/bin but users can provide custom bin locations

	if binPathExist == true { //if bin path exist - check if we can write to to it

		binPathWritable := false //assume bin path is not writable
		if runtime.GOOS != "windows" {
			binPathWritable = CheckDirWritable(binDir) //check if is writable on ( only works on LINUX)
		}

		// IF: "/usr/local/bin" or `custom bin path` provided by user is non-writable, (binPathWritable == false), we will attempt to install terragrunt at the ~/bin location. See ELSE
		if binPathWritable == false {

			homeBinExist := CheckDirExist(filepath.Join(usr.HomeDir, "bin")) //check to see if ~/bin exist
			if homeBinExist {                                                //if ~/bin exist, install at ~/bin/terragrunt
				fmt.Printf("Installing terragrunt at %s\n", filepath.Join(usr.HomeDir, "bin"))
				return filepath.Join(usr.HomeDir, "bin", "terragrunt")
			} else { //if ~/bin directory does not exist, create ~/bin for terragrunt installation
				fmt.Printf("Unable to write to: %s\n", userBinPath)
				fmt.Printf("Creating bin directory at: %s\n", filepath.Join(usr.HomeDir, "bin"))
				CreateDirIfNotExist(filepath.Join(usr.HomeDir, "bin")) //create ~/bin
				fmt.Printf("RUN `export PATH=$PATH:%s` to append bin to $PATH\n", filepath.Join(usr.HomeDir, "bin"))
				return filepath.Join(usr.HomeDir, "bin", "terragrunt")
			}
		} else { // ELSE: the "/usr/local/bin" or custom path provided by user is writable, we will return installable location
			return filepath.Join(userBinPath)
		}
	}
	fmt.Printf("[Error] : Binary path does not exist: %s\n", userBinPath)
	fmt.Printf("[Error] : Manually create bin directory at: %s and try again.\n", binDir)
	os.Exit(1)
	return ""
}

func PrintCreateDirStmt(unableDir string, writable string) {
	fmt.Printf("Unable to write to: %s\n", unableDir)
	fmt.Printf("Creating bin directory at: %s\n", writable)
	fmt.Printf("RUN `export PATH=$PATH:%s` to append bin to $PATH\n", writable)
}

// ConvertExecutableExt : convert excutable with local OS extension
func ConvertExecutableExt(fpath string) string {
	switch runtime.GOOS {
	case "windows":
		if filepath.Ext(fpath) == ".exe" {
			return fpath
		}
		return fpath + ".exe"
	default:
		return fpath
	}
}
