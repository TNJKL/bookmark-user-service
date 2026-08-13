package utils

import (
	"testing"
	"time"

	"github.com/TNJKL/bookmark-pkg/pkg/jwtutils"
	"github.com/TNJKL/bookmark-user-service/internal/api"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// SetupTestJWT loads the test RSA keys from local test utils folder
func SetupTestJWT(t *testing.T) (jwtutils.JWTGenerator, jwtutils.JWTValidator) {
	jwtGen, err := jwtutils.NewJWTGenerator("../utils/test.private.key")
	assert.NoError(t, err)
	jwtVal, err := jwtutils.NewJWTValidator("../utils/test.public.key")
	assert.NoError(t, err)
	return jwtGen, jwtVal
}

// generateTestToken creates and signs a valid JWT token for test purposes
func GenerateTestToken(t *testing.T, jwtGen jwtutils.JWTGenerator, sub, email string) string {
	tokenContent := jwt.MapClaims{
		"sub":   sub,
		"email": email,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	}
	token, err := jwtGen.GenerateJWT(tokenContent)
	assert.NoError(t, err)
	return token
}

// buildTestAPI instantiates the Gin engine with mocking dependencies
func BuildTestAPI(db *gorm.DB, redisClient *redis.Client, jwtGen jwtutils.JWTGenerator, jwtVal jwtutils.JWTValidator) api.Engine {
	return api.NewEngine(&api.EngineOpts{
		App:         gin.New(),
		Cfg:         &api.Config{},
		RedisClient: redisClient,
		Db:          db,
		JWTGen:      jwtGen,
		JWTVal:      jwtVal,
	})
}
