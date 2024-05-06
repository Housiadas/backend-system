package commands

import (
	"context"

	confkafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"

	"github.com/Housiadas/backend-system/business/config"
	"github.com/Housiadas/backend-system/business/data/sqldb"
	"github.com/Housiadas/backend-system/business/domain/userbus"
	"github.com/Housiadas/backend-system/foundation/kafka"
	"github.com/Housiadas/backend-system/foundation/logger"
)

func UserEvents(log *logger.Logger, dbConfig sqldb.Config, cfg config.Kafka) error {
	ctx := context.Background()
	log.Info(ctx, "userevents", "status", "initializing users events worker")

	consumer, err := kafka.NewConsumer(log, kafka.ConsumerConfig{
		Brokers:          cfg.Brokers,
		GroupId:          "foo",
		AddressFamily:    cfg.AddressFamily,
		SecurityProtocol: cfg.SecurityProtocol,
		SessionTimeout:   cfg.SessionTimeout,
	})
	if err != nil {
		log.Error(ctx, "userevents", "ERROR", "failed to initialize consumer")
		return err
	}

	err = consumer.Subscribe(userbus.UserUpdatedEvent)
	if err != nil {
		log.Error(ctx, "userevents", "ERROR", "failed to subscribe to topic")
		return err
	}

	err = consumer.Consume(ctx, func(msg *confkafka.Message) error {
		println("test", msg)
		return nil
	})
	if err != nil {
		log.Error(ctx, "userevents", "ERROR", "failed to consume user events")
		return err
	}

	return nil
}
