package main

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/lib/pq"
)

const linkLength = 10
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"

func generateShortLink(originalLink string) (string, error) {
	shortLink := ""
	exist := false
	err := error(nil)

	if shortLink, exist, err = isOriginalLinkExists(originalLink); exist == false && err == nil {
		// Генерация сокращенной ссылки
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

			if isShortLinkUnique(shortLink) {
				saveLinkMapping(shortLink, originalLink)
				return shortLink, nil
			}

			shortLink = ""
		}
	} else if err != nil {
		return "", err
	}

	return shortLink, nil
}

func isOriginalLinkExists(originalLink string) (string, bool, error) {
	// Устанавливаем соединение с базой данных
	db, err := sql.Open("postgres", "host=link.postgres port=5432 user=admin password=admin dbname=LinkDB sslmode=disable")
	if err != nil {
		fmt.Println("Ошибка соединения с базой данных в isOriginalLinkExists:", err)
		return "", false, err
	}
	defer db.Close()

	// Проверяем наличие оригинальной ссылки в базе данных
	query := "SELECT short_link FROM link_mapping WHERE original_link = $1"
	var shortLink string
	err = db.QueryRow(query, originalLink).Scan(&shortLink)
	if err != nil {
		return "", false, nil
	}

	return shortLink, true, nil
}

func isShortLinkUnique(shortLink string) bool {
	// Устанавливаем соединение с базой данных
	db, err := sql.Open("postgres", "host=link.postgres port=5432 user=admin password=admin dbname=LinkDB sslmode=disable")
	if err != nil {
		fmt.Println("Ошибка соединения с базой данных в isShortLinkUnique:", err)
		return false
	}
	defer db.Close()

	// Проверяем уникальность сокращенной ссылки в базе данных
	query := "SELECT COUNT(*) FROM link_mapping WHERE short_link = $1"
	var count int
	err = db.QueryRow(query, shortLink).Scan(&count)
	if err != nil {
		fmt.Println("Ошибка выполнения запроса в isShortLinkUnique:", err)
		return false
	}

	return count == 0
}

func saveLinkMapping(shortLink, originalLink string) error {
	// Устанавливаем соединение с базой данных
	db, err := sql.Open("postgres", "host=link.postgres port=5432 user=admin password=admin dbname=LinkDB sslmode=disable")
	if err != nil {
		fmt.Println("Ошибка соединения с базой данных в saveLinkMapping:", err)
		return err
	}
	defer db.Close()

	// Выполняем запись соответствия в базу данных
	query := "INSERT INTO link_mapping (short_link, original_link) VALUES ($1, $2)"
	_, err = db.Exec(query, shortLink, originalLink)
	if err != nil {
		fmt.Println("Ошибка выполнения запроса в saveLinkMapping:", err)
		return err
	}

	fmt.Printf("Сохранено соответствие: Оригинальная: %s, Сокращенная: %s\n", originalLink, shortLink)
	return nil
}

func handleLink(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		shortLink := r.FormValue("shortlink")

		originalLink, err := getOriginalLink(shortLink)
		if err != nil {
			if err == sql.ErrNoRows {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, "404 Not Found")
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "Ошибка получения оригинальной ссылки: %s", err)
			}
			return
		}

		fmt.Fprintf(w, "Оригинальная ссылка: %s", originalLink)
		// http.Redirect(w, r, originalLink, http.StatusFound)

	case http.MethodPost:
		originalLink := r.FormValue("url")

		shortLink, err := generateShortLink(originalLink)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Ошибка генерации сокращенной ссылки: %s", err)
			return
		}

		fmt.Fprintf(w, "Сокращенная ссылка: %s", shortLink)

	default:
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
	}
}

func getOriginalLink(shortLink string) (string, error) {
	// Устанавливаем соединение с базой данных
	db, err := sql.Open("postgres", "host=link.postgres port=5432 user=admin password=admin dbname=LinkDB sslmode=disable")
	if err != nil {
		fmt.Println("Ошибка соединения с базой данных в getOriginalLink:", err)
		return "", err
	}
	defer db.Close()
	// Получаем оригинальную ссылку из базы данных
	query := "SELECT original_link FROM link_mapping WHERE short_link = $1"
	var originalLink string
	err = db.QueryRow(query, shortLink).Scan(&originalLink)
	if err != nil {
		return "", err
	}

	return originalLink, nil
}

func main() {
	http.HandleFunc("/shorten", handleLink)
	http.ListenAndServe(":11000", nil)
}
