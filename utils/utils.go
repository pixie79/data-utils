// Description: Generic utils functions
// Author: Pixie79
// ============================================================================
// package utils

package utils

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

// Die prints an error message and exits
func Die(err error, msg string) {
	Logger.Error(fmt.Sprintf("%+v: %+v", msg, err))
	os.Exit(1)
}

// MaybeDie prints an error message and exits if the error is not nil
func MaybeDie(err error, msg string) {
	if err != nil {
		Die(err, msg)
	}
}

// GetEnv Simple helper function to read an environment or return a default value
func GetEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

// GetEnvOrDie Simple helper function to read an environment or die
func GetEnvOrDie(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		Die(Err, fmt.Sprintf("missing environment variable %s", key))
	}
	return value
}

// LinesFromReader reads lines from a reader
func LinesFromReader(r io.Reader) []string {
	var lines []string

	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	err := scanner.Err()
	MaybeDie(err, "could not parse lines")

	return lines
}

// UrlToLines reads lines from a url
func UrlToLines(url string, username string, password string) []string {
	client := http.Client{Timeout: 5 * time.Second}

	req, err := http.NewRequest(http.MethodGet, url, http.NoBody)
	MaybeDie(err, "could not create http request")

	// Add basic auth if username and password are set
	if username != "" && password != "" {
		req.SetBasicAuth(username, password)
	}

	res, err := client.Do(req)
	MaybeDie(err, "could not authenticate")

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		MaybeDie(err, "error closing connection")
	}(res.Body)

	if !InBetween(res.StatusCode, 200, 299) {
		Die(fmt.Errorf("%d", res.StatusCode), fmt.Sprintf("url access error %s", url))
	}

	return LinesFromReader(res.Body)
}

// InBetween checks if a number is in between two other numbers
func InBetween(i, min, max int) bool {
	if (i >= min) && (i <= max) {
		return true
	} else {
		return false
	}
}

// ChunkBy splits a slice into chunks of a given size
func ChunkBy[T any](items []T, chunkSize int) (chunks [][]T) {
	for chunkSize < len(items) {
		items, chunks = items[chunkSize:], append(chunks, items[0:chunkSize:chunkSize])
	}
	return append(chunks, items)
}

// B64DecodeMsg decodes a base64 encoded string
func B64DecodeMsg(b64Key string, offsetF ...int) ([]byte, error) {
	offset := 7
	if len(offsetF) > 0 {
		offset = offsetF[0]
	}
	//logger.Debug(fmt.Sprintf("base64 Encoded String: %s", b64Key))
	var key []byte
	var err error
	if len(b64Key)%4 != 0 {
		key, err = base64.RawStdEncoding.DecodeString(b64Key)
	} else {
		key, err = base64.StdEncoding.DecodeString(b64Key)
	}
	if err != nil {
		return []byte{}, err
	}
	result := key[offset:]
	//logger.Debug(fmt.Sprintf("base64 Decoded String: %s", result))
	return result, nil
}

// Contains does the list contain the matching string?
func Contains(s []string, str string) bool {
	for _, v := range s {
		if strings.EqualFold(v, str) {
			return true
		}
	}
	return false
}

// DifferenceInSlices Returns
// missing from List1 but in list 2
// missing from List2 but in list 1
// common in both
func DifferenceInSlices(l1, l2 []string) ([]string, []string, []string) {
	var missingL1, missingL2, common []string
	sort.Strings(l1)
	sort.Strings(l2)
	for _, v := range l1 {
		if !Contains(l2, v) {
			missingL2 = append(missingL2, v)
		}
	}
	for _, v := range l2 {
		if !Contains(l1, v) {
			missingL1 = append(missingL1, v)
		}
	}
	for _, v := range l1 {
		if Contains(l2, v) {
			common = append(common, v)
		}
	}
	return missingL1, missingL2, common
}

// CreateBytes creates a byte array from any data
func CreateBytes(data any) []byte {
	var envBuffer bytes.Buffer
	encData := gob.NewEncoder(&envBuffer)
	err := encData.Encode(data)
	MaybeDie(err, "encoding to bytes failed")
	return envBuffer.Bytes()
}

func TimePtr(t time.Time) time.Time {
	return t
}

// CreateKey creates a key from a byte array
func CreateKey(key []byte) []byte {
	// If key is empty, use hostname as key
	if len(key) < 1 {
		return []byte(Hostname)
	} else {
		return key
	}
}
