package main

import (
	"Service/internal/handler"
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Start app")
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("assets"))))
	http.HandleFunc("/shorten", handler.HandleLink)
	http.ListenAndServe(":11000", nil)
}
