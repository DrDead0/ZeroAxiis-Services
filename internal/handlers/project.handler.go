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

func GetProjects(c *gin.Context) {

	cachedData, err := utils.GetCache("project")

	if err == nil {

		var projects []models.Project

		err = json.Unmarshal(
			[]byte(cachedData),
			&projects,
		)

		if err == nil {

			pkg.Log.Info(
				"Projects Fetched From Redis Cache",
			)

			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"data":    projects,
			})

			return
		}

		pkg.Log.Warn(
			"Failed To Decode Cached Projects",
			zap.Error(err),
		)
	}

	projectCollection := database.MongoClient.
		Database("zeroaxiiscom").
		Collection("project")

	cursor, err := projectCollection.Find(
		context.Background(),
		bson.M{},
	)

	if err != nil {

		pkg.Log.Error(
			"Failed To Fetch Projects",
			zap.Error(err),
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed To Fetch Projects",
		})

		return
	}

	var projects []models.Project

	err = cursor.All(
		context.Background(),
		&projects,
	)

	if err != nil {

		pkg.Log.Error(
			"Failed To Decode Projects",
			zap.Error(err),
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed To Decode Projects",
		})

		return
	}

	err = utils.SetCache(
		"project",
		projects,
		6*time.Hour,
	)

	if err != nil {

		pkg.Log.Warn(
			"Failed To Cache Projects",
			zap.Error(err),
		)
	}

	pkg.Log.Info(
		"Projects Successfully Fetched",
		zap.Any("data", projects),
	)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    projects,
	})
}
func CreateProject(c *gin.Context) {

	var request models.CreateProjectRequest

	err := c.ShouldBind(&request)
	if err != nil {

		pkg.Log.Warn(
			"All Project Fields Are Required",
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
			"Project Thumbnail Is Required",
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
		"Projects",
	)

	if err != nil {

		pkg.Log.Error(
			"Failed To Upload Project Thumbnail",
			zap.Error(err),
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed To Upload Thumbnail",
		})

		return
	}

	project := models.Project{
		Title:         request.Title,
		Organization:  request.Organization,
		Description:   request.Description,
		ProjectURL:    request.ProjectURL,
		ImageURL:      imageResponse.ImageURL,
		ImagePublicID: imageResponse.ImagePublicID,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	projectCollection := database.MongoClient.
		Database("zeroaxiiscom").
		Collection("project")

	_, err = projectCollection.InsertOne(
		context.Background(),
		project,
	)

	if err != nil {

		_ = utils.DeleteImage(
			imageResponse.ImagePublicID,
		)

		pkg.Log.Error(
			"Failed To Create Project",
			zap.Error(err),
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed To Create Project",
		})

		return
	}

	err = utils.DeleteCache("project")
	if err != nil {

		pkg.Log.Warn(
			"Failed To Delete Project Cache",
			zap.Error(err),
		)

	}

	pkg.Log.Info(
		"Project Created Successfully",
		zap.String("title", project.Title),
	)

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Project Created Successfully",
	})

}
func UpdateProject(c *gin.Context) {

	id := c.Param("id")

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {

		pkg.Log.Warn(
			"Invalid Project ID",
			zap.Error(err),
		)

		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid Project ID",
		})

		return
	}

	projectCollection := database.MongoClient.
		Database("zeroaxiiscom").
		Collection("project")

	var project models.Project

	err = projectCollection.FindOne(
		context.Background(),
		bson.M{
			"_id": objectID,
		},
	).Decode(&project)

	if err != nil {

		pkg.Log.Warn(
			"Project Not Found",
			zap.String("id", id),
		)

		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Project Not Found",
		})

		return
	}

	var request models.UpdateProjectRequest

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

	if request.Organization != "" {
		update["organization"] = request.Organization
	}

	if request.Description != "" {
		update["description"] = request.Description
	}

	if request.ProjectURL != "" {
		update["project_url"] = request.ProjectURL
	}

	var imageResponse models.ImageUploadReponse

	fileHeader, err := c.FormFile("image")
	if err == nil {

		imageResponse, err = utils.UploadImage(
			fileHeader,
			"Projects",
		)

		if err != nil {

			pkg.Log.Error(
				"Failed To Upload New Project Thumbnail",
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

	_, err = projectCollection.UpdateOne(
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
			"Failed To Update Project",
			zap.Error(err),
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed To Update Project",
		})

		return
	}

	if imageResponse.ImagePublicID != "" {

		err = utils.DeleteImage(
			project.ImagePublicID,
		)

		if err != nil {

			pkg.Log.Warn(
				"Failed To Delete Old Project Thumbnail",
				zap.Error(err),
			)

		}
	}

	err = utils.DeleteCache("project")
	if err != nil {

		pkg.Log.Warn(
			"Failed To Delete Project Cache",
			zap.Error(err),
		)

	}

	pkg.Log.Info(
		"Project Updated Successfully",
		zap.String("id", id),
	)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Project Updated Successfully",
	})

}
func DeleteProject(c *gin.Context) {

	id := c.Param("id")

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {

		pkg.Log.Warn(
			"Invalid Project ID",
			zap.Error(err),
		)

		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid Project ID",
		})

		return
	}

	projectCollection := database.MongoClient.
		Database("zeroaxiiscom").
		Collection("project")

	var project models.Project

	err = projectCollection.FindOne(
		context.Background(),
		bson.M{
			"_id": objectID,
		},
	).Decode(&project)

	if err != nil {

		pkg.Log.Warn(
			"Project Not Found",
			zap.String("id", id),
		)

		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Project Not Found",
		})

		return
	}

	_, err = projectCollection.DeleteOne(
		context.Background(),
		bson.M{
			"_id": objectID,
		},
	)

	if err != nil {

		pkg.Log.Error(
			"Failed To Delete Project",
			zap.Error(err),
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed To Delete Project",
		})

		return
	}

	err = utils.DeleteImage(
		project.ImagePublicID,
	)

	if err != nil {

		pkg.Log.Warn(
			"Failed To Delete Project Thumbnail",
			zap.Error(err),
		)

	}

	err = utils.DeleteCache("project")

	if err != nil {

		pkg.Log.Warn(
			"Failed To Delete Project Cache",
			zap.Error(err),
		)

	}

	pkg.Log.Info(
		"Project Deleted Successfully",
		zap.String("id", id),
	)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Project Deleted Successfully",
	})

}