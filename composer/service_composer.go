package composer

import (
	"github.com/gin-gonic/gin"
	sctx "github.com/viettranx/service-context"
	"golang-ai-management/common"
	authBusiness "golang-ai-management/service/coordinate"
	authSQLRepository "golang-ai-management/service/repository/mysql"
	authAPI "golang-ai-management/transport/api"
)

type AuthService interface {
	LoginHdl() func(*gin.Context)
	RegisterHdl() func(*gin.Context)
}

func ComposeAuthAPIService(serviceCtx sctx.ServiceContext) AuthService {
	db := serviceCtx.MustGet(common.KeyCompMySQL).(common.GormComponent)
	jwtComp := serviceCtx.MustGet(common.KeyCompJWT).(common.JWTProvider)

	authRepo := authSQLRepository.NewMySQLRepository(db.GetDB())
	hasher := new(common.Hasher)

	biz := authBusiness.NewBusiness(authRepo, jwtComp, hasher)
	serviceAPI := authAPI.NewAPI(serviceCtx, biz)

	return serviceAPI
}
