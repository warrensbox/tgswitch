package lib

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

//Repo : properties
type Repo struct {
	URL             string    `json:"url"`
	AssetsURL       string    `json:"assets_url"`
	UploadURL       string    `json:"upload_url"`
	HTMLURL         string    `json:"html_url"`
	ID              int       `json:"id"`
	NodeID          string    `json:"node_id"`
	TagName         string    `json:"tag_name"`
	TargetCommitish string    `json:"target_commitish"`
	Name            string    `json:"name"`
	Draft           bool      `json:"draft"`
	Prerelease      bool      `json:"prerelease"`
	CreatedAt       time.Time `json:"created_at"`
	PublishedAt     time.Time `json:"published_at"`
	TarballURL      string    `json:"tarball_url"`
	ZipballURL      string    `json:"zipball_url"`
	Body            string    `json:"body"`
	Author          Author    `json:"author"`
	Assets          []Assets  `json:"assets"`
}

//Author : git owner properties
type Author struct {
	Login             string `json:"login"`
	ID                int    `json:"id"`
	NodeID            string `json:"node_id"`
	AvatarURL         string `json:"avatar_url"`
	GravatarID        string `json:"gravatar_id"`
	URL               string `json:"url"`
	HTMLURL           string `json:"html_url"`
	FollowersURL      string `json:"followers_url"`
	FollowingURL      string `json:"following_url"`
	GistsURL          string `json:"gists_url"`
	StarredURL        string `json:"starred_url"`
	SubscriptionsURL  string `json:"subscriptions_url"`
	OrganizationsURL  string `json:"organizations_url"`
	ReposURL          string `json:"repos_url"`
	EventsURL         string `json:"events_url"`
	ReceivedEventsURL string `json:"received_events_url"`
	Type              string `json:"type"`
	SiteAdmin         bool   `json:"site_admin"`
}

