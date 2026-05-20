package handlers

import (
	"banking-system-backend/constants"
	"banking-system-backend/internal/dto"
	"banking-system-backend/internal/requestctx"
	serviceInterfaces "banking-system-backend/internal/services/interfaces"
	"banking-system-backend/pkg/utils"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

type AccountNomineeHandler struct {
	// accountNomineeService *services.AccountNomineeService
	accountNomineeService serviceInterfaces.AccountNomineeServiceInterface
}

func NewAccountNomineeHandler(accountNomineeService serviceInterfaces.AccountNomineeServiceInterface) *AccountNomineeHandler {
	return &AccountNomineeHandler{accountNomineeService: accountNomineeService}
}

func (h *AccountNomineeHandler) AddAccountNominee(c *gin.Context) {
	ctx := c.Request.Context()
	log := requestctx.GetLogger(ctx)
	log.Info("create account nominee request received")

	var req dto.AccountNomineeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn("failed to bind create account nominee request",
			zap.Error(err),
		)
		utils.Error400(c, err)
		return
	}

	/*tenantIDHex := c.GetString("tenantID")
	tenantID, err := primitive.ObjectIDFromHex(tenantIDHex)
	if err != nil {
		log.Warn("failed to get tenant ID",
			zap.Error(err),
		)
		utils.Error401(c, err)
		return
	}*/

	authCtx := requestctx.MustGetAuth(c)

	accountNumber := strings.TrimSpace(c.Param("accountNumber"))
	if accountNumber == "" {
		log.Warn("account number required for adding nominee")
		utils.Error400(c, constants.ErrAccNumRequired)
		return
	}

	requestID, err := h.accountNomineeService.AddAccountNominee(ctx, authCtx.UserID, authCtx.Role, accountNumber, req)
	if err != nil {
		switch {

		case errors.Is(err, constants.ErrInvalidNomineePercentage),
			errors.Is(err, constants.ErrInvalidNomineeID),
			errors.Is(err, constants.ErrNomineeIDRequired),
			errors.Is(err, constants.ErrPrimaryNomineeAlreadyExists),
			errors.Is(err, constants.ErrTotalNomineePercentageExceeded):

			log.Warn("invalid nominee data while adding account nominee",
				zap.Error(err),
			)
			utils.Error400(c, err)
			return

		case errors.Is(err, constants.ErrNomineeNotFound),
			errors.Is(err, constants.ErrAccNotFound):

			log.Warn("account or nominee not found while adding account nominee",
				zap.String("account_number", accountNumber),
				zap.Error(err),
			)
			utils.Error404(c, err)
			return

		case errors.Is(err, constants.ErrNomineeNotActive):
			log.Warn("failed to add account nominee due to inactive nominee",
				zap.String("account_number", accountNumber),
				zap.Error(err),
			)
			utils.Error409(c, err)
			return

		case errors.Is(err, constants.ErrUnauthorizedAccountAccess):
			log.Warn("unauthorized account access while adding nominee",
				zap.String("account_number", accountNumber),
				zap.Error(err),
			)
			utils.Error403(c, err)
			return

		case errors.Is(err, constants.ErrNomineeAlreadyMappedToAccount):
			log.Warn("nominee already mapped to account",
				zap.String("account_number", accountNumber),
				zap.Error(err),
			)

			utils.Error409(c, err)
			return

		default:
			log.Warn("failed to add account nominee",
				zap.String("account_number", accountNumber),
				zap.Error(err),
			)
			utils.Error500(c, err)
			return
		}
	}

	log.Info("account nominee add request created successfully", zap.String("request_id", requestID.Hex()))
	utils.SuccessMessage(c, http.StatusOK, "Nominee add request created with ID: "+requestID.Hex())
}

