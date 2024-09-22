package routes

import (
	"github.com/gin-gonic/gin"
	controller "golang-ai-management/controllers"
)

func FaceRegRoutes(incomingRoutes *gin.Engine) {

	faceController := &controller.FaceController{}
	facesRoutes := incomingRoutes.Group("/faces")
	facesRoutes.GET("/list", faceController.GetIdentities())       // Get all identities
	facesRoutes.GET("/:face_name", faceController.GetIdentity())   // Get a specific identity
	facesRoutes.POST("/enroll", faceController.FaceEnroll())       // Enroll a face
	facesRoutes.POST("/recognize", faceController.FaceRegconize()) // Recognize a face
}
