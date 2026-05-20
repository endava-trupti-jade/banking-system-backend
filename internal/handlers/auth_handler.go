package handlers

import (
	"banking-system-backend/constants"
	"banking-system-backend/internal/dto"
	"banking-system-backend/internal/requestctx"
	// "banking-system-backend/internal/services"
	serviceInterfaces "banking-system-backend/internal/services/interfaces"
	"banking-system-backend/pkg/utils"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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
	//log.Println("AuthHandler Register() started")
	ctx := c.Request.Context()
	log := requestctx.GetLogger(ctx)
	log.Info("create user request received")

	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn("invalid register request payload",
			zap.Error(err),
		)
		utils.Error400(c, err)
		return
	}

	if err := h.authService.Register(ctx, req); err != nil {
		if errors.Is(err, constants.ErrEmailExists) {
			log.Warn("registration failed email already exists",
				zap.String("email", req.Email),
			)

			utils.Error409(c, constants.ErrEmailExists)
			return
		}

		log.Error("registration failed",
			zap.String("email", req.Email),
			zap.Error(err),
		)
		utils.Error500(c, err)

		return
	}

	log.Info("user registered successfully",
		zap.String("email", req.Email),
	)
	utils.SuccessMessage(c, http.StatusCreated, constants.MsgUserRegistration)

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
	ctx := c.Request.Context()
	log := requestctx.GetLogger(ctx)
	log.Info("login request received")

	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn("invalid login request payload",
			zap.Error(err),
		)
		utils.Error400(c, err)
		return
	}

	res, err := h.authService.Login(ctx, req)
	if err != nil {
		log.Warn("login failed invalid credentials",
			zap.String("email", utils.NormalizeEmail(req.Email)),
			zap.String("client_ip", c.ClientIP()),
		)
		utils.Error401(c, constants.ErrInvalidCredentials)
		return
	}

	response := dto.LoginResponse{
		Email: utils.NormalizeEmail(req.Email),
		Role:  res.Role,
		Token: res.Token,
	}

	log.Info("user login successful",
		zap.String("email", utils.NormalizeEmail(req.Email)),
	)
	utils.Success(c, http.StatusOK, response)

}
