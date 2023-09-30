package shortlink

import (
	"Service/internal/database"
	"crypto/rand"
)

const linkLength = 10
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"

func GenerateShortLink(originalLink string) (string, error) {
	shortLink := ""
	exist := false
	err := error(nil)

	if shortLink, exist, err = database.IsOriginalLinkExists(originalLink); exist == false && err == nil {
		// Generating a shortened link
		for {
			randomBytes := make([]byte, linkLength)
			_, err = rand.Read(randomBytes)
			if err != nil {
				return "", err
			}

			for i := 0; i < linkLength; i++ {
				index := int(randomBytes[i]) % len(charset)
				shortLink += string(charset[index])
			}

			if database.IsShortLinkUnique(shortLink) {
				database.SaveLinkMapping(shortLink, originalLink)
				return shortLink, nil
			}

			shortLink = ""
		}
	} else if err != nil {
		return "", err
	}

	return shortLink, nil
}
