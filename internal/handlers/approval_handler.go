package handlers

import (
	"banking-system-backend/constants"
	approveInterfaces "banking-system-backend/internal/approval/interfaces"
	"banking-system-backend/internal/requestctx"
	"banking-system-backend/pkg/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

type ApprovalHandler struct {
	approvalService approveInterfaces.ApprovalServiceInterface
}

func NewApprovalHandler(approvalService approveInterfaces.ApprovalServiceInterface) *ApprovalHandler {
	return &ApprovalHandler{approvalService: approvalService}
}

func (h *ApprovalHandler) ListRequests(c *gin.Context) {
	ctx := c.Request.Context()
	log := requestctx.GetLogger(ctx)
	log.Info("list approval requests received")

	entity := c.Query("entity") // optional filter
	status := c.Query("status")

	requests, err := h.approvalService.List(c.Request.Context(), entity, status)
	if err != nil {
		log.Error("failed to list approval requests", zap.Error(err))
		utils.Error500(c, err)
		return
	}

	log.Info("approval requests listed successfully", zap.Int("count", len(requests)))
	utils.Success(c, http.StatusOK, requests)
}

func (h *ApprovalHandler) GetRequest(c *gin.Context) {
	ctx := c.Request.Context()
	log := requestctx.GetLogger(ctx)
	log.Info("get approval request received")

	requestIDHex := strings.TrimSpace(c.Param("requestID"))

	requestID, err := primitive.ObjectIDFromHex(requestIDHex)
	if err != nil {
		log.Warn("invalid request ID",
			zap.String("request_id", requestIDHex),
			zap.Error(err),
		)
		utils.Error400(c, constants.ErrInvalidRequestID)
		return
	}

	req, err := h.approvalService.GetRequestByID(c.Request.Context(), requestID)
	if err != nil {
		log.Warn("approval request not found",
			zap.String("request_id", requestID.Hex()),
			zap.Error(err),
		)
		utils.Error404(c, err)
		return
	}

	log.Info("approval request retrieved successfully",
		zap.String("request_id", req.ID.Hex()),
	)
	utils.Success(c, http.StatusOK, req)
}

func (h *ApprovalHandler) Approve(c *gin.Context) {
	ctx := c.Request.Context()
	log := requestctx.GetLogger(ctx)
	log.Info("ApprovalHandler Approve() started")

	requestIDHex := strings.TrimSpace(c.Param("requestID"))
	if requestIDHex == "" {
		log.Warn("request ID is required")
		utils.Error400(c, constants.ErrRequestIDRequired)
		return
	}

	requestID, err := primitive.ObjectIDFromHex(requestIDHex)
	if err != nil {
		log.Warn("invalid request ID",
			zap.String("request_id", requestIDHex),
			zap.Error(err),
		)
		utils.Error400(c, constants.ErrInvalidRequestID)
		return
	}

	log.Info("ApprovalHandler Approve requestID=%s", zap.String("request_id", requestIDHex))

	// checkerID := c.MustGet("userID").(primitive.ObjectID)
	authCtx := requestctx.MustGetAuth(c)
	checkerID := authCtx.UserID
	if checkerID == primitive.NilObjectID {
		log.Warn("invalid checker ID")
		utils.Error400(c, constants.ErrInvalidCheckerID)
		return
	}

	//role := strings.TrimSpace(c.GetString("role"))
	rolePolicies := c.GetStringSlice("policies")

	err = h.approvalService.Decide(
		c.Request.Context(),
		requestID,
		checkerID,
		authCtx.Role,
		rolePolicies,
		constants.ActionApprove,
		"",
	)

	if err != nil {
		switch {

		case err == constants.ErrMakerCheckerRequestNotFound:
			log.Warn("approval request not found",
				zap.String("request_id", requestID.Hex()),
				zap.Error(err),
			)
			utils.Error404(c, err)

		case err == constants.ErrInvalidMakerCheckerStatus:
			log.Warn("invalid approval request status",
				zap.String("request_id", requestID.Hex()),
				zap.Error(err),
			)
			utils.Error409(c, err)

		case err == constants.ErrMakerCheckerViolation:
			log.Warn("approval request violation",
				zap.String("request_id", requestID.Hex()),
				zap.Error(err),
			)
			utils.Error403(c, err)

		case err == constants.ErrInvalidAction:
			log.Warn("invalid action",
				zap.String("request_id", requestID.Hex()),
				zap.Error(err),
			)
			utils.Error400(c, err)

		default:
			log.Error("unexpected error",
				zap.String("request_id", requestID.Hex()),
				zap.Error(err),
			)
			utils.Error500(c, err)
		}
		return
	}

	log.Info("Request approved successfully", zap.String("request_id", requestID.Hex()))
	utils.SuccessMessage(c, http.StatusOK, "Request approved successfully")
}

type RejectRequest struct {
	Reason string `json:"reason"`
}

func (h *ApprovalHandler) Reject(c *gin.Context) {
	ctx := c.Request.Context()
	log := requestctx.GetLogger(ctx)
	log.Info("ApprovalHandler Reject() started")

	requestIDHex := strings.TrimSpace(c.Param("requestID"))
	if requestIDHex == "" {
		log.Warn("request ID is required")
		utils.Error400(c, constants.ErrRequestIDRequired)
		return
	}

	requestID, err := primitive.ObjectIDFromHex(requestIDHex)
	if err != nil {
		log.Warn("invalid request ID",
			zap.String("request_id", requestIDHex),
			zap.Error(err),
		)
		utils.Error400(c, constants.ErrInvalidRequestID)
		return
	}

	checkerID := c.MustGet("userID").(primitive.ObjectID)
	if checkerID == primitive.NilObjectID {
		log.Warn("invalid checker ID")
		utils.Error400(c, constants.ErrInvalidCheckerID)
		return
	}

	role := strings.TrimSpace(c.GetString("role"))
	rolePolicies := c.GetStringSlice("policies")

	var req RejectRequest
	if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.Reason) == "" {
		log.Warn("reject reason is required")
		utils.Error400(c, constants.ErrRejectReasonRequired)
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
			log.Warn("approval request not found",
				zap.String("request_id", requestID.Hex()),
				zap.Error(err),
			)
			utils.Error404(c, err)

		case err == constants.ErrInvalidMakerCheckerStatus:
			log.Warn("invalid approval request status",
				zap.String("request_id", requestID.Hex()),
				zap.Error(err),
			)
			utils.Error409(c, err)

		case err == constants.ErrMakerCheckerViolation:
			log.Warn("approval request violation",
				zap.String("request_id", requestID.Hex()),
				zap.Error(err),
			)
			utils.Error403(c, err)

		case err == constants.ErrInvalidAction:
			log.Warn("invalid action",
				zap.String("request_id", requestID.Hex()),
				zap.Error(err),
			)
			utils.Error400(c, err)

		case err == constants.ErrRejectReasonRequired:
			log.Warn("reject reason is required",
				zap.String("request_id", requestID.Hex()),
				zap.Error(err),
			)
			utils.Error400(c, err)

		default:
			log.Error("unexpected error",
				zap.String("request_id", requestID.Hex()),
				zap.Error(err),
			)
			utils.Error500(c, err)
		}
		return
	}

	log.Info("Request rejected successfully", zap.String("request_id", requestID.Hex()))
	utils.SuccessMessage(c, http.StatusOK, "Request rejected successfully")
}
