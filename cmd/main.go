package main

import (
	"Service/internal/database"
	"Service/internal/handler"
	"fmt"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	fmt.Println("Start app")

	db, err := database.ConnectToDB()
	if err != nil {
		log.Fatal()
	}
	defer db.Close()
	r := gin.Default() // Creating a Gin Router

	r.Use(static.Serve("/", static.LocalFile("./assets", true)))
	r.POST("/shorten", func(c *gin.Context) {
		handler.CreateLink(c, db)
	})
	r.GET("/shorten", func(c *gin.Context) {
		handler.GetLink(c, db)
	})

	r.Run(":11000")
}
