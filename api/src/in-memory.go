package main

import (
	"fmt"
)


func IsShortLinkUniqueMemory(shortLink string) bool {
    for _, link := range linkMapping {
        if link == shortLink {
            return false
        }
    }
    return true
}

func SaveLinkMappingMemory(shortLink, originalLink string) error {
    linkMapping[originalLink] = shortLink
    fmt.Printf("Сохранено соответствие: Оригинальная: %s, Сокращенная: %s\n", originalLink, shortLink)
    return nil
}

func IsOriginalLinkExistsMemory(originalLink string) (string, bool, error) {
    shortLink, exist := linkMapping[originalLink]
    return shortLink, exist, nil
}

func GetOriginalLinkMemory(shortLink string) (string, error) {
    for originalLink, link := range linkMapping {
        if link == shortLink {
            return originalLink, nil
        }
    }
    return "", fmt.Errorf("сокращенная ссылка не найдена")
}