package services

import (
	"context"
	"encoding/base64"
	"strings"
	"time"

	"backend-dragonhak/config"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

// UploadImage uploads a base64-encoded image to Cloudinary and returns the URL and public ID
func UploadImage(base64Image string, folder string) (string, string, error) {
	// Remove the data URL prefix if present
	base64Image = strings.Split(base64Image, ",")[1]

	// Decode the base64 string
	imageBytes, err := base64.StdEncoding.DecodeString(base64Image)
	if err != nil {
		return "", "", err
	}

	// Create a unique filename
	filename := time.Now().Format("20060102150405") + ".jpg"

	// Upload the image
	ctx := context.Background()
	uploadResult, err := config.Cld.Upload.Upload(
		ctx,
		imageBytes,
		uploader.UploadParams{
			PublicID: folder + "/" + filename,
			Folder:   folder,
		},
	)
	if err != nil {
		return "", "", err
	}

	return uploadResult.SecureURL, uploadResult.PublicID, nil
}

// GetImageURL returns the URL of an image given its public ID
func GetImageURL(publicID string) (string, error) {
	asset, err := config.Cld.Image(publicID)
	if err != nil {
		return "", err
	}
	url, err := asset.String()
	if err != nil {
		return "", err
	}
	return url, nil
}

// DeleteImage deletes an image from Cloudinary using its public ID
func DeleteImage(publicID string) error {
	ctx := context.Background()
	_, err := config.Cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: publicID,
	})
	return err
}