// GetNominees godoc
// @Summary Get nominees for an account
// @Description Get nominees associated with an account
// @Tags Nominees
// @Security BearerAuth
// @Param accountNumber path string true "Account Number"
// @Success 200 {array} dto.AccountNomineeResponse
// @Failure 400 {object} map[string]string
// @Router /accounts/{accountNumber}/nominees [GET]
func (h *AccountNomineeHandler) ListAccountNomineesByAccountNumber(c *gin.Context) {
	ctx := c.Request.Context()
	log := requestctx.GetLogger(ctx)
	log.Info("list account nominees request received")

	accountNumber := strings.TrimSpace(c.Param("accountNumber"))
	if accountNumber == "" {
		log.Warn("account number required for listing nominees")
		utils.Error400(c, constants.ErrAccNumRequired)
		return
	}

	// Gin will automatically serialize the DTO slice.
	nominees, err := h.accountNomineeService.ListAccountNomineesByAccountNumber(ctx, accountNumber)
	if err != nil {
		if errors.Is(err, constants.ErrAccNotFound) {
			log.Warn("account not found while listing nominees",
				zap.String("account_number", accountNumber),
				zap.Error(err),
			)
			utils.Error404(c, err)
			return
		}
		utils.Error500(c, err)
		return
	}

	log.Info("account nominees retrieved successfully", zap.String("account_number", accountNumber))
	utils.Success(c, http.StatusOK, nominees)
}

func (h *AccountNomineeHandler) UpdateAccountNominee(c *gin.Context) {
	ctx := c.Request.Context()
	log := requestctx.GetLogger(ctx)
	log.Info("update account nominee request received")

	var req dto.AccountNomineeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn("failed to bind update account nominee request",
			zap.Error(err),
		)
		utils.Error400(c, err)
		return
	}

	mappingIDHex := strings.TrimSpace(c.Param("mappingId"))
	mappingID, err := primitive.ObjectIDFromHex(mappingIDHex)
	if err != nil {
		log.Warn("invalid mapping ID format while updating account nominee",
			zap.String("mapping_id", mappingIDHex),
			zap.Error(err),
		)
		utils.Error400(c, err)
		return
	}
	//		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid mapping id"})

	/*tenantIDHex := c.GetString("tenantID")
	tenantID, err := primitive.ObjectIDFromHex(tenantIDHex)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid tenant"})
		return
	}*/

	authCtx := requestctx.MustGetAuth(c)

	accountNumber := strings.TrimSpace(c.Param("accountNumber"))
	if accountNumber == "" {
		log.Warn("account number required while updating nominee",
			zap.String("mapping_id", mappingIDHex),
		)
		utils.Error400(c, constants.ErrAccNumRequired)
		return
	}

	requestID, err := h.accountNomineeService.UpdateAccountNominee(ctx, authCtx.UserID, authCtx.Role, accountNumber, mappingID, req)
	if err != nil {
		switch {
		case errors.Is(err, constants.ErrInvalidNomineePercentage),
			errors.Is(err, constants.ErrInvalidNomineeID),
			errors.Is(err, constants.ErrNomineeIDRequired),
			errors.Is(err, constants.ErrPrimaryNomineeAlreadyExists),
			errors.Is(err, constants.ErrTotalNomineePercentageExceeded):

			log.Warn("invalid nominee data while updating account nominee",
				zap.String("account_number", accountNumber),
				zap.String("mapping_id", mappingIDHex),
				zap.Error(err),
			)
			utils.Error400(c, err)
			return

		case errors.Is(err, constants.ErrNomineeNotFound),
			errors.Is(err, constants.ErrAccNotFound):

			log.Warn("account or nominee not found while updating account nominee",
				zap.String("account_number", accountNumber),
				zap.String("mapping_id", mappingIDHex),
				zap.Error(err),
			)
			utils.Error404(c, err)
			return

		case errors.Is(err, constants.ErrNomineeNotActive):

			log.Warn("nominee is not active while updating account nominee",
				zap.String("account_number", accountNumber),
				zap.String("mapping_id", mappingIDHex),
				zap.Error(err),
			)
			utils.Error409(c, err)
			return

		case errors.Is(err, constants.ErrUnauthorizedAccountAccess):
			log.Warn("unauthorized account access while updating nominee",
				zap.String("account_number", accountNumber),
				zap.String("mapping_id", mappingIDHex),
				zap.Error(err),
			)
			utils.Error403(c, err)
			return

		case errors.Is(err, constants.ErrNomineeAlreadyMappedToAccount):
			log.Warn("nominee already mapped to account while updating account nominee",
				zap.String("account_number", accountNumber),
				zap.String("mapping_id", mappingIDHex),
				zap.Error(err),
			)
			utils.Error409(c, err)
			return

		default:
			log.Error("failed to update account nominee",
				zap.String("account_number", accountNumber),
				zap.String("mapping_id", mappingIDHex),
				zap.Error(err),
			)
			utils.Error500(c, err)
			return
		}
	}

	log.Info("account nominee updated successfully", zap.String("account_number", accountNumber), zap.String("mapping_id", mappingIDHex))
	utils.SuccessMessage(c, http.StatusOK, "Nominee update request created with ID: "+requestID.Hex())
}

