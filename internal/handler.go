package main

import (
	"fmt"
	"net/http"
	"database/sql"
)

func HandleLink(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		shortLink := r.FormValue("shortlink")

		originalLink, err := GetOriginalLink(shortLink)
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

		shortLink, err := GenerateShortLink(originalLink)
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