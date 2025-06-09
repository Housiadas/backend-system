package commands

import (
	"context"

	confkafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"

	"github.com/Housiadas/backend-system/internal/core/service/userbus"
	"github.com/Housiadas/backend-system/pkg/kafka"
)

func (cmd *Command) UserEvents() error {
	ctx := context.Background()
	cmd.Log.Info(ctx, "userevents", "status", "initializing users events worker")

	consumer, err := kafka.NewConsumer(cmd.Log, kafka.ConsumerConfig{
		Brokers:          cmd.Kafka.Brokers,
		GroupId:          "foo",
		AddressFamily:    cmd.Kafka.AddressFamily,
		SecurityProtocol: cmd.Kafka.SecurityProtocol,
		SessionTimeout:   cmd.Kafka.SessionTimeout,
	})
	if err != nil {
		cmd.Log.Error(ctx, "userevents", "ERROR", "failed to initialize consumer")
		return err
	}

	err = consumer.Subscribe(userbus.UserUpdatedEvent)
	if err != nil {
		cmd.Log.Error(ctx, "userevents", "ERROR", "failed to subscribe to topic")
		return err
	}

	err = consumer.Consume(ctx, func(msg *confkafka.Message) error {
		println("test", msg)
		return nil
	})
	if err != nil {
		cmd.Log.Error(ctx, "userevents", "ERROR", "failed to consume user events")
		return err
	}

	return nil
}
