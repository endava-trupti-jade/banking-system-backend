package consumer

import (
	"banking-system-backend/internal/kafka/events"
	"banking-system-backend/pkg/logger"
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type NomineeConsumer struct {
	reader *kafka.Reader
}

func NewNomineeConsumer() *NomineeConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"kafka:9092"},
		Topic:   "nominee-events",
		GroupID: "notification-group",
	})

	return &NomineeConsumer{
		reader: reader,
	}
}

func (c *NomineeConsumer) Start(ctx context.Context) {
	logger.Log.Info("nominee consumer started")

	for {
		logger.Log.Info("waiting for nominee kafka message")

		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			logger.Log.Error("nominee consumer kafka read error", zap.Error(err))
			continue
		}

		var event events.NomineeCreatedEvent
		err = json.Unmarshal(msg.Value, &event)
		if err != nil {
			logger.Log.Error("nominee consumer kafka unmarshal error", zap.Error(err))
			continue
		}

		logger.Log.Info("nominee created event received",
			zap.String("event", event.EventID),
			zap.String("nominee_id", event.NomineeID))
	}
}
