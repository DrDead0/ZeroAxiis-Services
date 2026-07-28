package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zeroaxiis/ZeroAxiis-Services/internal/database"
	"github.com/zeroaxiis/ZeroAxiis-Services/internal/models"
	"github.com/zeroaxiis/ZeroAxiis-Services/internal/pkg"
	"github.com/zeroaxiis/ZeroAxiis-Services/internal/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.uber.org/zap"
)

func GetTestimonials(c *gin.Context) {
	cachedData, err := utils.GetCache("testimonial")

	if err == nil {
		var testimonials []models.Testimonial

		err = json.Unmarshal(
			[]byte(cachedData),
			&testimonials,
		)
		if err == nil {
			pkg.Log.Info(
				"Testimonial Fetched From Redis Cache",
			)
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"data":    testimonials,
			})
			return
		}
		pkg.Log.Warn(
			"Failed to Decode Cached Testimonials",
			zap.Error(err),
		)
	}
	testimonialCollection := database.MongoClient.Database("zeroaxiiscom").Collection("testimonial")

	cursor, err := testimonialCollection.Find(
		context.Background(),
		bson.M{},
	)
	if err != nil {
		pkg.Log.Error(
			"Failed to Fetch Testimonials",
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to Fetch Testimonials",
		})
		return
	}
	var testimonials []models.Testimonial
	err = cursor.All(
		context.Background(),
		&testimonials,
	)
	if err != nil {
		pkg.Log.Error(
			"Failed to Decode the Testimonials",
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to Decode Testimonils",
		})
		return 
	}
	err = utils.SetCache(
		"testimonial",
		testimonials,
		12*time.Hour,
	)
	if err != nil {
		pkg.Log.Warn(
			"Failed to Cache Testimonial",
			zap.Error(err),
		)
	}

	pkg.Log.Info(
		"Testimonial Successfully fetched",
		zap.Any("data", testimonials),
	)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    testimonials,
	})

}
func CreateTestimonial(c *gin.Context) {
	var request models.CreateTestimonialRequest
	err := c.ShouldBind(&request)
	if err != nil {
		pkg.Log.Warn(
			"All fields are Required for testimonials",
			zap.Error(err),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "All Fields are required",
		})
		return
	}

	testimonials := models.Testimonial{
		Name:     request.Name,
		Role:     request.Role,
		Company:  request.Company,
		Comment:  request.Comment,
		CreateAt: time.Now(),
		UpdateAt: time.Now(),
	}

	testimonialCollection := database.MongoClient.Database("zeroaxiiscom").Collection("testimonial")

	_, err = testimonialCollection.InsertOne(
		context.Background(),
		testimonials,
	)
	if err != nil {
		pkg.Log.Error(
			"Failed to create testimonials",
			zap.Error(err),
		)
		c.JSON(
			http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Failed to Create Testimonials",
			})
		return
	}
	err = utils.DeleteCache("testimonial")
	if err != nil {
		pkg.Log.Warn(
			"Failed to Delete Testimonial Cache",
			zap.Error(err),
		)
	}
	pkg.Log.Info(
		"Testimonial Successfully Created",
		zap.Any("data", testimonials),
	)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Testimonial Created Successfully",
	})
}
func UpdateTestimonial(c *gin.Context) {
	id := c.Param("id")

	objectID, err := bson.ObjectIDFromHex(id)

	if err != nil {
		pkg.Log.Warn(
			"Invalid TestimonialID",
			zap.Error(err),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid TestimonialID",
		})
		return
	}
	testimonialCollection := database.MongoClient.Database("zeroaxiiscom").Collection("testimonial")
	var testimonials models.Testimonial

	err = testimonialCollection.FindOne(
		context.Background(),
		bson.M{
			"_id": objectID,
		},
	).Decode(&testimonials)

	if err != nil {
		pkg.Log.Warn(
			"Testimonial Not Found",
			zap.String("id", id),
		)
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Testimonial Not Found",
		})
		return
	}
	var request models.UpdateTestimonialRequest
	err = c.ShouldBind(&request)
	if err != nil {
		pkg.Log.Warn(
			"Invalid Update Request For Testimonial",
			zap.Error(err),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid Update Request For Testimonial",
		})
		return
	}
	update := bson.M{}
	if request.Name != "" {
		update["name"] = request.Name
	}
	if request.Role != "" {
		update["role"] = request.Role
	}
	if request.Company != "" {
		update["company"] = request.Company
	}
	if request.Comment != "" {
		update["comment"] = request.Comment
	}
	if len(update) == 0 {
		pkg.Log.Warn(
			"No Field Provided For Update",
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "No Field Provided For Update",
		})
		return
	}
	update["updated_at"] = time.Now()

	_, err = testimonialCollection.UpdateOne(
		context.Background(),
		bson.M{
			"_id": objectID,
		},
		bson.M{
			"$set": update,
		},
	)
	if err != nil {
		pkg.Log.Error(
			"Failed to Update Testimonial",
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to Update Testimonials",
		})
		return
	}
	err = utils.DeleteCache("testimonial")
	if err != nil {
		pkg.Log.Warn(
			"Failed To Delete Testimonial Cache",
			zap.Error(err),
		)
	}
	pkg.Log.Info(
		"Testimonial Update Successfully",
		zap.String("id", id),
	)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Testimonial Updated Successfully",
	})
}
func DeleteTestimonial(c *gin.Context) {
	id := c.Param("id")
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		pkg.Log.Warn(
			"Invalid Testimonial ID",
			zap.Error(err),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid Testimonial ID",
		})
		return
	}
	testimonialCollection := database.MongoClient.Database("zeroaxiiscom").Collection("testimonial")

	var testimonials models.Testimonial
	err = testimonialCollection.FindOne(
		context.Background(),
		bson.M{
			"_id": objectID,
		},
	).Decode(&testimonials)

	if err != nil {
		pkg.Log.Warn(
			"Testimonial not Found",
			zap.Error(err),
		)
		c.JSON(
			http.StatusNotFound, gin.H{
				"success": false,
				"message": "Testimonial Not Found",
			},
		)
		return
	}
	_, err = testimonialCollection.DeleteOne(
		context.Background(),
		bson.M{
			"_id": objectID,
		},
	)
	if err != nil {
		pkg.Log.Error(
			"Failed to Delete Testimonial",
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to Delete Testimonial",
		})
		return
	}
	err = utils.DeleteCache("testimonial")
	if err != nil {
		pkg.Log.Warn(
			"Failed to Delete Testimonial Cache",
			zap.Error(err),
		)
	}
	pkg.Log.Info(
		"Testimonial Deleted Successfully",
		zap.String("id", id),
	)
	c.JSON(http.StatusOK , gin.H{
		"success":true,
		"Message":"Testimonial Deleted Successfully",
	})
}
