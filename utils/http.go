// Description: Generic utils functions
// Author: Pixie79
// ============================================================================
// package utils

package utils

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"time"

	tuUtils "github.com/pixie79/tiny-utils/utils"
)

// UrlToLines retrieves the contents of a URL and returns them line by line.
//
// Parameters:
// - url: the URL to retrieve the contents from.
// - username: the username for basic authentication. If not needed, leave it empty.
// - password: the password for basic authentication. If not needed, leave it empty.
//
// Returns:
// - lines: an array of strings containing the lines of the retrieved content.
func UrlToLines(url string, username string, password string) []string {
	client := http.Client{Timeout: 5 * time.Second}

	req, err := http.NewRequest(http.MethodGet, url, http.NoBody)
	tuUtils.MaybeDie(err, "could not create http request")

	// Add basic auth if username and password are set
	if username != "" && password != "" {
		req.SetBasicAuth(username, password)
	}

	res, err := client.Do(req)
	tuUtils.MaybeDie(err, "could not authenticate")

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		tuUtils.MaybeDie(err, "error closing connection")
	}(res.Body)

	if !tuUtils.InBetween(res.StatusCode, 200, 299) {
		tuUtils.Die(fmt.Sprintf("url access error %s, Status Code: %d", url, res.StatusCode))
	}

	return LinesFromReader(res.Body)
}

// LinesFromReader returns an array of strings representing each line read from the provided io.Reader.
//
// The function takes an io.Reader as a parameter and scans it line by line using a bufio.Scanner.
// Each line is then appended to the `lines` array.
// After scanning is complete, the function checks for any errors and calls the MaybeDie function if there is any error.
// Finally, the `lines` array is returned.
func LinesFromReader(r io.Reader) []string {
	var lines []string

	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	err := scanner.Err()
	tuUtils.MaybeDie(err, "could not parse lines")

	return lines
}
