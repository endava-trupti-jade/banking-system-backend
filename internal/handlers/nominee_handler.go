package handlers

import (
	"banking-system-backend/constants"
	"banking-system-backend/internal/dto"
	serviceInterfaces "banking-system-backend/internal/services/interfaces"
	"banking-system-backend/pkg/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	log.Println("NomineeHandler CreateNominee() started")

	var req dto.NomineeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("userID").(primitive.ObjectID)
	if userID == primitive.NilObjectID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": constants.ErrUnauthorized.Error()})
		return
	}

	nominee, err := h.nomineeService.CreateNominee(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Println("NomineeHandler CreateNominee() end")

	c.JSON(http.StatusCreated, gin.H{"message": constants.MsgNomineeCreated, "id": nominee.ID})
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
	log.Println("NomineeHandler UpdateNominee() started")

	var req dto.NomineeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	nomineeIDHex := c.Param("nomineeId")
	if nomineeIDHex == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nominee ID is required"})
		return
	}

	nomineeId, err := primitive.ObjectIDFromHex(nomineeIDHex)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Nominee ID"})
		return
	}

	userID := c.MustGet("userID").(primitive.ObjectID)
	if userID == primitive.NilObjectID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": constants.ErrUnauthorized.Error()})
		return
	}

	nominee, err := h.nomineeService.UpdateNominee(c.Request.Context(), nomineeId, userID, req)
	if err != nil {
		if err == constants.ErrNomineeNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		if err == constants.ErrUnauthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}

		utils.Error500(c, err)
		return
	}

	log.Println("NomineeHandler UpdateNominee() end")
	c.JSON(http.StatusOK, nominee)
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
	log.Println("NomineeHandler DeleteNominee() started")
	nomineeIDHex := c.Param("nomineeId")
	if nomineeIDHex == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nominee ID is required"})
		return
	}

	nomineeID, err := primitive.ObjectIDFromHex(nomineeIDHex)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Nominee ID"})
		return
	}

	userID := c.MustGet("userID").(primitive.ObjectID)
	if userID == primitive.NilObjectID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": constants.ErrUnauthorized.Error()})
		return
	}

	if err := h.nomineeService.SoftDeleteNominee(c.Request.Context(), nomineeID, userID); err != nil {
		if err == constants.ErrNomineeNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		if err == constants.ErrUnauthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}

		if err == constants.ErrNomineeAlreadyMappedToAccount {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		utils.Error500(c, err)
		return
	}

	log.Println("NomineeHandler DeleteNominee() end")
	c.JSON(http.StatusOK, gin.H{"message": constants.MsgNomineeDeleted})
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
	log.Println("NomineeHandler GetNominee() started")

	nomineeIDHex := c.Param("nomineeId")
	if nomineeIDHex == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nominee ID is required"})
		return
	}

	nomineeID, err := primitive.ObjectIDFromHex(nomineeIDHex)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Nominee ID"})
		return
	}

	userID := c.MustGet("userID").(primitive.ObjectID)
	if userID == primitive.NilObjectID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": constants.ErrUnauthorized.Error()})
		return
	}

	nominee, err := h.nomineeService.GetNomineeByID(c.Request.Context(), nomineeID, userID)
	if err != nil {
		if err == constants.ErrNomineeNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err == constants.ErrUnauthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		utils.Error500(c, err)
		return
	}

	log.Println("NomineeHandler GetNominee() end")
	c.JSON(http.StatusOK, nominee)

}

func (h *NomineeHandler) ListNominees(c *gin.Context) {

	userID := c.MustGet("userID").(primitive.ObjectID)
	if userID == primitive.NilObjectID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	nominees, err := h.nomineeService.ListNominees(c.Request.Context(), userID)
	if err != nil {
		utils.Error500(c, err)
		return
	}

	c.JSON(http.StatusOK, nominees)
}
