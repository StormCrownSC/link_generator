package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var (
	dataSourceName = "host=link.postgres port=5432 user=admin password=admin dbname=LinkDB sslmode=disable"
)

func ConnectToDB() (*sql.DB, error) {
	// Инициализация базы данных PostgresSQL
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func IsShortLinkUnique(shortLink string, db *sql.DB) bool {
	query := "SELECT COUNT(*) FROM link_mapping WHERE short_link = $1"
	var count int
	err := db.QueryRow(query, shortLink).Scan(&count)
	if err != nil {
		fmt.Println("Failed to execute query:", err)
		return false
	}

	return count == 0
}

func SaveLinkMapping(shortLink, originalLink string, db *sql.DB) error {
	query := "INSERT INTO link_mapping (short_link, original_link) VALUES ($1, $2)"
	pgxErr := db.QueryRow(query, shortLink, originalLink)
	if pgxErr != nil {
		fmt.Println("Failed to execute query:", pgxErr)
		return pgxErr.Err()
	}

	fmt.Printf("Saved link mapping: Original: %s, Short: %s\n", originalLink, shortLink)
	return nil
}

func IsOriginalLinkExists(originalLink string, db *sql.DB) (string, bool, error) {
	// Проверяем наличие оригинальной ссылки в базе данных
	query := "SELECT short_link FROM link_mapping WHERE original_link = $1"
	var shortLink string
	err := db.QueryRow(query, originalLink).Scan(&shortLink)
	if err != nil {
		return "", false, nil
	}

	return shortLink, true, nil
}

func GetOriginalLink(shortLink string, db *sql.DB) (string, error) {
	// Получаем оригинальную ссылку из базы данных
	query := "SELECT original_link FROM link_mapping WHERE short_link = $1"
	var originalLink string
	err := db.QueryRow(query, shortLink).Scan(&originalLink)
	if err != nil {
		return "", err
	}

	return originalLink, nil
}

func GetAllLink(db *sql.DB) ([]string, error) {
	// Получаем оригинальную ссылку из базы данных
	query := "SELECT original_link FROM link_mapping"
	var originalLink []string
	err := db.QueryRow(query).Scan(&originalLink)
	if err != nil {
		return nil, err
	}

	return originalLink, nil
}

func DeleteOriginalLink(Link string, db *sql.DB) {
	// Получаем оригинальную ссылку из базы данных
	query := "DELETE FROM link_mapping WHERE short_link = $1"
	err := db.QueryRow(query).Scan(&Link)
	if err != nil {
		fmt.Println(err)
	}

}
