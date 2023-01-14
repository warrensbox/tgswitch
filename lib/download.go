package lib

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/google/go-github/v49/github"
)

// DownloadFromURL : Downloads the binary from the source url
func DownloadFromURL(ctx context.Context, ghClient *github.Client, installLocation string, asset *github.ReleaseAsset) (string, error) {

	fileName := *asset.Name
	url := *asset.BrowserDownloadURL
	repoOwner := ctx.Value("repoOwner").(string)
	repoName := ctx.Value("repoName").(string)
	fmt.Println("Downloading", url, "to", fileName)
	fmt.Println("Downloading ...")

	output, err := os.Create(installLocation + fileName)
	if err != nil {
		fmt.Println("Error while creating", installLocation+fileName, "-", err)
		return "", err
	}
	defer output.Close()

	rc, _, err := ghClient.Repositories.DownloadReleaseAsset(ctx, repoOwner, repoName, *asset.ID, http.DefaultClient)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return "", err
	}
	defer rc.Close()

	n, errCopy := io.Copy(output, rc)
	if errCopy != nil {
		fmt.Println("Error while writing to disk", url, "-", errCopy)
		return "", errCopy
	}

	fmt.Println(n, "bytes downloaded.")
	return installLocation + fileName, nil
}
