package database

import (
	"errors"
	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"testing"
)

func TestIsShortLinkUnique_UniqueShortLink(t *testing.T) {
	// Создаем подключение к SQL и мок базы данных
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Создаем тестовые данные и ожидаемый результат
	shortLink := "abc123"
	expectedCount := 0

	// Ожидаем, что запрос вернет ожидаемое количество записей
	mock.ExpectQuery("SELECT COUNT(.+)").WithArgs(shortLink).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

	// Вызываем тестируемую функцию
	isUnique := IsShortLinkUnique(shortLink, db)

	// Проверяем, что результат соответствует ожиданиям
	if isUnique != true {
		t.Errorf("Expected true, but got false")
	}

	// Проверяем, что все ожидаемые запросы выполнены
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Expectations were not met: %s", err)
	}
}

func TestIsShortLinkUnique_NonUniqueShortLink(t *testing.T) {
	// Создаем подключение к SQL и мок базы данных
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Создаем тестовые данные и ожидаемый результат
	shortLink := "abc123"
	expectedCount := 1

	// Ожидаем, что запрос вернет ожидаемое количество записей
	mock.ExpectQuery("SELECT COUNT(.+)").WithArgs(shortLink).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

	// Вызываем тестируемую функцию
	isUnique := IsShortLinkUnique(shortLink, db)

	// Проверяем, что результат соответствует ожиданиям
	if isUnique != false {
		t.Errorf("Expected false, but got true")
	}

	// Проверяем, что все ожидаемые запросы выполнены
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Expectations were not met: %s", err)
	}
}

func TestIsShortLinkUnique_QueryError(t *testing.T) {
	// Создаем подключение к SQL и мок базы данных
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Создаем тестовые данные
	shortLink := "abc123"

	// Ожидаем, что запрос вызовет ошибку
	mock.ExpectQuery("SELECT COUNT(.+)").WithArgs(shortLink).
		WillReturnError(errors.New("database error"))

	// Вызываем тестируемую функцию
	isUnique := IsShortLinkUnique(shortLink, db)

	// Проверяем, что функция возвращает false (ошибка при выполнении запроса)
	if isUnique != false {
		t.Errorf("Expected false, but got true")
	}

	// Проверяем, что все ожидаемые запросы выполнены
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Expectations were not met: %s", err)
	}
}

func TestIsShortLinkUnique_ScanError(t *testing.T) {
	// Создаем подключение к SQL и мок базы данных
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Создаем тестовые данные
	shortLink := "abc123"

	// Ожидаем, что запрос вернет некорректные данные
	mock.ExpectQuery("SELECT COUNT(.+)").WithArgs(shortLink).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow("invalid"))

	// Вызываем тестируемую функцию
	isUnique := IsShortLinkUnique(shortLink, db)

	// Проверяем, что функция возвращает false (ошибка при сканировании результата)
	if isUnique != false {
		t.Errorf("Expected false, but got true")
	}

	// Проверяем, что все ожидаемые запросы выполнены
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Expectations were not met: %s", err)
	}
}

func TestSaveLinkMapping_Success(t *testing.T) {
	// Создаем подключение к SQL и мок базы данных
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Создаем тестовые данные
	shortLink := "abc123"
	originalLink := "https://example.com"

	// Ожидаем, что запрос на вставку будет выполнен
	mock.ExpectQuery("INSERT INTO link_mapping").WithArgs(shortLink, originalLink).
		WillReturnRows(sqlmock.NewRows([]string{"last_insert_id"}).AddRow(1))

	// Вызываем тестируемую функцию
	err = SaveLinkMapping(shortLink, originalLink, db)

	// Проверяем, что ошибка равна nil
	if err != nil {
		t.Errorf("Expected nil error, but got %v", err)
	}

	// Проверяем, что все ожидаемые запросы выполнены
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Expectations were not met: %s", err)
	}
}

func TestSaveLinkMapping_DBError(t *testing.T) {
	// Создаем мок базы данных, который вернет ошибку
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	db.Close() // Закрываем базу данных, чтобы вызвать ошибку при подключении

	// Вызываем тестируемую функцию
	err = SaveLinkMapping("abc123", "https://example.com", db)

	// Проверяем, что функция возвращает ошибку (ошибка при подключении)
	if err == nil {
		t.Errorf("Expected error, but got nil")
	}
}

func TestSaveLinkMapping_ScanError(t *testing.T) {
	// Создаем подключение к SQL и мок базы данных
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Создаем тестовые данные
	shortLink := "abc123"
	originalLink := "https://example.com"

	// Ожидаем, что запрос вернет некорректные данные
	mock.ExpectQuery("INSERT INTO link_mapping").WithArgs(shortLink, originalLink).
		WillReturnRows(sqlmock.NewRows([]string{"last_insert_id"}).AddRow("invalid"))

	// Вызываем тестируемую функцию
	err = SaveLinkMapping(shortLink, originalLink, db)

	// Проверяем, что функция не возвращает ошибку
	if err != nil {
		t.Errorf("Expected no error, but got: %s", err)
	}

	// Проверяем, что все ожидаемые запросы выполнены
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Expectations were not met: %s", err)
	}
}

func TestSaveLinkMapping_LastInsertIDError(t *testing.T) {
	// Создаем подключение к SQL и мок базы данных
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Создаем тестовые данные
	shortLink := "abc123"
	originalLink := "https://example.com"

	// Ожидаем, что запрос на вставку будет выполнен, но не вернет last_insert_id
	mock.ExpectQuery("INSERT INTO link_mapping").WithArgs(shortLink, originalLink).
		WillReturnRows(sqlmock.NewRows([]string{"last_insert_id"}))

	// Вызываем тестируемую функцию
	err = SaveLinkMapping(shortLink, originalLink, db)

	// Проверяем, что функция не возвращает ошибку
	if err != nil {
		t.Errorf("Expected no error, but got: %s", err)
	}

	// Проверяем, что все ожидаемые запросы выполнены
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Expectations were not met: %s", err)
	}
}
