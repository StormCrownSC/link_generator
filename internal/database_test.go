package main

import (
	"testing"
)

func init() {
	linkMapping = make(map[string]string)
}

func TestIsShortLinkUnique(t *testing.T) {
	// Test for in-memory DB
	dbType = "in-memory"
	result := IsShortLinkUnique("shortLink")
	if !result {
		t.Error("Expected true, got false")
	}

	// Test for postgres DB
	// dbType = "postgres"
	// result = IsShortLinkUnique("shortLink")
	// if !result {
	// 	t.Error("Expected true, got false")
	// }

	// Test for invalid DB type
	dbType = "invalid"
	result = IsShortLinkUnique("shortLink")
	if !result {
		t.Error("Expected true, got false")
	}
}

func TestSaveLinkMapping(t *testing.T) {
	// Test for in-memory DB
	dbType = "in-memory"
	err := SaveLinkMapping("shortLink", "originalLink")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Test for postgres DB
	// dbType = "postgres"
	// err = SaveLinkMapping("shortLink", "originalLink")
	// if err != nil {
	// 	t.Errorf("Expected no error, got %v", err)
	// }

	// Test for invalid DB type
	dbType = "invalid"
	err = SaveLinkMapping("shortLink", "originalLink")
	if err != nil {
		t.Errorf("Expected no error, got %v'", err)
	}
}

func TestIsOriginalLinkExists(t *testing.T) {
	// Test for in-memory DB
	dbType = "in-memory"
	originalLink := "originalLink"
	_, exists, err := IsOriginalLinkExists(originalLink)
	if !exists {
		t.Error("Expected true, got false")
	}
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Test for postgres DB
	// dbType = "postgres"
	// result, exists, err = IsOriginalLinkExists(originalLink)
	// if result != expectedResult {
	// 	t.Errorf("Expected '%s', got '%s'", expectedResult, result)
	// }
	// if !exists {
	// 	t.Error("Expected true, got false")
	// }
	// if err != nil {
	// 	t.Errorf("Expected no error, got %v", err)
	// }

	// Test for invalid DB type
	dbType = "invalid"
	_, exists, err = IsOriginalLinkExists(originalLink)
	if err != nil {
		t.Errorf("Expected no error, got %v'", err)
	}
}

func TestGetOriginalLink(t *testing.T) {
	// Test for in-memory DB
	dbType = "in-memory"
	shortLink := "shortLink"
	expectedResult := "originalLink"
	result, err := GetOriginalLink(shortLink)
	if result != expectedResult {
		t.Errorf("Expected '%s', got '%s'", expectedResult, result)
	}
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Test for postgres DB
	// dbType = "postgres"
	// result, err = GetOriginalLink(shortLink)
	// if result != expectedResult {
	// 	t.Errorf("Expected '%s', got '%s'", expectedResult, result)
	// }
	// if err != nil {
	// 	t.Errorf("Expected no error, got %v", err)
	// }

	// Test for invalid DB type
	dbType = "invalid"
	result, err = GetOriginalLink(shortLink)
	if err != nil {
		t.Errorf("Expected no error, got %v'", err)
	}
}
