package main

import (
	"handlers"
	"log"
	"net/http"
	"storage"
)

func main() {
	memStorage := storage.NewMemStorage()
	http.HandleFunc("/update/", handlers.SaveMetricHandler(memStorage))
	server := &http.Server{
		Addr: "127.0.0.1:8080",
	}
	log.Fatal(server.ListenAndServe())
}
