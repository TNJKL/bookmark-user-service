package api

import (
	"fmt"
	"net/http"

	"github.com/TNJKL/bookmark-pkg/pkg/jwtutils"
	"github.com/TNJKL/bookmark-pkg/pkg/utils"
	"github.com/TNJKL/bookmark-user-service/docs"
	_ "github.com/TNJKL/bookmark-user-service/docs" // Load tài liệu Swagger đã generate

	"github.com/TNJKL/bookmark-user-service/internal/app/handler/healthcheck"

	"github.com/TNJKL/bookmark-pkg/middlewares"
	"github.com/TNJKL/bookmark-pkg/ratelimit"
	userHandler "github.com/TNJKL/bookmark-user-service/internal/app/handler/user"
	"github.com/TNJKL/bookmark-user-service/internal/app/repository/ping"
	"github.com/TNJKL/bookmark-user-service/internal/app/repository/user"
	healthcheck2 "github.com/TNJKL/bookmark-user-service/internal/app/service/healthcheck"
	userSvc "github.com/TNJKL/bookmark-user-service/internal/app/service/user"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

// Interface định nghĩa "Engine có thể làm gì"
// Engine defines the contract for starting the server and handling HTTP requests.
type Engine interface {
	Start() error
	ServerHTTP(w http.ResponseWriter, req *http.Request)
}

// Struct thực tế implement interface
type engine struct {
	app         *gin.Engine
	cfg         *Config
	redisClient *redis.Client
	db          *gorm.DB
	jwtGen      jwtutils.JWTGenerator
	jwtVal      jwtutils.JWTValidator
}

// EngineOpts holds the configuration and dependencies required to initialize the API Engine
type EngineOpts struct {
	App         *gin.Engine
	Cfg         *Config
	RedisClient *redis.Client
	Db          *gorm.DB
	JWTGen      jwtutils.JWTGenerator
	JWTVal      jwtutils.JWTValidator
}

// NewEngine creates and configures a new HTTP API Engine.
func NewEngine(opts *EngineOpts) Engine {

	app := &engine{
		app:         opts.App, // Tạo Gin router
		cfg:         opts.Cfg,
		redisClient: opts.RedisClient,
		db:          opts.Db,
		jwtGen:      opts.JWTGen,
		jwtVal:      opts.JWTVal,
	}
	app.initRoutes() // Đăng ký các routes

	return app
}

// Start runs the HTTP server on the configured port.
func (e *engine) Start() error {
	return e.app.Run(fmt.Sprintf(":%s", e.cfg.Apport))
}

// Server HTTP to test the API endpoint
// ServerHTTP handles HTTP requests directly, primarily used for testing endpoints.
func (e *engine) ServerHTTP(w http.ResponseWriter, req *http.Request) {
	e.app.ServeHTTP(w, req)
}

// handlers aggregates all HTTP handler dependencies used to register
// the application's routes.
type handlers struct {
	healthCheckHandler healthcheck.HealthCheck
	userHandler        userHandler.Handler
}

// initHandlers initializes the api handlers
func (e *engine) initHandlers() *handlers {
	pingRepo := ping.NewHealthRepository(e.redisClient)
	healthCheckSvc := healthcheck2.NewHealthCheck(e.cfg.ServiceName, e.cfg.InstanceID, pingRepo)

	//init user handler
	userRepo := user.NewSQLRepository(e.db)
	hasher := utils.NewHasher()
	userService := userSvc.NewService(userRepo, hasher, e.jwtGen)

	return &handlers{
		healthCheckHandler: healthcheck.NewHealthCheck(healthCheckSvc),
		userHandler:        userHandler.NewHandler(userService),
	}
}

// middlewares represents the struct for containing all the necessary middlewares for API
type middlewares struct {
	jwtAuth   middleware.JWTAuth
	ratelimit middleware.RateLimit
}

// initMiddlewares initializes the api middlewares
func (e *engine) initMiddlewares() middlewares {
	rateLimitRepo := ratelimit.NewRedisRepo(e.redisClient)
	return middlewares{
		jwtAuth:   middleware.NewJWTAuth(e.jwtVal),
		ratelimit: middleware.NewRateLimit(rateLimitRepo),
	}
}

// initRoutes initializes the api routes
func (e *engine) initRoutes() {
	allHandler := e.initHandlers()
	allMiddlewares := e.initMiddlewares()

	//health-check
	e.app.GET("/health-check", allHandler.healthCheckHandler.HealthCheck)

	//Init swagger routes
	docs.SwaggerInfo.BasePath = e.cfg.BasePath
	e.app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//public routes
	v1Routes := e.app.Group("/v1")
	{
		//user
		v1Routes.POST("/users/register", allHandler.userHandler.Register)
		v1Routes.POST("/users/login", allHandler.userHandler.Login)

	}
	//private routes (need Auth)
	privateRoutes := e.app.Group("")
	privateRoutes.Use(allMiddlewares.jwtAuth.JWTAuth())
	privateRoutes.Use(allMiddlewares.ratelimit.RateLimit())
	{
		privateV1Routes := privateRoutes.Group("/v1")
		{
			//self endpoints
			privateV1Routes.GET("self/info", allHandler.userHandler.GetSelfInfo)
			privateV1Routes.PUT("self/info", allHandler.userHandler.UpdateSelfInfo)

		}

	}
}
