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