func (h *AccountNomineeHandler) SoftDeleteAccountNominee(c *gin.Context) {
	ctx := c.Request.Context()
	log := requestctx.GetLogger(ctx)
	log.Info("soft delete account nominee request received")

	mappingIDHex := strings.TrimSpace(c.Param("mappingId"))
	mappingID, err := primitive.ObjectIDFromHex(mappingIDHex)
	if err != nil {
		log.Warn("invalid mapping ID format while soft deleting account nominee",
			zap.String("mapping_id", mappingIDHex),
			zap.Error(err),
		)
		utils.Error400(c, err)
		return
	}

	/*tenantIDHex := c.GetString("tenantID")
	tenantID, err := primitive.ObjectIDFromHex(tenantIDHex)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid tenant"})
		return
	}*/

	authCtx := requestctx.MustGetAuth(c)

	accountNumber := strings.TrimSpace(c.Param("accountNumber"))
	if accountNumber == "" {
		log.Warn("account number required while soft deleting nominee",
			zap.String("mapping_id", mappingIDHex),
		)
		utils.Error400(c, constants.ErrAccNumRequired)
		return
	}

	requestID, err := h.accountNomineeService.SoftDeleteAccountNominee(ctx, authCtx.UserID, authCtx.Role, accountNumber, mappingID)
	if err != nil {
		switch {

		case errors.Is(err, constants.ErrInvalidNomineePercentage),
			errors.Is(err, constants.ErrInvalidNomineeID),
			errors.Is(err, constants.ErrNomineeIDRequired),
			errors.Is(err, constants.ErrPrimaryNomineeAlreadyExists),
			errors.Is(err, constants.ErrTotalNomineePercentageExceeded):

			log.Warn("invalid nominee data while soft deleting account nominee",
				zap.String("account_number", accountNumber),
				zap.String("mapping_id", mappingIDHex),
				zap.Error(err),
			)
			utils.Error400(c, err)
			return

		case errors.Is(err, constants.ErrNomineeNotFound),
			errors.Is(err, constants.ErrAccNotFound):

			log.Warn("account or nominee not found while soft deleting account nominee",
				zap.String("account_number", accountNumber),
				zap.String("mapping_id", mappingIDHex),
				zap.Error(err),
			)
			utils.Error404(c, err)
			return

		case errors.Is(err, constants.ErrNomineeNotActive):
			log.Warn("nominee is not active while soft deleting account nominee",
				zap.String("account_number", accountNumber),
				zap.String("mapping_id", mappingIDHex),
				zap.Error(err),
			)
			utils.Error409(c, err)
			return

		case errors.Is(err, constants.ErrUnauthorizedAccountAccess):
			log.Warn("unauthorized account access while soft deleting nominee",
				zap.String("account_number", accountNumber),
				zap.String("mapping_id", mappingIDHex),
				zap.Error(err),
			)
			utils.Error403(c, err)
			return

		case errors.Is(err, constants.ErrNomineeAlreadyMappedToAccount):
			log.Warn("nominee already mapped to account while soft deleting account nominee",
				zap.String("account_number", accountNumber),
				zap.String("mapping_id", mappingIDHex),
				zap.Error(err),
			)
			utils.Error409(c, err)
			return

		default:
			log.Error("failed to soft delete account nominee",
				zap.String("account_number", accountNumber),
				zap.String("mapping_id", mappingIDHex),
				zap.Error(err),
			)
			utils.Error500(c, err)
			return
		}
	}

	log.Info("soft delete account nominee completed", zap.String("account_number", accountNumber), zap.String("mapping_id", mappingIDHex))
	utils.SuccessMessage(c, http.StatusOK, "Nominee delete mapping request created with ID: "+requestID.Hex())
}
