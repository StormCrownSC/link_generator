package database

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"os"
)

var (
	dataSourceName = "host=link.postgres port=5432 user=admin password=admin dbname=LinkDB sslmode=disable"
)

func connectToDB() (*sql.DB, error) {
	// Инициализация базы данных PostgresSQL
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func connectDBPoll() (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(context.Background(), os.Getenv(dataSourceName))
	if err != nil {
		return nil, err
	}
	return pool, nil

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
	db, err := connectDBPoll()
	if err != nil {
		fmt.Println("Failed to connect to the database:", err)
		return err
	}
	defer db.Close()

	query := "INSERT INTO link_mapping (short_link, original_link) VALUES ($1, $2)"
	pgxErr := db.QueryRow(context.Background(), query, shortLink, originalLink)
	if pgxErr != nil {
		fmt.Println("Failed to execute query:", err)
		return err
	}

	fmt.Printf("Saved link mapping: Original: %s, Short: %s\n", originalLink, shortLink)
	return nil
}

func IsOriginalLinkExists(originalLink string) (string, bool, error) {
	// Устанавливаем соединение с базой данных
	db, err := connectDBPoll()
	if err != nil {
		fmt.Println("Ошибка соединения с базой данных в isOriginalLinkExists:", err)
		return "", false, err
	}
	defer db.Close()

	// Проверяем наличие оригинальной ссылки в базе данных
	query := "SELECT short_link FROM link_mapping WHERE original_link = $1"
	var shortLink string
	err = db.QueryRow(context.Background(), query, originalLink).Scan(&shortLink)
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
