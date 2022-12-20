package lib

import (
	"fmt"
	"sort"

	semver "github.com/hashicorp/go-version"
)

// GetSemver : returns version that will be installed based on server constaint provided
func GetSemver(tgconstraint *string, proxyUrl string) (string, error) {

	tglist := GetAppList(proxyUrl) //get list of versions
	fmt.Printf("Reading required version from constraint: %s\n", *tgconstraint)
	tfversion, err := SemVerParser(tgconstraint, tglist)
	return tfversion, err
}

// ValidateSemVer : Goes through the list of terragrunt version, return a valid tf version for contraint provided
func SemVerParser(tfconstraint *string, tflist []string) (string, error) {
	tfversion := ""
	constraints, err := semver.NewConstraint(*tfconstraint) //NewConstraint returns a Constraints instance that a Version instance can be checked against
	if err != nil {
		return "", fmt.Errorf("error parsing constraint: %s", err)
	}
	versions := make([]*semver.Version, len(tflist))
	//put tfversion into semver object
	for i, tfvals := range tflist {
		version, err := semver.NewVersion(tfvals) //NewVersion parses a given version and returns an instance of Version or an error if unable to parse the version.
		if err != nil {
			return "", fmt.Errorf("error parsing constraint: %s", err)
		}
		versions[i] = version
	}

	sort.Sort(sort.Reverse(semver.Collection(versions)))

	for _, element := range versions {
		if constraints.Check(element) { // Validate a version against a constraint
			tfversion = element.String()
			fmt.Printf("Matched version: %s\n", tfversion)
			if ValidVersionFormat(tfversion) { //check if version format is correct
				return tfversion, nil
			}
		}
	}

	PrintInvalidTFVersion()
	return "", fmt.Errorf("error parsing constraint: %s", *tfconstraint)
}

// Print invalid TF version
func PrintInvalidTFVersion() {
	fmt.Println("Version does not exist or invalid terraform version format.\n Format should be #.#.# or #.#.#-@# where # are numbers and @ are word characters.\n For example, 0.11.7 and 0.11.9-beta1 are valid versions")
}

// Function that check constraint for darwin M1. Terragrunt started release arm64 versions for linux and darwin OS from version 0.28.12 included.
// However, amd64 versions work on darwin arm64. To be tested on linux platforms.
func CheckDarwinArm64VersionConstraint(tgversion string, goarch string, goos string) (bool, error) {
	version, err := semver.NewVersion(tgversion)
	darwinM1constraint, err := semver.NewConstraint("< 0.28.12")

	return darwinM1constraint.Check(version) && goarch == "arm64" && goos == "darwin", err
}
