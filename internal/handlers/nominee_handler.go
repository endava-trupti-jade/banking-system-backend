package handlers

import (
	"banking-system-backend/constants"
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

type NomineeHandler struct {
	// nomineeService *services.NomineeService
	nomineeService serviceInterfaces.NomineeServiceInterface
}

func NewNomineeHandler(service serviceInterfaces.NomineeServiceInterface) *NomineeHandler {
	return &NomineeHandler{nomineeService: service}
}

// CreateNominee godoc
// @Summary Create a nominee
// @Description Create a nominee for an account
// @Tags Nominees
// @Security BearerAuth
// @Param request body dto.NomineeRequest true "Nominee Payload"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /api/nominees [POST]
func (h *NomineeHandler) CreateNominee(c *gin.Context) {
	ctx := c.Request.Context()
	log := requestctx.GetLogger(ctx)
	log.Info("create nominee request received")

	var req dto.NomineeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn("failed to bind create nominee request",
			zap.Error(err),
		)
		utils.Error400(c, err)
		return
	}

	authCtx := requestctx.MustGetAuth(c)

	nominee, err := h.nomineeService.CreateNominee(ctx, authCtx.UserID, req)
	if err != nil {
		log.Warn("failed to create nominee",
			zap.Error(err),
		)
		utils.HandleServiceError(c, err)
		return
	}

	log.Info("nominee created successfully",
		zap.String("nominee_id", nominee.ID.Hex()),
	)

	utils.SuccessMessage(c, http.StatusCreated, constants.MsgNomineeCreated)
}

// UpdateNominee godoc
// @Summary Update a nominee
// @Description Update nominee details
// @Tags Nominees
// @Security BearerAuth
// @Param nomineeId path string true "Nominee ID"
// @Param request body dto.NomineeRequest true "Nominee Payload"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /api/nominees/{nomineeId} [PUT]
func (h *NomineeHandler) UpdateNominee(c *gin.Context) {
	ctx := c.Request.Context()
	log := requestctx.GetLogger(ctx)
	log.Info("update nominee request received")

	var req dto.NomineeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn("failed to bind update nominee request",
			zap.Error(err),
		)
		utils.Error400(c, err)
		return
	}

	nomineeIDHex := c.Param("nomineeId")
	if nomineeIDHex == "" {
		log.Warn("nominee ID is required")
		utils.Error400(c, constants.ErrNomineeIDRequired)
		return
	}

	nomineeId, err := primitive.ObjectIDFromHex(nomineeIDHex)
	if err != nil {
		log.Warn("invalid nominee ID",
			zap.String("nominee_id", nomineeIDHex),
			zap.Error(err),
		)
		utils.Error400(c, constants.ErrInvalidNomineeID)
		return
	}

	authCtx := requestctx.MustGetAuth(c)

	nominee, err := h.nomineeService.UpdateNominee(ctx, nomineeId, authCtx.UserID, req)
	if err != nil {
		if errors.Is(err, constants.ErrNomineeNotFound) {
			log.Warn("nominee not found",
				zap.String("nominee_id", nomineeId.Hex()),
			)
			utils.Error404(c, err)
			return
		}

		if errors.Is(err, constants.ErrUnauthorized) {
			log.Warn("unauthorized request")
			utils.Error403(c, err)
			return
		}

		log.Warn("failed to update nominee",
			zap.String("nominee_id", nomineeId.Hex()),
			zap.Error(err),
		)
		utils.Error500(c, err)
		return
	}

	log.Info("nominee updated successfully",
		zap.String("nominee_id", nominee.ID.Hex()),
	)
	utils.Success(c, http.StatusOK, nominee)
}

