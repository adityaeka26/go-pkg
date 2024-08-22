package kafka

import (
	"context"
	"strings"

	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"
)

type KafkaConsumer struct {
	writer   *kafka.Writer
	sasl     bool
	brokers  []string
	username string
	password string
}

func NewKafkaConsumer(sasl bool, hosts, username, password string) (*KafkaConsumer, error) {
	config := kafka.WriterConfig{}
	config.Brokers = strings.Split(hosts, ",")
	config.Balancer = &kafka.LeastBytes{}
	if sasl {
		mechanism, err := scram.Mechanism(scram.SHA512, username, password)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		config.Dialer = &kafka.Dialer{SASLMechanism: mechanism}
	}

	writer := kafka.NewWriter(config)
	writer.AllowAutoTopicCreation = true

	return &KafkaConsumer{
		writer:   writer,
		sasl:     sasl,
		brokers:  strings.Split(hosts, ","),
		username: username,
		password: password,
	}, nil
}

func (k *KafkaConsumer) ConsumeMessage(ctx context.Context, groupId, topic string, messages chan<- *kafka.Message, errChan chan<- error) (*kafka.Reader, error) {
	config := kafka.ReaderConfig{}
	config.Brokers = k.brokers
	config.GroupID = groupId
	config.Topic = topic

	if k.sasl {
		mechanism, err := scram.Mechanism(scram.SHA512, k.username, k.password)
		if err != nil {
			errChan <- err
			return nil, err
		}

		config.Dialer = &kafka.Dialer{SASLMechanism: mechanism}
	}

	reader := kafka.NewReader(config)

	go func() {
		defer close(messages)
		defer close(errChan)
		for {
			m, err := reader.FetchMessage(ctx)
			if err != nil {
				errChan <- err
				continue
			}
			messages <- &m
		}
	}()

	return reader, nil
}

func (k *KafkaConsumer) Close(ctx context.Context) error {
	return k.writer.Close()
}
