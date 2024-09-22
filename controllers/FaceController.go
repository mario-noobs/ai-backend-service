package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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
		identities, err := f.service.ListIdentities()
		if err != nil {
			// If an error occurs, return a JSON error response with 500 Internal Server Error
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		// If successful, return the list of identities as a JSON response
		c.JSON(http.StatusOK, gin.H{
			"identities": identities,
		})
	}
}

func (f *FaceController) GetIdentity() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, "GetIdentity")
	}
}

func (f *FaceController) FaceEnroll() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, "FaceEnroll")
	}
}

func (f *FaceController) FaceRegconize() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, "FaceRegconize")
	}
}
