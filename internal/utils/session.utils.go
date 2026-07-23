package utils

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/zeroaxiis/ZeroAxiis-Services/internal/database"
)

const (
	SessionPrefix = "session:"
	SessionTTL    = 15 * time.Minute
)

func CreateSession(adminID string) (string, error) {
	sessionID := uuid.NewString()

	key := SessionPrefix + sessionID

	err := database.RedisClient.Set(
		context.Background(),
		key,
		adminID,
		SessionTTL,
	).Err()

	if err != nil {
		return "", err
	}

	return sessionID, nil
}

func DeleteSession(sessionID string) error {
	key := SessionPrefix + sessionID

	err := database.RedisClient.Del(
		context.Background(),
		key,
	).Err()

	if err != nil {
		return err
	}

	return nil
}

func RefreshSession(sessionID string) error {

	key := SessionPrefix + sessionID

	err := database.RedisClient.Expire(
		context.Background(),
		key,
		SessionTTL,
	).Err()

	if err != nil {
		return err
	}

	return nil
}

func GetSession(sessionID string) (string, error) {

	key := SessionPrefix + sessionID

	adminID, err := database.RedisClient.Get(
		context.Background(),
		key,
	).Result()

	if err != nil {
		return "", err
	}

	return adminID, nil
}
