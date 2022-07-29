package events

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/project1/post-service/config"
	handler "github.com/project1/post-service/events/handler"
	"github.com/project1/post-service/pkg/logger"
	"github.com/project1/post-service/storage"
	kafka "github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	kafkaConsumer *kafka.Reader
	eventHandler  *handler.EventHandler
	log           logger.Logger
}

func NewKafkaConsumer(db *sqlx.DB, conf *config.Config, log logger.Logger, topic string) *KafkaConsumer {
	connString := fmt.Sprintf("%s:%d", conf.KafkaHost, conf.KafkaPort)
	return &KafkaConsumer{
		kafkaConsumer: kafka.NewReader(kafka.ReaderConfig{
			Brokers:        []string{connString},
			Topic:          topic,
			MinBytes:       10e3,
			MaxBytes:       10e6,
			Partition:      0,
			CommitInterval: 0,
		}),
		eventHandler: handler.NewEventHandlerFunc(*conf, storage.NewStoragePg(db), log),
		log:          log,
	}
}

func (k KafkaConsumer) Start() {
	fmt.Println(">>> Kafka consumer started")
	for {
		m, err := k.kafkaConsumer.ReadMessage(context.Background())
		if err != nil {
			k.log.Error("Error on consuming a message", logger.Error(err))
			break
		}
		err = k.eventHandler.Handler(m.Value)
		if err != nil {
			k.log.Error("failed to handle consumed topic:",
				logger.String("on toc]pic", m.Topic), logger.Error(err))
		} else {
			fmt.Println()
			k.log.Info("Successfully consumed message",
				logger.String("on topic", m.Topic),
				logger.String("message", "success"))
			fmt.Println()
		}
	}

	err := k.kafkaConsumer.Close()
	if err != nil {
		k.log.Error("error while closing kafka reader", logger.Error(err))
	}
}
