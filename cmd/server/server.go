package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/tiraill/go_collect_metrics/internal/handlers"
	"github.com/tiraill/go_collect_metrics/internal/storage"
	"github.com/tiraill/go_collect_metrics/internal/utils"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	address       *string
	restore       *bool
	storeInterval *time.Duration
	storeFile     *string
)

func init() {
	address = flag.String("a", "127.0.0.1:8080", "server address")
	restore = flag.Bool("r", true, "restore flag")
	storeInterval = flag.Duration("i", 30*time.Second, "store interval")
	storeFile = flag.String("f", "devops-metrics-db.json", "store file")
}

func saveStorage(storage *storage.MemStorage, config utils.ServerConfig) {
	ticker := time.NewTicker(config.StoreInterval)
	defer ticker.Stop()

	for range ticker.C {
		err := storage.SaveToFile(config.StoreFile)
		if err != nil {
			log.Print("Failed save to file", err)
		} else {
			log.Print("Save storage to file")
		}
	}
}

func main() {
	flag.Parse()
	config := utils.MakeServerConfig(*address, *restore, *storeInterval, *storeFile)
	fmt.Println(config)
	memStorage := storage.NewMemStorage()
	router := handlers.GetRouter(memStorage)
	srv := &http.Server{
		Addr:    config.Address,
		Handler: router,
	}

	if config.Restore {
		err := memStorage.LoadFromFile(config.StoreFile)
		if err != nil {
			log.Print("Failed to load memStorage from file: ", err)
		} else {
			log.Print("Load memStorage from file")
		}
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	go saveStorage(memStorage, config)
	log.Print("Server Started")

	s := <-done
	log.Print("Server Stopped. Signal: ", s)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer func() {
		err := memStorage.SaveToFile(config.StoreFile)
		if err != nil {
			log.Print("Failed save to file", err)
		} else {
			log.Print("Save storage to file")
		}
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Print("Server Exited Properly")
}
