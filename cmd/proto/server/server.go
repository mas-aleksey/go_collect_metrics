package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/tiraill/go_collect_metrics/internal/storage"
	"github.com/tiraill/go_collect_metrics/internal/utils"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	// импортируем пакет со сгенерированными protobuf-файлами
	pb "github.com/tiraill/go_collect_metrics/cmd/proto"
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
	trustedSubnet *string
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
	trustedSubnet = flag.String("t", "", "trusted subnet")
	// postgresql://ml_platform_orchestrator_admin:pwd@localhost:5467/yandex
}

// MetricsServer поддерживает все необходимые методы сервера.
type MetricsServer struct {
	// нужно встраивать тип pb.Unimplemented<TypeName>
	// для совместимости с будущими версиями
	pb.UnimplementedMetricsServer
	config utils.ServerConfig
	db     storage.Storage
}

func pbMetricToJsonMetric(m *pb.Metric) utils.JSONMetric {
	return utils.JSONMetric{
		ID:    m.Id,
		MType: m.Type,
		Delta: &m.Delta,
		Value: &m.Value,
		Hash:  &m.Hash,
	}
}

func (s *MetricsServer) SaveMetric(ctx context.Context, in *pb.SaveMetricRequest) (*pb.SaveMetricResponse, error) {
	log.Print("Handle SaveMetric")
	var response pb.SaveMetricResponse
	metric := pbMetricToJsonMetric(in.Metric)
	err := metric.ValidatesAll(s.config.HashKey)
	if err != nil {
		return nil, fmt.Errorf("ошибка валидации метрики: %v", err)
	}
	updatedMetric, err := s.db.UpdateJSONMetric(ctx, metric)
	if err != nil {
		return nil, fmt.Errorf("ошибка записи метрики в Storage: %v", err)
	}
	updatedMetric.Hash = utils.CalcHash(metric.String(), s.config.HashKey)
	response.Metric = utils.JSONMetricToPbMetric(&updatedMetric)
	return &response, nil
}

func (s *MetricsServer) SaveBatchMetrics(ctx context.Context, in *pb.SaveBatchMetricRequest) (*pb.SaveBatchMetricResponse, error) {
	log.Print("Handle SaveBatchMetrics")
	var response pb.SaveBatchMetricResponse
	metrics := make([]utils.JSONMetric, len(in.Metrics))
	for i, metric := range in.Metrics {
		metrics[i] = pbMetricToJsonMetric(metric)
		err := metrics[i].ValidatesAll(s.config.HashKey)
		if err != nil {
			return nil, fmt.Errorf("ошибка валидации метрики: %v", err)
		}
	}
	updatedMetrics, err := s.db.UpdateJSONMetrics(ctx, metrics)
	if err != nil {
		return nil, fmt.Errorf("ошибка записи метрик в Storage: %v", err)
	}
	pbMetrics := make([]*pb.Metric, 0, len(updatedMetrics))
	for i, metric := range updatedMetrics {
		metrics[i].Hash = utils.CalcHash(metric.String(), s.config.HashKey)
		pbMetric := utils.JSONMetricToPbMetric(&metric)
		pbMetrics = append(pbMetrics, pbMetric)
	}
	log.Println(pbMetrics)
	response.Metrics = pbMetrics
	return &response, nil
}

func (s *MetricsServer) GetMetric(ctx context.Context, in *pb.GetMetricRequest) (*pb.GetMetricResponse, error) {
	log.Print("Handle GetMetric")
	var response pb.GetMetricResponse
	metric := pbMetricToJsonMetric(in.Metric)
	if !metric.IsValidType() {
		return nil, fmt.Errorf("ошибка валидации метрики")
	}
	selectedMetric, err := s.db.GetJSONMetric(ctx, metric.ID, metric.MType)
	if err != nil {
		return nil, fmt.Errorf("метрика не найдена %v", err)
	}
	selectedMetric.Hash = utils.CalcHash(metric.String(), s.config.HashKey)
	response.Metric = utils.JSONMetricToPbMetric(&selectedMetric)
	return &response, nil
}

func (s *MetricsServer) GetListMetrics(ctx context.Context, _ *pb.ListMetricRequest) (*pb.ListMetricResponse, error) {
	log.Print("Handle GetListMetrics")
	var response pb.ListMetricResponse
	metrics, err := s.db.GetAllMetrics(ctx)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения метрик из Storage: %v", err)
	}
	for _, metric := range metrics {
		metric.Hash = utils.CalcHash(metric.String(), s.config.HashKey)
		response.Metrics = append(response.Metrics, utils.JSONMetricToPbMetric(&metric))
	}
	return &response, nil
}

func (s *MetricsServer) Ping(ctx context.Context, _ *pb.PingRequest) (*pb.PingResponse, error) {
	log.Print("Handle Ping")
	var response pb.PingResponse
	ok := s.db.Ping(ctx)
	if !ok {
		return nil, fmt.Errorf("ошибка соединения с БД")
	}
	return &response, nil
}

func main() {
	fmt.Println("Build version:", buildVersion)
	fmt.Println("Build date:", buildDate)
	fmt.Println("Build commit:", buildCommit)
	flag.Parse()
	serverConfig, err := utils.MakeServerConfig(*configFile, *address, *hashKey, *cryptoKey, *trustedSubnet)
	if err != nil {
		log.Fatal(err)
	}
	storageConfig, err := utils.MakeStorageConfig(*configFile, *restore, *storeInterval, *storeFile, *databaseDSN)
	if err != nil {
		log.Fatal(err)
	}
	// определяем порт для сервера
	listen, err := net.Listen("tcp", ":3200")
	if err != nil {
		log.Fatal(err)
	}
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	dbCtx, dbCancel := context.WithCancel(context.Background())
	db := storage.NewStorage(&storageConfig)
	err = db.Init(dbCtx)
	if err != nil {
		log.Printf("Error init db: %s", err)
	} else {
		log.Print("Init db success")
	}
	// создаём gRPC-сервер без зарегистрированной службы
	srv := grpc.NewServer()
	metricServer := &MetricsServer{
		config: serverConfig,
		db:     storage.NewStorage(&storageConfig),
	}

	pb.RegisterMetricsServer(srv, metricServer)

	go func() {
		if err := srv.Serve(listen); err != nil {
			log.Fatal(err)
		}
	}()
	log.Print("gRPC Server Started")

	sig := <-done
	log.Print("Server Stopped. Signal: ", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer func() {
		dbCancel()
		db.Close(ctx)
		cancel()
	}()

	srv.GracefulStop()
	log.Print("Server Exited Properly")
}
