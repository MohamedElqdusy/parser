package process

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"parser/logger"
	"strings"
)

// downloadFile downloads the file from the url and returns the filepath locally
func downloadFile(url string) (string, error) {
	var (
		fileName string
		out      *os.File
		resp     *http.Response
		err      error
	)

	// extract the file name from the url
	if fileName, err = extractFileName(url); err != nil {
		return "", err
	}

	// Create the file with its name
	if out, err = os.Create(fileName); err != nil {
		return "", err
	}
	defer out.Close()

	logger.Info("started downloading file ", fileName)
	// Get the data
	if resp, err = http.Get(url); err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Check the download response
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("couldn't download the file with url %s [bad status: %s]", url, resp.Status)
	}

	// Writer the response body to file
	if _, err = io.Copy(out, resp.Body); err != nil {
		return "", err
	}
	logger.Info("finished downloading file ", fileName)
	return fileName, nil
}

// extractFileName creates the fileName from the fullPath
func extractFileName(fullURL string) (string, error) {
	var (
		fileURL *url.URL
		err     error
	)

	if fileURL, err = url.Parse(fullURL); err != nil {
		return "", err
	}
	path := fileURL.Path
	segments := strings.Split(path, "/")
	return segments[len(segments)-1], nil
}
