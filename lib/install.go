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

// AddRecent : add to recent file
func AddRecent(requestedVersion string) {

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

	if !semverRegex.MatchString(version) {
		fmt.Println("Invalid terragrunt version format. Format should be #.#.# or #.#.#-@# where # is numbers and @ is word characters. For example, 0.11.7 and 0.11.9-beta1 are valid versions")
	}

	return semverRegex.MatchString(version)
}

//Install : Install the provided version in the argument
func Install2(tgversion string, usrBinPath string, mirrorURL string) string {

	if !ValidVersionFormat(tgversion) {
		fmt.Printf("The provided terraform version format does not exist - %s. Try `tfswitch -l` to see all available versions.\n", tgversion)
		os.Exit(1)
	}

	/* Check to see if user has permission to the default bin location which is  "/usr/local/bin/terraform"
	 * If user does not have permission to default bin location, proceed to create $HOME/bin and install the tfswitch there
	 * Inform user that they dont have permission to default location, therefore tfswitch was installed in $HOME/bin
	 * Tell users to add $HOME/bin to their path
	 */
	binPath := InstallableBinLocation(usrBinPath)

	initialize()                           //initialize path
	installLocation = GetInstallLocation() //get installation location -  this is where we will put our terraform binary file

	goarch := runtime.GOARCH
	goos := runtime.GOOS

	installFileVersionPath := ConvertExecutableExt(filepath.Join(installLocation, installVersion+tgversion))

	/* check if selected version already downloaded */
	fileExist := CheckFileExist(installLocation + installVersion + tgversion)
	if fileExist {
		installLocation := ChangeSymlink(binPath, tgversion)
		return installLocation
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

	/* rename unzipped file to terraform version name - terraform_x.x.x */
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
	//AddRecent(tgversion) //add to recent file for faster lookup
	os.Exit(0)
	return ""
}

//InstallableBinLocation : Checks if terraform is installable in the location provided by the user.
//If not, create $HOME/bin. Ask users to add  $HOME/bin to $PATH
//Return $HOME/bin as install location
func InstallableBinLocation(binLocation string) string {

	usr, errCurr := user.Current()
	if errCurr != nil {
		log.Fatal(errCurr)
	}
	pathDir := Path(binLocation)              //get path directory from binary path
	existDefaultBin := CheckDirExist(pathDir) //the default is /usr/local/bin but users can provide custom bin locations
	if existDefaultBin {                      //if exist - now see if we can write to to it

		writableToDefault := false
		if runtime.GOOS != "windows" {
			writableToDefault = CheckDirWritable(pathDir) //check if is writable on ( only works on LINUX)
		}

		if !writableToDefault {
			exisHomeBin := CheckDirExist(filepath.Join(usr.HomeDir, "bin"))
			if exisHomeBin {
				fmt.Printf("Installing terraform at %s\n", filepath.Join(usr.HomeDir, "bin"))
				return filepath.Join(usr.HomeDir, "bin", "terraform")
			}
			PrintCreateDirStmt(pathDir, filepath.Join(usr.HomeDir, "bin"))
			CreateDirIfNotExist(filepath.Join(usr.HomeDir, "bin"))
			return filepath.Join(usr.HomeDir, "bin", "terraform")
		}
		return binLocation
	}
	fmt.Printf("[Error] : Binary path does not exist: %s\n", binLocation)
	fmt.Printf("[Error] : Manually create bin directory at: %s and try again.\n", binLocation)
	os.Exit(1)
	return ""
}

func PrintCreateDirStmt(unableDir string, writable string) {
	fmt.Printf("Unable to write to: %s\n", unableDir)
	fmt.Printf("Creating bin directory at: %s\n", writable)
	fmt.Printf("RUN `export PATH=$PATH:%s` to append bin to $PATH\n", writable)
}

//ConvertExecutableExt : convert excutable with local OS extension
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
