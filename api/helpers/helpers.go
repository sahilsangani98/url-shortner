package helpers

import (
	"os"
	"strings"
)

// EnforceHTTP enforces HTTP prototcol
func EnforceHTTP(url string) string {
	if url[:4] != "http" {
		return "http://" + url
	}
	return url
}

// RemoveDomainError validates domain and returns false if any domain abuse found
func RemoveDomainError(url string) bool {
	// Remove common URL prefixes such as "http://", "https://", and "www."
	// Then extract the domain part of the URL.
	// If the resulting domain matches the one defined in the environment variable "DOMAIN",
	// return false to indicate a match; otherwise, return true.

	// If the provided URL is the same as the one defined in the "DOMAIN" environment variable,
	// there is a domain error (locahost given as an input to shorten URL), so return false.
	if url == os.Getenv("DOMAIN") {
		return false
	}
	newURL := strings.Replace(url, "http://", "", 1)
	newURL = strings.Replace(newURL, "https://", "", 1)
	newURL = strings.Replace(newURL, "www.", "", 1)
	newURL = strings.Split(newURL, "/")[0]

	if newURL == os.Getenv("DOMAIN") {
		return false
	}
	return true
}
