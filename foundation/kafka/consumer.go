package kafka

import (
	"context"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Consumer interface {
	Subscribe(topic string) error
	Consume(ctx context.Context, msg *kafka.Message) error
	Close()
}

type ConsumerConfig struct {
	Broker           string
	GroupId          string
	AddressFamily    string
	SecurityProtocol string
	SessionTimeout   int
}

type ConsumerClient struct {
	consumer *kafka.Consumer
}

func NewConsumer(cfg ConsumerConfig) (*ConsumerClient, error) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":        cfg.Broker,
		"group.id":                 cfg.GroupId,
		"broker.address.family":    cfg.AddressFamily,
		"session.timeout.ms":       cfg.SessionTimeout,
		"auto.offset.reset":        "earliest",
		"enable.auto.offset.store": true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}

	return &ConsumerClient{
		consumer: consumer,
	}, nil
}

func (c *ConsumerClient) Close() {
	c.consumer.Close()
}

func (c *ConsumerClient) Subscribe(topic string) error {
	err := c.consumer.Subscribe(topic, nil)
	return err
}

func (c *ConsumerClient) Consume(fn func() error) error {
	run := true
	for run {
		ev := c.consumer.Poll(100)
		switch e := ev.(type) {
		case *kafka.Message:
			// Callback, application specific
			err := fn()
			if err != nil {
				return err
			}
			run = false
		case kafka.Error:
			return fmt.Errorf("kakfa consumer error: %v\n", e)
		default:
			return fmt.Errorf("kafka consumer uknown: %v\n", e)
		}
	}
	return nil
}
