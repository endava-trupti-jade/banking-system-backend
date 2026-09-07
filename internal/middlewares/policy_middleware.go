package middlewares

import (
	"banking-system-backend/constants"
	"banking-system-backend/internal/requestctx"
	"banking-system-backend/pkg/logger"
	"banking-system-backend/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func PolicyMiddleware(policyCodes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		log := requestctx.GetLogger(ctx)
		if log == nil {
			log = logger.Log
		}

		if log == nil {
			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				gin.H{"error": "logger unavailable"},
			)
			return
		}

		// Get role from centralized request context
		role, ok := requestctx.GetRole(ctx)
		if !ok || role == "" {
			log.Warn("role missing in request context")
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": constants.ErrRoleRequired.Error()})
			return
		}

		// userID, exists := c.Get("userID")
		// if !exists {
		// 	c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": constants.ErrUserIDRequired.Error()})
		// 	return
		// }

		// Get policies for this role
		rolePolicies, exists := utils.RolePolicies[role]
		if !exists || len(rolePolicies) == 0 {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": constants.ErrPolicyRequired.Error()})
			return
		}

		c.Set("policies", rolePolicies)

		// Convert role policies to a set for faster lookup
		policySet := make(map[string]struct{}, len(rolePolicies))

		for _, p := range rolePolicies {
			policySet[p] = struct{}{}
		}

		//Check if any of the allowed policies match the user's role policies
		log.Info(
			"checking authorization policies",
			zap.String("role", role),
			zap.Strings("required_policies", policyCodes),
			zap.Strings("role_policies", rolePolicies),
		)

		// Check permissions
		for _, policyCode := range policyCodes {
			// Check role has this policy
			if _, ok := policySet[policyCode]; ok {
				log.Info(
					"authorization successful",
					zap.String("matched_policy", policyCode),
				)
				c.Next()
				return
			}
		}

		log.Warn(
			"permission denied",
		)
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": constants.ErrPermissionDenied.Error()})
	}
}
