package handlers

import (
	"banking-system-backend/constants"
	"banking-system-backend/internal/dto"
	// "banking-system-backend/internal/services"
	serviceInterfaces "banking-system-backend/internal/services/interfaces"
	"banking-system-backend/pkg/utils"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type AuthHandler struct {
	authService serviceInterfaces.AuthServiceInterface
}

func NewAuthHandler(service serviceInterfaces.AuthServiceInterface) *AuthHandler {
	return &AuthHandler{authService: service}
}

// Regsiter godoc
// @Summary Register user
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Register Payload"
// @Success 200 {object} models.User
// @Router /auth/register [POST]
func (h *AuthHandler) Register(c *gin.Context) {
	log.Println("AuthHandler Register() started")

	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.authService.Register(c.Request.Context(), req); err != nil {
		if err != nil {
			if errors.Is(err, constants.ErrEmailExists) {
				c.JSON(http.StatusConflict, gin.H{"error": constants.ErrEmailExists.Error()})
				return
			}
		}
		utils.Error500(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": constants.MsgUserRegistration})

	log.Println("AuthHandler Register() end")
}

// Login godoc
// @Summary Login user
// @Tags Auth
// @Accept json
// @Produce json
// @param request body dto.LoginRequest true "Login Payload"
// @Success 200 {object} dto.LoginResponse
// @Router /auth/login [POST]
func (h *AuthHandler) Login(c *gin.Context) {
	log.Println("AuthHandler Login() started")

	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": constants.ErrInvalidCredentials.Error()})
		return
	}

	var response dto.LoginResponse
	response.Email = req.Email
	response.Role = res.Role
	response.Token = res.Token

	c.JSON(http.StatusOK, response)

	log.Println("AuthHandler Login() end")
}
