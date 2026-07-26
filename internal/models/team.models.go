package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type TeamMember struct {
	ID            bson.ObjectID `bson:"_id,omitempty" json:"id" binding:"required"`
	Name          string        `bson:"name" json:"name" binding:"required"`
	Role          string        `bson:"role" json:"role" binding:"required"`
	Description   string        `bson:"description" json:"description" binding:"required"`
	ImageURL      string        `bson:"image_url" json:"image_url" binding:"required"`
	ImagePublicID string        `bson:"image_public_id" json:"image_public_id" binding:"required"`
	CreatedAt     time.Time     `bson:"created_at" json:"created_at" binding:"required"`
	UpdatedAt     time.Time     `bson:"updated_at" json:"updated_at" binding:"required"`
}

type CreateTeamMemberRequest struct {
	Name          string `form:"name" binding:"required"`
	Role          string `form:"role" binding:"required"`
	Description   string `form:"description" binding:"required"`
	// ImageURL      string `json:"image_url" binding:"required"`
	// ImagePublicID string `json:"image_public_id" binding:"required"`
}

type UpdateTeamMemberRequest struct {
	Name          string `json:"name" binding:"required"`
	Role          string `json:"role" binding:"required"`
	Description   string `json:"description" binding:"required"`
	ImageURL      string `json:"image_url" binding:"required"`
	ImagePublicID string `json:"image_public_id" binding:"required"`
}