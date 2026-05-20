package utils

import (
	"banking-system-backend/constants"
	"banking-system-backend/internal/config"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SuccessResponse struct {
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

func Success(c *gin.Context, status int, data any) {
	c.JSON(status, SuccessResponse{
		Success: true,
		Data:    data,
	})
}

func SuccessMessage(c *gin.Context, status int, message string) {
	c.JSON(status, SuccessResponse{
		Success: true,
		Message: message,
	})
}

func Fail(c *gin.Context, status int, message string) {
	c.JSON(status, ErrorResponse{
		Success: false,
		Error:   message,
	})
}

func buildErrorMessage(defaultMsg string, err error) string {
	if config.AppConfig.AppEnv == "development" {
		return err.Error()
	}
	return defaultMsg
}

func Error400(c *gin.Context, err error) {
	Fail(c, http.StatusBadRequest, buildErrorMessage(constants.ErrBadRequest.Error(), err))
}

func Error401(c *gin.Context, err error) {
	Fail(c, http.StatusUnauthorized, buildErrorMessage(constants.ErrUnauthorized.Error(), err))
}

func Error403(c *gin.Context, err error) {
	Fail(c, http.StatusForbidden, buildErrorMessage(constants.ErrForbidden.Error(), err))
}

func Error404(c *gin.Context, err error) {
	Fail(c, http.StatusNotFound, buildErrorMessage(constants.ErrNotFound.Error(), err))
}

func Error409(c *gin.Context, err error) {
	Fail(c, http.StatusConflict, buildErrorMessage(constants.ErrConflict.Error(), err))
}

func Error500(c *gin.Context, err error) {
	Fail(c, http.StatusInternalServerError, buildErrorMessage(constants.StatusInternalServerError, err))
}

func HandleServiceError(c *gin.Context, err error) {
	switch {

	// 400
	case errors.Is(err, constants.ErrInvalidBeneficiaryID),
		errors.Is(err, constants.ErrInvalidNomineeID),
		errors.Is(err, constants.ErrNomineeIDRequired),
		errors.Is(err, constants.ErrAccNumRequired),
		errors.Is(err, constants.ErrInvalidNomineePercentage),
		errors.Is(err, constants.ErrPrimaryNomineeAlreadyExists),
		errors.Is(err, constants.ErrTotalNomineePercentageExceeded):

		Error400(c, err)

	// 401
	case errors.Is(err, constants.ErrUnauthorized):

		Error401(c, err)

	// 403
	case errors.Is(err, constants.ErrOwnershipViolation),
		errors.Is(err, constants.ErrUnauthorizedAccountAccess):

		Error403(c, err)

	// 404
	case errors.Is(err, constants.ErrAccNotFound),
		errors.Is(err, constants.ErrNomineeNotFound),
		errors.Is(err, constants.ErrBeneficiaryNotFound):

		Error404(c, err)

	// 409
	case errors.Is(err, constants.ErrNomineeAlreadyMappedToAccount),
		errors.Is(err, constants.ErrNomineeNotActive):

		Error409(c, err)

	// default
	default:
		Error500(c, err)
	}
}
