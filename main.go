package main

/*
* Version 0.3.0
* Compatible with Mac OS X ONLY
 */

/*** OPERATION WORKFLOW ***/
/*
* 1- Create /usr/local/terraform directory if does not exist
* 2- Download zip file from url to /usr/local/terraform
* 3- Unzip the file to /usr/local/terraform
* 4- Rename the file from `terraform` to `terraform_version`
* 5- Remove the downloaded zip file
* 6- Read the existing symlink for terraform (Check if it's a homebrew symlink)
* 7- Remove that symlink (Check if it's a homebrew symlink)
* 8- Create new symlink to binary  `terraform_version`
 */

import (
	"fmt"
	"log"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/pborman/getopt"
	lib "github.com/warrensbox/terragrunt-switcher/lib"
)

const (
	hashiURL = "https://api.github.com/repos/gruntwork-io/terragrunt/releases"
)

var version = "0.3.0\n"

func main() {

	getopt.Parse()
	args := getopt.Args()

	if len(args) == 0 {

		tflist, _ := lib.GetTGList(hashiURL)
		recentVersions, _ := lib.GetRecentVersions() //get recent versions from RECENT file
		tflist = append(recentVersions, tflist...)   //append recent versions to the top of the list
		tflist = lib.RemoveDuplicateVersions(tflist) //remove duplicate version

		/* prompt user to select version of terraform */
		prompt := promptui.Select{
			Label: "Select Terraform version",
			Items: tflist,
		}

		_, tfversion, errPrompt := prompt.Run()

		if errPrompt != nil {
			log.Printf("Prompt failed %v\n", errPrompt)
			os.Exit(1)
		}

		fmt.Printf("Terraform version %q selected\n", tfversion)
		lib.AddRecent(tfversion) //add to recent file for faster lookup
		lib.Install(tfversion)

	} else {
		usageMessage()
	}

}

func usageMessage() {
	fmt.Print("\n\n")
	getopt.PrintUsage(os.Stderr)
	fmt.Println("Supply the terraform version as an argument, or choose from a menu")
}
