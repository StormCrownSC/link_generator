package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func connectToDB() (*sql.DB, error) {
	// Инициализация базы данных PostgreSQL
	db, err := sql.Open("postgres", "host=link.postgres port=5432 user=admin password=admin dbname=LinkDB sslmode=disable")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func IsShortLinkUnique(shortLink string) bool {
	db, err := connectToDB()
	if err != nil {
		fmt.Println("Failed to connect to the database:", err)
		return false
	}
	defer db.Close()

	query := "SELECT COUNT(*) FROM link_mapping WHERE short_link = $1"
	var count int
	err = db.QueryRow(query, shortLink).Scan(&count)
	if err != nil {
		fmt.Println("Failed to execute query:", err)
		return false
	}

	return count == 0
}

func SaveLinkMapping(shortLink, originalLink string) error {
	db, err := connectToDB()
	if err != nil {
		fmt.Println("Failed to connect to the database:", err)
		return err
	}
	defer db.Close()

	query := "INSERT INTO link_mapping (short_link, original_link) VALUES ($1, $2)"
	_, err = db.Exec(query, shortLink, originalLink)
	if err != nil {
		fmt.Println("Failed to execute query:", err)
		return err
	}

	fmt.Printf("Saved link mapping: Original: %s, Short: %s\n", originalLink, shortLink)
	return nil
}

func IsOriginalLinkExists(originalLink string) (string, bool, error) {
	// Устанавливаем соединение с базой данных
	db, err := connectToDB()
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

func GetOriginalLink(shortLink string) (string, error) {
	// Устанавливаем соединение с базой данных
	db, err := connectToDB()
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
