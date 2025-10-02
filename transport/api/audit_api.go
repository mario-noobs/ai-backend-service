package api

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"golang-ai-management/common"
	"golang-ai-management/models"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	sctx "github.com/viettranx/service-context"
	"github.com/viettranx/service-context/core"
)

type AuditBusiness interface {
	GetUserAuditLogs(ctx context.Context, userID string, limit, offset int) ([]*models.AuditLog, error)
}

type AuditAPI struct {
	serviceCtx    sctx.ServiceContext
	auditBusiness AuditBusiness
}

// JWT Claims for parsing user ID from token
type AuditJWTClaims struct {
	Sub string `json:"sub"`
	Exp int64  `json:"exp"`
	Iat int64  `json:"iat"`
}

// parseJWTTokenForAudit extracts user ID from JWT token
func parseJWTTokenForAudit(tokenString string) (*AuditJWTClaims, error) {
	// Remove "Bearer " prefix if present
	if strings.HasPrefix(tokenString, "Bearer ") {
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	}

	// Split JWT token into parts
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return nil, nil
	}

	// Decode payload (second part)
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, err
	}

	// Parse JSON claims
	var claims AuditJWTClaims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, err
	}

	return &claims, nil
}

func NewAuditAPI(serviceCtx sctx.ServiceContext, auditBusiness AuditBusiness) *AuditAPI {
	return &AuditAPI{
		serviceCtx:    serviceCtx,
		auditBusiness: auditBusiness,
	}
}

func (api *AuditAPI) ListAuditLogsHdl() func(*gin.Context) {
	return func(c *gin.Context) {
		// Get JWT token from context (set by auth middleware)
		tokenInterface, exists := c.Get("token")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication token not found"})
			return
		}

		tokenStr, ok := tokenInterface.(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			return
		}

		// Parse JWT to get user ID
		claims, err := parseJWTTokenForAudit(tokenStr)
		if err != nil || claims == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid JWT token"})
			return
		}

		// Get pagination parameters
		limit := 50 // default limit
		offset := 0 // default offset

		if limitStr := c.Query("limit"); limitStr != "" {
			if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
				limit = parsedLimit
				if limit > 100 { // max limit
					limit = 100
				}
			}
		}

		if offsetStr := c.Query("offset"); offsetStr != "" {
			if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
				offset = parsedOffset
			}
		}

		// Get user's audit logs
		auditLogs, err := api.auditBusiness.GetUserAuditLogs(context.Background(), claims.Sub, limit, offset)
		if err != nil {
			common.WriteErrorResponse(c, err)
			return
		}

		c.JSON(http.StatusOK, core.ResponseData(gin.H{
			"audit_logs": auditLogs,
			"user_id":    claims.Sub,
			"limit":      limit,
			"offset":     offset,
		}))
	}
}
