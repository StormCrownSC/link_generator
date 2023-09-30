package shortlink

import (
	"testing"
)

var (
	dbType      string
	linkMapping map[string]string
)

func TestGenerateShortLink(t *testing.T) {
	// Mock functions for dependencies
	dbType = "in-memory"
	linkMapping = make(map[string]string)
	var secondShortLink string

	// Test case 1: Generate short link for a new original link
	originalLink := "https://example.com/original"
	shortLink, err := GenerateShortLink(originalLink)
	if err != nil {
		t.Errorf("Test case 1 failed: Unexpected error: %v", err)
	}

	// Test case 2: Generate short link for the same original link
	secondShortLink, err = GenerateShortLink(originalLink)
	if err != nil {
		t.Errorf("Test case 2 failed: Unexpected error: %v", err)
	}

	if secondShortLink != shortLink {
		t.Errorf("Test case 2 failed: Expected short link, but got: %s", shortLink)
	}

	// Test case 3: Generate short link for a different original link
	originalLink = "https://stormcrown.ru/"
	shortLink, err = GenerateShortLink(originalLink)
	if err != nil {
		t.Errorf("Test case 3 failed: Unexpected error: %v", err)
	}

	// Test case 4: Generate short link for a single character original link
	originalLink = "1"
	shortLink, err = GenerateShortLink(originalLink)
	if err != nil {
		t.Errorf("Test case 4 failed: Unexpected error: %v", err)
	}

	// Test case 5: Generate short link for a two-character original link
	originalLink = "ht"
	shortLink, err = GenerateShortLink(originalLink)
	if err != nil {
		t.Errorf("Test case 5 failed: Unexpected error: %v", err)
	}

}
