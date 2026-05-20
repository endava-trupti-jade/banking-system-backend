package handlers

import (
	"banking-system-backend/internal/dto"
	"banking-system-backend/internal/requestctx"
	serviceInterfaces "banking-system-backend/internal/services/interfaces"
	"banking-system-backend/pkg/utils"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

type AccountBeneficiaryHandler struct {
	// accountBeneficiaryService *services.AccountBeneficiaryService
	accountBeneficiaryService serviceInterfaces.AccountBeneficiaryServiceInterface
}

func NewAccountBeneficiaryHandler(s serviceInterfaces.AccountBeneficiaryServiceInterface) *AccountBeneficiaryHandler {
	return &AccountBeneficiaryHandler{accountBeneficiaryService: s}
}

func (h *AccountBeneficiaryHandler) AddAccountBeneficiary(c *gin.Context) {
	ctx := c.Request.Context()
	log := requestctx.GetLogger(ctx)
	log.Info("add account beneficiary request received")

	accountNumber := c.Param("accountNumber")

	var req dto.AccountBeneficiaryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn("failed to bind add account beneficiary request",
			zap.Error(err),
		)
		utils.Error400(c, err)
		return
	}

	authCtx := requestctx.MustGetAuth(c)

	requestID, err := h.accountBeneficiaryService.AddAccountBeneficiary(
		c.Request.Context(),
		authCtx.UserID,
		authCtx.Role,
		accountNumber,
		req,
	)
	if err != nil {
		log.Warn("failed to add account beneficiary",
			zap.Error(err),
		)
		utils.Error400(c, err)
		return
	}

	log.Info("account beneficiary added successfully", zap.String("request_id", requestID.Hex()))
	utils.Success(c, http.StatusOK, gin.H{"request_id": requestID})
}

func (h *AccountBeneficiaryHandler) ListAccountBeneficiariesByAccountNumber(c *gin.Context) {
	ctx := c.Request.Context()
	log := requestctx.GetLogger(ctx)
	log.Info("list account beneficiaries by account number request received")

	accountNumber := c.Query("account_number")
	if accountNumber == "" {
		log.Warn("account number required")
		utils.Error400(c, errors.New("account_number is required"))
		return
	}

	result, err := h.accountBeneficiaryService.ListAccountBeneficiariesByAccountNumber(ctx, accountNumber)
	if err != nil {
		log.Warn("failed to list account beneficiaries",
			zap.Error(err),
		)
		utils.Error400(c, err)
		return
	}

	utils.Success(c, http.StatusOK, result)
}

func (h *AccountBeneficiaryHandler) UpdateAccountBeneficiary(c *gin.Context) {
	ctx := c.Request.Context()
	log := requestctx.GetLogger(ctx)
	log.Info("update account beneficiary request received")

	accountNumber := c.Param("accountNumber")
	mappingID, _ := primitive.ObjectIDFromHex(c.Param("mappingID"))

	var req dto.AccountBeneficiaryRequest
	_ = c.ShouldBindJSON(&req)

	authCtx := requestctx.MustGetAuth(c)

	requestID, err := h.accountBeneficiaryService.UpdateAccountBeneficiary(
		c.Request.Context(),
		authCtx.UserID,
		authCtx.Role,
		accountNumber,
		mappingID,
		req,
	)
	if err != nil {
		log.Warn("failed to update account beneficiary",
			zap.Error(err),
		)
		utils.Error400(c, err)
		return
	}

	log.Info("account beneficiary updated successfully", zap.String("mapping_id", mappingID.Hex()))
	utils.Success(c, http.StatusOK, gin.H{"request_id": requestID})
}

func (h *AccountBeneficiaryHandler) SoftDeleteAccountBeneficiary(c *gin.Context) {
	ctx := c.Request.Context()
	log := requestctx.GetLogger(ctx)
	log.Info("soft delete account beneficiary request received")

	mappingIDParam := c.Param("id")
	mappingID, err := primitive.ObjectIDFromHex(mappingIDParam)
	if err != nil {
		log.Warn("invalid mapping id",
			zap.String("mapping_id", mappingIDParam),
			zap.Error(err),
		)
		utils.Error400(c, errors.New("invalid mapping id"))
		return
	}

	accountNumber := c.Query("account_number")
	if accountNumber == "" {
		log.Warn("account number required")
		utils.Error400(c, errors.New("account_number is required"))
		return
	}

	authCtx := requestctx.MustGetAuth(c)

	requestID, err := h.accountBeneficiaryService.SoftDeleteAccountBeneficiary(
		ctx,
		authCtx.UserID,
		authCtx.Role,
		accountNumber,
		mappingID,
	)
	if err != nil {
		log.Warn("failed to soft delete account beneficiary",
			zap.Error(err),
		)
		utils.Error400(c, err)
		return
	}

	log.Info("account beneficiary soft deleted successfully", zap.String("mapping_id", mappingID.Hex()))
	utils.Success(c, http.StatusOK, gin.H{"request_id": requestID})
}
