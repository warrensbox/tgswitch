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

func main() {

	custBinPath := getopt.StringLong("bin", 'b', defaultBin, "Custom binary path. For example: /Users/username/bin/terragrunt")
	versionFlag := getopt.BoolLong("version", 'v', "displays the version of tgswitch")
	helpFlag := getopt.BoolLong("help", 'h', "displays help message")
	_ = versionFlag

	getopt.Parse()
	args := getopt.Args()

	dir, err := os.Getwd()
	if err != nil {
		log.Printf("Failed to get current directory %v\n", err)
		os.Exit(1)
	}

	tgvfile := dir + fmt.Sprintf("/%s", tgvFilename) //settings for .terragrunt-version file in current directory (tgenv compatible)
	rcfile := dir + fmt.Sprintf("/%s", rcFilename)   //settings for .tgswitchrc file in current directory

	if *versionFlag {
		fmt.Printf("\nVersion: %v\n", version)
	} else if *helpFlag {
		usageMessage()
	} else {
		installLocation := lib.GetInstallLocation()
		if _, err := os.Stat(rcfile); err == nil && len(args) == 0 { //if there is a .tgswitchrc file, and no commmand line arguments
			fmt.Printf("Reading required terragrunt version %s \n", rcFilename)

			fileContents, err := ioutil.ReadFile(rcfile)
			if err != nil {
				fmt.Printf("Failed to read %s file. Follow the README.md instructions for setup. https://github.com/warrensbox/tgswitch/blob/master/README.md\n", rcFilename)
				fmt.Printf("Error: %s\n", err)
				os.Exit(1)
			}
			tgversion := strings.TrimSuffix(string(fileContents), "\n")
			fileExist := lib.CheckFileExist(installLocation + installVersion + tgversion)
			if fileExist {
				lib.ChangeSymlink(*custBinPath, string(tgversion))
				os.Exit(0)
			}
			listOfVersions := lib.GetAppList(proxyUrl)

			if lib.ValidVersionFormat(tgversion) && lib.VersionExist(tgversion, listOfVersions) { //check if version format is correct && if version exist
				lib.Install(tgversion, *custBinPath, terragruntURL)
			} else {
				os.Exit(1)
			}

		} else if _, err := os.Stat(tgvfile); err == nil && len(args) == 0 {
			fmt.Printf("Reading required terragrunt version %s \n", tgvFilename)

			fileContents, err := ioutil.ReadFile(tgvfile)
			if err != nil {
				fmt.Printf("Failed to read %s file. Follow the README.md instructions for setup. https://github.com/warrensbox/tgswitch/blob/master/README.md\n", tgvFilename)
				fmt.Printf("Error: %s\n", err)
				os.Exit(1)
			}
			tgversion := strings.TrimSuffix(string(fileContents), "\n")
			fileExist := lib.CheckFileExist(installLocation + installVersion + string(tgversion))
			if fileExist {
				lib.ChangeSymlink(*custBinPath, string(tgversion))
				os.Exit(0)
			}
			listOfVersions := lib.GetAppList(proxyUrl)

			if lib.ValidVersionFormat(tgversion) && lib.VersionExist(tgversion, listOfVersions) { //check if version format is correct && if version exist
				lib.Install(tgversion, *custBinPath, terragruntURL)
			} else {
				os.Exit(1)
			}

		} else if len(args) == 1 {
			requestedVersion := args[0]

			if lib.ValidVersionFormat(requestedVersion) {

				fileExist := lib.CheckFileExist(installLocation + installVersion + string(requestedVersion))
				if fileExist {
					lib.ChangeSymlink(*custBinPath, string(requestedVersion))
					os.Exit(0)
				}

				//check if version exist before downloading it
				listOfVersions := lib.GetAppList(proxyUrl)
				exist := lib.VersionExist(requestedVersion, listOfVersions)

				if exist {
					installLocation := lib.Install(requestedVersion, *custBinPath, terragruntURL)
					fmt.Println("remove later - installLocation:", installLocation)
				}

			} else {
				fmt.Println("Args must be a valid terragrunt version")
				usageMessage()
			}

		} else if len(args) == 0 {

			listOfVersions := lib.GetAppList(proxyUrl)
			recentVersions, _ := lib.GetRecentVersions()                 //get recent versions from RECENT file
			listOfVersions = append(recentVersions, listOfVersions...)   //append recent versions to the top of the list
			listOfVersions = lib.RemoveDuplicateVersions(listOfVersions) //remove duplicate version

			/* prompt user to select version of terragrunt */
			prompt := promptui.Select{
				Label: "Select terragrunt version",
				Items: listOfVersions,
			}

			_, tgversion, errPrompt := prompt.Run()
			tgversion = strings.Trim(tgversion, " *recent")

			if errPrompt != nil {
				log.Printf("Prompt failed %v\n", errPrompt)
				os.Exit(1)
			}

			lib.Install(tgversion, *custBinPath, terragruntURL)
			os.Exit(0)
		} else {
			usageMessage()
		}
	}

}

func usageMessage() {
	fmt.Print("\n\n")
	getopt.PrintUsage(os.Stderr)
	fmt.Println("Supply the terragrunt version as an argument, or choose from a menu")
}
