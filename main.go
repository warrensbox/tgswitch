package main

/*
* Version 0.3.0
* Compatible with Mac OS X ONLY
 */

/*** OPERATION WORKFLOW ***/
/*
* 1- Create /usr/local/terragrunt directory if does not exist
* 2- Download binary file from url to /usr/local/terragrunt
* 3- Rename the file from `terragrunt` to `terragrunt_version`
* 4- Read the existing symlink for terragrunt (Check if it's a homebrew symlink)
* 6- Remove that symlink (Check if it's a homebrew symlink)
* 7- Create new symlink to binary  `terragrunt_version`
 */

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/pborman/getopt"
	lib "github.com/warrensbox/tgswitch/lib"
)

const (
	terragruntURL  = "https://github.com/gruntwork-io/terragrunt/releases/download/"
	defaultBin     = "/usr/local/bin/terragrunt" //default bin installation dir
	rcFilename     = ".tgswitchrc"
	tgvFilename    = ".terragrunt-version"
	installVersion = "terragrunt_"
	proxyUrl       = "https://warrensbox.github.io/terragunt-versions-list/index.json"
)

var version = "0.5.0\n"

// findVersionFile searches recursively upwards from the current working directory for a terragrunt
// version file (ie. ".tgswitchrc" or ".terragrunt-version"). It then parses this file and returns
// the version string.
func findVersionFile() string {
	var tgVersion string
	var fileContents []byte
	var err error

	cwd, err := os.Getwd()
	if err != nil {
		log.Printf("Failed to get current directory %v\n", err)
		os.Exit(1)
	}

	path := strings.Split(cwd, string(filepath.Separator))
	pathSegments := len(path)

	// range over the segemnts of the full cwd path
	for i := range path {
		// get the current segment position in the path based on loop iterator
		currentSegment := pathSegments - i

		// create an array of all path segments based on iterator position
		// since we are working backwards, we only want the segments _up to_ the
		// current iterator position
		var currentPathSegments = path[:currentSegment]

		// check current segment position for either version file type
		versionFiles := [2]string{rcFilename, tgvFilename}
		for _, fileName := range versionFiles {
			// constuct the full file path to check from the array of segments
			filePath := append(currentPathSegments, fileName)
			absPath := string(filepath.Separator) + filepath.Join(filePath...)

			// check for the version file, if found break the loop
			fileContents, _ = ioutil.ReadFile(absPath)
			if fileContents != nil {
				fmt.Printf("Found version file at %s\n", absPath)
				tgVersion = strings.TrimSuffix(string(fileContents), "\n")
				break
			}
		}

		// if found a version break the parent loop
		if tgVersion != "" {
			break
		}
	}

	return tgVersion
}

func main() {
	custBinPath := getopt.StringLong("bin", 'b', defaultBin, "Custom binary path. For example: /Users/username/bin/terragrunt")
	versionFlag := getopt.BoolLong("version", 'v', "displays the version of tgswitch")
	helpFlag := getopt.BoolLong("help", 'h', "displays help message")

	getopt.Parse()
	args := getopt.Args()

	var requestedVersion string
	var listOfVersions []string

	if *versionFlag {
		fmt.Printf("\nVersion: %v\n", version)
	} else if *helpFlag {
		usageMessage()
	} else {
		installLocation := lib.GetInstallLocation()
		if len(args) == 0 {
			requestedVersion = findVersionFile()
		} else if len(args) == 1 {
			requestedVersion = args[0]
		}

		if requestedVersion != "" && lib.ValidVersionFormat(requestedVersion) {
			fileExist := lib.CheckFileExist(installLocation + installVersion + requestedVersion)
			if fileExist {
				lib.ChangeSymlink(*custBinPath, requestedVersion)
				os.Exit(0)
			}

			// check if version exists locally before downloading it
			listOfVersions = lib.GetAppList(proxyUrl)

			// create a regex that will match a version string that itself contains a regex
			latestRegex := regexp.MustCompile(`^latest\:(.*)$`)

			// first check if the version is simply "latest" and use the first version
			// next check if the version is a valid regex that matches one or more versions in
			// the list of versions.
			if requestedVersion == "latest" {
				requestedVersion = listOfVersions[0]
			} else if latestRegex.MatchString(requestedVersion) {
				versionRegex, err := regexp.Compile(latestRegex.FindStringSubmatch(requestedVersion)[1])
				if err != nil {
					fmt.Printf("Version regex %q is not valid\n", requestedVersion)
				} else {
					for _, v := range listOfVersions {
						if versionRegex.MatchString(v) {
							requestedVersion = v
							break
						}
					}
				}
			}

			exist := lib.VersionExist(requestedVersion, listOfVersions)
			if exist {
				lib.Install(requestedVersion, *custBinPath, terragruntURL)
				os.Exit(0)
			} else {
				usageMessage()
			}
		}

		if len(listOfVersions) == 0 {
			listOfVersions = lib.GetAppList(proxyUrl)
		}
		recentVersions, _ := lib.GetRecentVersions()                 //get recent versions from RECENT file
		listOfVersions = append(recentVersions, listOfVersions...)   //append recent versions to the top of the list
		listOfVersions = lib.RemoveDuplicateVersions(listOfVersions) //remove duplicate version

		/* prompt user to select version of terragrunt */
		prompt := promptui.Select{
			Label: "Select terragrunt version",
			Items: listOfVersions,
		}

		_, requestedVersion, errPrompt := prompt.Run()
		requestedVersion = strings.TrimSuffix(requestedVersion, " *recent")

		if errPrompt != nil {
			log.Printf("Prompt failed %v\n", errPrompt)
			os.Exit(1)
		}

		lib.Install(requestedVersion, *custBinPath, terragruntURL)
		os.Exit(0)
	}
}

func usageMessage() {
	fmt.Print("\n\n")
	getopt.PrintUsage(os.Stderr)
	fmt.Println("Supply the terragrunt version as an argument, or choose from a menu")
	os.Exit(0)
}
