package handlers

import (
	"banking-system-backend/constants"
	"banking-system-backend/internal/dto"
	"banking-system-backend/internal/requestctx"
	serviceInterfaces "banking-system-backend/internal/services/interfaces"
	"banking-system-backend/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

type BeneficiaryHandler struct {
	// beneficiaryService *services.BeneficiaryService
	beneficiaryService serviceInterfaces.BeneficiaryServiceInterface
}

func NewBeneficiaryHandler(s serviceInterfaces.BeneficiaryServiceInterface) *BeneficiaryHandler {
	return &BeneficiaryHandler{beneficiaryService: s}
}

func (h *BeneficiaryHandler) CreateBeneficiary(c *gin.Context) {
	ctx := c.Request.Context()
	log := requestctx.GetLogger(ctx)
	log.Info("create beneficiary request received")

	var req dto.BeneficiaryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn("failed to bind create beneficiary request",
			zap.Error(err),
		)
		utils.Error400(c, err)
		return
	}

	authCtx := requestctx.MustGetAuth(c)

	res, err := h.beneficiaryService.CreateBeneficiary(ctx, authCtx.UserID, req)
	if err != nil {
		log.Warn("failed to create beneficiary",
			zap.Error(err),
		)
		utils.HandleServiceError(c, err)
		return
	}

	log.Info("beneficiary created successfully",
		zap.String("beneficiary_id", res.ID.Hex()),
	)
	utils.Success(c, http.StatusCreated, res)
}

func (h *BeneficiaryHandler) DeleteBeneficiary(c *gin.Context) {
	ctx := c.Request.Context()
	log := requestctx.GetLogger(ctx)
	log.Info("delete beneficiary request received")

	idParam := c.Param("id")

	beneficiaryID, err := utils.ParseObjectID(idParam, "beneficiary ID")
	if err != nil {
		log.Warn("invalid beneficiary ID",
			zap.String("beneficiary_id", idParam),
			zap.Error(err),
		)
		utils.HandleServiceError(c, constants.ErrInvalidBeneficiaryID)
		return
	}
	authCtx := requestctx.MustGetAuth(c)

	err = h.beneficiaryService.SoftDeleteBeneficiary(ctx, beneficiaryID, authCtx.UserID)
	if err != nil {
		log.Warn("failed to delete beneficiary",
			zap.String("beneficiary_id", idParam),
			zap.Error(err),
		)
		utils.HandleServiceError(c, err)
		return
	}

	log.Info("beneficiary deleted successfully",
		zap.String("beneficiary_id", idParam),
	)
	utils.SuccessMessage(c, http.StatusOK, "beneficiary deleted successfully")
}

func (h *BeneficiaryHandler) GetBeneficiary(c *gin.Context) {
	ctx := c.Request.Context()
	log := requestctx.GetLogger(ctx)
	log.Info("get beneficiary request received")

	beneficiaryIDParam := c.Param("id")
	beneficiaryID, err := utils.ParseObjectID(beneficiaryIDParam, "beneficiary ID")
	if err != nil {
		log.Warn("invalid beneficiary ID",
			zap.String("beneficiary_id", beneficiaryIDParam),
			zap.Error(err),
		)
		utils.HandleServiceError(c, constants.ErrInvalidBeneficiaryID)
		return
	}

	authCtx := requestctx.MustGetAuth(c)

	beneficiary, err := h.beneficiaryService.GetBeneficiaryByID(ctx, beneficiaryID, authCtx.UserID)
	if err != nil {
		log.Warn("failed to get beneficiary",
			zap.String("beneficiary_id", beneficiaryIDParam),
			zap.Error(err),
		)
		utils.HandleServiceError(c, err)
		return
	}

	log.Info("beneficiary retrieved successfully",
		zap.String("beneficiary_id", beneficiary.ID.Hex()),
	)
	utils.Success(c, http.StatusOK, beneficiary)
}

func (h *BeneficiaryHandler) UpdateBeneficiary(c *gin.Context) {
	ctx := c.Request.Context()
	log := requestctx.GetLogger(ctx)
	log.Info("update beneficiary request received")

	beneficiaryIDParam := c.Param("id")
	beneficiaryID, err := primitive.ObjectIDFromHex(beneficiaryIDParam)
	if err != nil {
		log.Warn("invalid beneficiary ID",
			zap.String("beneficiary_id", beneficiaryIDParam),
			zap.Error(err),
		)
		utils.Error400(c, constants.ErrInvalidBeneficiaryID)
		return
	}

	var req dto.BeneficiaryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn("failed to bind update beneficiary request",
			zap.Error(err),
		)
		utils.Error400(c, err)
		return
	}

	authCtx := requestctx.MustGetAuth(c)

	updated, err := h.beneficiaryService.UpdateBeneficiary(ctx, beneficiaryID, authCtx.UserID, req)
	if err != nil {
		log.Warn("failed to update beneficiary",
			zap.String("beneficiary_id", beneficiaryIDParam),
			zap.Error(err),
		)
		utils.HandleServiceError(c, err)
		return
	}

	log.Info("beneficiary updated successfully",
		zap.String("beneficiary_id", updated.ID.Hex()),
	)
	utils.Success(c, http.StatusOK, updated)
}

func (h *BeneficiaryHandler) ListBeneficiaries(c *gin.Context) {
	ctx := c.Request.Context()
	log := requestctx.GetLogger(ctx)
	log.Info("list beneficiaries request received")

	authCtx := requestctx.MustGetAuth(c)

	res, err := h.beneficiaryService.ListBeneficiaries(ctx, authCtx.UserID)
	if err != nil {
		log.Warn("failed to list beneficiaries",
			zap.Error(err),
		)
		utils.Error500(c, err)
		return
	}

	log.Info("beneficiaries listed successfully",
		zap.Int("count", len(res)),
	)
	utils.Success(c, http.StatusOK, res)
}