//Author : author properties
type Assets struct {
	URL                string    `json:"url"`
	ID                 int       `json:"id"`
	NodeID             string    `json:"node_id"`
	Name               string    `json:"name"`
	Label              string    `json:"label"`
	ContentType        string    `json:"content_type"`
	State              string    `json:"state"`
	Size               int       `json:"size"`
	DownloadCount      int       `json:"download_count"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	BrowserDownloadURL string    `json:"browser_download_url"`
	Uploader           Uploader  `json:"uploader"`
}

//Uploader : repo uploader properties
type Uploader struct {
	Login             string `json:"login"`
	ID                int    `json:"id"`
	NodeID            string `json:"node_id"`
	AvatarURL         string `json:"avatar_url"`
	GravatarID        string `json:"gravatar_id"`
	URL               string `json:"url"`
	HTMLURL           string `json:"html_url"`
	FollowersURL      string `json:"followers_url"`
	FollowingURL      string `json:"following_url"`
	GistsURL          string `json:"gists_url"`
	StarredURL        string `json:"starred_url"`
	SubscriptionsURL  string `json:"subscriptions_url"`
	OrganizationsURL  string `json:"organizations_url"`
	ReposURL          string `json:"repos_url"`
	EventsURL         string `json:"events_url"`
	ReceivedEventsURL string `json:"received_events_url"`
	Type              string `json:"type"`
	SiteAdmin         bool   `json:"site_admin"`
}

type tgVersionList struct {
	tglist []string
}

var wg = sync.WaitGroup{}

var numPages = 5

//GetTGList :  Get the list of available terraform version given the hashicorp url
func GetTGList(gruntURL string) ([]string, error) {

	gswitch := http.Client{
		Timeout: time.Second * 2, // Maximum of 2 secs [decresing this seem to fail]
	}

	req, err := http.NewRequest(http.MethodGet, gruntURL, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "terragrunt-switcher")

	resp, _ := gswitch.Do(req)
	links := resp.Header.Get("Link")
	link := strings.Split(links, ",")

	for _, pagNum := range link {
		if strings.Contains(pagNum, "last") {
			strPage := inBetween(pagNum, "page=", ">")
			page, err := strconv.Atoi(strPage)
			if err != nil {
				fmt.Println(err)
				os.Exit(2)
			}
			numPages = page
		}
	}

	tglist := getTGVersion(gruntURL, numPages)

	return tglist, nil
}

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

	return exists
}

//RemoveDuplicateVersions : remove duplicate version
func RemoveDuplicateVersions(elements []string) []string {
	// Use map to record duplicates as we find them.
	encountered := map[string]bool{}
	result := []string{}

	for v := range elements {
		if encountered[elements[v]] == true {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[elements[v]] = true
			// Append to result slice.
			result = append(result, elements[v])
		}
	}
	// Return the new slice.
	return result
}

func getVersions(gruntURLPage string) {

	gswitch := http.Client{
		Timeout: time.Second * 2, // Maximum of 2 secs [decresing this seem to fail]
	}

	req, err := http.NewRequest(http.MethodGet, gruntURLPage, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "terragrunt-switcher")

	res, getErr := gswitch.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	var repo []Repo
	jsonErr := json.Unmarshal(body, &repo)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	var tgVersionList tgVersionList

	for _, num := range repo {
		if num.Prerelease == false && num.Draft == false {
			semverRegex := regexp.MustCompile(`\Av\d+(\.\d+){2}\z`)
			if semverRegex.MatchString(num.Name) {
				trimstr := strings.Trim(num.Name, "v")
				tgVersionList.tglist = append(tgVersionList.tglist, trimstr)
			}
		}

	}
	fmt.Println(tgVersionList.tglist)
}

func inBetween(value string, a string, b string) string {
	// Get substring between two strings.
	posFirst := strings.Index(value, a)
	if posFirst == -1 {
		return ""
	}
	posLast := strings.Index(value, b)
	if posLast == -1 {
		return ""
	}
	posFirstAdjusted := posFirst + len(a)
	if posFirstAdjusted >= posLast {
		return ""
	}
	return value[posFirstAdjusted:posLast]
}

func getTGVersion(gruntURLPage string, numPages int) []string {
	version := make([]string, 0)
	ch := make(chan []string, 10)

	for i := 1; i <= numPages; i++ {
		page := strconv.Itoa(i)
		api := gruntURLPage + "?page=" + page
		wg.Add(1)
		go getTGBody(api, ch)
	}

	go func(ch chan<- []string) {
		defer close(ch)
		wg.Wait()
	}(ch)

	for i := range ch {
		version = append(version, i...)
	}

	semvers := []*Version{}

	var sortedVersion []string

	for _, v := range version {
		sv, err := NewVersion(v)
		if err != nil {
			fmt.Println(err)
		}
		semvers = append(semvers, sv)
	}

	Sort(semvers)

	for _, sv := range semvers {
		sortedVersion = append(sortedVersion, sv.String())
	}

	return sortedVersion
}

func getTGBody(gruntURLPage string, ch chan<- []string) {
	defer wg.Done()

	gswitch := http.Client{
		Timeout: time.Second * 2, // Maximum of 2 secs [decresing this seem to fail]
	}

	req, err := http.NewRequest(http.MethodGet, gruntURLPage, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "terragrunt-switcher")

	res, getErr := gswitch.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	var repo []Repo
	jsonErr := json.Unmarshal(body, &repo)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	var tgVersionList tgVersionList

	for _, num := range repo {
		if num.Prerelease == false && num.Draft == false {
			semverRegex := regexp.MustCompile(`\Av\d+(\.\d+){2}\z`)
			if semverRegex.MatchString(num.Name) {
				trimstr := strings.Trim(num.Name, "v")
				tgVersionList.tglist = append(tgVersionList.tglist, trimstr)
			}
		}

	}
	ch <- tgVersionList.tglist
}

type Version struct {
	Major      int64
	Minor      int64
	Patch      int64
	PreRelease PreRelease
	Metadata   string
}

type PreRelease string

func splitOff(input *string, delim string) (val string) {
	parts := strings.SplitN(*input, delim, 2)

	if len(parts) == 2 {
		*input = parts[0]
		val = parts[1]
	}

	return val
}

func New(version string) *Version {
	return Must(NewVersion(version))
}

func NewVersion(version string) (*Version, error) {
	v := Version{}

	if err := v.Set(version); err != nil {
		return nil, err
	}

	return &v, nil
}

// Must is a helper for wrapping NewVersion and will panic if err is not nil.
func Must(v *Version, err error) *Version {
	if err != nil {
		panic(err)
	}
	return v
}

// Set parses and updates v from the given version string. Implements flag.Value
func (v *Version) Set(version string) error {
	metadata := splitOff(&version, "+")
	preRelease := PreRelease(splitOff(&version, "-"))
	dotParts := strings.SplitN(version, ".", 3)

	if len(dotParts) != 3 {
		return fmt.Errorf("%s is not in dotted-tri format", version)
	}

	if err := validateIdentifier(string(preRelease)); err != nil {
		return fmt.Errorf("failed to validate pre-release: %v", err)
	}

	if err := validateIdentifier(metadata); err != nil {
		return fmt.Errorf("failed to validate metadata: %v", err)
	}

	parsed := make([]int64, 3, 3)

	for i, v := range dotParts[:3] {
		val, err := strconv.ParseInt(v, 10, 64)
		parsed[i] = val
		if err != nil {
			return err
		}
	}

	v.Metadata = metadata
	v.PreRelease = preRelease
	v.Major = parsed[0]
	v.Minor = parsed[1]
	v.Patch = parsed[2]
	return nil
}

func (v Version) String() string {
	var buffer bytes.Buffer

	fmt.Fprintf(&buffer, "%d.%d.%d", v.Major, v.Minor, v.Patch)

	if v.PreRelease != "" {
		fmt.Fprintf(&buffer, "-%s", v.PreRelease)
	}

	if v.Metadata != "" {
		fmt.Fprintf(&buffer, "+%s", v.Metadata)
	}

	return buffer.String()
}

func (v *Version) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var data string
	if err := unmarshal(&data); err != nil {
		return err
	}
	return v.Set(data)
}

func (v Version) MarshalJSON() ([]byte, error) {
	return []byte(`"` + v.String() + `"`), nil
}

func (v *Version) UnmarshalJSON(data []byte) error {
	l := len(data)
	if l == 0 || string(data) == `""` {
		return nil
	}
	if l < 2 || data[0] != '"' || data[l-1] != '"' {
		return errors.New("invalid semver string")
	}
	return v.Set(string(data[1 : l-1]))
}

// Compare tests if v is less than, equal to, or greater than versionB,
// returning -1, 0, or +1 respectively.
func (v Version) Compare(versionB Version) int {
	if cmp := recursiveCompare(v.Slice(), versionB.Slice()); cmp != 0 {
		return cmp
	}
	return preReleaseCompare(v, versionB)
}

// Equal tests if v is equal to versionB.
func (v Version) Equal(versionB Version) bool {
	return v.Compare(versionB) == 0
}

// LessThan tests if v is less than versionB.
func (v Version) LessThan(versionB Version) bool {
	return v.Compare(versionB) < 0
}

// Slice converts the comparable parts of the semver into a slice of integers.
func (v Version) Slice() []int64 {
	return []int64{v.Major, v.Minor, v.Patch}
}

func (p PreRelease) Slice() []string {
	preRelease := string(p)
	return strings.Split(preRelease, ".")
}

func preReleaseCompare(versionA Version, versionB Version) int {
	a := versionA.PreRelease
	b := versionB.PreRelease

	/* Handle the case where if two versions are otherwise equal it is the
	 * one without a PreRelease that is greater */
	if len(a) == 0 && (len(b) > 0) {
		return 1
	} else if len(b) == 0 && (len(a) > 0) {
		return -1
	}

	// If there is a prerelease, check and compare each part.
	return recursivePreReleaseCompare(a.Slice(), b.Slice())
}

func recursiveCompare(versionA []int64, versionB []int64) int {
	if len(versionA) == 0 {
		return 0
	}

	a := versionA[0]
	b := versionB[0]

	if a > b {
		return 1
	} else if a < b {
		return -1
	}

	return recursiveCompare(versionA[1:], versionB[1:])
}

func recursivePreReleaseCompare(versionA []string, versionB []string) int {
	// A larger set of pre-release fields has a higher precedence than a smaller set,
	// if all of the preceding identifiers are equal.
	if len(versionA) == 0 {
		if len(versionB) > 0 {
			return -1
		}
		return 0
	} else if len(versionB) == 0 {
		// We're longer than versionB so return 1.
		return 1
	}

	a := versionA[0]
	b := versionB[0]

	aInt := false
	bInt := false

	aI, err := strconv.Atoi(versionA[0])
	if err == nil {
		aInt = true
	}

	bI, err := strconv.Atoi(versionB[0])
	if err == nil {
		bInt = true
	}

	// Numeric identifiers always have lower precedence than non-numeric identifiers.
	if aInt && !bInt {
		return -1
	} else if !aInt && bInt {
		return 1
	}

	// Handle Integer Comparison
	if aInt && bInt {
		if aI > bI {
			return 1
		} else if aI < bI {
			return -1
		}
	}

	// Handle String Comparison
	if a > b {
		return 1
	} else if a < b {
		return -1
	}

	return recursivePreReleaseCompare(versionA[1:], versionB[1:])
}

// BumpMajor increments the Major field by 1 and resets all other fields to their default values
func (v *Version) BumpMajor() {
	v.Major += 1
	v.Minor = 0
	v.Patch = 0
	v.PreRelease = PreRelease("")
	v.Metadata = ""
}

// BumpMinor increments the Minor field by 1 and resets all other fields to their default values
func (v *Version) BumpMinor() {
	v.Minor += 1
	v.Patch = 0
	v.PreRelease = PreRelease("")
	v.Metadata = ""
}

// BumpPatch increments the Patch field by 1 and resets all other fields to their default values
func (v *Version) BumpPatch() {
	v.Patch += 1
	v.PreRelease = PreRelease("")
	v.Metadata = ""
}

// validateIdentifier makes sure the provided identifier satisfies semver spec
func validateIdentifier(id string) error {
	if id != "" && !reIdentifier.MatchString(id) {
		return fmt.Errorf("%s is not a valid semver identifier", id)
	}
	return nil
}

// reIdentifier is a regular expression used to check that pre-release and metadata
// identifiers satisfy the spec requirements
var reIdentifier = regexp.MustCompile(`^[0-9A-Za-z-]+(\.[0-9A-Za-z-]+)*$`)

type Versions []*Version

func (s Versions) Len() int {
	return len(s)
}

func (s Versions) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s Versions) Less(i, j int) bool {
	return s[i].LessThan(*s[j])
}

// Sort sorts the given slice of Version
func Sort(versions []*Version) {
	sort.Sort(sort.Reverse(Versions(versions)))
}
