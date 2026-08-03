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

func GetCreatives(c *gin.Context) {

	cachedData, err := utils.GetCache("creative")

	if err == nil {

		var creatives []models.Creative

		err = json.Unmarshal(
			[]byte(cachedData),
			&creatives,
		)

		if err == nil {

			pkg.Log.Info(
				"Creatives Fetched From Redis Cache",
			)

			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"data":    creatives,
			})

			return
		}

		pkg.Log.Warn(
			"Failed To Decode Cached Creatives",
			zap.Error(err),
		)
	}

	creativeCollection := database.MongoClient.
		Database("zeroaxiiscom").
		Collection("creative")

	cursor, err := creativeCollection.Find(
		context.Background(),
		bson.M{},
	)

	if err != nil {

		pkg.Log.Error(
			"Failed To Fetch Creatives",
			zap.Error(err),
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed To Fetch Creatives",
		})

		return
	}

	var creatives []models.Creative

	err = cursor.All(
		context.Background(),
		&creatives,
	)

	if err != nil {

		pkg.Log.Error(
			"Failed To Decode Creatives",
			zap.Error(err),
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed To Decode Creatives",
		})

		return
	}

	err = utils.SetCache(
		"creative",
		creatives,
		6*time.Hour,
	)

	if err != nil {

		pkg.Log.Warn(
			"Failed To Cache Creatives",
			zap.Error(err),
		)
	}

	pkg.Log.Info(
		"Creatives Successfully Fetched",
		zap.Any("data", creatives),
	)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    creatives,
	})
}

