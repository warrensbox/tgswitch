package lib_test

import (
	"log"
	"testing"

	"github.com/warrensbox/tgswitch/lib"
)

const (
	gruntURL = "https://api.github.com/repos/gruntwork-io/terragrunt/releases"
)

//TestRemoveDuplicateVersions :  test to removed duplicate
func TestRemoveDuplicateVersions(t *testing.T) {

	test_array := []string{"0.0.1", "0.0.2", "0.0.3", "0.0.1"}

	list := lib.RemoveDuplicateVersions(test_array)

	if len(list) == len(test_array) {
		log.Fatalf("Not able to remove duplicate: %s\n", test_array)
	} else {
		t.Log("Write versions exist (expected)")
	}
}

//TestValidVersionFormat : test if func returns valid version format
// more regex testing at https://rubular.com/r/UvWXui7EU2icSb
func TestValidVersionFormat(t *testing.T) {

	var version string
	version = "0.11.8"

	valid := lib.ValidVersionFormat(version)

	if valid == true {
		t.Logf("Valid version format : %s (expected)", version)
	} else {
		log.Fatalf("Failed to verify version format: %s\n", version)
	}

	version = "1.11.9"

	valid = lib.ValidVersionFormat(version)

	if valid == true {
		t.Logf("Valid version format : %s (expected)", version)
	} else {
		log.Fatalf("Failed to verify version format: %s\n", version)
	}

	version = "1.11.a"

	valid = lib.ValidVersionFormat(version)

	if valid == false {
		t.Logf("Invalid version format : %s (expected)", version)
	} else {
		log.Fatalf("Failed to verify version format: %s\n", version)
	}

	version = "22323"

	valid = lib.ValidVersionFormat(version)

	if valid == false {
		t.Logf("Invalid version format : %s (expected)", version)
	} else {
		log.Fatalf("Failed to verify version format: %s\n", version)
	}

	version = "@^&*!)!"

	valid = lib.ValidVersionFormat(version)

	if valid == false {
		t.Logf("Invalid version format : %s (expected)", version)
	} else {
		log.Fatalf("Failed to verify version format: %s\n", version)
	}

	version = "latest"

	valid = lib.ValidVersionFormat(version)

	if valid == true {
		t.Logf("Valid version format : %s (expected)", version)
	} else {
		log.Fatalf("Failed to verify version format: %s\n", version)
	}

	version = "latest:^0.35"

	valid = lib.ValidVersionFormat(version)

	if valid == true {
		t.Logf("Valid version format : %s (expected)", version)
	} else {
		log.Fatalf("Failed to verify version format: %s\n", version)
	}

	version = "latest:^0.35.13"

	valid = lib.ValidVersionFormat(version)

	if valid == true {
		t.Logf("Valid version format : %s (expected)", version)
	} else {
		log.Fatalf("Failed to verify version format: %s\n", version)
	}
}
