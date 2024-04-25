package kafka

import (
	"context"
	"fmt"

	"github.com/Housiadas/backend-system/foundation/logger"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

const (
	MinCommitCount = 4
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
	log      *logger.Logger
}

func NewConsumer(cfg ConsumerConfig, log *logger.Logger) (*ConsumerClient, error) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":        cfg.Broker,
		"group.id":                 cfg.GroupId,
		"broker.address.family":    cfg.AddressFamily,
		"session.timeout.ms":       cfg.SessionTimeout,
		"auto.offset.reset":        "earliest",
		"enable.auto.offset.store": false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}

	return &ConsumerClient{
		consumer: consumer,
		log:      log,
	}, nil
}

func (c *ConsumerClient) Close() {
	c.consumer.Close()
}

func (c *ConsumerClient) Subscribe(topic string) error {
	err := c.consumer.Subscribe(topic, nil)
	return err
}

func (c *ConsumerClient) Consume(ctx context.Context, fn func() error) error {
	msgCount := 0
	run := true
	for run == true {
		ev := c.consumer.Poll(100)
		switch e := ev.(type) {
		case *kafka.Message:
			msgCount += 1
			if msgCount%MinCommitCount == 0 {
				go func() {
					_, err := c.consumer.Commit()
					c.log.Error(ctx, fmt.Sprintf("consumer: Commiting%v\n", err))
				}()
			}
			// Callback, application specific
			err := fn()
			if err != nil {
				c.log.Error(ctx, fmt.Sprintf("consumer: %v\n", e))
			}
			fmt.Printf("%% Message on %s:\n%s\n", e.TopicPartition, string(e.Value))
		case kafka.PartitionEOF:
			c.log.Info(ctx, fmt.Sprintf("consumer: EOF Reached %v\n", e))
		case kafka.Error:
			c.log.Error(ctx, fmt.Sprintf("consumer: %v\n", e))
			run = false
		default:
			c.log.Info(ctx, fmt.Sprintf("consumer: Ignored %v\n", e))
		}
	}
	return nil
}
