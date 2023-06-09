package main

import (
	"os"
	"context"
	"fmt"
	"net"
	"database/sql"
	"knocker/proto"
	"google.golang.org/grpc"
	"github.com/joho/godotenv"
)

var (
	dbType      string
	linkMapping map[string]string
)

type LinkServiceServer struct {
	proto.UnimplementedLinkServiceServer
}

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

func (s *LinkServiceServer) ShortenLink(ctx context.Context, req *proto.LinkRequest) (*proto.LinkResponse, error) {
	// Логика обработки запроса на сокращение ссылки
	// Возвращаем заполненный LinkResponse с сокращенной ссылкой
	shortenedURL, err := GenerateShortLink(req.OriginalLink)
	if err != nil {
		fmt.Println("Ошибка при генерации сокращённой ссылки")
		return nil, err // Возвращаем ошибку клиенту
	}
	return &proto.LinkResponse{
		ShortLink: shortenedURL,
	}, nil
}

func (s *LinkServiceServer) ExpandLink(ctx context.Context, req *proto.LinkRequest) (*proto.LinkResponse, error) {
	// Логика обработки запроса на получение оригинальной ссылки
	// Возвращаем заполненный LinkResponse с оригинальной ссылкой
	originalLink, err := GetOriginalLink(req.OriginalLink)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Сокращённой ссылки не существует")
			return nil, err // Возвращаем ошибку клиенту
		} else {
			fmt.Println("Ошибка получения оригинальной ссылки: %s", err)
			return nil, err // Возвращаем ошибку клиенту
		}
	}
	return &proto.LinkResponse{
		ShortLink: originalLink,
	}, nil
}

func main() {
	listenAddr := ":11000"

	// Создаем gRPC сервер
	grpcServer := grpc.NewServer()

	// Регистрируем LinkServiceServer
	linkServiceServer := &LinkServiceServer{}
	proto.RegisterLinkServiceServer(grpcServer, linkServiceServer)

	// Запускаем сервер на указанном адресе
	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		fmt.Printf("Failed to listen: %v", err)
		return
	}
	fmt.Printf("Server listening on %s\n", listenAddr)
	err = grpcServer.Serve(lis)
	if err != nil {
		fmt.Printf("Failed to serve: %v", err)
		return
	}
}
