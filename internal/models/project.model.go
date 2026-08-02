package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Project struct {
	ID bson.ObjectID `bson:"_id,omitempty" json:"id"`

	Title string `bson:"title" json:"title"`

	Organization string `bson:"organization" json:"organization"`

	Description string `bson:"description" json:"description"`

	ProjectURL string `bson:"project_url" json:"project_url"`

	ImageURL string `bson:"image_url" json:"image_url"`

	ImagePublicID string `bson:"image_public_id" json:"image_public_id"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`

	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

type CreateProjectRequest struct {
	Title string `form:"title" binding:"required"`

	Organization string `form:"organization" binding:"required"`

	Description string `form:"description" binding:"required"`

	ProjectURL string `form:"project_url" binding:"required"`
}

type UpdateProjectRequest struct {
	Title string `form:"title"`

	Organization string `form:"organization"`

	Description string `form:"description"`

	ProjectURL string `form:"project_url"`
}