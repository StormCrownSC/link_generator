package handler

import (
	"Service/internal/database"
	"Service/internal/shortlink"
	"context"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/rsocket/rsocket-go"
	"github.com/rsocket/rsocket-go/payload"
	"github.com/rsocket/rsocket-go/rx/flux"
	"github.com/rsocket/rsocket-go/rx/mono"
	"log"
	"sync"
)

type RequestBody struct {
	Url string `json:"url"`
}

func StartServer() {
	db, err := database.ConnectToDB()
	if err != nil {
		return
	}
	defer db.Close()

	wg := &sync.WaitGroup{}
	wg.Add(2)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		defer wg.Done()
		r := gin.Default()
		r.Use(static.Serve("/", static.LocalFile("./assets", true)))
		r.Run(":12000")
	}()

	// Rsocket
	go func() {
		defer wg.Done()
		err := rsocket.Receive().OnStart(func() {
			log.Println("Server Started")
		}).Acceptor(func(_ context.Context, _ payload.SetupPayload, _ rsocket.CloseableRSocket) (rsocket.RSocket, error) {
			return rsocket.NewAbstractSocket(
				// Request-Response
				rsocket.RequestResponse(func(c payload.Payload) mono.Mono {
					// Извлечение параметра "url" из тела запроса
					originalLink := c.DataUTF8()
					if originalLink == "" {
						return mono.Error(fmt.Errorf("invalid url format"))
					}

					shortLink, err := shortlink.GenerateShortLink(originalLink, db)
					if err != nil {
						return mono.Error(err)
					}

					return mono.Just(payload.NewString(shortLink, ""))
				}),

				// Request-Stream
				rsocket.RequestStream(func(c payload.Payload) flux.Flux {
					allLink, err := database.GetAllLink(db)
					if err != nil {
						if err == sql.ErrNoRows {
							return flux.Error(fmt.Errorf("url not found"))
						} else {
							return flux.Error(fmt.Errorf("with get original link"))
						}
					}

					// Return a successful response
					return flux.Create(func(_ context.Context, s flux.Sink) {
						for _, link := range allLink {
							s.Next(payload.NewString(link, ""))
						}
						s.Complete()
					})
				}),

				// Request-Channel
				rsocket.RequestChannel(func(c flux.Flux) flux.Flux {
					shorts := make(chan string)
					originals := make(chan string)

					c.DoOnComplete(func() {
						close(shorts)
					}).DoOnNext(func(msg payload.Payload) error {
						short := msg.DataUTF8()
						if err != nil {
							log.Fatalln(err)
							return nil
						}
						shorts <- short
						return nil
					}).Subscribe(context.Background())

					go func() {
						for short := range shorts {
							original, _ := database.GetOriginalLink(short, db)
							originals <- original
						}
						close(originals)
					}()
					return flux.Create(func(_ context.Context, s flux.Sink) {
						for original := range originals {
							s.Next(payload.NewString(original, ""))
						}
						s.Complete()
					})
				}),

				// Fire-and-forget
				rsocket.FireAndForget(func(c payload.Payload) {
					originalLink := c.DataUTF8()
					if originalLink != "" {
						database.DeleteOriginalLink(originalLink, db)
					}
				}),
			), nil
		}).Transport(rsocket.TCPServer().SetAddr(":11000").Build()).Serve(ctx)

		if err != nil {
			log.Fatalln(err)
		}
	}()
	wg.Wait()
	cancel()
}
