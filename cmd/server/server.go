// Сервер сбора метрик
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tiraill/go_collect_metrics/internal/handlers"
	"github.com/tiraill/go_collect_metrics/internal/storage"
	"github.com/tiraill/go_collect_metrics/internal/utils"
)

var (
	address       *string
	restore       *bool
	storeInterval *time.Duration
	storeFile     *string
	hashKey       *string
	databaseDSN   *string
	buildVersion  = "N/A"
	buildDate     = "N/A"
	buildCommit   = "N/A"
)

func init() {
	address = flag.String("a", "127.0.0.1:8080", "server address")
	restore = flag.Bool("r", true, "restore flag")
	storeInterval = flag.Duration("i", 30*time.Second, "store interval")
	storeFile = flag.String("f", "/tmp/devops-metrics-db.json", "store file")
	hashKey = flag.String("k", "", "hash key")
	databaseDSN = flag.String("d", "", "database connection string")
	// postgresql://ml_platform_orchestrator_admin:pwd@localhost:5467/yandex
}

func main() {
	fmt.Println("Build version:", buildVersion)
	fmt.Println("Build date:", buildDate)
	fmt.Println("Build commit:", buildCommit)
	flag.Parse()
	serverConfig := utils.MakeServerConfig(*address, *hashKey)
	storageConfig, err := utils.MakeStorageConfig(*restore, *storeInterval, *storeFile, *databaseDSN)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	db := storage.NewStorage(&storageConfig)
	err = db.Init(ctx)
	if err != nil {
		log.Printf("Error init db: %s", err)
	} else {
		log.Print("Init db success")
	}

	router := handlers.GetRouter(db, serverConfig)
	srv := &http.Server{
		Addr:    serverConfig.Address,
		Handler: router,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Print("Server Started")

	s := <-done
	log.Print("Server Stopped. Signal: ", s)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer func() {
		db.Close(ctx)
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server Shutdown Failed:%+v", err)
	} else {
		log.Print("Server Exited Properly")
	}
}
