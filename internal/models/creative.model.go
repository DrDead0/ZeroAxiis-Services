package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Creative struct {
	ID bson.ObjectID `bson:"_id,omitempty" json:"id"`

	VideoURL string `bson:"video_url" json:"video_url"`
	VideoID string `bson:"video_id" json:"video_id"`

	Title string `bson:"title" json:"title"`
	Description string `bson:"description" json:"description"`
	ThumbnailURL string `bson:"thumbnail_url" json:"thumbnail_url"`

	ChannelTitle string `bson:"channel_title" json:"channel_title"`
	Duration string `bson:"duration" json:"duration"`
	PublishedAt time.Time `bson:"published_at" json:"published_at"`

	Summary string `bson:"summary" json:"summary"`

	Category string `bson:"category" json:"category"`

	Featured bool `bson:"featured" json:"featured"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

type CreateCreativeRequest struct {
	VideoURL string `form:"video_url" binding:"required"`

	Summary string `form:"summary" binding:"required"`

	Category string `form:"category" binding:"required"`

	Featured bool `form:"featured"`
}

type UpdateCreativeRequest struct {
	VideoURL string `form:"video_url"`

	Summary string `form:"summary"`

	Category string `form:"category"`

	Featured *bool `form:"featured"`
}