package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/sanda-bunescu/ExploRO/services"
)

func FirebaseAuthMiddleware(firebaseService services.FirebaseServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		firebaseUID, err := firebaseService.GetAndVerifyIDToken(c)
		if err != nil {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		// Store firebaseUID in the context for later use
		c.Set("firebaseUID", firebaseUID)
		c.Next()
	}
}
