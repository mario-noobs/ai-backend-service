package composer

import (
	"github.com/gin-gonic/gin"
	sctx "github.com/viettranx/service-context"
	"golang-ai-management/common"
	authService "golang-ai-management/service/auth"
	business "golang-ai-management/service/coordinate"
	authAPI "golang-ai-management/transport/api"
)

type AuthService interface {
	LoginHdl() func(*gin.Context)
	RegisterHdl() func(*gin.Context)
}

func ComposeAuthAPIService(serviceCtx sctx.ServiceContext) AuthService {
	jwtComp := serviceCtx.MustGet(common.KeyCompJWT).(common.JWTProvider)

	auth := authService.NewClient(ComposeUserAuthRPCClient(serviceCtx))

	hasher := new(common.Hasher)

	biz := business.NewBusiness(auth, jwtComp, hasher)
	serviceAPI := authAPI.NewAPI(serviceCtx, biz)

	return serviceAPI
}
