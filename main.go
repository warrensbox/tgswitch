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
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl2/gohcl"
	"github.com/hashicorp/hcl2/hclparse"
	"github.com/manifoldco/promptui"
	"github.com/pborman/getopt"
	"github.com/spf13/viper"
	lib "github.com/warrensbox/tgswitch/lib"
)

const (
	terragruntURL = "https://github.com/gruntwork-io/terragrunt/releases/download/"
	defaultBin    = "/usr/local/bin/terragrunt" //default bin installation dir
	rcFilename    = ".tgswitchrc"
	tgvFilename   = ".terragrunt-version"
	versionPrefix = "terragrunt_"
	proxyUrl      = "https://warrensbox.github.io/terragunt-versions-list/index.json"
	tomlFilename  = ".tgswitch.toml"
	tgHclFilename = "terragrunt.hcl"
)

var version = "0.5.0\n"

func main() {

	dir := lib.GetCurrentDirectory()
	homedir := lib.GetHomeDirectory()

	custBinPath := getopt.StringLong("bin", 'b', lib.ConvertExecutableExt(defaultBin), "Custom binary path. Ex: tgswitch -b "+lib.ConvertExecutableExt("/Users/username/bin/terragrunt"))
	installPath := getopt.StringLong("install", 'i', homedir, "Custom install path. Ex: tfswitch -i /Users/username")
	versionFlag := getopt.BoolLong("version", 'v', "displays the version of tgswitch")
	helpFlag := getopt.BoolLong("help", 'h', "displays help message")
	chDirPath := getopt.StringLong("chdir", 'c', dir, "Switch to a different working directory before executing the given command. Ex: tgswitch --chdir terragrunt dir will run tgswitch in the directory")
	_ = versionFlag

	getopt.Parse()
	args := getopt.Args()

	dir, err := os.Getwd()
	if err != nil {
		log.Printf("Failed to get current directory %v\n", err)
		os.Exit(1)
	}

	TOMLConfigFile := filepath.Join(*chDirPath, tomlFilename)  //settings for .tgswitch.toml file in current directory (option to specify bin directory)
	HomeTOMLConfigFile := filepath.Join(homedir, tomlFilename) //settings for .tgswitch.toml file in home directory (option to specify bin directory)
	RCFile := filepath.Join(*chDirPath, rcFilename)            //settings for .tgswitchrc file in current directory (backward compatible purpose)
	TGVersionFile := filepath.Join(*chDirPath, tgvFilename)    //settings for .terragrunt-version file in current directory (tfenv compatible)
	TGHACLFile := filepath.Join(*chDirPath, tgHclFilename)     //settings for terragrunt.hcl file in current directory
	switch {
	case *versionFlag:
		fmt.Printf("\nVersion: %v\n", version)
	case *helpFlag:
		usageMessage()
	/* Checks if the .tgswitch.toml file exist in home or current directory
	 * This block checks to see if the tgswitch toml file is provided in the current path.
	 * If the .tgswitch.toml file exist, it has a higher precedence than the .tgswitchrc file
	 * You can specify the custom binary path and the version you desire
	 * If you provide a custom binary path with the -b option, this will override the bin value in the toml file
	 * If you provide a version on the command line, this will override the version value in the toml file
	 */
	case lib.FileExists(TOMLConfigFile) || lib.FileExists(HomeTOMLConfigFile):
		version := ""
		binPath := *custBinPath
		if lib.FileExists(TOMLConfigFile) { //read from toml from current directory
			version, binPath, *installPath = GetParamsTOML(binPath, *installPath, *chDirPath)
		} else { // else read from toml from home directory
			version, binPath, *installPath = GetParamsTOML(binPath, *installPath, homedir)
		}

		/* GIVEN A TOML FILE, */
		switch {
		case len(args) == 1:
			requestedVersion := args[0]
			if lib.ValidVersionFormat(requestedVersion) {
				//check if version exist before downloading it
				listOfVersions := lib.GetAppList(proxyUrl)
				exist := lib.VersionExist(requestedVersion, listOfVersions)

				if exist {
					installLocation := lib.Install(requestedVersion, binPath, *installPath, terragruntURL)
					fmt.Println("Install Location:", installLocation)
				}
			} else {
				fmt.Println("Args must be a valid terragrunt version")
				usageMessage()
			}
		/* provide a tgswitchrc file (IN ADDITION TO A TOML FILE) */
		case lib.FileExists(RCFile) && len(args) == 0:
			lib.ReadingFileMsg(rcFilename)
			tgversion := lib.RetrieveFileContents(RCFile)
			installVersion(tgversion, &binPath, installPath)
		/* if .terragrunt-version file found (IN ADDITION TO A TOML FILE) */
		case lib.FileExists(TGVersionFile) && len(args) == 0:
			lib.ReadingFileMsg(TGVersionFile)
			tgversion := lib.RetrieveFileContents(TGVersionFile)
			installVersion(tgversion, &binPath, installPath)
		/* if terragrunt.hcl file found (IN ADDITION TO A TOML FILE) */
		case lib.FileExists(TGHACLFile) && checkVersionDefinedHCL(&TGHACLFile) && len(args) == 0:
			installTGHclFile(&TGHACLFile, binPath, *installPath, proxyUrl)
		/* if terragrunt Version environment variable is set  (IN ADDITION TO A TOML FILE)*/
		case checkTGEnvExist() && len(args) == 0 && version == "":
			tgversion := os.Getenv("TG_VERSION")
			fmt.Printf("Terragrunt version environment variable: %s\n", tgversion)
			installVersion(tgversion, &binPath, installPath)
		/* if version is specified in the .toml file */
		case version != "":
			lib.Install(version, binPath, *installPath, terragruntURL)
		/* show dropdown */
		default:
			installFromList(&binPath, installPath)
		}

	case len(args) == 1:
		requestedVersion := args[0]
		if lib.ValidVersionFormat(requestedVersion) {
			//check if version exist before downloading it
			listOfVersions := lib.GetAppList(proxyUrl)
			exist := lib.VersionExist(requestedVersion, listOfVersions)

			if exist {
				installLocation := lib.Install(requestedVersion, *custBinPath, *installPath, terragruntURL)
				fmt.Println("Install Location:", installLocation)
			}
		} else {
			fmt.Println("Args must be a valid terragrunt version")
			usageMessage()
		}
	/* provide a tgswitchrc file */
	case lib.FileExists(RCFile) && len(args) == 0:
		lib.ReadingFileMsg(rcFilename)
		tgversion := lib.RetrieveFileContents(RCFile)
		installVersion(tgversion, custBinPath, installPath)
	case lib.FileExists(TGVersionFile) && len(args) == 0:
		lib.ReadingFileMsg(TGVersionFile)
		tgversion := lib.RetrieveFileContents(TGVersionFile)
		installVersion(tgversion, custBinPath, installPath)
	/* if terragrunt.hcl file found */
	case lib.FileExists(TGHACLFile) && checkVersionDefinedHCL(&TGHACLFile) && len(args) == 0:
		installTGHclFile(&TGHACLFile, *custBinPath, *installPath, proxyUrl)
	/* if terragrunt Version environment variable is set*/
	case checkTGEnvExist() && len(args) == 0:
		tgversion := os.Getenv("TG_VERSION")
		fmt.Printf("Terragrunt version environment variable: %s\n", tgversion)
		installVersion(tgversion, custBinPath, installPath)
	/* show dropdown */
	default:
		installFromList(custBinPath, installPath)
		os.Exit(0)
	}

}

