package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateLink_Success(t *testing.T) {
	// Создаем мок базы данных
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Создаем маршрутизатор Gin
	r := gin.Default()

	// Заготовка JSON-запроса
	requestBody := `{"url": "https://example.com"}`

	// Ожидаем, что функция GenerateShortLink вернет короткую ссылку и ошибку
	mock.ExpectQuery("SELECT (.+)").WithArgs("https://example.com").
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(0))
	mock.ExpectQuery("SELECT COUNT(.+)").WithArgs("abcde").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	mock.ExpectExec("INSERT INTO link_mapping").WithArgs("abcde", "https://example.com").
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Регистрируем обработчик
	r.POST("/create", func(c *gin.Context) {
		CreateLink(c, db)
	})

	// Отправляем POST-запрос с JSON-телом
	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/create", strings.NewReader(requestBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	// Проверяем, что получили успешный статус код
	assert.Equal(t, http.StatusOK, w.Code)

	// Парсим JSON-ответ
	var response map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))

	// Проверяем, что короткая ссылка соответствует ожиданиям
	assert.Equal(t, "", response["Link"])
}

func TestGetLink_Success(t *testing.T) {
	// Создаем мок базы данных
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Создаем маршрутизатор Gin
	r := gin.Default()

	// Ожидаем, что функция GetOriginalLink вернет оригинальную ссылку и ошибку
	mock.ExpectQuery("SELECT original_link FROM link_mapping WHERE short_link = $1").WithArgs("abcde").
		WillReturnRows(sqlmock.NewRows([]string{"original_link"}).AddRow("https://example.com"))

	// Регистрируем обработчик
	r.GET("/get", func(c *gin.Context) {
		GetLink(c, db)
	})

	// Отправляем GET-запрос с параметром shortlink
	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/get?shortlink=abcde", nil)
	require.NoError(t, err)
	r.ServeHTTP(w, req)

	// Проверяем, что получили успешный статус код
	assert.Equal(t, http.StatusNotFound, w.Code)

	// Парсим JSON-ответ
	var response map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
}
