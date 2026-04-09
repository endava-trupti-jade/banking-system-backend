package handlers

import (
	"banking-system-backend/internal/dto"
	serviceInterfaces "banking-system-backend/internal/services/interfaces"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AccountBeneficiaryHandler struct {
	// accountBeneficiaryService *services.AccountBeneficiaryService
	accountBeneficiaryService serviceInterfaces.AccountBeneficiaryServiceInterface
}

func NewAccountBeneficiaryHandler(s serviceInterfaces.AccountBeneficiaryServiceInterface) *AccountBeneficiaryHandler {
	return &AccountBeneficiaryHandler{accountBeneficiaryService: s}
}

func (h *AccountBeneficiaryHandler) AddAccountBeneficiary(c *gin.Context) {
	accountNumber := c.Param("accountNumber")

	var req dto.AccountBeneficiaryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("userID").(primitive.ObjectID)
	role := c.MustGet("role").(string)

	requestID, err := h.accountBeneficiaryService.AddAccountBeneficiary(
		c.Request.Context(),
		userID,
		role,
		accountNumber,
		req,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"request_id": requestID})
}

func (h *AccountBeneficiaryHandler) ListAccountBeneficiariesByAccountNumber(c *gin.Context) {
	ctx := c.Request.Context()

	accountNumber := c.Query("account_number")
	if accountNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "account_number is required"})
		return
	}

	result, err := h.accountBeneficiaryService.ListAccountBeneficiariesByAccountNumber(ctx, accountNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *AccountBeneficiaryHandler) UpdateAccountBeneficiary(c *gin.Context) {
	accountNumber := c.Param("accountNumber")
	mappingID, _ := primitive.ObjectIDFromHex(c.Param("mappingID"))

	var req dto.AccountBeneficiaryRequest
	_ = c.ShouldBindJSON(&req)

	userID := c.MustGet("userID").(primitive.ObjectID)
	role := c.MustGet("role").(string)

	requestID, err := h.accountBeneficiaryService.UpdateAccountBeneficiary(
		c.Request.Context(),
		userID,
		role,
		accountNumber,
		mappingID,
		req,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"request_id": requestID})
}

func (h *AccountBeneficiaryHandler) SoftDeleteAccountBeneficiary(c *gin.Context) {
	ctx := c.Request.Context()

	mappingIDParam := c.Param("id")
	mappingID, err := primitive.ObjectIDFromHex(mappingIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid mapping id"})
		return
	}

	accountNumber := c.Query("account_number")
	if accountNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "account_number is required"})
		return
	}

	userID := c.MustGet("userID").(primitive.ObjectID)
	role := c.MustGet("role").(string)

	requestID, err := h.accountBeneficiaryService.SoftDeleteAccountBeneficiary(
		ctx,
		userID,
		role,
		accountNumber,
		mappingID,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"request_id": requestID})
}