func usageMessage() {
	fmt.Print("\n\n")
	getopt.PrintUsage(os.Stderr)
	fmt.Println("Supply the terragrunt version as an argument, or choose from a menu")
}

/* parses everything in the toml file, return required version and bin path */
func GetParamsTOML(binPath string, installPath string, dir string) (string, string, string) {
	path := lib.GetHomeDirectory()
	if dir == path {
		path = "home directory"
	} else {
		path = "current directory"
	}
	fmt.Printf("Reading configuration from %s\n", path+" for "+tomlFilename) //takes the default bin (defaultBin) if user does not specify bin path
	configfileName := lib.GetFileName(tomlFilename)                          //get the config file
	viper.SetConfigType("toml")
	viper.SetConfigName(configfileName)
	viper.AddConfigPath(dir)

	errs := viper.ReadInConfig() // Find and read the config file
	if errs != nil {
		log.Fatalf("Error: %s\nUnable to read %s provided\n", errs, tomlFilename) // Handle errors reading the config file
	}

	bin := viper.Get("bin")                                            // read custom binary location
	if binPath == lib.ConvertExecutableExt(defaultBin) && bin != nil { // if the bin path is the same as the default binary path and if the custom binary is provided in the toml file (use it)
		binPath = os.ExpandEnv(bin.(string))
	}
	//fmt.Println(binPath) //uncomment this to debug
	version := viper.Get("version") //attempt to get the version if it's provided in the toml
	if version == nil {
		version = ""
	}

	install := viper.Get("install") //attempt to get the install path if it's provided in the toml
	if install != nil {
		installPath = os.ExpandEnv(install.(string))
	}

	return version.(string), binPath, installPath
}

