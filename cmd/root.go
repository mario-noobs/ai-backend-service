package cmd

import (
	"fmt"
	"golang-ai-management/common"
	"golang-ai-management/composer"
	"golang-ai-management/middleware"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	sctx "github.com/viettranx/service-context"
	"github.com/viettranx/service-context/component/ginc"
	smdlw "github.com/viettranx/service-context/component/ginc/middleware"
	"github.com/viettranx/service-context/component/gormc"
	// Removed jwtc import - not needed in ai-backend-service
)

func newServiceCtx() sctx.ServiceContext {
	return sctx.NewServiceContext(
		sctx.WithName("Coordinate Server"),
		sctx.WithComponent(ginc.NewGin(common.KeyCompGIN)),
		sctx.WithComponent(gormc.NewGormDB(common.KeyCompMySQL, "")),
		// Removed JWT component - all JWT logic handled by auth-service
		sctx.WithComponent(NewConfig()),
	)
}

var rootCmd = &cobra.Command{
	Use:   "app",
	Short: "Start service",
	Run: func(cmd *cobra.Command, args []string) {
		serviceCtx := newServiceCtx()

		logger := sctx.GlobalLogger().GetLogger("service")

		// Make some delay for DB ready (migration)
		// remove it if you already had your own DB
		time.Sleep(time.Second * 5)
		if err := serviceCtx.Load(); err != nil {
			logger.Fatal(err)
		}

		ginComp := serviceCtx.MustGet(common.KeyCompGIN).(common.GINComponent)

		router := ginComp.GetRouter()
		router.Use(gin.Recovery(), gin.Logger(), smdlw.Recovery(serviceCtx))

		router.Use(middleware.CORSMiddleware())
		router.Use(middleware.RequestID())
		router.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"data": "pong"})
		})

		//go StartGRPCServices(serviceCtx)

		v1User := router.Group("/api/v1/user")
		SetupUserRoutes(v1User, serviceCtx)

		v1Face := router.Group("/api/v1/face")
		v1Face.Use(middleware.Authentication())
		SetupFaceRoutes(v1Face, serviceCtx)

		v1Profile := router.Group("/")
		SetupProfileRoutes(v1Profile, serviceCtx)

		if err := router.Run(fmt.Sprintf(":%d", ginComp.GetPort())); err != nil {
			logger.Fatal(err)
		}
	},
}

func SetupUserRoutes(router *gin.RouterGroup, serviceCtx sctx.ServiceContext) {
	authAPIService := composer.ComposeAuthAPIService(serviceCtx)

	router.POST("/authenticate", authAPIService.LoginHdl())
	router.POST("/register", authAPIService.RegisterHdl())
	router.POST("/logout", authAPIService.LogoutHdl())
}

func SetupFaceRoutes(router *gin.RouterGroup, serviceCtx sctx.ServiceContext) {
	faceAPIService := composer.ComposeFaceAPIService(serviceCtx)

	router.POST("/register-identity", faceAPIService.RegisterFaceHdl())
	router.POST("/recognize-identity", faceAPIService.RecognizeFaceHdl())
	router.POST("/delete-identity", faceAPIService.DeleteFaceHdl())
	router.GET("/is-registered", faceAPIService.IsRegisteredHdl())
}

func SetupProfileRoutes(router *gin.RouterGroup, serviceCtx sctx.ServiceContext) {
	profileAPIService := composer.ComposeProfileAPIService(serviceCtx)

	router.POST("/profile", profileAPIService.GetProfileHdl())
}

func Execute() {
	rootCmd.AddCommand(outEnvCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
