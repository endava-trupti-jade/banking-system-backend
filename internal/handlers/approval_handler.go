package handlers

import (
	"banking-system-backend/constants"
	approveInterfaces "banking-system-backend/internal/approval/interfaces"
	"banking-system-backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"strings"
)

type ApprovalHandler struct {
	approvalService approveInterfaces.ApprovalServiceInterface
}

func NewApprovalHandler(approvalService approveInterfaces.ApprovalServiceInterface) *ApprovalHandler {
	return &ApprovalHandler{approvalService: approvalService}
}

func (h *ApprovalHandler) ListRequests(c *gin.Context) {
	entity := c.Query("entity") // optional filter
	status := c.Query("status")

	requests, err := h.approvalService.List(c.Request.Context(), entity, status)
	if err != nil {
		utils.Error500(c, err)
		return
	}

	c.JSON(http.StatusOK, requests)
}

func (h *ApprovalHandler) GetRequest(c *gin.Context) {
	requestIDHex := strings.TrimSpace(c.Param("requestID"))

	requestID, err := primitive.ObjectIDFromHex(requestIDHex)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid requestID"})
		return
	}

	req, err := h.approvalService.GetRequestByID(c.Request.Context(), requestID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, req)
}

func (h *ApprovalHandler) Approve(c *gin.Context) {
	log.Println("ApprovalHandler Approve() started")

	requestIDHex := strings.TrimSpace(c.Param("requestID"))
	if requestIDHex == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request ID is required"})
		return
	}

	requestID, err := primitive.ObjectIDFromHex(requestIDHex)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
		return
	}

	log.Printf("ApprovalHandler Approve requestID=%s", requestIDHex)

	checkerID := c.MustGet("userID").(primitive.ObjectID)
	if checkerID == primitive.NilObjectID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid checker ID"})
		return
	}

	role := strings.TrimSpace(c.GetString("role"))
	rolePolicies := c.GetStringSlice("policies")

	err = h.approvalService.Decide(
		c.Request.Context(),
		requestID,
		checkerID,
		role,
		rolePolicies,
		constants.ActionApprove,
		"",
	)

	if err != nil {
		switch {

		case err == constants.ErrMakerCheckerRequestNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})

		case err == constants.ErrInvalidMakerCheckerStatus:
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})

		case err == constants.ErrMakerCheckerViolation:
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})

		case err == constants.ErrInvalidAction:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		default:
			//utils.Error500(c, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	log.Println("ApprovalHandler Approve() end")
	c.JSON(http.StatusOK, gin.H{"message": "Request approved successfully"})
}

type RejectRequest struct {
	Reason string `json:"reason"`
}

func (h *ApprovalHandler) Reject(c *gin.Context) {
	log.Println("ApprovalHandler Reject() started")

	requestIDHex := strings.TrimSpace(c.Param("requestID"))
	if requestIDHex == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request ID is required"})
		return
	}

	requestID, err := primitive.ObjectIDFromHex(requestIDHex)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
		return
	}

	checkerID := c.MustGet("userID").(primitive.ObjectID)
	if checkerID == primitive.NilObjectID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid checker ID"})
		return
	}

	role := strings.TrimSpace(c.GetString("role"))
	rolePolicies := c.GetStringSlice("policies")

	var req RejectRequest
	if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.Reason) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.ErrRejectReasonRequired.Error()})
		return
	}

	err = h.approvalService.Decide(
		c.Request.Context(),
		requestID,
		checkerID,
		role,
		rolePolicies,
		constants.ActionReject,
		req.Reason,
	)

	if err != nil {
		switch {

		case err == constants.ErrMakerCheckerRequestNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})

		case err == constants.ErrInvalidMakerCheckerStatus:
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})

		case err == constants.ErrMakerCheckerViolation:
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})

		case err == constants.ErrInvalidAction:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		case err == constants.ErrRejectReasonRequired:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		default:
			utils.Error500(c, err)
		}
		return
	}

	log.Println("ApprovalHandler Reject() end")
	c.JSON(http.StatusOK, gin.H{"message": "Request rejected successfully"})
}
