package middleware

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"golang-ai-management/common"
	"golang-ai-management/models"
	"golang-ai-management/service/repository/mysql"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	sctx "github.com/viettranx/service-context"
)

type JWTClaims struct {
	Sub string `json:"sub"`
	Exp int64  `json:"exp"`
	Iat int64  `json:"iat"`
}

// parseJWTToken extracts claims from JWT token without verification (for audit purposes only)
func parseJWTToken(tokenString string) (*JWTClaims, error) {
	// Remove "Bearer " prefix if present
	if strings.HasPrefix(tokenString, "Bearer ") {
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	}

	// Split JWT token into parts
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return nil, nil // Invalid JWT format
	}

	// Decode payload (second part)
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, err
	}

	// Parse JSON claims
	var claims JWTClaims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, err
	}

	return &claims, nil
}

// AuditLogging middleware logs all API calls to the database
func AuditLogging(serviceCtx sctx.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		// Calculate response time in milliseconds
		responseTime := time.Since(start).Milliseconds()

		// Get request details
		clientIP := c.ClientIP()
		apiCall := c.Request.URL.Path
		method := c.Request.Method
		status := c.Writer.Status()

		// Try to get user ID with priority order:
		// 1. audit_user_id (set by login/register after successful auth)
		// 2. JWT token from context
		var userID *string

		// First priority: Check for audit_user_id set by successful login/register
		if auditUserID, exists := c.Get("audit_user_id"); exists {
			if uid, ok := auditUserID.(string); ok {
				userID = &uid
			}
		}

		// Second priority: Try to get user ID from JWT token in context
		if userID == nil {
			if tokenInterface, exists := c.Get("token"); exists {
				if tokenStr, ok := tokenInterface.(string); ok {
					if claims, err := parseJWTToken(tokenStr); err == nil && claims != nil {
						userID = &claims.Sub
					}
				}
			}
		}

		// Create and save audit log asynchronously
		go func() {
			auditLog := models.NewAuditLog(clientIP, apiCall, method, status, responseTime, userID)

			mysqlComp := serviceCtx.MustGet(common.KeyCompMySQL)
			if gormComp, ok := mysqlComp.(common.GormComponent); ok {
				if db := gormComp.GetDB(); db != nil {
					repo := mysql.NewMySQLRepository(db)
					repo.AddAuditLog(context.Background(), &auditLog)
				}
			}
		}()
	}
}