func CreateCreative(c *gin.Context) {

	var request models.CreateCreativeRequest

	err := c.ShouldBind(&request)
	if err != nil {

		pkg.Log.Warn(
			"All Creative Fields Are Required",
			zap.Error(err),
		)

		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "All Fields Are Required",
		})

		return
	}

	video, err := utils.GetYoutubeVideo(
		request.VideoURL,
	)

	if err != nil {

		pkg.Log.Error(
			"Failed To Fetch Video Details",
			zap.Error(err),
		)

		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid YouTube Video",
		})

		return
	}

	if len(video.Items) == 0 {

		pkg.Log.Warn(
			"No Video Found On YouTube",
		)

		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Video Not Found",
		})

		return
	}

	videoID, err := utils.ExtractVideoID(
		request.VideoURL,
	)

	if err != nil {

		pkg.Log.Error(
			"Failed To Extract Video ID",
			zap.Error(err),
		)

		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid YouTube URL",
		})

		return
	}

	creativeCollection := database.MongoClient.
		Database("zeroaxiiscom").
		Collection("creative")

	if request.Featured {

		_, err = creativeCollection.UpdateMany(
			context.Background(),
			bson.M{},
			bson.M{
				"$set": bson.M{
					"featured": false,
				},
			},
		)

		if err != nil {

			pkg.Log.Error(
				"Failed To Update Featured Creative",
				zap.Error(err),
			)

			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Failed To Create Creative",
			})

			return
		}
	}

	item := video.Items[0]

	publishedAt, err := time.Parse(
		time.RFC3339,
		item.Snippet.PublishedAt,
	)

	if err != nil {
		publishedAt = time.Now()
	}

	creative := models.Creative{
		VideoURL: request.VideoURL,
		VideoID:  videoID,

		Title:        item.Snippet.Title,
		Description:  item.Snippet.Description,
		ThumbnailURL: item.Snippet.Thumbnails.High.URL,

		ChannelTitle: item.Snippet.ChannelTitle,
		Duration:     item.ContentDetails.Duration,
		PublishedAt:  publishedAt,

		Summary:  request.Summary,
		Category: request.Category,
		Featured: request.Featured,

		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = creativeCollection.InsertOne(
		context.Background(),
		creative,
	)

	if err != nil {

		pkg.Log.Error(
			"Failed To Create Creative",
			zap.Error(err),
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed To Create Creative",
		})

		return
	}

	err = utils.DeleteCache("creative")
	if err != nil {

		pkg.Log.Warn(
			"Failed To Delete Creative Cache",
			zap.Error(err),
		)
	}

	pkg.Log.Info(
		"Creative Created Successfully",
		zap.String("title", creative.Title),
	)

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Creative Created Successfully",
	})
}
func UpdateCreative(c *gin.Context) {

	id := c.Param("id")

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {

		pkg.Log.Warn(
			"Invalid Creative ID",
			zap.Error(err),
		)

		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid Creative ID",
		})

		return
	}

	creativeCollection := database.MongoClient.
		Database("zeroaxiiscom").
		Collection("creative")

	var creative models.Creative

	err = creativeCollection.FindOne(
		context.Background(),
		bson.M{
			"_id": objectID,
		},
	).Decode(&creative)

	if err != nil {

		pkg.Log.Warn(
			"Creative Not Found",
			zap.String("id", id),
		)

		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Creative Not Found",
		})

		return
	}

	var request models.UpdateCreativeRequest

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

	if request.Summary != "" {
		update["summary"] = request.Summary
	}

	if request.Category != "" {
		update["category"] = request.Category
	}

	if request.Featured != nil {

		if *request.Featured {

			_, err = creativeCollection.UpdateMany(
				context.Background(),
				bson.M{
					"_id": bson.M{
						"$ne": objectID,
					},
				},
				bson.M{
					"$set": bson.M{
						"featured": false,
					},
				},
			)

			if err != nil {

				pkg.Log.Error(
					"Failed To Update Featured Creative",
					zap.Error(err),
				)

				c.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": "Failed To Update Creative",
				})

				return
			}
		}

		update["featured"] = *request.Featured
	}

	if request.VideoURL != "" {

		video, err := utils.GetYoutubeVideo(
			request.VideoURL,
		)

		if err != nil {

			pkg.Log.Error(
				"Failed To Fetch YouTube Video",
				zap.Error(err),
			)

			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Invalid YouTube Video",
			})

			return
		}

		videoID, err := utils.ExtractVideoID(
			request.VideoURL,
		)

		if err != nil {

			pkg.Log.Error(
				"Failed To Extract Video ID",
				zap.Error(err),
			)

			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Invalid YouTube URL",
			})

			return
		}
		if len(video.Items) == 0 {

			pkg.Log.Warn(
				"No Video Found On YouTube",
			)

			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Video Not Found",
			})

			return
		}

		item := video.Items[0]

		publishedAt, err := time.Parse(
			time.RFC3339,
			item.Snippet.PublishedAt,
		)

		if err != nil {
			publishedAt = time.Now()
		}

		update["video_url"] = request.VideoURL
		update["video_id"] = videoID
		update["title"] = item.Snippet.Title
		update["description"] = item.Snippet.Description
		update["thumbnail_url"] = item.Snippet.Thumbnails.High.URL
		update["channel_title"] = item.Snippet.ChannelTitle
		update["duration"] = item.ContentDetails.Duration
		update["published_at"] = publishedAt
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

	_, err = creativeCollection.UpdateOne(
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
			"Failed To Update Creative",
			zap.Error(err),
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed To Update Creative",
		})

		return
	}

	err = utils.DeleteCache("creative")
	if err != nil {

		pkg.Log.Warn(
			"Failed To Delete Creative Cache",
			zap.Error(err),
		)

	}

	pkg.Log.Info(
		"Creative Updated Successfully",
		zap.String("id", id),
	)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Creative Updated Successfully",
	})

}
func DeleteCreative(c *gin.Context) {

	id := c.Param("id")

	objectId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		pkg.Log.Warn(
			"Invalid Creative ID Format",
			zap.Error(err),
		)

		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid Creative ID Format",
		})
		return
	}

	creativeCollection := database.MongoClient.
		Database("zeroaxiiscom").
		Collection("creative")

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	result, err := creativeCollection.DeleteOne(
		ctx,
		bson.M{"_id": objectId},
	)

	if err != nil {
		pkg.Log.Error(
			"Failed To Delete Creative",
			zap.Error(err),
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed To Delete Creative",
		})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Creative Not Found",
		})
		return
	}

	err = utils.DeleteCache("creative")
	if err != nil {
		pkg.Log.Warn(
			"Failed To Delete Creative Cache",
			zap.Error(err),
		)
	}

	pkg.Log.Info(
		"Creative Deleted Successfully",
		zap.String("id", id),
	)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Creative Deleted Successfully",
	})
}
