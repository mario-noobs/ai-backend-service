package routes

import (
	controller "golang-ai-management/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	userRoutes := incomingRoutes.Group("api/v1/user")
	userRoutes.GET("", controller.GetUsers())
	userRoutes.GET("/:user_id", controller.GetUser())
	userRoutes.POST("/signup", controller.SignUp())
	userRoutes.POST("/login", controller.Login())
}
