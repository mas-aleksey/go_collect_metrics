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
	cryptoKey     *string
	databaseDSN   *string
	configFile    *string
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
	cryptoKey = flag.String("crypto-key", "", "private crypto key")
	databaseDSN = flag.String("d", "", "database connection string")
	configFile = flag.String("config", "", "config file")
	// postgresql://ml_platform_orchestrator_admin:pwd@localhost:5467/yandex
}

func main() {
	fmt.Println("Build version:", buildVersion)
	fmt.Println("Build date:", buildDate)
	fmt.Println("Build commit:", buildCommit)
	flag.Parse()
	serverConfig, err := utils.MakeServerConfig(*configFile, *address, *hashKey, *cryptoKey)
	if err != nil {
		log.Fatal(err)
	}
	storageConfig, err := utils.MakeStorageConfig(*configFile, *restore, *storeInterval, *storeFile, *databaseDSN)
	if err != nil {
		log.Fatal(err)
	}
	privateKey, err := utils.LoadPrivateKey(serverConfig.CryptoKey)
	if err != nil {
		log.Fatal("Failed to load private key: ", err)
	}
	dbCtx, dbCancel := context.WithCancel(context.Background())

	db := storage.NewStorage(&storageConfig)
	err = db.Init(dbCtx)
	if err != nil {
		log.Printf("Error init db: %s", err)
	} else {
		log.Print("Init db success")
	}

	router := handlers.GetRouter(db, serverConfig, privateKey)
	srv := &http.Server{
		Addr:    serverConfig.Address,
		Handler: router,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

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
		dbCancel()
		db.Close(ctx)
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server Shutdown Failed:%+v", err)
	} else {
		log.Print("Server Exited Properly")
	}
}
