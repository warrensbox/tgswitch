package lib

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// DownloadFromURL : Downloads the binary from the source url
func DownloadFromURL(installLocation string, url string) (string, error) {

	tokens := strings.Split(url, "/")
	fileName := tokens[len(tokens)-1]
	fmt.Println("Downloading", url, "to", fileName)
	fmt.Println("Downloading ...")

	output, err := os.Create(installLocation + fileName)
	if err != nil {
		fmt.Println("Error while creating", installLocation+fileName, "-", err)
		return "", err
	}
	defer output.Close()

	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return "", err
	}
	defer response.Body.Close()

	n, errCopy := io.Copy(output, response.Body)
	if errCopy != nil {
		fmt.Println("Error while downloading", url, "-", errCopy)
		return "", errCopy
	}

	fmt.Println(n, "bytes downloaded.")
	return installLocation + fileName, nil
}
