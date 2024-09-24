package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang-ai-management/models"
	"golang-ai-management/service"
	"net/http"
)

type FaceController struct {
	service service.MarioFaceService
}

var validate = validator.New()

func (f *FaceController) GetIdentities() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Call the listIdentities function
		response := f.service.ListIdentities()

		// If successful, return the list of identities as a JSON response
		c.JSON(http.StatusOK, response)
	}
}

func (f *FaceController) GetIdentity() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, "GetIdentity")
	}
}

func (f *FaceController) FaceEnroll() gin.HandlerFunc {
	return func(c *gin.Context) {
		var faceModel models.Face
		if err := c.ShouldBindJSON(&faceModel); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		response := f.service.Enroll(faceModel)
		c.JSON(http.StatusOK, response)
	}
}

func (f *FaceController) FaceRegconize() gin.HandlerFunc {
	return func(c *gin.Context) {
		var faceModel models.Face
		if err := c.ShouldBindJSON(&faceModel); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		response := f.service.Recognize(faceModel)
		c.JSON(http.StatusOK, response)
	}
}
