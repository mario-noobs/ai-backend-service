package composer

import (
	"github.com/gin-gonic/gin"
	sctx "github.com/viettranx/service-context"
	"golang-ai-management/common"
	helper "golang-ai-management/helpers"
	authService "golang-ai-management/service/auth"
	faceBusiness "golang-ai-management/service/face"
	profileBusiness "golang-ai-management/service/profile"
	profileService "golang-ai-management/service/profile"
	authAPI "golang-ai-management/transport/api"
	profileAPI "golang-ai-management/transport/api"
)

type AuthService interface {
	LoginHdl() func(*gin.Context)
	RegisterHdl() func(*gin.Context)
}

type FaceServiceHandler interface {
	RegisterFaceHdl() func(*gin.Context)
	RecognizeFaceHdl() func(*gin.Context)
}

type ProfileHandler interface {
	GetProfileHdl() func(*gin.Context)
}

func ComposeAuthAPIService(serviceCtx sctx.ServiceContext) AuthService {
	jwtComp := serviceCtx.MustGet(common.KeyCompJWT).(common.JWTProvider)

	auth := authService.NewClient(ComposeUserAuthRPCClient(serviceCtx))

	hasher := new(common.Hasher)
	time := new(helper.Timer)

	faceService := new(faceBusiness.FaceService)
	faceServiceConfig := new(faceBusiness.MarioFaceServiceConfig)

	authBiz := authService.NewBusiness(auth, jwtComp, hasher)

	faceBiz := faceBusiness.NewFaceBusiness(*faceService, *faceServiceConfig, *time)

	serviceAPI := authAPI.NewAPI(serviceCtx, authBiz, faceBiz)

	return serviceAPI
}

func ComposeFaceAPIService(serviceCtx sctx.ServiceContext) FaceServiceHandler {
	jwtComp := serviceCtx.MustGet(common.KeyCompJWT).(common.JWTProvider)

	auth := authService.NewClient(ComposeUserAuthRPCClient(serviceCtx))

	hasher := new(common.Hasher)
	time := new(helper.Timer)

	faceService := new(faceBusiness.FaceService)
	faceServiceConfig := new(faceBusiness.MarioFaceServiceConfig)

	authBiz := authService.NewBusiness(auth, jwtComp, hasher)

	faceBiz := faceBusiness.NewFaceBusiness(*faceService, *faceServiceConfig, *time)

	serviceAPI := authAPI.NewAPI(serviceCtx, authBiz, faceBiz)

	return serviceAPI
}

func ComposeProfileAPIService(serviceCtx sctx.ServiceContext) ProfileHandler {
	jwtComp := serviceCtx.MustGet(common.KeyCompJWT).(common.JWTProvider)

	profile := profileService.NewClient(ComposeProfileRPCClient(serviceCtx))

	profileBiz := profileBusiness.NewBusiness(profile, jwtComp)

	serviceAPI := profileAPI.NewProfileAPI(serviceCtx, profileBiz)

	return serviceAPI
}
