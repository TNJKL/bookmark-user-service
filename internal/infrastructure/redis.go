package infrastructure

import (
	redisPkg "github.com/TNJKL/bookmark-pkg/pkg/redis"
	"github.com/TNJKL/bookmark-pkg/pkg/utils"
	"github.com/redis/go-redis/v9"
)

// CreateRedisConn creates a new redis connection
func CreateRedisConn() *redis.Client {
	// Create redis db connection
	redisClient, err := redisPkg.NewClient("")
	utils.NoErr(err)

	return redisClient
}
