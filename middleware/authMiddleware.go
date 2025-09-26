package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

// Authentication middleware for ai-backend-service
// This service doesn't validate JWT tokens directly - it relies on auth-service via gRPC
func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("Authorization")
		if clientToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No Authorization header provided"})
			c.Abort()
			return
		}

		// TODO: Call auth-service via gRPC to validate token and get user info
		// For now, just pass the token through
		c.Set("token", clientToken)

		c.Next()
	}
}

const appId = "AiBackend"

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		xRequestID := uuid.New().String()
		timestamp := strconv.FormatInt(time.Now().Unix(), 10)
		c.Set("requestId", xRequestID+"_"+appId+"_"+timestamp)
		c.Next()
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "false")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
