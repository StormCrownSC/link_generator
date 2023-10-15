package shortlink

import (
	"Service/internal/database"
	"crypto/rand"
	"database/sql"
)

const linkLength = 10
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"

func GenerateShortLink(originalLink string, db *sql.DB) (string, error) {
	var (
		shortLink string
		exist     bool = false
		err       error
	)
	if shortLink, exist, err = database.IsOriginalLinkExists(originalLink, db); exist == false && err == nil {
		// Generating a shortened link
		for i := 0; i < 5; i++ {
			randomBytes := make([]byte, linkLength)
			_, err = rand.Read(randomBytes)
			if err != nil {
				return "", err
			}

			for i := 0; i < linkLength; i++ {
				index := int(randomBytes[i]) % len(charset)
				shortLink += string(charset[index])
			}

			if database.IsShortLinkUnique(shortLink, db) {
				database.SaveLinkMapping(shortLink, originalLink, db)
				return shortLink, nil
			}

			shortLink = ""
		}
	} else if err != nil {
		return "", err
	}

	return shortLink, nil
}
