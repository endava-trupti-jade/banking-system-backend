package handlers

import (
	"banking-system-backend/constants"
	"banking-system-backend/internal/dto"
	serviceInterfaces "banking-system-backend/internal/services/interfaces"
	"banking-system-backend/pkg/utils"
	"errors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"strings"
)

type AccountNomineeHandler struct {
	// accountNomineeService *services.AccountNomineeService
	accountNomineeService serviceInterfaces.AccountNomineeServiceInterface
}

func NewAccountNomineeHandler(accountNomineeService serviceInterfaces.AccountNomineeServiceInterface) *AccountNomineeHandler {
	return &AccountNomineeHandler{accountNomineeService: accountNomineeService}
}

func (h *AccountNomineeHandler) AddAccountNominee(c *gin.Context) {
	log.Println("AccountNomineeHandler AddAccountNominee() started")

	var req dto.AccountNomineeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	/*tenantIDHex := c.GetString("tenantID")
	tenantID, err := primitive.ObjectIDFromHex(tenantIDHex)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid tenant"})
		return
	}*/

	userID := c.MustGet("userID").(primitive.ObjectID)
	if userID == primitive.NilObjectID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": constants.ErrUnauthorized.Error()})
		return
	}

	role := c.GetString("role")

	accountNumber := strings.TrimSpace(c.Param("accountNumber"))
	if accountNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Account number required"})
		return
	}

	requestID, err := h.accountNomineeService.AddAccountNominee(c.Request.Context(), userID, role, accountNumber, req)
	if err != nil {
		log.Println("AccountNomineeHandler error :", err)
		switch {

		case errors.Is(err, constants.ErrInvalidNomineePercentage),
			errors.Is(err, constants.ErrInvalidNomineeID),
			errors.Is(err, constants.ErrNomineeIDRequired),
			errors.Is(err, constants.ErrPrimaryNomineeAlreadyExists),
			errors.Is(err, constants.ErrTotalNomineePercentageExceeded):

			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return

		case errors.Is(err, constants.ErrNomineeNotFound),
			errors.Is(err, constants.ErrAccNotFound):

			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return

		case errors.Is(err, constants.ErrNomineeNotActive):

			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return

		case errors.Is(err, constants.ErrUnauthorizedAccountAccess):

			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return

		case errors.Is(err, constants.ErrNomineeAlreadyMappedToAccount):

			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return

		default:
			utils.Error500(c, err)
			return
		}
	}

	log.Println("AccountNomineeHandler AddAccountNominee() end")
	c.JSON(http.StatusCreated, gin.H{"message": "Nominee request created", "request_id": requestID})
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
	log.Println("AccountNomineeHandler ListAccountNomineesByAccountNumber() started")

	accountNumber := strings.TrimSpace(c.Param("accountNumber"))
	if accountNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Account Number is required"})
		return
	}

	// Gin will automatically serialize the DTO slice.
	nominees, err := h.accountNomineeService.ListAccountNomineesByAccountNumber(c.Request.Context(), accountNumber)
	if err != nil {
		if errors.Is(err, constants.ErrAccNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		utils.Error500(c, err)
		return
	}

	log.Println("AccountNomineeHandler ListAccountNomineesByAccountNumber() end")
	c.JSON(http.StatusOK, nominees)
}

func (h *AccountNomineeHandler) UpdateAccountNominee(c *gin.Context) {
	log.Println("AccountNomineeHandler UpdateAccountNominee() started")

	var req dto.AccountNomineeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mappingIDHex := strings.TrimSpace(c.Param("mappingId"))
	mappingID, err := primitive.ObjectIDFromHex(mappingIDHex)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid mapping id"})
		return
	}

	/*tenantIDHex := c.GetString("tenantID")
	tenantID, err := primitive.ObjectIDFromHex(tenantIDHex)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid tenant"})
		return
	}*/

	userID := c.MustGet("userID").(primitive.ObjectID)
	if userID == primitive.NilObjectID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": constants.ErrUnauthorized.Error()})
		return
	}

	role := c.GetString("role")

	accountNumber := strings.TrimSpace(c.Param("accountNumber"))
	if accountNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Account number required"})
		return
	}

	requestID, err := h.accountNomineeService.UpdateAccountNominee(c.Request.Context(), userID, role, accountNumber, mappingID, req)
	if err != nil {
		switch {

		case errors.Is(err, constants.ErrInvalidNomineePercentage),
			errors.Is(err, constants.ErrInvalidNomineeID),
			errors.Is(err, constants.ErrNomineeIDRequired),
			errors.Is(err, constants.ErrPrimaryNomineeAlreadyExists),
			errors.Is(err, constants.ErrTotalNomineePercentageExceeded):

			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return

		case errors.Is(err, constants.ErrNomineeNotFound),
			errors.Is(err, constants.ErrAccNotFound):

			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return

		case errors.Is(err, constants.ErrNomineeNotActive):

			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return

		case errors.Is(err, constants.ErrUnauthorizedAccountAccess):

			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return

		case errors.Is(err, constants.ErrNomineeAlreadyMappedToAccount):

			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return

		default:
			utils.Error500(c, err)
			return
		}
	}

	log.Println("AccountNomineeHandler UpdateAccountNominee() end")
	c.JSON(http.StatusOK, gin.H{"message": "Nominee update mapping request created", "request_id": requestID})
}

func (h *AccountNomineeHandler) SoftDeleteAccountNominee(c *gin.Context) {
	log.Println("AccountNomineeHandler SoftDeleteAccountNominee() started")

	mappingIDHex := strings.TrimSpace(c.Param("mappingId"))
	mappingID, err := primitive.ObjectIDFromHex(mappingIDHex)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid mapping id"})
		return
	}

	/*tenantIDHex := c.GetString("tenantID")
	tenantID, err := primitive.ObjectIDFromHex(tenantIDHex)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid tenant"})
		return
	}*/

	userID := c.MustGet("userID").(primitive.ObjectID)
	if userID == primitive.NilObjectID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": constants.ErrUnauthorized.Error()})
		return
	}

	role := c.GetString("role")

	accountNumber := strings.TrimSpace(c.Param("accountNumber"))
	if accountNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Account number required"})
		return
	}

	requestID, err := h.accountNomineeService.SoftDeleteAccountNominee(c.Request.Context(), userID, role, accountNumber, mappingID)
	if err != nil {
		switch {

		case errors.Is(err, constants.ErrInvalidNomineePercentage),
			errors.Is(err, constants.ErrInvalidNomineeID),
			errors.Is(err, constants.ErrNomineeIDRequired),
			errors.Is(err, constants.ErrPrimaryNomineeAlreadyExists),
			errors.Is(err, constants.ErrTotalNomineePercentageExceeded):

			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return

		case errors.Is(err, constants.ErrNomineeNotFound),
			errors.Is(err, constants.ErrAccNotFound):

			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return

		case errors.Is(err, constants.ErrNomineeNotActive):

			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return

		case errors.Is(err, constants.ErrUnauthorizedAccountAccess):

			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return

		case errors.Is(err, constants.ErrNomineeAlreadyMappedToAccount):

			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return

		default:
			utils.Error500(c, err)
			return
		}
	}

	log.Println("AccountNomineeHandler SoftDeleteAccountNominee() end")
	c.JSON(http.StatusOK, gin.H{"message": "Nominee delete mapping request created", "request_id": requestID})
}
