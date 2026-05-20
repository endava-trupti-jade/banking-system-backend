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
	"go.uber.org/zap"
)

type AccountHandler struct {
	// accountService *services.AccountService
	accountService serviceInterfaces.AccountServiceInterface
}

func NewAccountHandler(service serviceInterfaces.AccountServiceInterface) *AccountHandler {
	return &AccountHandler{accountService: service}
}

// CreateAccount godoc
// @Summary Create new bank account
// @Description Create account for logged-in user
// @Tags Account
// @Security BearerAuth
// @Success 201 {object} models.Account
// @Failure 401 {object} map[string]string
// @Router /account/ [POST]
func (h *AccountHandler) CreateAccount(c *gin.Context) {
	ctx := c.Request.Context()
	log := requestctx.GetLogger(ctx)
	log.Info("create account request received")

	var req dto.CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn("failed binding request",
			zap.Error(err),
		)
		utils.Error400(c, err)
		return
	}

	authCtx := requestctx.MustGetAuth(c)

	account, err := h.accountService.CreateAccount(ctx, authCtx.Role, authCtx.UserID, req)
	if err != nil {
		if errors.Is(err, constants.ErrAccCreationOwnershipDenied) {
			log.Warn("account creation ownership denied",
				zap.Error(err),
			)
			utils.Error403(c, err)
			return
		}
		log.Error("failed to create account",
			zap.Error(err),
		)
		utils.Error500(c, err)
		return
	}

	log.Info("account created successfully", zap.String("account_number", account.AccountNumber))
	utils.Success(c, http.StatusCreated, account)
}

// GetAccount godoc
// @Summary Get account details
// @Description Fetch account by account number
// @Tags account
// @Security BearerAuth
// @Param accountNumber path string true "Account Number"
// @Success 200 {object} models.Account
// @Failure 404 {object} map[string]string
// @Router /account/{accountNumber} [GET]
func (h *AccountHandler) GetAccount(c *gin.Context) {
	ctx := c.Request.Context()
	log := requestctx.GetLogger(ctx)
	log.Info("get account request received")

	accountNumber := strings.TrimSpace(c.Param("accountNumber"))
	if accountNumber == "" {
		utils.Error400(c, constants.ErrAccNumRequired)
		log.Warn("account number required")
		return
	}

	authCtx := requestctx.MustGetAuth(c)

	account, err := h.accountService.GetAccount(ctx, accountNumber, authCtx.Role, authCtx.UserID)
	if err != nil {
		if errors.Is(err, constants.ErrOwnershipViolation) {
			log.Warn("ownership violation for account",
				zap.String("account_number", accountNumber),
				zap.Error(err),
			)

			utils.Error403(c, err)
			return
		}

		if errors.Is(err, constants.ErrAccNotFound) {
			log.Warn("account not found",
				zap.String("account_number", accountNumber),
				zap.Error(err),
			)

			utils.Error404(c, err)
			return
		}

		log.Error("failed to fetch account",
			zap.String("account_number", accountNumber),
			zap.Error(err),
		)
		utils.Error500(c, err)
		return
	}

	log.Info("account fetched successfully", zap.String("account_number", accountNumber))
	utils.Success(c, http.StatusOK, account)
}

//UpdateAccount godoc
// @Summary Update account details
// @Description Update account information by account number
// @Tags Account
// @Security BearerAuth
// @Param accountNumber path string true "Account Number"
// @Success 200 {object} models.Account
// @Failure 404 {object} map[string]string
// @Router /account/{accountNumber} [PUT]
func (h *AccountHandler) UpdateAccount(c *gin.Context) {
	ctx := c.Request.Context()
	log := requestctx.GetLogger(ctx)
	log.Info("account update request received")

	accountNumber := strings.TrimSpace(c.Param("accountNumber"))
	if accountNumber == "" {
		utils.Error400(c, constants.ErrAccNumRequired)
		log.Warn("account number required")
		return
	}

	authCtx := requestctx.MustGetAuth(c)

	var req dto.UpdateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn("failed to bind update account request",
			zap.Error(err),
		)
		utils.Error400(c, err)
		return
	}

	account, err := h.accountService.UpdateAccount(ctx, accountNumber, authCtx.Role, authCtx.UserID, req)
	if err != nil {
		if errors.Is(err, constants.ErrOwnershipViolation) {
			log.Warn("ownership violation for account",
				zap.String("account_number", accountNumber),
				zap.Error(err),
			)

			utils.Error403(c, err)
			return
		}

		if errors.Is(err, constants.ErrAccNotFound) {
			log.Warn("account not found",
				zap.String("account_number", accountNumber),
				zap.Error(err),
			)

			utils.Error404(c, err)
			return
		}

		log.Error("failed to update account",
			zap.String("account_number", accountNumber),
			zap.Error(err),
		)

		utils.Error500(c, err)
		return
	}

	log.Info("account updated successfully", zap.String("account_number", accountNumber))
	utils.Success(c, http.StatusOK, account)
}

// DeleteAccount godoc
// @Summary Delete account
// @Description Remove account by account number
// @Tags Account
// @Security BearerAuth
// @Param accountNumber path string true "Account Number"
// @Success 200 {object} map[string]string
// @Failure 404  {object} map[string]string
// @Router /account/{accountNumber} [DELETE]
func (h *AccountHandler) DeleteAccount(c *gin.Context) {
	ctx := c.Request.Context()
	log := requestctx.GetLogger(ctx)
	log.Info("account delete request received")

	accountNumber := strings.TrimSpace(c.Param("accountNumber"))
	if accountNumber == "" {
		log.Warn("account number required")
		utils.Error400(c, constants.ErrAccNumRequired)
		return
	}

	authCtx := requestctx.MustGetAuth(c)

	err := h.accountService.DeleteAccount(ctx, accountNumber, authCtx.Role, authCtx.UserID)
	if err != nil {
		if errors.Is(err, constants.ErrOwnershipViolation) {
			log.Warn("ownership violation for account",
				zap.String("account_number", accountNumber),
				zap.Error(err),
			)

			utils.Error403(c, err)
			return
		}

		if errors.Is(err, constants.ErrAccNotFound) {
			log.Warn("account not found",
				zap.String("account_number", accountNumber),
				zap.Error(err),
			)

			utils.Error404(c, err)
			return
		}

		log.Error("failed to delete account",
			zap.String("account_number", accountNumber),
			zap.Error(err),
		)

		utils.Error500(c, err)
		return
	}

	log.Info("account deleted successfully", zap.String("account_number", accountNumber))
	utils.SuccessMessage(c, http.StatusOK, constants.MsgAccountDeletion)
}
