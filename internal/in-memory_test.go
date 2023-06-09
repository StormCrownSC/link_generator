package main

import (
	"testing"
)

func TestIsShortLinkUniqueMemory(t *testing.T) {
	// Initialize linkMapping
	linkMapping = map[string]string{
		"https://example.com/original1": "short1",
		"https://example.com/original2": "short2",
	}

	// Test case 1: Check unique short link
	shortLink := "short3"
	unique := IsShortLinkUniqueMemory(shortLink)
	if !unique {
		t.Errorf("Test case 1 failed: Expected unique short link, but got non-unique")
	}

	// Test case 2: Check non-unique short link
	shortLink = "short1"
	unique = IsShortLinkUniqueMemory(shortLink)
	if unique {
		t.Errorf("Test case 2 failed: Expected non-unique short link, but got unique")
	}
}

func TestSaveLinkMappingMemory(t *testing.T) {
	// Initialize linkMapping
	linkMapping = make(map[string]string)

	// Test case: Save link mapping
	shortLink := "short1"
	originalLink := "https://example.com/original"
	err := SaveLinkMappingMemory(shortLink, originalLink)
	if err != nil {
		t.Errorf("Test case failed: Unexpected error: %v", err)
	}
}

func TestIsOriginalLinkExistsMemory(t *testing.T) {
	// Initialize linkMapping
	linkMapping = map[string]string{
		"https://example.com/original1": "short1",
		"https://example.com/original2": "short2",
	}

	// Test case 1: Check existing original link
	originalLink := "https://example.com/original1"
	shortLink, exist, err := IsOriginalLinkExistsMemory(originalLink)
	if err != nil {
		t.Errorf("Test case 1 failed: Unexpected error: %v", err)
	}
	if !exist {
		t.Errorf("Test case 1 failed: Expected existing original link, but got non-existing")
	}
	if shortLink != "short1" {
		t.Errorf("Test case 1 failed: Expected short link 'short1', but got '%s'", shortLink)
	}

	// Test case 2: Check non-existing original link
	originalLink = "https://example.com/nonexistent"
	shortLink, exist, err = IsOriginalLinkExistsMemory(originalLink)
	if err != nil {
		t.Errorf("Test case 2 failed: Unexpected error: %v", err)
	}
	if exist {
		t.Errorf("Test case 2 failed: Expected non-existing original link, but got existing")
	}
	if shortLink != "" {
		t.Errorf("Test case 2 failed: Expected empty short link, but got '%s'", shortLink)
	}
}

func TestGetOriginalLinkMemory(t *testing.T) {
	// Initialize linkMapping
	linkMapping = map[string]string{
		"https://example.com/original1": "short1",
		"https://example.com/original2": "short2",
	}

	// Test case 1: Get original link for existing short link
	shortLink := "short1"
	originalLink, err := GetOriginalLinkMemory(shortLink)
	if err != nil {
		t.Errorf("Test case 1 failed: Unexpected error: %v", err)
	}
	if originalLink != "https://example.com/original1" {
		t.Errorf("Test case 1 failed: Expected original link 'https://example.com/original1', but got '%s'", originalLink)
	}

	// Test case 2: Get original link for non-existing short link
	shortLink = "nonexistent"
	originalLink, err = GetOriginalLinkMemory(shortLink)
	if err == nil {
		t.Errorf("Test case 2 failed: Expected error, but got no error")
	}
	expectedErrMsg := "сокращенная ссылка не найдена"
	if err.Error() != expectedErrMsg {
		t.Errorf("Test case 2 failed: Expected error message '%s', but got '%v'", expectedErrMsg, err.Error())
	}
	if originalLink != "" {
		t.Errorf("Test case 2 failed: Expected empty original link, but got '%s'", originalLink)
	}
}
