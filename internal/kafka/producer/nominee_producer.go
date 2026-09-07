package producer

import (
	"banking-system-backend/pkg/logger"
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type NomineeProducer struct {
	writer *kafka.Writer
}

func NewNomineeProducer() *NomineeProducer {
	writer := &kafka.Writer{
		Addr:     kafka.TCP("kafka:9092"),
		Topic:    "nominee-events",
		Balancer: &kafka.LeastBytes{},
	}

	return &NomineeProducer{
		writer: writer,
	}
}

func (p *NomineeProducer) Publish(ctx context.Context, event interface{}) error {
	logger.Log.Info("publishing nominee kafka event",
		zap.Any("event", event),
	)

	payload, err := json.Marshal(event)
	if err != nil {
		logger.Log.Info("nominee producer kafka marshal error", zap.Error(err))
		return err
	}

	err = p.writer.WriteMessages(
		ctx,
		kafka.Message{
			Value: payload,
		},
	)
	if err != nil {
		logger.Log.Error("failed to publish nominee event", zap.Error(err))
		return err
	}

	logger.Log.Info("nominee event published")
	return nil
}
