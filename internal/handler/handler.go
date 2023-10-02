package handler

import (
	"Service/internal/database"
	"Service/internal/shortlink"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type RequestBody struct {
	Url string `json:"url"`
}

func HandleLink(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		shortLink := r.FormValue("shortlink")
		if shortLink == "" {
			http.Error(w, "Ошибка, ссылка не может быть пустой строкой", http.StatusBadRequest)
			return
		}

		originalLink, err := database.GetOriginalLink(shortLink)
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
		body, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, "Ошибка чтения тела запроса", http.StatusInternalServerError)
			return
		}
		// Распаковка данных в структуру RequestBody
		var requestBody RequestBody
		if err := json.Unmarshal(body, &requestBody); err != nil {
			http.Error(w, "Ошибка разбора JSON", http.StatusBadRequest)
			return
		}

		// Извлечение параметра "url" из тела запроса
		originalLink := requestBody.Url
		if originalLink == "" {
			http.Error(w, "Ошибка, ссылка не может быть пустой строкой", http.StatusBadRequest)
			return
		}

		shortLink, err := shortlink.GenerateShortLink(originalLink)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			http.Error(w, fmt.Sprintf("Ошибка генерации сокращенной ссылки: %s", err), http.StatusBadRequest)
			return
		}

		fmt.Fprintf(w, "Сокращенная ссылка: %s", shortLink)

	default:
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
	}
}
