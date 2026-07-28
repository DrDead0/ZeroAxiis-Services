package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Testimonial struct {
	ID bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Name string `bson:"name" json:"name"`
	Role string `bson:"role" json:"role"`
	Company string `bson:"company" json:"company"`
	Comment string `bson:"comment" json:"comment"`
	CreateAt time.Time `bson:"created_at" json:"created_at"`
	UpdateAt time.Time `bson:"update_at" json:"update_at"`
}

type CreateTestimonialRequest struct{
	Name string `form:"name" binding:"required"`
	Role string `form:"role" binding:"required"`
	Company string `form:"company" binding:"required"`
	Comment string `form:"comment" binding:"required"`
}

type UpdateTestimonialRequest struct{
	Name string `form:"name"`
	Role string `form:"role"`
	Company string `form:"company"`
	Comment string `form:"comment"`
}
