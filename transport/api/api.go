package api

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"golang-ai-management/common"
	helper "golang-ai-management/helpers"
	"golang-ai-management/models"
	"golang-ai-management/models/response"
	"golang-ai-management/proto"
	"golang-ai-management/proto/pb"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/ptypes/empty"
	sctx "github.com/viettranx/service-context"
	"github.com/viettranx/service-context/core"
	"google.golang.org/grpc/metadata"
)

var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

// JWT Claims structure to parse the token
type JWTClaims struct {
	Sub string `json:"sub"`
	Exp int64  `json:"exp"`
	Iat int64  `json:"iat"`
}

// parseJWTTokenFromAPI extracts claims from JWT token without verification (for audit purposes only)
func parseJWTTokenFromAPI(tokenString string) (*JWTClaims, error) {
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

type AuthBusiness interface {
	Login(ctx context.Context, data *proto.AuthEmailPassword) (*proto.TokenResponse, error)
	Register(ctx context.Context, in *proto.AuthRegister) (*empty.Empty, error)
	Logout(ctx context.Context, data *proto.LogoutRequest) (*empty.Empty, error)
	RefreshToken(ctx context.Context, data *proto.RefreshTokenRequest) (*proto.TokenResponse, error)
}

type FaceBusiness interface {
	Enroll(ctx context.Context, data models.Face, jwt string) response.FaceRegResponse
	Recognize(ctx context.Context, data models.Face, jwt string) response.FaceRegResponse
	Delete(ctx context.Context, data models.Face, jwt string) response.FaceRegResponse
	IsUserRegistered(ctx context.Context, userId string, jwt string) (bool, error)
}

type ProfileBusiness interface {
	GetUserProfile(context.Context) (*pb.User, error)
}

type api struct {
	serviceCtx      sctx.ServiceContext
	authBusiness    AuthBusiness
	faceBusiness    FaceBusiness
	profileBusiness ProfileBusiness
	time            helper.Timer
}

func (api *api) GetProfileHdl() func(*gin.Context) {
	return func(c *gin.Context) {
		transactionId, exists := c.Get("requestId")
		if !exists {
			logger.Error("response", "method", "GetProfileHdl", "error", "TransactionId is null", "ms", api.time.End())
			c.JSON(http.StatusOK, gin.H{"error": "TransactionId is null"})
			return
		}

		api.time.Start()

		logger.Info("request", "requestId", transactionId, "method", "GetProfileHdl")
		token := c.Request.Header.Get("Authorization")
		if token == "" {
			logger.Error("response", "requestId", transactionId, "method", "GetProfileHdl", "error", "JWT not found", "ms", api.time.End())
			c.JSON(http.StatusUnauthorized, gin.H{"error": "JWT not found"})
			return
		}
		// Step 2: Create metadata with the token
		md := metadata.Pairs("authorization", token)

		// Step 3: Create a new outgoing context with the metadata
		ctx := metadata.NewOutgoingContext(c.Request.Context(), md)

		resp, err := api.profileBusiness.GetUserProfile(ctx)

		if err != nil {
			logger.Error("response", "requestId", transactionId, "method", "GetProfileHdl", "err", err, "ms", api.time.End())
			errMsg := err.Error()
			if errMsg == "invalid token" ||
				strings.Contains(errMsg, "invalid token") {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed: " + errMsg})
				return
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user profile: " + errMsg})
				return
			}
		}

		logger.Info("response", "requestId", transactionId, "method", "GetProfileHdl", "data", resp, "ms", api.time.End())
		c.JSON(http.StatusOK, core.ResponseData(resp))
	}
}

func NewAPI(serviceCtx sctx.ServiceContext, authBusiness AuthBusiness, faceBusiness FaceBusiness) *api {
	return &api{serviceCtx: serviceCtx, authBusiness: authBusiness, faceBusiness: faceBusiness}
}

func NewProfileAPI(serviceCtx sctx.ServiceContext, business ProfileBusiness) *api {
	return &api{serviceCtx: serviceCtx, profileBusiness: business}
}

func (api *api) LoginHdl() func(*gin.Context) {
	return func(c *gin.Context) {
		var data proto.AuthEmailPassword
		var method = "LoginHdl"
		transactionId, exists := c.Get("requestId")
		if !exists {
			logger.Error("response", "method", "GetProfileHdl", "requestId", transactionId, "error", "TransactionId is null", "ms", api.time.End())
			c.JSON(http.StatusOK, gin.H{"error": "TransactionId is null"})
			return
		}

		api.time.Start()
		logger.Info("request", "requestId", transactionId, "method", method)
		if err := c.ShouldBind(&data); err != nil {
			logger.Error("response", "requestId", transactionId, "method", method, "error", err, "ms", api.time.End())
			common.WriteErrorResponse(c, core.ErrBadRequest.WithError(err.Error()))
			return
		}

		data.TransactionId = transactionId.(string)

		resp, err := api.authBusiness.Login(c.Request.Context(), &data)

		if err != nil {
			logger.Error("response", "requestId", transactionId, "method", method, "error", err, "ms", api.time.End())
			common.WriteErrorResponse(c, err)
			return
		}

		if resp != nil && resp.AccessToken != nil && resp.AccessToken.Token != "" {
			if claims, parseErr := parseJWTTokenFromAPI(resp.AccessToken.Token); parseErr == nil && claims != nil {
				c.Set("audit_user_id", claims.Sub)
			}
		}

		logger.Info("response", "requestId", transactionId, "method", method, "data", true, "ms", api.time.End())
		c.JSON(http.StatusOK, core.ResponseData(resp))
	}
}

func (api *api) RegisterHdl() func(*gin.Context) {
	return func(c *gin.Context) {
		var data proto.AuthRegister
		var method = "RegisterHdl"
		transactionId, exists := c.Get("requestId")
		if !exists {
			logger.Error("response", "method", method, "error", "TransactionId is null", "ms", api.time.End())
			c.JSON(http.StatusOK, gin.H{"error": "TransactionId is null"})
			return
		}
		api.time.Start()
		logger.Info("request", "method", method)
		if err := c.ShouldBind(&data); err != nil {
			logger.Error("response", "requestId", transactionId, "method", method, "error", err, "ms", api.time.End())
			common.WriteErrorResponse(c, core.ErrBadRequest.WithError(err.Error()))
			return
		}

		data.AuthEmailPassword.TransactionId = transactionId.(string)

		_, err := api.authBusiness.Register(c.Request.Context(), &data)

		if err != nil {
			logger.Error("response", "requestId", transactionId, "method", method, "error", err, "ms", api.time.End())
			common.WriteErrorResponse(c, err)
			return
		}
		logger.Info("response", "requestId", transactionId, "method", method, "data", true, "ms", api.time.End())
		c.JSON(http.StatusOK, core.ResponseData(true))
	}
}

func (api *api) LogoutHdl() func(*gin.Context) {
	return func(c *gin.Context) {
		var data proto.LogoutRequest
		var method = "LogoutHdl"
		transactionId, exists := c.Get("requestId")
		if !exists {
			logger.Error("response", "method", method, "error", "TransactionId is null", "ms", api.time.End())
			c.JSON(http.StatusOK, gin.H{"error": "TransactionId is null"})
			return
		}
		api.time.Start()
		logger.Info("request", "method", method)

		tokenStr := c.Request.Header.Get("Authorization")
		if tokenStr == "" {
			logger.Error("response", "requestId", transactionId, "method", method, "error", "JWT not found", "ms", api.time.End())
			c.JSON(http.StatusUnauthorized, gin.H{"error": "JWT not found"})
			return
		}

		// Remove "Bearer " prefix if present (middleware passes full Authorization header)
		if len(tokenStr) > 7 && tokenStr[:7] == "Bearer " {
			tokenStr = tokenStr[7:]
		}

		data.AccessToken = tokenStr

		_, err := api.authBusiness.Logout(c.Request.Context(), &data)

		if err != nil {
			logger.Error("response", "requestId", transactionId, "method", method, "error", err, "ms", api.time.End())
			common.WriteErrorResponse(c, err)
			return
		}
		logger.Info("response", "requestId", transactionId, "method", method, "data", "logout successful", "ms", api.time.End())
		c.JSON(http.StatusOK, core.ResponseData(gin.H{"message": "logout successful"}))
	}
}

func (api *api) RefreshTokenHdl() func(*gin.Context) {
	return func(c *gin.Context) {
		var data proto.RefreshTokenRequest
		var method = "RefreshTokenHdl"
		transactionId, exists := c.Get("requestId")
		if !exists {
			logger.Error("response", "method", method, "error", "TransactionId is null", "ms", api.time.End())
			c.JSON(http.StatusOK, gin.H{"error": "TransactionId is null"})
			return
		}
		api.time.Start()
		logger.Info("request", "method", method)

		if err := c.ShouldBind(&data); err != nil {
			logger.Error("response", "requestId", transactionId, "method", method, "error", err, "ms", api.time.End())
			common.WriteErrorResponse(c, core.ErrBadRequest.WithError(err.Error()))
			return
		}

		resp, err := api.authBusiness.RefreshToken(c.Request.Context(), &data)

		if err != nil {
			logger.Error("response", "requestId", transactionId, "method", method, "error", err, "ms", api.time.End())
			common.WriteErrorResponse(c, err)
			return
		}
		logger.Info("response", "requestId", transactionId, "method", method, "data", "tokens refreshed", "ms", api.time.End())
		c.JSON(http.StatusOK, core.ResponseData(resp))
	}
}

func (api *api) RegisterFaceHdl() func(*gin.Context) {
	return func(c *gin.Context) {
		transactionId, exists := c.Get("requestId")
		if !exists {
			logger.Error("response", "method", "RegisterFaceHdl", "requestId", transactionId, "error", "TransactionId is null", "ms", api.time.End())
			c.JSON(http.StatusOK, gin.H{"error": "TransactionId is null"})
			return
		}
		var data models.Face

		data.TransactionId = transactionId.(string)
		api.time.Start()
		logger.Info("request", "requestId", transactionId, "method", "RegisterFaceHdl")
		jwtToken, exists := c.Get("token")
		if !exists {
			logger.Error("response", "requestId", transactionId, "method", "RegisterFaceHdl", "error", "JWT not found", "ms", api.time.End())
			c.JSON(http.StatusUnauthorized, gin.H{"error": "JWT not found"})
			return
		}

		if err := c.ShouldBind(&data); err != nil {
			logger.Error("response", "requestId", transactionId, "method", "RegisterFaceHdl", "error", err, "ms", api.time.End())
			common.WriteErrorResponse(c, core.ErrBadRequest.WithError(err.Error()))
			return
		}

		resp := api.faceBusiness.Enroll(c.Request.Context(), data, jwtToken.(string))
		logger.Info("response", "requestId", transactionId, "method", "RegisterFaceHdl", "data", resp, "ms", api.time.End())
		c.JSON(http.StatusOK, core.ResponseData(resp))
	}
}

func (api *api) RecognizeFaceHdl() func(*gin.Context) {
	return func(c *gin.Context) {
		transactionId, exists := c.Get("requestId")
		if !exists {
			logger.Error("response", "method", "GetProfileHdl", "error", "TransactionId is null", "ms", api.time.End())
			c.JSON(http.StatusOK, gin.H{"error": "TransactionId is null"})
			return
		}
		var data models.Face

		data.TransactionId = transactionId.(string)
		api.time.Start()

		logger.Info("request", "requestId", transactionId, "method", "RecognizeFaceHdl")
		jwtToken, exists := c.Get("token")
		if !exists {
			logger.Error("response", "requestId", transactionId, "method", "RecognizeFaceHdl", "error", "JWT not found", "ms", api.time.End())
			c.JSON(http.StatusUnauthorized, gin.H{"error": "JWT not found"})
			return
		}

		if err := c.ShouldBind(&data); err != nil {
			logger.Error("response", "requestId", transactionId, "method", "RecognizeFaceHdl", "error", err, "ms", api.time.End())
			common.WriteErrorResponse(c, core.ErrBadRequest.WithError(err.Error()))
			return
		}

		resp := api.faceBusiness.Recognize(c.Request.Context(), data, jwtToken.(string))
		logger.Info("response", "requestId", transactionId, "method", "RecognizeFaceHdl", "data", resp, "ms", api.time.End())
		c.JSON(http.StatusOK, core.ResponseData(resp))
	}
}

func (api *api) DeleteFaceHdl() func(*gin.Context) {
	return func(c *gin.Context) {
		transactionId, exists := c.Get("requestId")
		if !exists {
			logger.Error("response", "method", "DeleteFaceHdl", "error", "TransactionId is null", "ms", api.time.End())
			c.JSON(http.StatusOK, gin.H{"error": "TransactionId is null"})
			return
		}
		var data models.Face

		data.TransactionId = transactionId.(string)
		api.time.Start()

		logger.Info("request", "requestId", transactionId, "method", "DeleteFaceHdl")
		jwtToken, exists := c.Get("token")
		if !exists {
			logger.Error("response", "requestId", transactionId, "method", "DeleteFaceHdl", "error", "JWT not found", "ms", api.time.End())
			c.JSON(http.StatusUnauthorized, gin.H{"error": "JWT not found"})
			return
		}

		if err := c.ShouldBind(&data); err != nil {
			logger.Error("response", "requestId", transactionId, "method", "DeleteFaceHdl", "error", err, "ms", api.time.End())
			common.WriteErrorResponse(c, core.ErrBadRequest.WithError(err.Error()))
			return
		}

		resp := api.faceBusiness.Delete(c.Request.Context(), data, jwtToken.(string))
		logger.Info("response", "requestId", transactionId, "method", "DeleteFaceHdl", "data", resp, "ms", api.time.End())
		c.JSON(http.StatusOK, core.ResponseData(resp))
	}
}

func (api *api) IsRegisteredHdl() func(*gin.Context) {
	return func(c *gin.Context) {
		userId := c.Query("userId")
		if userId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "userId is required"})
			return
		}
		jwtToken, exists := c.Get("token")
		if !exists || jwtToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "JWT not found"})
			return
		}
		registered, err := api.faceBusiness.IsUserRegistered(c.Request.Context(), userId, jwtToken.(string))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"registered": registered})
	}
}
