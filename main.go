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
	"regexp"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/pborman/getopt"
	lib "github.com/warrensbox/tgswitch/lib"
	"github.com/warrensbox/tgswitch/modal"
)

const (
	terragruntURL = "https://api.github.com/repos/gruntwork-io/terragrunt/releases?"
	defaultBin    = "/usr/local/bin/terragrunt" //default bin installation dir
	rcFilename    = ".tgswitchrc"
	tgvFilename   = ".terragrunt-version"
)

var version = "0.2.0\n"

var CLIENT_ID = "xxx"
var CLIENT_SECRET = "xxx"

func main() {

	var client modal.Client

	client.ClientID = CLIENT_ID
	client.ClientSecret = CLIENT_SECRET

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

		if _, err := os.Stat(rcfile); err == nil && len(args) == 0 { //if there is a .tgswitchrc file, and no commmand line arguments
			fmt.Printf("Reading required terragrunt version %s \n", rcFilename)

			fileContents, err := ioutil.ReadFile(rcfile)
			if err != nil {
				fmt.Printf("Failed to read %s file. Follow the README.md instructions for setup. https://github.com/warrensbox/tgswitch/blob/master/README.md\n", rcFilename)
				fmt.Printf("Error: %s\n", err)
				os.Exit(1)
			}
			tgversion := strings.TrimSuffix(string(fileContents), "\n")
			_, assets := lib.GetAppList(terragruntURL, &client)

			if lib.ValidVersionFormat(tgversion) { //check if version is correct
				lib.Install(terragruntURL, string(tgversion), assets, *custBinPath)
			} else {
				fmt.Println("Invalid terragrunt version format. Format should be #.#.# or #.#.#-@# where # is numbers and @ is word characters. For example, 0.11.7 and 0.11.9-beta1 are valid versions")
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
			_, assets := lib.GetAppList(terragruntURL, &client)

			if lib.ValidVersionFormat(tgversion) { //check if version is correct
				lib.Install(terragruntURL, string(tgversion), assets, *custBinPath)
			} else {
				fmt.Println("Invalid terragrunt version format. Format should be #.#.# or #.#.#-@# where # is numbers and @ is word characters. For example, 0.11.7 and 0.11.9-beta1 are valid versions")
				os.Exit(1)
			}

		} else if len(args) == 1 {

			semverRegex := regexp.MustCompile(`\A\d+(\.\d+){2}\z`)
			if semverRegex.MatchString(args[0]) {
				requestedVersion := args[0]

				//check if version exist before downloading it
				tflist, assets := lib.GetAppList(terragruntURL, &client)
				exist := lib.VersionExist(requestedVersion, tflist)

				if exist {
					installLocation := lib.Install(terragruntURL, requestedVersion, assets, *custBinPath)
					lib.AddRecent(requestedVersion, installLocation) //add to recent file for faster lookup
				} else {
					fmt.Println("Not a valid terragrunt version")
				}

			} else {
				fmt.Println("Not a valid terragrunt version")
				fmt.Println("Args must be a valid terragrunt version")
				usageMessage()
			}

		} else if len(args) == 0 {

			tglist, assets := lib.GetAppList(terragruntURL, &client)
			recentVersions, _ := lib.GetRecentVersions() //get recent versions from RECENT file
			tglist = append(recentVersions, tglist...)   //append recent versions to the top of the list
			tglist = lib.RemoveDuplicateVersions(tglist) //remove duplicate version

			/* prompt user to select version of terragrunt */
			prompt := promptui.Select{
				Label: "Select terragrunt version",
				Items: tglist,
			}

			_, tgversion, errPrompt := prompt.Run()
			tgversion = strings.Trim(tgversion, " *recent")

			if errPrompt != nil {
				log.Printf("Prompt failed %v\n", errPrompt)
				os.Exit(1)
			}

			installLocation := lib.Install(terragruntURL, tgversion, assets, *custBinPath)
			lib.AddRecent(tgversion, installLocation) //add to recent file for faster lookup
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
