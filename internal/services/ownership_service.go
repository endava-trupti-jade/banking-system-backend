package services

import (
	"banking-system-backend/constants"
	"banking-system-backend/internal/repositories/mongorepo"
	"banking-system-backend/internal/repositories/redisrepo"
	"banking-system-backend/internal/requestctx"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"time"
)

type OwnershipService struct {
	accountRepo *mongorepo.AccountRepository
	cacheRepo   *redisrepo.CacheRepository
}

func NewOwnershipService(accountRepo *mongorepo.AccountRepository, cacheRepo *redisrepo.CacheRepository) *OwnershipService {
	return &OwnershipService{
		accountRepo: accountRepo,
		cacheRepo:   cacheRepo,
	}
}

const accountOwnerKeyPrefix = "account_owner:"

func (s *OwnershipService) SaveAccountOwner(ctx context.Context, accountNumber, userID string) error {
	log := requestctx.GetLogger(ctx).With(
		zap.String("account_number", accountNumber),
		zap.String("user_id", userID),
	)
	log.Info("OwnershipService SaveAccountOwner() started")

	cacheKey := accountOwnerKeyPrefix + accountNumber

	err := s.cacheRepo.Set(ctx, cacheKey, userID, time.Hour)
	if err != nil {
		log.Error("Failed to save account owner in cache", zap.Error(err))
		return err
	}

	return nil
}

func (s *OwnershipService) ResolveAccountOwner(ctx context.Context, accountNumber string) (string, error) {
	log := requestctx.GetLogger(ctx).With(zap.String("account_number", accountNumber))
	log.Info("OwnershipService ResolveAccountOwner() started")

	cacheKey := accountOwnerKeyPrefix + accountNumber

	// Get cache value
	if ownerID, err := s.cacheRepo.Get(ctx, cacheKey); err == nil && ownerID != "" {
		log.Info("Found ownerID in cache", zap.String("owner_id", ownerID))
		return ownerID, nil
	}

	log.Info("OwnerID not found in cache")

	account, err := s.accountRepo.FindByAccountNumber(ctx, accountNumber)
	if err != nil {
		log.Error("OwnershipService ResolveAccountOwner() error", zap.Error(err))
		return "", err
	}

	ownerID := account.CustomerID.Hex()

	log.Info(" accountNumber OwnershipService -> ResolveAccountOwner : ", zap.String("account_number", accountNumber))
	// Cache in redis
	_ = s.cacheRepo.Set(ctx, cacheKey, ownerID, time.Hour)

	log.Info("OwnershipService ResolveAccountOwner() end")
	return ownerID, nil
}

func (s *OwnershipService) EnforceAccountOwnership(ctx context.Context, accountNumber, role string, userID primitive.ObjectID) error {
	log := requestctx.GetLogger(ctx).With(zap.String("account_number", accountNumber), zap.String("user_id", userID.Hex()))
	log.Info("OwnershipService EnforceAccountOwnership() started")

	// Admin bypass
	if role == constants.RoleAdmin {
		return nil
	}

	ownerID, err := s.ResolveAccountOwner(ctx, accountNumber)
	if err != nil {
		log.Error("Ownership lookup failed", zap.Error(err))
		return constants.ErrOwnershipCheckFailed
	}

	if ownerID != userID.Hex() {
		log.Warn("Ownership violation detected",
			zap.String("owner_id", ownerID),
		)
		return constants.ErrOwnershipViolation
	}

	log.Info("OwnershipService EnforceAccountOwnership() end")
	return nil
}