/* installFromList : displays & installs tf version */
func installFromList(custBinPath *string, installPath *string) {

	listOfVersions := lib.GetAppList(proxyUrl)
	recentVersions, _ := lib.GetRecentVersions(*installPath)     //get recent versions from RECENT file
	listOfVersions = append(recentVersions, listOfVersions...)   //append recent versions to the top of the list
	listOfVersions = lib.RemoveDuplicateVersions(listOfVersions) //remove duplicate version

	prompt := promptui.Select{
		Label: "Select Terragrunt version",
		Items: listOfVersions,
	}
	_, tgversion, errPrompt := prompt.Run()
	tgversion = strings.Trim(tgversion, " *recent")

	if errPrompt != nil {
		log.Printf("Prompt failed %v\n", errPrompt)
		os.Exit(1)
	}

	lib.Install(tgversion, *custBinPath, *installPath, terragruntURL)
	os.Exit(0)
}

// install with provided version as argument
func installVersion(arg string, custBinPath *string, installPath *string) {
	if lib.ValidVersionFormat(arg) {
		requestedVersion := arg

		//check to see if the requested version has been downloaded before
		installLocation := lib.GetInstallLocation(*installPath)
		installFileVersionPath := lib.ConvertExecutableExt(filepath.Join(installLocation, versionPrefix+requestedVersion))
		recentDownloadFile := lib.CheckFileExist(installFileVersionPath)
		if recentDownloadFile {
			lib.ChangeSymlink(installFileVersionPath, *custBinPath)
			fmt.Printf("Switched terragrunt to version %q \n", requestedVersion)
			lib.AddRecent(requestedVersion, *installPath) //add to recent file for faster lookup
			os.Exit(0)
		}

		//if the requested version had not been downloaded before
		//go get the list of versions
		listOfVersions := lib.GetAppList(proxyUrl)
		//check if version exist before downloading it
		exist := lib.VersionExist(requestedVersion, listOfVersions)

		if exist {
			installLocation := lib.Install(requestedVersion, *custBinPath, *installPath, terragruntURL)
			fmt.Println("Install Location:", installLocation)
		}

	} else {
		lib.PrintInvalidTGVersion()
		usageMessage()
		log.Fatalln("Args must be a valid terragrunt version")
	}
}

// checkTGEnvExist - checks if the TG_VERSION environment variable is set
func checkTGEnvExist() bool {
	return os.Getenv("TG_VERSION") != ""
}

// install using a version constraint
func installFromConstraint(tgconstraint *string, custBinPath, installPath, mirrorURL string) {

	tgversion, err := lib.GetSemver(tgconstraint, mirrorURL)
	if err == nil {
		lib.Install(tgversion, custBinPath, installPath, mirrorURL)
	}
	fmt.Println(err)
	fmt.Println("No version found to match constraint. Follow the README.md instructions for setup. https://github.com/warrensbox/tgswitch/blob/master/README.md")
	os.Exit(1)
}

// Install using version constraint from terragrunt file
func installTGHclFile(tgFile *string, custBinPath, installPath, mirrorURL string) {
	fmt.Printf("Terragrunt file found: %s\n", *tgFile)
	parser := hclparse.NewParser()
	file, diags := parser.ParseHCLFile(*tgFile) //use hcl parser to parse HCL file
	if diags.HasErrors() {
		fmt.Println("Unable to parse HCL file")
		os.Exit(1)
	}
	var version terragruntVersionConstraints
	gohcl.DecodeBody(file.Body, nil, &version)
	installFromConstraint(&version.TerragruntVersionConstraint, custBinPath, installPath, mirrorURL)
}

// check if version is defined in hcl file /* lazy-emergency fix - will improve later */
func checkVersionDefinedHCL(tgFile *string) bool {
	parser := hclparse.NewParser()
	file, diags := parser.ParseHCLFile(*tgFile) //use hcl parser to parse HCL file
	if diags.HasErrors() {
		log.Fatalf("Unable to parse HCL file: %q\n", diags.Error())
	}
	var version terragruntVersionConstraints
	gohcl.DecodeBody(file.Body, nil, &version)
	return version != (terragruntVersionConstraints{})
}

// Install using version constraint from terragrunt file
// func installTGHclFile(tgFile *string, custBinPath *string, installPath *string) {
// 	fmt.Printf("Terragrunt file found: %s\n", *tgFile)
// 	parser := hclparse.NewParser()
// 	file, diags := parser.ParseHCLFile(*tgFile) //use hcl parser to parse HCL file
// 	if diags.HasErrors() {
// 		fmt.Println("Unable to parse HCL file")
// 		os.Exit(1)
// 	}
// 	var version terragruntVersionConstraints
// 	gohcl.DecodeBody(file.Body, nil, &version)
// 	//TODO figure out sermver
// 	//installFromConstraint(&version.TerragruntVersionConstraint, custBinPath, installPath)
// }

type terragruntVersionConstraints struct {
	TerragruntVersionConstraint string `hcl:"terragrunt_version_constraint"`
}
