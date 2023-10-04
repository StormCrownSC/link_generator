package main

import (
	"Service/internal/handler"
	"fmt"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Start app")

	r := gin.Default() // Creating a Gin Router

	r.Use(static.Serve("/", static.LocalFile("./assets", true)))
	r.POST("/shorten", handler.CreateLink)
	r.GET("/shorten", handler.GetLink)

	r.Run(":11000")
}
