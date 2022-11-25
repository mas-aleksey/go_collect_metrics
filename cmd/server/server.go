package main

import (
	"github.com/tiraill/go_collect_metrics/internal/handlers"
	"github.com/tiraill/go_collect_metrics/internal/storage"
	"log"
	"net/http"
)

func main() {
	memStorage := storage.NewMemStorage()
	router := handlers.GetRouter(memStorage)
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", router))
}
