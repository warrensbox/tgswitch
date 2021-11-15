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
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
				exists = true
				return exists
			}
		}
	}

	if !exists {
		fmt.Println("Requested version does not exist")
	}

	return exists
}

//RemoveDuplicateVersions : remove duplicate version
func RemoveDuplicateVersions(elements []string) []string {
	// Use map to record duplicates as we find them.
	encountered := map[string]bool{}
	result := []string{}

	for _, val := range elements {
		versionOnly := strings.Trim(val, " *recent")
		if encountered[versionOnly] == true {
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
		Timeout: time.Second * 10, // Maximum of 10 secs [decresing this seem to fail]
	}

	req, err := http.NewRequest(http.MethodGet, gruntURLPage, nil)
	if err != nil {
		log.Fatal("Unable to make request. Please try again.")
	}

	req.Header.Set("User-Agent", "github-appinstaller")

	res, getErr := gswitch.Do(req)
	if getErr != nil {
		log.Fatal("Unable to make request Please try again.")
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Println("Unable to get release from repo ", string(body))
		log.Fatal(readErr)
	}

	var repo ListVersion
	jsonErr := json.Unmarshal(body, &repo)
	if jsonErr != nil {
		log.Println("Unable to get release from repo ", string(body))
		log.Fatal(jsonErr)
	}

	return repo.Versions
}

type ListVersion struct {
	Versions []string `json:"Versions"`
}
