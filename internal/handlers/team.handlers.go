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

func GetTeamMembers(c *gin.Context) {
	cachedData, err := utils.GetCache("team")
	if err == nil {
		var members []models.TeamMember
		err = json.Unmarshal(
			[]byte(cachedData),
			&members,
		)

		if err == nil {
			pkg.Log.Info(
				"Team Members Featched From Redis Cache",
			)
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"data":    members,
			})
			return
		}
		pkg.Log.Warn(
			"Failed to Decode Cached Team data",
			zap.Error(err),
		)
	}
	teamCollection := database.MongoClient.Database("zeroaxiiscom").Collection("team")

	cursor, err := teamCollection.Find(
		context.Background(),
		bson.M{},
	)
	if err != nil {
		pkg.Log.Error(
			"Failed to Fetch team members",
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Internal Server Error",
		})
		return
	}
	var members []models.TeamMember
	err = cursor.All(
		context.Background(),
		&members,
	)
	if err != nil {
		pkg.Log.Error(
			"Failed to decode Team Member",
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Internal Server Error",
		})
		return
	}
	err = utils.SetCache(
		"team",
		members,
		time.Hour,
	)
	if err != nil {
		pkg.Log.Warn(
			"Failed To Cache Team members",
			zap.Error(err),
		)
	}

	pkg.Log.Info(
		"Team Members Successfully Fetched",
		zap.Any("data", members),
	)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    members,
	})

}

func CreateTeamMember(c *gin.Context) {
	var request models.CreateTeamMemberRequest

	err := c.ShouldBind(&request)
	if err != nil {
		pkg.Log.Warn(
			"User has not filled all the fields",
			zap.Error(err),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "All fields are required",
		})
		return
	}

	teamCollection := database.MongoClient.Database("zeroaxiiscom").Collection("team")
	var existingMembers models.TeamMember
	err = teamCollection.FindOne(
		context.Background(),
		bson.M{
			"name": request.Name,
		},
	).Decode(&existingMembers)

	if err == nil {
		pkg.Log.Warn(
			"Team Member Already Exists",
			zap.String("name", request.Name),
		)
		c.JSON(http.StatusConflict, gin.H{
			"success": false,
			"message": "Team Member Already Exists",
		})
		return
	}

	fileHeader, err := c.FormFile("image")
	if err != nil {
		pkg.Log.Warn(
			"User Has Not Uploaded The Image",
			zap.Error(err),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Image is required..!",
		})
		return
	}
	imageResponse, err := utils.UploadImage(fileHeader, "Team")
	if err != nil {
		pkg.Log.Warn(
			"Failed to Upload Image",
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to upload image",
		})
		return
	}
	member := models.TeamMember{
		Name:          request.Name,
		Role:          request.Role,
		Description:   request.Description,
		ImageURL:      imageResponse.ImageURL,
		ImagePublicID: imageResponse.ImagePublicID,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	_, err = teamCollection.InsertOne(
		context.Background(),
		member,
	)
	if err != nil {
		pkg.Log.Error(
			"Failed to Create Team Member",
			zap.Error(err),
		)
		_ = utils.DeleteImage(imageResponse.ImagePublicID)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to create team member",
		})

		return
	}
	err = utils.DeleteCache("team")
	if err != nil {
		pkg.Log.Warn(
			"Failed to Delete Team Cache",
			zap.Error(err),
		)
	}
	pkg.Log.Info(
		"Team Member Created Successfully",
		zap.String("name", member.Name),
	)
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Team Member Created Successfully",
	})

}

func UpdateTeamMember(c *gin.Context) {

	id := c.Param("id")

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		pkg.Log.Warn(
			"Invlaid Team Member ID",
			zap.Error(err),
		)

		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid Team Member ID",
		})
		return
	}

	teamCollection := database.MongoClient.Database("zeroaxiiscom").Collection("team")

	var member models.TeamMember

	err = teamCollection.FindOne(
		context.Background(),
		bson.M{
			"_id": objectID,
		},
	).Decode(&member)

	if err != nil {
		pkg.Log.Warn(
			"Team Memeber Not Found",
			zap.String("id", id),
		)

		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Team Member Not Found",
		})
		return
	}

	var request models.UpdateTeamMemberRequest

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

	if request.Name != "" {
		update["name"] = request.Name
	}

	if request.Role != "" {
		update["role"] = request.Role
	}

	if request.Description != "" {
		update["description"] = request.Description
	}

	var imageResponse models.ImageUploadReponse

	fileHeader, err := c.FormFile("image")
	if err == nil {

		imageResponse, err = utils.UploadImage(
			fileHeader,
			"Team",
		)

		if err != nil {

			pkg.Log.Error(
				"Failed To Upload New Team Member Image",
				zap.Error(err),
			)

			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Failed To Upload Image",
			})

			return
		}

		update["image_url"] = imageResponse.ImageURL
		update["image_public_id"] = imageResponse.ImagePublicID
	}

	if len(update) == 0 {

		pkg.Log.Warn(
			"No Feild Provided For Update",
		)

		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "No Field Provided For Update",
		})

		return
	}

	update["updated_at"] = time.Now()
	_, err = teamCollection.UpdateOne(
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
			"Failed To Update Team Member",
			zap.Error(err),
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed To Update Team Member",
		})

		return
	}

	if imageResponse.ImagePublicID != "" {

		err = utils.DeleteImage(
			member.ImagePublicID,
		)

		if err != nil {

			pkg.Log.Warn(
				"Failed To Delete Old Team Member Image",
				zap.Error(err),
			)

		}
	}

	err = utils.DeleteCache("team")
	if err != nil {

		pkg.Log.Warn(
			"Failed To Delete Team Cache",
			zap.Error(err),
		)

	}

	pkg.Log.Info(
		"Team Member Updated Successfully",
		zap.String("id", id),
	)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Team Member Updated Successfully",
	})

}

func DeleteTeamMember(c *gin.Context) {
	id := c.Param("id")
	objectID, err := bson.ObjectIDFromHex(id)

	if err != nil {
		pkg.Log.Warn(
			"Invalid Team Member ID",
			zap.Error(err),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid Team Member ID",
		})
		return
	}
	teamCollection := database.MongoClient.Database("zeroaxiiscom").Collection("team")
	var members models.TeamMember

	err = teamCollection.FindOne(
		context.Background(),
		bson.M{
			"_id": objectID,
		},
	).Decode(&members)
	if err != nil {
		pkg.Log.Warn(
			"Team Member Not Found",
			zap.String("id", id),
		)
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Team Member not Found",
		})
		return
	}
	_, err = teamCollection.DeleteOne(
		context.Background(),
		bson.M{
			"_id": objectID,
		},
	)
	if err != nil {
		pkg.Log.Error(
			"Failed To Delete Team Member",
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed To Delete Team Member",
		})
		return
	}
	err = utils.DeleteCache("team")
	if err != nil {
		pkg.Log.Warn(
			"Failed to Delete Team Cache",
			zap.Error(err),
		)
	}
	err = utils.DeleteImage(members.ImagePublicID)
	if err != nil {
		pkg.Log.Warn(
			"Failed to Delete Team Member Image",
			zap.Error(err),
		)
	}

	pkg.Log.Info(
		"Team Member Deleted Successfully",
		zap.String("id", id),
	)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Team Memeber Deleted Successfully",
	})

}
