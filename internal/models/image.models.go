package models

type ImageUploadReponse struct {
	ImageURL string `json:"image_url"`
	ImagePublicID string `json:"public_id"`
}