package main

import (
	"github.com/tiraill/go_collect_metrics/internal/handlers"
	"github.com/tiraill/go_collect_metrics/internal/storage"
	"log"
	"net/http"
)

func main() {
	memStorage := storage.NewMemStorage()
	http.HandleFunc("/update/", handlers.SaveMetricHandler(memStorage))
	server := &http.Server{
		Addr: "127.0.0.1:8080",
	}
	log.Fatal(server.ListenAndServe())
}
