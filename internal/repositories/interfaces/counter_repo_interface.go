package interfaces

import (
	"context"
)

type CounterRepositoryInterface interface {
	GetNextAccountNumber(ctx context.Context) (string, error)
}
