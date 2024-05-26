package helper

import (
	"context"
	"encoding/json"
	"gis_map_info/support/redis_support"
	"time"
)

func SaveParameter[T any](ctx context.Context, c T, expired time.Duration) string {
	key := RandStringBytes(10)
	ggString, _ := json.Marshal(c)
	redis_support.RedisClient.Set(ctx, key, ggString, expired)
	return string(key)
}

func GetParameter[T any](ctx context.Context, key string, c *T) error {
	rawRedisData, err := redis_support.RedisClient.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(rawRedisData), &c)
	if err != nil {
		return err
	}
	return nil
}
