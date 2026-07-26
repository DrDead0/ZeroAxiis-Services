package utils

import (
	"context"
	"encoding/json"
	"time"

	"github.com/zeroaxiis/ZeroAxiis-Services/internal/database"
)

func SetCache(key string, data interface{}, expiration time.Duration) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = database.RedisClient.Set(
		context.Background(),
		key,
		jsonData,
		expiration,
	).Err()

	if err != nil {
		return err
	}
	return nil
}

func GetCache(key string) (string, error) {
	data, err := database.RedisClient.Get(
		context.Background(),
		key,
	).Result()
	if err != nil {
		return "", err
	}
	return data, nil
}
func DeleteCache(key string) error {
	err := database.RedisClient.Del(
		context.Background(),
		key,
	).Err()

	if err != nil {
		return err
	}
	return nil
}
