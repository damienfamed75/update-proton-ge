package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"io"
	"strings"

	"github.com/rs/zerolog/log"
)

// Response contains release information from github. Most importantly the names
// of files to download and their corresponding download URLs.
type Response struct {
	Assets []struct{
		Name string `json:"name"`
		DownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

// FileNames returns a slice of file names belonging to this response's assets.
func (r *Response) FileNames() []string {
	names := make([]string, len(r.Assets))
	for i := range r.Assets {
		names[i] = r.Assets[i].Name
	}
	return names
}

// TarName returns the tar ball file name.
func (r *Response) TarName() string {
	return r.findName("tar.gz")
}

// Sha512SumName returns the sha512sum file name.
func (r *Response) Sha512SumName() string {
	return r.findName("sha512sum")
}

func (r *Response) findName(containing string) string {
	for i := range r.Assets {
		if strings.Contains(r.Assets[i].Name, containing) {
			return r.Assets[i].Name
		}
	}
	return ""
}

// Downloads the latest release information such as where to download its tar ball and sha512sum files.
func downloadLatestReleaseInfo() (*Response, error) {
	res, err := http.Get("https://api.github.com/repos/GloriousEggroll/proton-ge-custom/releases/latest")
	if err != nil {
		return nil, err
	}

	var response Response
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return &response, nil
}

// Returns the file names that were downloaded and an error
func downloadLatestProtonRelease(releaseInfo *Response) error {
	var sha512FileName string // save the file name of our sha512sum so we can verify what we downloaded.
	for _, asset := range releaseInfo.Assets {
		log.Debug().Str("name", asset.Name).Msg("downloading...")
		if strings.Contains(asset.Name, "sha512") {
			sha512FileName = asset.Name
		}

		log.Trace().Str("fileName", asset.Name).Msg("creating file")
		outputFile, err := os.Create(asset.Name) // touch the output file
		if err != nil {
			return fmt.Errorf("create %s: %w", asset.Name, err)
		}
		defer outputFile.Close()

		log.Trace().Str("url", asset.DownloadURL).Msg("downloading file contents")
		res, err := http.Get(asset.DownloadURL)
		if err != nil {
			return fmt.Errorf("get %s: %w", asset.DownloadURL, err)
		}

		n, err := io.Copy(outputFile, res.Body)
		if err != nil {
			return fmt.Errorf("copy response body: %w", err)
		}
		log.Trace().Int64("bytes", n).Msg("copied response body")
	}

	// Validate the sha512 sum
	if bb, err := runCommand("sha512sum", "-c", sha512FileName); err != nil {
		fmt.Printf("%s\n", string(bb))
		return err
	}

	return nil
}

