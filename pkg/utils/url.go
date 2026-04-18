package utils

import (
	"net/url"
	"strings"
)

func NormalizeURL(base string, href string) (string, error) {
	// ignore empty
	if href == "" {
		return "", nil
	}

	// ignore mailto, tel, javascript
	if strings.HasPrefix(href, "mailto:") ||
		strings.HasPrefix(href, "tel:") ||
		strings.HasPrefix(href, "javascript:") {
		return "", nil
	}

	baseURL, err := url.Parse(base)
	if err != nil {
		return "", err
	}

	parsed, err := url.Parse(href)
	if err != nil {
		return "", err
	}

	resolved := baseURL.ResolveReference(parsed)

	// remove fragments (#section)
	resolved.Fragment = ""

	return resolved.String(), nil
}

func IsSameDomain(base string, target string) bool {
	baseURL, err1 := url.Parse(base)
	targetURL, err2 := url.Parse(target)

	if err1 != nil || err2 != nil {
		return false
	}

	return baseURL.Host == targetURL.Host
}