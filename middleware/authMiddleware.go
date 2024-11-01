package middleware

import (
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("Authorization")
		if clientToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("No Authorization header provided")})
			c.Abort()
			return
		}

		//claims, err := helper.ValidateToken(clientToken)
		//if err != "" {
		//	c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		//	c.Abort()
		//	return
		//}
		//
		c.Set("token", clientToken)
		//c.Set("first_name", claims.First_name)
		//c.Set("last_name", claims.Last_name)
		//c.Set("uid", claims.Uid)

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
