package main

import (
	"os"
	"net/http"
	"fmt"

	"github.com/joho/godotenv"
)

var (
	dbType string
	linkMapping map[string]string
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Failed to load .env file:", err)
		dbType = "in-memory"
	} else {
		dbType = os.Getenv("DB_TYPE")
	}
	linkMapping = make(map[string]string)
}

func main() {
	http.HandleFunc("/shorten", HandleLink)
	http.ListenAndServe(":11000", nil)
}
