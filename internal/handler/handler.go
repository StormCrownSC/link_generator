package handler

import (
	"Service/internal/database"
	"Service/internal/shortlink"
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"server/pg"
	"sync"
	"utils"

	"github.com/gin-gonic/gin"
	"github.com/rsocket/rsocket-go"
	"github.com/rsocket/rsocket-go/payload"
	"github.com/rsocket/rsocket-go/rx/flux"
	"github.com/rsocket/rsocket-go/rx/mono"
)

type RequestBody struct {
	Url string `json:"url"`
}

func StartServer() {
	db := pg.СonnectToDB()

	// Rsocket
	wg := &sync.WaitGroup{}
	wg.Add(1)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		err := rsocket.Receive().
			OnStart(func() {
				log.Println("Server Started")
			}).
			Acceptor(func(_ context.Context, _ payload.SetupPayload, _ rsocket.CloseableRSocket) (rsocket.RSocket, error) {
				return rsocket.NewAbstractSocket(
					// Request-Response
					rsocket.RequestResponse(func(c payload.Payload) mono.Mono {
						// Получить ID - вернуть робота с таким ID
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

						shortLink, err := shortlink.GenerateShortLink(originalLink, db)
						if err != nil {
							c.JSON(http.StatusInternalServerError, gin.H{"Error": "Generate shorten url"})
							return
						}

						// Return a successful response
						c.JSON(http.StatusOK, gin.H{
							"Shorten": shortLink,
						})
					}),

					// Request-Stream
					rsocket.RequestStream(func(c payload.Payload) flux.Flux {
						allLink, err := database.GetAllLink(db)
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
							"Link": allLink,
						}),
					}),

					// Request-Channel
					rsocket.RequestChannel(func(c flux.Flux) flux.Flux {
						shortLink := c.Query("shortlink")
						if shortLink == "" {
							c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid url format"})
							return
						}

						originalLink, err := database.GetOriginalLink(shortLink, db)
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
						}),
					}),

					// Fire-and-forget
					rsocket.FireAndForget(func(с payload.Payload) {
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

						database.DeleteOriginalLink(originalLink, db)

					}),

					////////////////////////////////////////////////////////////////////////
				), nil
			}).
			Transport(rsocket.TCPServer().SetAddr(":7878").Build()).
			Serve(ctx)

		if err != nil {
			log.Fatalln(err)
		}
		wg.Done()
	}()
	utils.EnterExit()
	cancel()
	wg.Wait()
}

func GetLink(c *gin.Context, db *sql.DB) {
	// Extracting the request parameters
	shortLink := c.Query("shortlink")
	if shortLink == "" {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid url format"})
		return
	}

	originalLink, err := database.GetOriginalLink(shortLink, db)
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
