package consumer

import (
	"banking-system-backend/internal/kafka/events"
	"banking-system-backend/pkg/logger"
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type BeneficiaryConsumer struct {
	reader *kafka.Reader
}

func NewBeneficiaryConsumer() *BeneficiaryConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"kafka:9092"},
		Topic:   "beneficiary-events",
		GroupID: "notification-group",
	})
	return &BeneficiaryConsumer{
		reader: reader,
	}
}

func (c *BeneficiaryConsumer) Start(ctx context.Context) {
	logger.Log.Info("beneficiary consumer started")

	for {
		logger.Log.Info("waiting for beneficiary kafka message")

		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			logger.Log.Error("beneficiary consumer kafka read error", zap.Error(err))
			continue
		}

		var event events.BeneficiaryCreatedEvent
		err = json.Unmarshal(msg.Value, &event)
		if err != nil {
			logger.Log.Error("beneficiary consumer kafka unmarshal error", zap.Error(err))
			continue
		}

		logger.Log.Info("beneficiary created event received",
			zap.String("event_id", event.EventID),
			zap.String("beneficiary_id", event.BeneficiaryID),
		)
	}
}
