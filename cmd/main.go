package main

import (
	"Service/internal/handler"
	"fmt"
)

func main() {
	fmt.Println("Start app")

	handler.StartServer()

}
