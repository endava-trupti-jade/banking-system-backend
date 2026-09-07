package producer

import (
	"banking-system-backend/pkg/logger"
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type BeneficiaryProducer struct {
	writer *kafka.Writer
}

func NewBeneficiaryProducer() *BeneficiaryProducer {
	writer := &kafka.Writer{
		Addr:     kafka.TCP("kafka:9092"),
		Topic:    "beneficiary-events",
		Balancer: &kafka.LeastBytes{},
	}

	return &BeneficiaryProducer{
		writer: writer,
	}
}

func (p *BeneficiaryProducer) Publish(ctx context.Context, event interface{}) error {
	logger.Log.Info("publishing beneficiary kafka event",
		zap.Any("event", event),
	)
	payload, err := json.Marshal(event)
	if err != nil {
		logger.Log.Error("beneficiary producer kafka marshal error", zap.Error(err))
		return err
	}

	err = p.writer.WriteMessages(
		ctx,
		kafka.Message{
			Value: payload,
		},
	)

	if err != nil {
		logger.Log.Error("failed to publish beneficiary event", zap.Error(err))
		return err
	}

	logger.Log.Info("beneficiary event published")

	return nil
}
