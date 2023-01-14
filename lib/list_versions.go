package lib

import (
	"context"
	"log"
	"reflect"
	"regexp"
	"strings"

	"github.com/google/go-github/v49/github"
)

// VersionExist : check if requested version exist
func VersionExist(val interface{}, array interface{}) (exists bool) {

	exists = false
	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) {
				exists = true
				return exists
			}
		}
	}

	return exists
}

// RemoveDuplicateVersions : remove duplicate version
func RemoveDuplicateVersions(elements []string) []string {
	// Use map to record duplicates as we find them.
	encountered := map[string]bool{}
	result := []string{}

	for _, val := range elements {
		versionOnly := strings.TrimSuffix(val, " *recent")
		if !encountered[versionOnly] {
			// Record this element as an encountered element.
			encountered[versionOnly] = true
			// Append to result slice.
			result = append(result, val)
		}
	}
	// Return the new slice.
	return result
}

func GetAppList(ctx context.Context, ghClient *github.Client) []string {
	repoOwner := ctx.Value("repoOwner").(string)
	repoName := ctx.Value("repoName").(string)

	opt := &github.ListOptions{Page: 1, PerPage: 100}
	releases, _, err := ghClient.Repositories.ListReleases(ctx, repoOwner, repoName, opt)
	if err != nil {
		log.Fatal("Unable to make request. Please try again.")
	}

	var re = regexp.MustCompile(`^v([0-9]+\.[0-9]+\.[0-9]+)$`)

	values := []string{}
	for _, v := range releases {
		name := *v.Name
		res := re.MatchString(name)
		if !res {
			continue
		}
		version := re.ReplaceAllString(name, `$1`)
		values = append(values, version)
	}

	return values
}

type ListVersion struct {
	Name string `json:"name"`
}
