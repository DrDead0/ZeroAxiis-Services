package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zeroaxiis/ZeroAxiis-Services/internal/database"
	"github.com/zeroaxiis/ZeroAxiis-Services/internal/models"
	"github.com/zeroaxiis/ZeroAxiis-Services/internal/pkg"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.uber.org/zap"
)

func GetTeamMembers(c *gin.Context) {
	teamCollection:= database.MongoClient.Database("zeroaxiiscom").Collection("team")

	cursor, err := teamCollection.Find(
		context.Background(),
		bson.M{},
	)
	if err != nil{
		pkg.Log.Error(
			"Failed to Fetch team members",
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":false,
			"message":"Internal Server Error",
		})
		return
	}
	var members []models.TeamMember
	err = cursor.All(
		context.Background(),
		&members,
	)
	if err != nil{
		pkg.Log.Error(
			"Failed to decode Team Member",
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":false,
			"message":"Internal Server Error",
		})
		return
	}
	pkg.Log.Info(
		"Team Members Successfully Fetched",
		zap.Any("data",members),
	)
	c.JSON(http.StatusOK, gin.H{
		"success":true,
		"data":members,
	})

	
}

func CreateTeamMember(c *gin.Context) {
	var request models.CreateTeamMemberRequest

	err := c.ShouldBindJSON(&request)
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
	member := models.TeamMember{
		Name:          request.Name,
		Role:          request.Role,
		Description:   request.Description,
		ImageURL:      request.ImageURL,
		ImagePublicID: request.ImagePublicID,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	teamCollection := database.MongoClient.Database("zeroaxiiscom").Collection("team")
	var existingMembers models.TeamMember
	err = teamCollection.FindOne(
		context.Background(),
		bson.M{
			"name": member.Name,
		},
	).Decode(&existingMembers)

	if err == nil{
		pkg.Log.Warn(
			"Team Member Already Exists",
			zap.String("name", member.Name),
		)
		c.JSON(http.StatusConflict, gin.H{
			"success":false,
			"message":"Team Member Already Exists",
		})
		return 
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

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to Create Team Memeber",
		})
		return
	}
	pkg.Log.Info(
		"Team Member Created Successfully",
		zap.String("name", member.Name),
	)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Team Member Created Successfully",
	})

}

func UpdateTeamMember(c *gin.Context) {

}

func DeleteTeamMember(c *gin.Context) {

}
