package components

import (
	"app-websocket/internal/broker/kafka"
	"app-websocket/internal/config"
	"app-websocket/internal/ports"
	grpcServer "app-websocket/internal/ports/grpc"
	"app-websocket/internal/ports/ws"
	"app-websocket/internal/services/message_cache"
	"app-websocket/internal/services/message_online"
	"app-websocket/internal/services/rooms"
	"app-websocket/internal/storage/erdtree"
	"app-websocket/internal/storage/pg"
	"app-websocket/internal/storage/redis"
	"app-websocket/internal/storage/scylla"
	"app-websocket/pkg/logger/slogpretty"
	"log/slog"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

type Components struct {
	HttpServer         *ports.Server
	Postgres           *pg.Postgres
	Scylla             *scylla.Scylla
	Redis              *redis.Redis
	KafkaProducer      *kafka.KafkaProducer
	KafkaConsumerGroup *kafka.ConsumerGroup
	GrpcServer         *grpcServer.GrpcServer
}

func InitComponents(cfg *config.Config, logger *slog.Logger) (*Components, error) {
	postgres, err := pg.New(cfg.Postgres.PostgresURL)
	if err != nil {
		return nil, err
	}

	scylladb, err := scylla.New(cfg.Scylla.ScyllaURL)
	if err != nil {
		return nil, err
	}

	erdtree, err := erdtree.New()
	if err != nil {
		return nil, err
	}

	grpcServer, _ := grpcServer.New(logger, erdtree)

	rds, err := redis.New(&cfg.Redis, logger)
	if err != nil {
		return nil, err
	}

	kafkaProducer, err := kafka.NewProducer(&cfg.Kafka, logger)
	if err != nil {
		return nil, err
	}

	kafkaConsumerGroup, err := kafka.NewConsumerGroup(&cfg.Kafka, logger)
	if err != nil {
		return nil, err
	}

	hub := ws.NewHub(kafkaConsumerGroup, logger)

	roomService := rooms.New(postgres)

	chatCache := message_cache.New(&cfg.Chat, rds, scylladb)

	chatOnline := message_online.New(kafkaProducer, kafkaConsumerGroup, rds, hub, erdtree)

	httpServer, err := ports.NewServer(&cfg.Http, chatCache, chatOnline, roomService, logger, hub)
	if err != nil {
		return nil, err
	}

	return &Components{
		HttpServer:         httpServer,
		Postgres:           postgres,
		Redis:              rds,
		KafkaProducer:      kafkaProducer,
		KafkaConsumerGroup: kafkaConsumerGroup,
		GrpcServer:         grpcServer,
	}, nil
}

func (c *Components) Shutdown() {
	c.Postgres.CloseConnection()
	c.Redis.Close()
	c.KafkaProducer.Close()
	c.KafkaConsumerGroup.Close()
	c.HttpServer.Stop()
	c.GrpcServer.Stop()
}

func SetupLogger(env string) *slog.Logger {
	var logger *slog.Logger

	switch env {
	case envLocal:
		logger = slogpretty.SetupPrettySlog()
	case envDev:
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return logger
}
