package handlers

import (
	"backend-dragonhak/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UploadImageRequest struct {
	Base64Image string `json:"base64_image" binding:"required"`
	Folder      string `json:"folder" binding:"required"`
}

// UploadImage handles image upload requests
func UploadImage(c *gin.Context) {
	var req UploadImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Upload the image to Cloudinary
	imageURL, publicID, err := services.UploadImage(req.Base64Image, req.Folder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"url":       imageURL,
		"public_id": publicID,
	})
}

// GetImage retrieves an image URL by its public ID
func GetImage(c *gin.Context) {
	publicID := c.Param("public_id")
	if publicID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Public ID is required"})
		return
	}

	imageURL, err := services.GetImageURL(publicID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get image URL: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"url": imageURL,
	})
}

// DeleteImage deletes an image by its public ID
func DeleteImage(c *gin.Context) {
	publicID := c.Param("public_id")
	if publicID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Public ID is required"})
		return
	}

	err := services.DeleteImage(publicID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete image: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Image deleted successfully",
	})
}
