package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cloudstorex/backend/internal/config"
	"github.com/cloudstorex/backend/internal/database"
	"github.com/cloudstorex/backend/internal/identity"
	"github.com/cloudstorex/backend/internal/middleware"
	"github.com/cloudstorex/backend/internal/policy"
	"github.com/cloudstorex/backend/internal/provider"
	"github.com/cloudstorex/backend/internal/provider/minio"
	"github.com/cloudstorex/backend/internal/shared/logger"
	"github.com/cloudstorex/backend/internal/shared/response"
	"github.com/cloudstorex/backend/internal/storage"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type App struct {
	Config          *config.Config
	DB              *gorm.DB
	RedisClient     *redis.Client
	StorageRegistry provider.Registry
	StorageRouter   storage.Router
	StorageService  storage.Service
	Router          *gin.Engine
}

func NewApp() (*App, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	logger.InitLogger(cfg.AppEnv)

	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		return nil, err
	}

	if err := database.RunMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	redisClient, err := database.NewRedisClient(cfg)
	if err != nil {
		return nil, err
	}

	// Initialize Storage Domain & MinIO Provider
	storageRegistry := provider.NewRegistry()
	minioCfg := &minio.Config{
		Endpoint:      cfg.MinioEndpoint,
		AccessKey:     cfg.MinioAccessKey,
		SecretKey:     cfg.MinioSecretKey,
		DefaultBucket: cfg.MinioBucket,
		UseSSL:        cfg.MinioUseSSL,
	}
	minioProv, err := minio.NewProvider(context.Background(), minioCfg, logger.Log)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize MinIO provider: %w", err)
	}
	if err := storageRegistry.Register("minio", minioProv); err != nil {
		return nil, fmt.Errorf("failed to register MinIO provider: %w", err)
	}

	policyEvaluator := policy.NewDefaultEvaluator("minio")
	storageRouter := storage.NewRouter(storageRegistry, policyEvaluator)
	storageService := storage.NewService(storageRouter)

	app := &App{
		Config:          cfg,
		DB:              db,
		RedisClient:     redisClient,
		StorageRegistry: storageRegistry,
		StorageRouter:   storageRouter,
		StorageService:  storageService,
		Router:          gin.New(), // Create without default middlewares
	}

	app.setupMiddlewares()
	app.setupRoutes()

	return app, nil
}

func (a *App) setupMiddlewares() {
	a.Router.Use(middleware.Recovery())
	a.Router.Use(middleware.RequestLogger())
	a.Router.Use(middleware.SecurityHeaders())
	a.Router.Use(middleware.CORS(a.Config))
}

func (a *App) setupRoutes() {
	v1 := a.Router.Group("/api/v1")
	
	// Health check
	v1.GET("/health", func(c *gin.Context) {
		dbStatus := "down"
		sqlDB, err := a.DB.DB()
		if err == nil && sqlDB.Ping() == nil {
			dbStatus = "up"
		}
		
		redisStatus := "down"
		if a.RedisClient.Ping(c.Request.Context()).Err() == nil {
			redisStatus = "up"
		}

		minioStatus := "down"
		if exists, err := a.StorageService.ObjectExists(c.Request.Context(), a.Config.MinioBucket, "health-check-dummy-key"); err == nil || errors.Is(err, storage.ErrObjectNotFound) || !exists {
			minioStatus = "up"
		}

		response.Success(c, http.StatusOK, gin.H{
			"status":   "up",
			"postgres": dbStatus,
			"redis":    redisStatus,
			"minio":    minioStatus,
			"version":  "1.0.0",
		})
	})

	// Identity dependencies
	tokenService := identity.NewTokenService(a.Config.JWTSecret)
	identityRepo := identity.NewRepository(a.DB)
	identityService := identity.NewService(identityRepo, tokenService)
	identityHandler := identity.NewHandler(identityService)

	// Auth routes
	auth := v1.Group("/auth")
	{
		auth.POST("/register", identityHandler.Register)
		auth.POST("/login", identityHandler.Login)
		
		// Protected test route
		protected := auth.Group("/me")
		protected.Use(middleware.AuthMiddleware(tokenService))
		protected.GET("", func(c *gin.Context) {
			userID, _ := c.Get("user_id")
			response.Success(c, http.StatusOK, gin.H{"user_id": userID})
		})
	}

	// Storage routes (protected by auth middleware)
	storageHandler := storage.NewHandler(a.StorageService, a.Config.MaxUploadSizeMB)
	storageGroup := v1.Group("/storage")
	storageGroup.Use(middleware.AuthMiddleware(tokenService))
	{
		// Bucket APIs
		storageGroup.POST("/buckets", storageHandler.CreateBucket)
		storageGroup.GET("/buckets", storageHandler.ListBuckets)
		storageGroup.DELETE("/buckets/:bucket", storageHandler.DeleteBucket)

		// Object APIs
		storageGroup.POST("/buckets/:bucket/objects", storageHandler.UploadObject)
		storageGroup.GET("/buckets/:bucket/objects", storageHandler.ListObjects)
		storageGroup.GET("/buckets/:bucket/objects/*key", storageHandler.DownloadObject)
		storageGroup.DELETE("/buckets/:bucket/objects/*key", storageHandler.DeleteObject)
		storageGroup.HEAD("/buckets/:bucket/objects/*key", storageHandler.ObjectExists)
	}

	// OpenAPI & Swagger UI routes (public)
	v1.GET("/openapi.json", func(c *gin.Context) {
		c.File("api/docs/swagger.json")
	})
	v1.GET("/docs", func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusOK, swaggerUIHTML)
	})
}

const swaggerUIHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <title>CloudStoreX API Docs</title>
  <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5/swagger-ui.css" />
</head>
<body>
<div id="swagger-ui"></div>
<script src="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5/swagger-ui-bundle.js" crossorigin></script>
<script>
  window.onload = () => {
    window.ui = SwaggerUIBundle({
      url: '/api/v1/openapi.json',
      dom_id: '#swagger-ui',
    });
  };
</script>
</body>
</html>`

func (a *App) Run() error {
	srv := &http.Server{
		Addr:    ":" + a.Config.Port,
		Handler: a.Router,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		slog.Info("Server listening", slog.String("port", a.Config.Port))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("listen", slog.String("error", err.Error()))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", slog.String("error", err.Error()))
		return err
	}

	slog.Info("Server exiting")
	return nil
}
