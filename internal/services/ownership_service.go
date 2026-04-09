package services

import (
	"banking-system-backend/constants"
	"banking-system-backend/internal/repositories/mongorepo"
	"banking-system-backend/internal/repositories/redisrepo"
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
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
	log.Println("OwnershipService SaveAccountOwner() started")
	cacheKey := accountOwnerKeyPrefix + accountNumber

	err := s.cacheRepo.Set(ctx, cacheKey, userID, time.Hour)
	if err != nil {
		log.Println("Failed to save account owner in cache:", err)
		return err
	}

	log.Println("OwnershipService SaveAccountOwner() end")
	return nil
}

func (s *OwnershipService) ResolveAccountOwner(ctx context.Context, accountNumber string) (string, error) {
	log.Println("OwnershipService ResolveAccountOwner() started")
	cacheKey := accountOwnerKeyPrefix + accountNumber

	// Get cache value
	if ownerID, err := s.cacheRepo.Get(ctx, cacheKey); err == nil && ownerID != "" {
		log.Println("Found ownerID in cache:", ownerID)
		return ownerID, nil
	}

	log.Println()

	account, err := s.accountRepo.FindByAccountNumber(ctx, accountNumber)
	if err != nil {
		log.Println("OwnershipService ResolveAccountOwner() error : ", err)
		return "", err
	}

	ownerID := account.CustomerID.Hex()

	log.Println(" accountNumber OwnershipService -> ResolveAccountOwner : ", accountNumber)
	// Cache in redis
	_ = s.cacheRepo.Set(ctx, cacheKey, ownerID, time.Hour)

	log.Println("OwnershipService ResolveAccountOwner() end")
	return ownerID, nil
}

func (s *OwnershipService) EnforceAccountOwnership(ctx context.Context, accountNumber, role string, userID primitive.ObjectID) error {
	log.Println("OwnershipService EnforceAccountOwnership() started")

	// Admin bypass
	if role == constants.RoleAdmin {
		return nil
	}

	ownerID, err := s.ResolveAccountOwner(ctx, accountNumber)
	if err != nil {
		log.Println("Ownership lookup failed:", err)
		return constants.ErrOwnershipCheckFailed
	}

	if ownerID != userID.Hex() {
		return constants.ErrOwnershipViolation
	}

	log.Println("OwnershipService EnforceAccountOwnership() end")
	return nil
}
