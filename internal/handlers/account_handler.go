package handlers

import (
	"banking-system-backend/constants"
	"banking-system-backend/internal/dto"
	serviceInterfaces "banking-system-backend/internal/services/interfaces"
	"banking-system-backend/pkg/utils"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	log.Println("AccountHandler CreateAccount() started")
	var req dto.CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("req : ", req, " err ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("userID").(primitive.ObjectID)
	if userID == primitive.NilObjectID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": constants.ErrUnauthorized.Error()})
		return
	}

	role := c.GetString("role")
	if role == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": constants.ErrRoleRequired.Error()})
		return
	}

	account, err := h.accountService.CreateAccount(c.Request.Context(), role, userID, req)
	if err != nil {
		log.Println("account : ", account, " err ", err)
		if err == constants.ErrAccCreationOwnershipDenied {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		utils.Error500(c, err)
		return
	}

	log.Println("AccountHandler CreateAccount() end")
	c.JSON(http.StatusCreated, account)
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
	log.Println("AccountHandler GetAccount() started")

	accountNumber := strings.TrimSpace(c.Param("accountNumber"))
	if accountNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.ErrAccNumRequired.Error()})
		return
	}

	userID := c.MustGet("userID").(primitive.ObjectID)
	if userID == primitive.NilObjectID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": constants.ErrUnauthorized.Error()})
		return
	}

	role := c.GetString("role")
	if role == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": constants.ErrRoleRequired.Error()})
		return
	}

	account, err := h.accountService.GetAccount(c.Request.Context(), accountNumber, role, userID)
	log.Println(" AccountHandler GetAccount : ", account, "=>", err)
	if err != nil {
		if err == constants.ErrOwnershipViolation {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		}
		return
	}

	log.Println("AccountHandler GetAccount() end")
	c.JSON(http.StatusOK, account)
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
	log.Println("AccountHandler UpdateAccount() started")
	accountNumber := strings.TrimSpace(c.Param("accountNumber"))
	if accountNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.ErrAccNumRequired.Error()})
		return
	}

	userID := c.MustGet("userID").(primitive.ObjectID)
	if userID == primitive.NilObjectID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": constants.ErrUnauthorized.Error()})
		return
	}

	role := c.GetString("role")
	if role == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": constants.ErrRoleRequired.Error()})
		return
	}

	var req dto.UpdateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	account, err := h.accountService.UpdateAccount(c.Request.Context(), accountNumber, role, userID, req)
	if err != nil {
		if err == constants.ErrOwnershipViolation {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": constants.ErrAccNotFound.Error()})
		}
		return
	}

	log.Println("AccountHandler UpdateAccount() end")
	c.JSON(http.StatusOK, account)
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
	log.Println("AccountHandler DeleteAccount() started")

	accountNumber := strings.TrimSpace(c.Param("accountNumber"))
	if accountNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.ErrAccNumRequired.Error()})
		return
	}

	userID := c.MustGet("userID").(primitive.ObjectID)
	if userID == primitive.NilObjectID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": constants.ErrUnauthorized.Error()})
		return
	}

	role := c.GetString("role")
	if role == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": constants.ErrRoleRequired.Error()})
		return
	}

	err := h.accountService.DeleteAccount(c.Request.Context(), accountNumber, role, userID)
	if err != nil {
		if err == constants.ErrOwnershipViolation {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		}
		return
	}

	log.Println("AccountHandler DeleteAccount() end")
	c.JSON(http.StatusOK, gin.H{"message": constants.MsgAccountDeletion})
}
