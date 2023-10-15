package shortlink

import (
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"testing"
)

func TestGenerateShortLink_Success(t *testing.T) {
	// Создаем подключение к SQL и мок базы данных
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Создаем тестовые данные
	originalLink := "https://example.com"
	shortLink := "abcde"

	// Ожидаем, что функция IsOriginalLinkExists вернет shortLink и true
	mock.ExpectQuery("SELECT short_link FROM link_mapping WHERE original_link = $1").WithArgs(originalLink).
		WillReturnRows(sqlmock.NewRows([]string{"short_link"}).AddRow(shortLink))
	// Ожидаем, что функция IsShortLinkUnique вернет true
	mock.ExpectQuery("SELECT COUNT(*) FROM link_mapping WHERE short_link = $1").WithArgs(shortLink).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	// Ожидаем, что функция SaveLinkMapping будет вызвана
	mock.ExpectExec("INSERT INTO link_mapping").WithArgs(shortLink, originalLink).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Вызываем тестируемую функцию
	resultShortLink, err := GenerateShortLink(originalLink, db)
	fmt.Println(resultShortLink, err)

	// Проверяем, что короткая ссылка соответствует ожиданиям
	if resultShortLink != "" {
		t.Errorf("Expected short link %s, but got %s", shortLink, resultShortLink)
	}

	// Проверяем, что ошибка равна nil
	if err != nil {
		t.Errorf("Expected nil error, but got %v", err)
	}
}

func TestGenerateShortLink_OriginalLinkExists(t *testing.T) {
	// Создаем подключение к SQL и мок базы данных
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Создаем тестовые данные
	originalLink := "https://example.com"
	shortLink := ""

	// Ожидаем, что функция IsOriginalLinkExists вернет true и ошибку
	mock.ExpectQuery("SELECT (.+)").WithArgs(originalLink).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(shortLink))

	// Вызываем тестируемую функцию
	resultShortLink, err := GenerateShortLink(originalLink, db)

	// Проверяем, что короткая ссылка равна пустой строке
	if resultShortLink != "" {
		t.Errorf("Expected empty short link, but got %s", resultShortLink)
	}

	// Проверяем, что ошибки нет
	if err != nil {
		t.Errorf("Expected no error, but got: %s", err)
	}

	// Проверяем, что все ожидаемые запросы выполнены
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Expectations were not met: %s", err)
	}
}

func TestGenerateShortLink_RandomGenerationError(t *testing.T) {
	// Создаем подключение к SQL и мок базы данных
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Создаем тестовые данные
	originalLink := "https://example.com"
	shortLink := ""

	// Ожидаем, что функция IsOriginalLinkExists вернет false и ошибку
	mock.ExpectQuery("SELECT (.+)").WithArgs(originalLink).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(shortLink))

	// Вызываем тестируемую функцию, преднамеренно вызывая ошибку при генерации случайной короткой ссылки
	resultShortLink, err := GenerateShortLink(originalLink, db)

	// Проверяем, что короткая ссылка равна пустой строке
	if resultShortLink != "" {
		t.Errorf("Expected empty short link, but got %s", resultShortLink)
	}

	// Проверяем, что ошибка не равна nil
	if err != nil {
		t.Errorf("Expected no error, but got: %s", err)
	}

	// Проверяем, что все ожидаемые запросы выполнены
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Expectations were not met: %s", err)
	}
}

func TestGenerateShortLink_SaveLinkMappingError(t *testing.T) {
	// Создаем подключение к SQL и мок базы данных
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Создаем тестовые данные
	originalLink := "https://example.com"
	shortLink := "abcde"

	// Ожидаем, что функция IsOriginalLinkExists вернет false и ошибку
	mock.ExpectQuery("SELECT (.+)").WithArgs(originalLink).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(shortLink))
	// Ожидаем, что функция IsShortLinkUnique вернет true
	mock.ExpectQuery("SELECT COUNT(.+)").WithArgs(shortLink).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	// Ожидаем, что функция SaveLinkMapping вызовет ошибку
	mock.ExpectQuery("INSERT INTO link_mapping").WithArgs(shortLink, originalLink).
		WillReturnError(errors.New("database error"))

	// Вызываем тестируемую функцию
	resultShortLink, err := GenerateShortLink(originalLink, db)

	// Проверяем, что короткая ссылка равна пустой строке
	if resultShortLink != shortLink {
		t.Errorf("Expected empty short link, but got %s", resultShortLink)
	}

	// Проверяем, что ошибка не равна nil
	if err != nil {
		t.Errorf("Expected no error, but got: %s", err)
	}
}
