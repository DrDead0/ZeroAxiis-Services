package models

import (
	"time"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Blog struct {
	ID  bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Title  string `bson:"title" json:"title"`
	Author string `bson:"author" json:"author"`
	Content string `bson:"content" json:"content"`
	ImageURL string `bson:"image_url" json:"image_url"`
	ImagePublicID string `bson:"image_public_id" json:"image_public_id"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

type CreateBlogRequest struct {
	Title string `form:"title" binding:"required"`
	Author string `form:"author" binding:"required"`
	Content string `form:"content" binding:"required"`
}

type UpdateBlogRequest struct {
	Title string `form:"title"`
	Author string `form:"author"`
	Content string `form:"content"`
}
