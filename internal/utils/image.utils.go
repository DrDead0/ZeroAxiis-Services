package utils

import (
	"context"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/zeroaxiis/ZeroAxiis-Services/internal/database"
	"github.com/zeroaxiis/ZeroAxiis-Services/internal/models"
)

// type ImageUploadReponse struct{
// 	ImageURL string `json:"image_url"`
// 	PublicID string `json:"public_id"`
// }

func UploadImage(fileHeader *multipart.FileHeader, folder string) (models.ImageUploadReponse, error) {
	var response models.ImageUploadReponse
	file, err := fileHeader.Open()
	if err != nil {
		return response, err
	}

	defer file.Close()
	result, err := database.CloudinaryClient.Upload.Upload(
		context.Background(),
		file,
		uploader.UploadParams{
			Folder: folder,
		},
	)

	if err != nil {
		return response, err
	}

	response.ImageURL = result.SecureURL
	response.ImagePublicID = result.PublicID

	return response, nil
}

func DeleteImage(publicID string) error {
	_, err := database.CloudinaryClient.Upload.Destroy(
		context.Background(),
		uploader.DestroyParams{
			PublicID: publicID,
		},
	)
	if err != nil {
		return err
	}

	return nil
}
