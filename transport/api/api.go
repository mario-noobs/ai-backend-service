package api

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/ptypes/empty"
	sctx "github.com/viettranx/service-context"
	"github.com/viettranx/service-context/core"
	"golang-ai-management/common"
	"golang-ai-management/models"
	"golang-ai-management/models/response"
	"golang-ai-management/proto/pb"
	"google.golang.org/grpc/metadata"
	"net/http"
)

type AuthBusiness interface {
	Login(ctx context.Context, data *pb.AuthEmailPassword) (*pb.TokenResponse, error)
	Register(ctx context.Context, in *pb.AuthRegister) (*empty.Empty, error)
}

type FaceBusiness interface {
	Enroll(ctx context.Context, data models.Face) response.FaceRegResponse
	Recognize(ctx context.Context, data models.Face) response.FaceRegResponse
}

type ProfileBusiness interface {
	GetUserProfile(context.Context) (*pb.User, error)
}

type api struct {
	serviceCtx      sctx.ServiceContext
	authBusiness    AuthBusiness
	faceBusiness    FaceBusiness
	profileBusiness ProfileBusiness
}

func (api api) GetProfileHdl() func(*gin.Context) {
	return func(c *gin.Context) {

		token := c.Request.Header.Get("Authorization")
		if token == "" {
			common.WriteErrorResponse(c, errors.New("missing authorization token"))
			return
		}
		// Step 2: Create metadata with the token
		md := metadata.Pairs("authorization", token)

		// Step 3: Create a new outgoing context with the metadata
		ctx := metadata.NewOutgoingContext(c.Request.Context(), md)

		response, err := api.profileBusiness.GetUserProfile(ctx)

		if err != nil {
			common.WriteErrorResponse(c, err)
			return
		}

		c.JSON(http.StatusOK, core.ResponseData(response))
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

		if err := c.ShouldBind(&data); err != nil {
			common.WriteErrorResponse(c, core.ErrBadRequest.WithError(err.Error()))
			return
		}

		response, err := api.authBusiness.Login(c.Request.Context(), &data)

		if err != nil {
			common.WriteErrorResponse(c, err)
			return
		}

		c.JSON(http.StatusOK, core.ResponseData(response))
	}
}

func (api *api) RegisterHdl() func(*gin.Context) {
	return func(c *gin.Context) {
		var data pb.AuthRegister

		if err := c.ShouldBind(&data); err != nil {
			common.WriteErrorResponse(c, core.ErrBadRequest.WithError(err.Error()))
			return
		}

		_, err := api.authBusiness.Register(c.Request.Context(), &data)

		if err != nil {
			common.WriteErrorResponse(c, err)
			return
		}

		c.JSON(http.StatusOK, core.ResponseData(true))
	}
}
func (api *api) RegisterFaceHdl() func(*gin.Context) {
	return func(c *gin.Context) {
		var data models.Face

		if err := c.ShouldBind(&data); err != nil {
			common.WriteErrorResponse(c, core.ErrBadRequest.WithError(err.Error()))
			return
		}

		resp := api.faceBusiness.Enroll(c.Request.Context(), data)

		c.JSON(http.StatusOK, core.ResponseData(resp))
	}
}

func (api *api) RecognizeFaceHdl() func(*gin.Context) {
	return func(c *gin.Context) {
		var data models.Face

		if err := c.ShouldBind(&data); err != nil {
			common.WriteErrorResponse(c, core.ErrBadRequest.WithError(err.Error()))
			return
		}

		resp := api.faceBusiness.Recognize(c.Request.Context(), data)

		c.JSON(http.StatusOK, core.ResponseData(resp))
	}
}