package main

import (
	"errors"

	_ "github.com/lib/pq"
)


func IsShortLinkUnique(shortLink string) bool {
	switch dbType {
	case "in-memory":
		return IsShortLinkUniqueMemory(shortLink)
	case "postgres":
		return IsShortLinkUniquePostgres(shortLink)
	default:
		return IsShortLinkUniqueMemory(shortLink)
	}
	return false
}

func SaveLinkMapping(shortLink, originalLink string) error {
	switch dbType {
	case "in-memory":
		return SaveLinkMappingMemory(shortLink, originalLink)
	case "postgres":
		return SaveLinkMappingPostgres(shortLink, originalLink)
	default:
		return SaveLinkMappingMemory(shortLink, originalLink)
	}
	return errors.New("Не выбран тип бд")
}

func IsOriginalLinkExists(originalLink string) (string, bool, error) {
	switch dbType {
	case "in-memory":
		return IsOriginalLinkExistsMemory(originalLink)
	case "postgres":
		return IsOriginalLinkExistsPostgres(originalLink)
	default:
		return IsOriginalLinkExistsMemory(originalLink)
	}
	return "", true, errors.New("Не выбран тип бд")
}

func GetOriginalLink(shortLink string) (string, error) {
	switch dbType {
	case "in-memory":
		return GetOriginalLinkMemory(shortLink)
	case "postgres":
		return GetOriginalLinkPostgres(shortLink)
	default:
		return GetOriginalLinkMemory(shortLink)
	}
	return "", errors.New("Не выбран тип бд")
}