// DeleteNominee godoc
// @Summary Delete a nominee
// @Description Remove a nominee from an account
// @Tags Nominees
// @Security BearerAuth
// @Param nomineeId path string true "Nominee ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/nominees/{nomineeId} [DELETE]
func (h *NomineeHandler) DeleteNominee(c *gin.Context) {
	ctx := c.Request.Context()
	log := requestctx.GetLogger(ctx)
	log.Info("delete nominee request received")

	nomineeIDHex := c.Param("nomineeId")
	if nomineeIDHex == "" {
		log.Warn("nominee ID is required")
		utils.Error400(c, constants.ErrNomineeIDRequired)
		return
	}

	nomineeID, err := utils.ParseObjectID(nomineeIDHex, "nominee ID")
	if err != nil {
		log.Warn("invalid nominee ID",
			zap.String("nominee_id", nomineeIDHex),
			zap.Error(err),
		)
		utils.Error400(c, constants.ErrInvalidNomineeID)
		return
	}

	authCtx := requestctx.MustGetAuth(c)

	if err := h.nomineeService.SoftDeleteNominee(ctx, nomineeID, authCtx.UserID); err != nil {
		if errors.Is(err, constants.ErrNomineeNotFound) {
			log.Warn("nominee not found",
				zap.String("nominee_id", nomineeID.Hex()),
			)
			utils.Error404(c, err)
			return
		}

		if errors.Is(err, constants.ErrUnauthorized) {
			log.Warn("unauthorized request")
			utils.Error403(c, err)
			return
		}

		if errors.Is(err, constants.ErrNomineeAlreadyMappedToAccount) {
			log.Warn("cannot delete nominee mapped to account",
				zap.String("nominee_id", nomineeID.Hex()),
				zap.Error(err),
			)
			utils.Error409(c, err)
			return
		}

		log.Warn("failed to delete nominee",
			zap.String("nominee_id", nomineeID.Hex()),
			zap.Error(err),
		)
		utils.Error500(c, err)
		return
	}

	log.Info("nominee deleted successfully",
		zap.String("nominee_id", nomineeID.Hex()),
	)
	utils.SuccessMessage(c, http.StatusOK, constants.MsgNomineeDeleted)
}

// GetNominee godoc
// @Summary Get nominee by ID
// @Description Retrieve nominee details by nominee ID
// @Tags Nominees
// @Security BearerAuth
// @Param nomineeId path string true "Nominee ID"
// @Success 200 {object} models.Nominee
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/nominees/{nomineeId} [GET]
func (h *NomineeHandler) GetNominee(c *gin.Context) {
	ctx := c.Request.Context()
	log := requestctx.GetLogger(ctx)
	log.Info("get nominee request received")

	nomineeIDHex := c.Param("nomineeId")
	if nomineeIDHex == "" {
		log.Warn("nominee ID is required")
		utils.Error400(c, constants.ErrNomineeIDRequired)
		return
	}

	nomineeID, err := utils.ParseObjectID(nomineeIDHex, "nominee ID")
	if err != nil {
		log.Warn("invalid nominee ID",
			zap.String("nominee_id", nomineeIDHex),
			zap.Error(err),
		)
		utils.Error400(c, constants.ErrInvalidNomineeID)
		return
	}

	authCtx := requestctx.MustGetAuth(c)

	nominee, err := h.nomineeService.GetNomineeByID(ctx, nomineeID, authCtx.UserID)
	if err != nil {
		if errors.Is(err, constants.ErrNomineeNotFound) {
			log.Warn("nominee not found",
				zap.String("nominee_id", nomineeID.Hex()),
			)
			utils.HandleServiceError(c, err)
			return
		}

		log.Warn("failed to get nominee",
			zap.String("nominee_id", nomineeID.Hex()),
			zap.Error(err),
		)
		utils.HandleServiceError(c, err)
		return
	}

	log.Info("nominee retrieved successfully",
		zap.String("nominee_id", nominee.ID.Hex()),
	)
	utils.Success(c, http.StatusOK, nominee)

}

func (h *NomineeHandler) ListNominees(c *gin.Context) {
	ctx := c.Request.Context()
	log := requestctx.GetLogger(ctx)
	log.Info("list nominees request received")

	authCtx := requestctx.MustGetAuth(c)

	nominees, err := h.nomineeService.ListNominees(ctx, authCtx.UserID)
	if err != nil {
		log.Warn("failed to list nominees",
			zap.Error(err),
		)
		utils.HandleServiceError(c, err)
		return
	}
	log.Info("nominees listed successfully",
		zap.Int("count", len(nominees)),
	)
	utils.Success(c, http.StatusOK, nominees)
}
