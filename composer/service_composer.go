package composer

import (
	"golang-ai-management/common"
	helper "golang-ai-management/helpers"
	auditService "golang-ai-management/service/audit"
	authService "golang-ai-management/service/auth"
	faceBusiness "golang-ai-management/service/face"
	profileBusiness "golang-ai-management/service/profile"
	profileService "golang-ai-management/service/profile"
	"golang-ai-management/service/repository/mysql"
	auditAPI "golang-ai-management/transport/api"
	authAPI "golang-ai-management/transport/api"
	profileAPI "golang-ai-management/transport/api"

	"github.com/gin-gonic/gin"
	sctx "github.com/viettranx/service-context"
)

type AuthService interface {
	LoginHdl() func(*gin.Context)
	RegisterHdl() func(*gin.Context)
	LogoutHdl() func(*gin.Context)
	RefreshTokenHdl() func(*gin.Context)
}

type FaceServiceHandler interface {
	RegisterFaceHdl() func(*gin.Context)
	RecognizeFaceHdl() func(*gin.Context)
	DeleteFaceHdl() func(*gin.Context)
	IsRegisteredHdl() func(*gin.Context)
}

type ProfileHandler interface {
	GetProfileHdl() func(*gin.Context)
}

type AuditHandler interface {
	ListAuditLogsHdl() func(*gin.Context)
}

func ComposeAuthAPIService(serviceCtx sctx.ServiceContext) AuthService {
	auth := authService.NewClient(ComposeUserAuthRPCClient(serviceCtx))

	hasher := new(common.Hasher)
	time := new(helper.Timer)

	faceService := new(faceBusiness.FaceService)
	faceServiceConfig := new(faceBusiness.MarioFaceServiceConfig)

	// Pass nil JWT provider - auth-service handles all JWT logic via gRPC
	authBiz := authService.NewBusiness(auth, nil, hasher)

	faceBiz := faceBusiness.NewFaceBusiness(*faceService, *faceServiceConfig, *time)

	serviceAPI := authAPI.NewAPI(serviceCtx, authBiz, faceBiz)

	return serviceAPI
}

func ComposeFaceAPIService(serviceCtx sctx.ServiceContext) FaceServiceHandler {
	auth := authService.NewClient(ComposeUserAuthRPCClient(serviceCtx))

	hasher := new(common.Hasher)
	time := new(helper.Timer)

	faceService := new(faceBusiness.FaceService)
	faceServiceConfig := new(faceBusiness.MarioFaceServiceConfig)

	// Pass nil JWT provider - auth-service handles all JWT logic via gRPC
	authBiz := authService.NewBusiness(auth, nil, hasher)

	faceBiz := faceBusiness.NewFaceBusiness(*faceService, *faceServiceConfig, *time)

	serviceAPI := authAPI.NewAPI(serviceCtx, authBiz, faceBiz)

	return serviceAPI
}

func ComposeProfileAPIService(serviceCtx sctx.ServiceContext) ProfileHandler {
	profile := profileService.NewClient(ComposeProfileRPCClient(serviceCtx))

	// Pass nil JWT provider - auth-service handles all JWT logic via gRPC
	profileBiz := profileBusiness.NewBusiness(profile, nil)

	serviceAPI := profileAPI.NewProfileAPI(serviceCtx, profileBiz)

	return serviceAPI
}

func ComposeAuditLogAPIService(serviceCtx sctx.ServiceContext) AuditHandler {
	// Get MySQL component from service context
	mysqlComp := serviceCtx.MustGet(common.KeyCompMySQL).(common.GormComponent)
	db := mysqlComp.GetDB()

	// Create repository
	repo := mysql.NewMySQLRepository(db)

	// Create business layer
	auditBiz := auditService.NewAuditBusiness(repo)

	// Create API layer
	auditAPIService := auditAPI.NewAuditAPI(serviceCtx, auditBiz)

	return auditAPIService
}
