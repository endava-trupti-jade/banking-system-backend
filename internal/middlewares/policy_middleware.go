package middlewares

import (
	"banking-system-backend/constants"
	"banking-system-backend/pkg/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func PolicyMiddleware(policyCodes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		currentRole, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": constants.ErrRoleRequired.Error()})
			return
		}

		// userID, exists := c.Get("userID")
		// if !exists {
		// 	c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": constants.ErrUserIDRequired.Error()})
		// 	return
		// }

		role, ok := currentRole.(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "invalid role type"})
			return
		}
		//userID = userID.(string)

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
		log.Println("required policies:", policyCodes)
		log.Println("user role policies:", rolePolicies)
		for _, policyCode := range policyCodes {
			// Check role has this policy
			if _, ok := policySet[policyCode]; ok {
				c.Next()
				return
			}
		}

		log.Println("no perm")
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": constants.ErrPermissionDenied.Error()})
	}
}
