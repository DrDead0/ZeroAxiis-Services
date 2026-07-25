package database

import (
	"context"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/zeroaxiis/ZeroAxiis-Services/internal/config"
)

var CloudinaryClient *cloudinary.Cloudinary

func ConnectCloudinary() error {

	cfg := config.MustLoad()

	client, err := cloudinary.NewFromParams(
		cfg.CloudinaryCloudName,
		cfg.CloudinaryAPIKey,
		cfg.CloudinaryAPISecret,
	)

	if err != nil {
		return err
	}

	CloudinaryClient = client

	_ = context.Background()

	return nil
}