package handlers

import (
	"banking-system-backend/constants"
	"banking-system-backend/internal/dto"
	serviceInterfaces "banking-system-backend/internal/services/interfaces"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

type BeneficiaryHandler struct {
	// beneficiaryService *services.BeneficiaryService
	beneficiaryService serviceInterfaces.BeneficiaryServiceInterface
}

func NewBeneficiaryHandler(s serviceInterfaces.BeneficiaryServiceInterface) *BeneficiaryHandler {
	return &BeneficiaryHandler{beneficiaryService: s}
}

func (h *BeneficiaryHandler) CreateBeneficiary(c *gin.Context) {
	var req dto.BeneficiaryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("userID").(primitive.ObjectID)
	if userID == primitive.NilObjectID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": constants.ErrUnauthorized.Error()})
		return
	}

	res, err := h.beneficiaryService.CreateBeneficiary(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *BeneficiaryHandler) DeleteBeneficiary(c *gin.Context) {
	idParam := c.Param("id")

	beneficiaryID, _ := primitive.ObjectIDFromHex(idParam)
	userID := c.MustGet("userID").(primitive.ObjectID)

	err := h.beneficiaryService.SoftDeleteBeneficiary(c.Request.Context(), beneficiaryID, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (h *BeneficiaryHandler) GetBeneficiary(c *gin.Context) {
	ctx := c.Request.Context()

	beneficiaryIDParam := c.Param("id")
	beneficiaryID, err := primitive.ObjectIDFromHex(beneficiaryIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid beneficiary id"})
		return
	}

	userID := c.MustGet("userID").(primitive.ObjectID)

	beneficiary, err := h.beneficiaryService.GetBeneficiaryByID(ctx, beneficiaryID, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, beneficiary)
}

func (h *BeneficiaryHandler) UpdateBeneficiary(c *gin.Context) {
	ctx := c.Request.Context()

	beneficiaryIDParam := c.Param("id")
	beneficiaryID, err := primitive.ObjectIDFromHex(beneficiaryIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid beneficiary id"})
		return
	}

	var req dto.BeneficiaryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("userID").(primitive.ObjectID)

	updated, err := h.beneficiaryService.UpdateBeneficiary(ctx, beneficiaryID, userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updated)
}

func (h *BeneficiaryHandler) ListBeneficiaries(c *gin.Context) {
	userID := c.MustGet("userID").(primitive.ObjectID)

	res, err := h.beneficiaryService.ListBeneficiaries(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}
