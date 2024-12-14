package components

import (
	"app-consumer/internal/broker/kafka"
	"app-consumer/internal/config"
	"app-consumer/internal/services/worker"
	"app-consumer/internal/storage/redis"
	"app-consumer/internal/storage/scylla"
	"app-consumer/pkg/logger/slogpretty"
	"log/slog"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

type Components struct {
	Scylla             *scylla.Scylla
	Redis              *redis.Redis
	KafkaConsumerGroup *kafka.ConsumerGroup
	Worker             *worker.Worker
}

func InitComponents(cfg *config.Config, logger *slog.Logger) (*Components, error) {
	scylladb, err := scylla.New(cfg.Scylla.ScyllaURL)
	if err != nil {
		return nil, err
	}

	rds, err := redis.New(&cfg.Redis, logger)
	if err != nil {
		return nil, err
	}

	kafkaConsumerGroup, err := kafka.NewConsumerGroup(&cfg.Kafka, logger)
	if err != nil {
		return nil, err
	}

	workerService := worker.New(logger, kafkaConsumerGroup, scylladb, rds)

	return &Components{
		Scylla:             scylladb,
		Redis:              rds,
		KafkaConsumerGroup: kafkaConsumerGroup,
		Worker:             workerService,
	}, nil
}

func (c *Components) Shutdown() {
	c.Redis.Close()
	c.KafkaConsumerGroup.Close()
	c.Scylla.Close()
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
