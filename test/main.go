package main

import (
	"context"
	"fmt"
	"log"
	"knocker/proto"
	"google.golang.org/grpc"
)

func main() {
	serverAddr := "localhost:11000"

	// Устанавливаем соединение с сервером gRPC
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	// Создаем клиент LinkService
	client := proto.NewLinkServiceClient(conn)

	// Отправляем запрос на сокращение ссылки
	shortenReq := &proto.LinkRequest{
		OriginalLink: "https://example.com",
	}
	shortenRes, err := client.ShortenLink(context.Background(), shortenReq)
	if err != nil {
		log.Fatalf("Failed to shorten link: %v", err)
	}
	fmt.Printf("Shortened URL: %s\n", shortenRes.ShortLink)

	// Отправляем запрос на получение оригинальной ссылки
	expandReq := &proto.LinkRequest{
		OriginalLink: shortenRes.ShortLink,
	}
	expandRes, err := client.ExpandLink(context.Background(), expandReq)
	if err != nil {
		log.Fatalf("Failed to original link: %v", err)
	}
	fmt.Printf("Expanded URL: %s\n", expandRes.ShortLink)
}
