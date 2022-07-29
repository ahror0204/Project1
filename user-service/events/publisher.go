package events

import (
	"context"
	"fmt"
	"time"

	"github.com/project1/user-service/config"
	"github.com/project1/user-service/pkg/logger"
	"github.com/project1/user-service/pkg/messagebroker"
	kafka "github.com/segmentio/kafka-go"
)

type KafkaPublisher struct {
	kafkaWriter *kafka.Writer
	log         logger.Logger
}

func NewKafkaPublisherBroker(conf config.Config, log logger.Logger, topic string) messagebroker.Publisher {
	conString := fmt.Sprintf("%s:%d", conf.KafkaHost, conf.KafkaPort)

	return &KafkaPublisher{
		kafkaWriter: &kafka.Writer{
			Addr:                   kafka.TCP(conString),
			Topic:                  topic,
			BatchTimeout:           10 * time.Millisecond,
			AllowAutoTopicCreation: true,
		},
		log: log,
	}

}

func (p *KafkaPublisher) Start() error {
	return nil
}

func (p *KafkaPublisher) Stop() error {
	err := p.kafkaWriter.Close()
	if err != nil {
		return err
	}
	return nil
}

func (p *KafkaPublisher) Publish(key, body []byte, logBody string) error {
	message := kafka.Message{
		Key:   key,
		Value: body,
	}

	err := p.kafkaWriter.WriteMessages(context.Background(), message)
	if err != nil {
		return err
	}

	p.log.Info("Message published(key/body): " + string(key) + "/" + logBody)
	return nil
}
