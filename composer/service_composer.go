package composer

import (
	"github.com/gin-gonic/gin"
	sctx "github.com/viettranx/service-context"
	"golang-ai-management/common"
	authService "golang-ai-management/service/auth"
	authBusiness "golang-ai-management/service/coordinate"
	faceBusiness "golang-ai-management/service/face"
	authAPI "golang-ai-management/transport/api"
)

type AuthService interface {
	LoginHdl() func(*gin.Context)
	RegisterHdl() func(*gin.Context)
}

type FaceServiceHandler interface {
	RegisterFaceHdl() func(*gin.Context)
	RecognizeFaceHdl() func(*gin.Context)
}

func ComposeAuthAPIService(serviceCtx sctx.ServiceContext) AuthService {
	jwtComp := serviceCtx.MustGet(common.KeyCompJWT).(common.JWTProvider)

	auth := authService.NewClient(ComposeUserAuthRPCClient(serviceCtx))

	hasher := new(common.Hasher)

	faceService := new(faceBusiness.FaceService)
	faceServiceConfig := new(faceBusiness.MarioFaceServiceConfig)

	authBiz := authBusiness.NewBusiness(auth, jwtComp, hasher)

	faceBiz := faceBusiness.NewFaceBusiness(*faceService, *faceServiceConfig)

	serviceAPI := authAPI.NewAPI(serviceCtx, authBiz, faceBiz)

	return serviceAPI
}

func ComposeFaceAPIService(serviceCtx sctx.ServiceContext) FaceServiceHandler {
	jwtComp := serviceCtx.MustGet(common.KeyCompJWT).(common.JWTProvider)

	auth := authService.NewClient(ComposeUserAuthRPCClient(serviceCtx))

	hasher := new(common.Hasher)

	faceService := new(faceBusiness.FaceService)
	faceServiceConfig := new(faceBusiness.MarioFaceServiceConfig)

	authBiz := authBusiness.NewBusiness(auth, jwtComp, hasher)

	faceBiz := faceBusiness.NewFaceBusiness(*faceService, *faceServiceConfig)

	serviceAPI := authAPI.NewAPI(serviceCtx, authBiz, faceBiz)

	return serviceAPI
}
