package api

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/ptypes/empty"
	sctx "github.com/viettranx/service-context"
	"github.com/viettranx/service-context/core"
	"golang-ai-management/common"
	helper "golang-ai-management/helpers"
	"golang-ai-management/models"
	"golang-ai-management/models/response"
	"golang-ai-management/proto/pb"
	"google.golang.org/grpc/metadata"
	"log/slog"
	"net/http"
	"os"
)

var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

type AuthBusiness interface {
	Login(ctx context.Context, data *pb.AuthEmailPassword) (*pb.TokenResponse, error)
	Register(ctx context.Context, in *pb.AuthRegister) (*empty.Empty, error)
}

type FaceBusiness interface {
	Enroll(ctx context.Context, data models.Face, jwt string) response.FaceRegResponse
	Recognize(ctx context.Context, data models.Face, jwt string) response.FaceRegResponse
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
			common.WriteErrorResponse(c, errors.New("missing authorization token"))
			return
		}
		// Step 2: Create metadata with the token
		md := metadata.Pairs("authorization", token)

		// Step 3: Create a new outgoing context with the metadata
		ctx := metadata.NewOutgoingContext(c.Request.Context(), md)

		resp, err := api.profileBusiness.GetUserProfile(ctx)

		if err != nil {
			logger.Error("response", "requestId", transactionId, "method", "GetProfileHdl", "err", err, "ms", api.time.End())
			common.WriteErrorResponse(c, err)
			return
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
		var data pb.AuthEmailPassword
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
		logger.Info("response", "requestId", transactionId, "method", method, "data", true, "ms", api.time.End())
		c.JSON(http.StatusOK, core.ResponseData(resp))
	}
}

func (api *api) RegisterHdl() func(*gin.Context) {
	return func(c *gin.Context) {
		var data pb.AuthRegister
		var method = "RegisterHdl"
		transactionId, exists := c.Get("requestId")
		if !exists {
			logger.Error("response", "method", method, "error", "TransactionId is null", "ms", api.time.End())
			c.JSON(http.StatusOK, gin.H{"error": "TransactionId is null"})
			return
		}
		data.AuthEmailPassword.TransactionId = transactionId.(string)
		api.time.Start()
		logger.Info("request", "method", method)
		if err := c.ShouldBind(&data); err != nil {
			logger.Error("response", "requestId", transactionId, "method", method, "error", err, "ms", api.time.End())
			common.WriteErrorResponse(c, core.ErrBadRequest.WithError(err.Error()))
			return
		}

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
