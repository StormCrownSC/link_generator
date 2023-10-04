package handler

import (
	"Service/internal/database"
	"Service/internal/shortlink"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type RequestBody struct {
	Url string `json:"url"`
}

func CreateLink(c *gin.Context) {
	var request RequestBody
	if err := c.BindJSON(&request); err != nil {
		fmt.Println(err, request)
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid request format"})
		return
	}

	// Извлечение параметра "url" из тела запроса
	originalLink := request.Url
	if originalLink == "" {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid url format"})
		return
	}

	shortLink, err := shortlink.GenerateShortLink(originalLink)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Generate shorten url"})
		return
	}

	// Return a successful response
	c.JSON(http.StatusOK, gin.H{
		"Shorten": shortLink,
	})

}

func GetLink(c *gin.Context) {
	// Extracting the request parameters
	shortLink := c.Query("shortlink")
	if shortLink == "" {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid url format"})
		return
	}

	originalLink, err := database.GetOriginalLink(shortLink)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"Error": "Url not found"})
		} else {
			c.JSON(http.StatusNotFound, gin.H{"Error": "with get original link"})
		}
		return
	}

	// Return a successful response
	c.JSON(http.StatusOK, gin.H{
		"Link": originalLink,
	})
}
