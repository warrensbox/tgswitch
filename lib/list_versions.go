package lib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strings"
	"time"
)

//VersionExist : check if requested version exist
func VersionExist(val interface{}, array interface{}) (exists bool) {

	exists = false
	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) {
				exists = true
			}
		}
	}

	if !exists {
		fmt.Printf("Requested version %q does not exist\n", val)
	}

	return exists
}

//RemoveDuplicateVersions : remove duplicate version
func RemoveDuplicateVersions(elements []string) []string {
	// Use map to record duplicates as we find them.
	encountered := map[string]bool{}
	result := []string{}

	for _, val := range elements {
		versionOnly := strings.TrimSuffix(val, " *recent")
		if encountered[versionOnly] {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[versionOnly] = true
			// Append to result slice.
			result = append(result, val)
		}
	}
	// Return the new slice.
	return result
}

func GetAppList(gruntURLPage string) []string {

	gswitch := http.Client{
		Timeout: time.Second * 10, // Maximum of 10 secs [decreasing this seems to fail]
	}

	req, err := http.NewRequest(http.MethodGet, gruntURLPage, nil)
	if err != nil {
		log.Fatal("Unable to make request. Please try again.")
	}

	req.Header.Set("User-Agent", "github-appinstaller")

	res, getErr := gswitch.Do(req)
	if getErr != nil {
		log.Fatal("Unable to make request. Please try again.")
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatalf("Unable to get release from repo %s:\n%s", body, readErr)
	}

	var repo ListVersion
	jsonErr := json.Unmarshal(body, &repo)
	if jsonErr != nil {
		log.Fatalf("Unable to get release from repo %s:\n%s", body, jsonErr)
	}

	return repo.Versions
}

type ListVersion struct {
	Versions []string `json:"Versions"`
}
