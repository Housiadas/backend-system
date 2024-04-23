package kafka

import (
	"context"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Producer interface {
	Produce(ctx context.Context, msg *kafka.Message) error
	Close()
}

type ProducerConfig struct {
	Broker                  string
	SecurityProtocol        string
	AddressFamily           string
	LogLevel                int
	ProducerMaxMessageBytes int
}

type ProducerClient struct {
	producer *kafka.Producer
}

func NewProducer(cfg ProducerConfig) (*ProducerClient, error) {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers":        cfg.Broker,
		"security.protocol":        cfg.SecurityProtocol,
		"broker.address.family":    cfg.AddressFamily,
		"log_level":                cfg.LogLevel,
		"allow.auto.create.topics": false,
		"go.logs.channel.enable":   true,
		"message.max.bytes":        cfg.ProducerMaxMessageBytes,
	})
	if err != nil {
		return nil, err
	}

	return &ProducerClient{
		producer: producer,
	}, nil
}

func (p *ProducerClient) Produce(_ context.Context, msg *kafka.Message) error {
	deliveryChan := make(chan kafka.Event)
	defer close(deliveryChan)

	if err := p.producer.Produce(msg, deliveryChan); err != nil {
		return err
	}

	e := <-deliveryChan
	m := e.(*kafka.Message)
	if m.TopicPartition.Error != nil {
		return fmt.Errorf("error delivering message to kafka : %w", m.TopicPartition.Error)
	}

	return nil
}

func (p *ProducerClient) Close() {
	p.producer.Close()
}
