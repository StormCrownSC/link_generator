package main

import (
	"Service/internal/database"
	"Service/internal/handler"
	"fmt"
)

func main() {
	fmt.Println("Start app")

	db, err := database.ConnectToDB()
	if err != nil {
		return
	}
	defer db.Close()
	handler.StartServer(db)
}
