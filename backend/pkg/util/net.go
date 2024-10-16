package util

import (
	"net/url"
)

func SameOrigin(urlA string, urlB string) bool {
	parsedURLA, err := url.Parse(urlA)
	if err != nil {
		return false
	}

	parsedURLB, err := url.Parse(urlB)
	if err != nil {
		return false
	}

	if parsedURLA.Scheme != parsedURLB.Scheme {
		return false
	}

	if parsedURLA.Host != parsedURLB.Host {
		return false
	}

	return true
}
