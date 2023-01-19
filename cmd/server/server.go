package main

import (
	"context"
	"flag"
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
	hashKey       *string
)

func init() {
	address = flag.String("a", "127.0.0.1:8080", "server address")
	restore = flag.Bool("r", true, "restore flag")
	storeInterval = flag.Duration("i", 30*time.Second, "store interval")
	storeFile = flag.String("f", "/tmp/devops-metrics-db.json", "store file")
	hashKey = flag.String("k", "", "hash key")
}

func saveStorage(storage *storage.MemStorage) {
	if storage.Config.StoreInterval == 0 {
		return
	}
	ticker := time.NewTicker(storage.Config.StoreInterval)
	defer ticker.Stop()

	for range ticker.C {
		storage.SaveToFileWithLog()
	}
}

func main() {
	flag.Parse()
	serverConfig := utils.MakeServerConfig(*address, *hashKey)
	storageConfig := utils.MakeMemStorageConfig(*restore, *storeInterval, *storeFile)

	memStorage := storage.NewMemStorage(storageConfig)
	router := handlers.GetRouter(memStorage, serverConfig)
	srv := &http.Server{
		Addr:    serverConfig.Address,
		Handler: router,
	}

	err := memStorage.LoadFromFile()
	if err != nil {
		log.Print("memStorage was not load from file: ", err)
	} else {
		log.Print("Load memStorage from file")
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	go saveStorage(memStorage)
	log.Print("Server Started")

	s := <-done
	log.Print("Server Stopped. Signal: ", s)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer func() {
		memStorage.SaveToFileWithLog()
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Print("Server Exited Properly")
}
