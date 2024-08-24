package kafka

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"
)

type KafkaProducer struct {
	writer   *kafka.Writer
	sasl     bool
	brokers  []string
	username string
	password string
}

func NewKafkaProducer(sasl bool, hosts, username, password string) (*KafkaProducer, error) {
	// producer
	config := kafka.WriterConfig{}
	config.Brokers = strings.Split(hosts, ",")
	config.Balancer = &kafka.LeastBytes{}
	if sasl {
		mechanism, err := scram.Mechanism(scram.SHA512, username, password)
		if err != nil {
			return nil, err
		}

		config.Dialer = &kafka.Dialer{SASLMechanism: mechanism}
	}

	writer := kafka.NewWriter(config)
	writer.AllowAutoTopicCreation = true

	return &KafkaProducer{
		writer:   writer,
		sasl:     sasl,
		brokers:  strings.Split(hosts, ","),
		username: username,
		password: password,
	}, nil
}

func (k *KafkaProducer) SendMessage(ctx context.Context, topic string, value []byte) error {
	fmt.Printf("send message to topic %s: %s\n", topic, string(value))
	return k.writer.WriteMessages(ctx, kafka.Message{
		Topic: topic,
		Value: value,
	})
}

func (k *KafkaProducer) SendMessageWithAutoTopicCreation(ctx context.Context, topic string, value []byte) error {
	var err error
	const retries = 3
	for i := 0; i < retries; i++ {
		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		err = k.writer.WriteMessages(ctx, kafka.Message{
			Topic: topic,
			Value: value,
		})
		if errors.Is(err, kafka.LeaderNotAvailable) || errors.Is(err, context.DeadlineExceeded) {
			time.Sleep(time.Millisecond * 250)
			continue
		}

		if err != nil {
			return err
		}
		break
	}

	fmt.Printf("send message (with topic creation) to topic %s: %s\n", topic, string(value))
	return nil
}

func (k *KafkaProducer) Close(ctx context.Context) error {
	return k.writer.Close()
}
