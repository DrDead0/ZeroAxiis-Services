package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/zeroaxiis/ZeroAxiis-Services/internal/config"
	"github.com/zeroaxiis/ZeroAxiis-Services/internal/models"
)

func ExtractVideoID(videoURL string) (string, error) {

	parsedURL, err := url.Parse(videoURL)
	if err != nil {
		return "", err
	}

	if strings.Contains(parsedURL.Host, "youtu.be") {

		videoID := strings.TrimPrefix(
			parsedURL.Path,
			"/",
		)

		if videoID == "" {
			return "", errors.New("invalid youtube url")
		}

		return videoID, nil
	}

	videoID := parsedURL.Query().Get("v")

	if videoID == "" {
		return "", errors.New("invalid youtube url")
	}

	return videoID, nil
}

func GetYoutubeVideo(videoURL string) (models.YouTubeResponse, error) {

	videoID, err := ExtractVideoID(videoURL)
	if err != nil {
		return models.YouTubeResponse{}, err
	}

	apiURL := fmt.Sprintf(
		"https://www.googleapis.com/youtube/v3/videos?id=%s&part=snippet,contentDetails&key=%s",
		videoID,
		config.MustLoad().YouTubeAPIKey,
	)

	response, err := http.Get(apiURL)
	if err != nil {
		return models.YouTubeResponse{}, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return models.YouTubeResponse{}, errors.New("youtube api request failed")
	}

	var youtubeResponse models.YouTubeResponse

	err = json.NewDecoder(response.Body).Decode(
		&youtubeResponse,
	)
	if err != nil {
		return models.YouTubeResponse{}, err
	}

	if len(youtubeResponse.Items) == 0 {
		return models.YouTubeResponse{}, errors.New("video not found")
	}

	return youtubeResponse, nil
}
