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

func GetBlogs(c *gin.Context) {

	cachedData, err := utils.GetCache("blog")

	if err == nil {

		var blogs []models.Blog

		err = json.Unmarshal(
			[]byte(cachedData),
			&blogs,
		)

		if err == nil {

			pkg.Log.Info(
				"Blogs Fetched From Redis Cache",
			)

			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"data":    blogs,
			})

			return
		}

		pkg.Log.Warn(
			"Failed To Decode Cached Blogs",
			zap.Error(err),
		)
	}

	blogCollection := database.MongoClient.
		Database("zeroaxiiscom").
		Collection("blog")

	cursor, err := blogCollection.Find(
		context.Background(),
		bson.M{},
	)

	if err != nil {

		pkg.Log.Error(
			"Failed To Fetch Blogs",
			zap.Error(err),
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed To Fetch Blogs",
		})

		return
	}

	var blogs []models.Blog

	err = cursor.All(
		context.Background(),
		&blogs,
	)

	if err != nil {

		pkg.Log.Error(
			"Failed To Decode Blogs",
			zap.Error(err),
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed To Decode Blogs",
		})

		return
	}

	err = utils.SetCache(
		"blog",
		blogs,
		6*time.Hour,
	)

	if err != nil {

		pkg.Log.Warn(
			"Failed To Cache Blogs",
			zap.Error(err),
		)
	}

	pkg.Log.Info(
		"Blogs Successfully Fetched",
		zap.Any("data", blogs),
	)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    blogs,
	})
}
func CreateBlog(c *gin.Context) {

	var request models.CreateBlogRequest

	err := c.ShouldBind(&request)
	if err != nil {

		pkg.Log.Warn(
			"All Blog Fields Are Required",
			zap.Error(err),
		)

		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "All Fields Are Required",
		})

		return
	}

	fileHeader, err := c.FormFile("image")
	if err != nil {

		pkg.Log.Warn(
			"Blog Thumbnail Is Required",
			zap.Error(err),
		)

		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Thumbnail Is Required",
		})

		return
	}

	imageResponse, err := utils.UploadImage(
		fileHeader,
		"Blogs",
	)

	if err != nil {

		pkg.Log.Error(
			"Failed To Upload Blog Thumbnail",
			zap.Error(err),
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed To Upload Thumbnail",
		})

		return
	}

	blog := models.Blog{
		Title:         request.Title,
		Author:        request.Author,
		Content:       request.Content,
		ImageURL:      imageResponse.ImageURL,
		ImagePublicID: imageResponse.ImagePublicID,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	blogCollection := database.MongoClient.
		Database("zeroaxiiscom").
		Collection("blog")

	_, err = blogCollection.InsertOne(
		context.Background(),
		blog,
	)

	if err != nil {

		pkg.Log.Error(
			"Failed To Create Blog",
			zap.Error(err),
		)

		_ = utils.DeleteImage(
			imageResponse.ImagePublicID,
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed To Create Blog",
		})

		return
	}

	err = utils.DeleteCache("blog")
	if err != nil {

		pkg.Log.Warn(
			"Failed To Delete Blog Cache",
			zap.Error(err),
		)

	}

	pkg.Log.Info(
		"Blog Created Successfully",
		zap.String("title", blog.Title),
	)

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Blog Created Successfully",
	})

}
func UpdateBlog(c *gin.Context) {

	id := c.Param("id")

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {

		pkg.Log.Warn(
			"Invalid Blog ID",
			zap.Error(err),
		)

		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid Blog ID",
		})

		return
	}

	blogCollection := database.MongoClient.
		Database("zeroaxiiscom").
		Collection("blog")

	var blog models.Blog

	err = blogCollection.FindOne(
		context.Background(),
		bson.M{
			"_id": objectID,
		},
	).Decode(&blog)

	if err != nil {

		pkg.Log.Warn(
			"Blog Not Found",
			zap.String("id", id),
		)

		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Blog Not Found",
		})

		return
	}

	var request models.UpdateBlogRequest

	err = c.ShouldBind(&request)
	if err != nil {

		pkg.Log.Warn(
			"Invalid Update Request",
			zap.Error(err),
		)

		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid Request",
		})

		return
	}

	update := bson.M{}

	if request.Title != "" {
		update["title"] = request.Title
	}

	if request.Author != "" {
		update["author"] = request.Author
	}

	if request.Content != "" {
		update["content"] = request.Content
	}

	var imageResponse models.ImageUploadReponse

	fileHeader, err := c.FormFile("image")
	if err == nil {

		imageResponse, err = utils.UploadImage(
			fileHeader,
			"Blogs",
		)

		if err != nil {

			pkg.Log.Error(
				"Failed To Upload New Blog Thumbnail",
				zap.Error(err),
			)

			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Failed To Upload Thumbnail",
			})

			return
		}

		update["image_url"] = imageResponse.ImageURL
		update["image_public_id"] = imageResponse.ImagePublicID
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

	_, err = blogCollection.UpdateOne(
		context.Background(),
		bson.M{
			"_id": objectID,
		},
		bson.M{
			"$set": update,
		},
	)

	if err != nil {

		if imageResponse.ImagePublicID != "" {

			_ = utils.DeleteImage(
				imageResponse.ImagePublicID,
			)

		}

		pkg.Log.Error(
			"Failed To Update Blog",
			zap.Error(err),
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed To Update Blog",
		})

		return
	}

	if imageResponse.ImagePublicID != "" {

		err = utils.DeleteImage(
			blog.ImagePublicID,
		)

		if err != nil {

			pkg.Log.Warn(
				"Failed To Delete Old Blog Thumbnail",
				zap.Error(err),
			)

		}
	}

	err = utils.DeleteCache("blog")
	if err != nil {

		pkg.Log.Warn(
			"Failed To Delete Blog Cache",
			zap.Error(err),
		)

	}

	pkg.Log.Info(
		"Blog Updated Successfully",
		zap.String("id", id),
	)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Blog Updated Successfully",
	})

}
func DeleteBlog(c *gin.Context) {

	id := c.Param("id")

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {

		pkg.Log.Warn(
			"Invalid Blog ID",
			zap.Error(err),
		)

		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid Blog ID",
		})

		return
	}

	blogCollection := database.MongoClient.
		Database("zeroaxiiscom").
		Collection("blog")

	var blog models.Blog

	err = blogCollection.FindOne(
		context.Background(),
		bson.M{
			"_id": objectID,
		},
	).Decode(&blog)

	if err != nil {

		pkg.Log.Warn(
			"Blog Not Found",
			zap.String("id", id),
		)

		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Blog Not Found",
		})

		return
	}

	_, err = blogCollection.DeleteOne(
		context.Background(),
		bson.M{
			"_id": objectID,
		},
	)

	if err != nil {

		pkg.Log.Error(
			"Failed To Delete Blog",
			zap.Error(err),
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed To Delete Blog",
		})

		return
	}

	err = utils.DeleteImage(
		blog.ImagePublicID,
	)

	if err != nil {

		pkg.Log.Warn(
			"Failed To Delete Blog Thumbnail",
			zap.Error(err),
		)

	}

	err = utils.DeleteCache("blog")

	if err != nil {

		pkg.Log.Warn(
			"Failed To Delete Blog Cache",
			zap.Error(err),
		)

	}

	pkg.Log.Info(
		"Blog Deleted Successfully",
		zap.String("id", id),
	)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Blog Deleted Successfully",
	})

}
