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

	"github.com/google/uuid"

	"github.com/cloudstorex/backend/internal/config"
	"github.com/cloudstorex/backend/internal/database"
	"github.com/cloudstorex/backend/internal/identity"
	"github.com/cloudstorex/backend/internal/middleware"
	"github.com/cloudstorex/backend/internal/metadata/events"
	"github.com/cloudstorex/backend/internal/metadata/repository"
	"github.com/cloudstorex/backend/internal/metadata/service"
	"github.com/cloudstorex/backend/internal/policy/engine"
	policyEvents "github.com/cloudstorex/backend/internal/policy/events"
	policyHandlerPkg "github.com/cloudstorex/backend/internal/policy/handler"
	policyRepo "github.com/cloudstorex/backend/internal/policy/repository"
	"github.com/cloudstorex/backend/internal/policy/rules"
	policySvc "github.com/cloudstorex/backend/internal/policy/service"
	"github.com/cloudstorex/backend/internal/provider"
	"github.com/cloudstorex/backend/internal/provider/aws"
	"github.com/cloudstorex/backend/internal/provider/dto"
	providerRepo "github.com/cloudstorex/backend/internal/provider/repository"
	providerSvc "github.com/cloudstorex/backend/internal/provider/service"
	"github.com/cloudstorex/backend/internal/provider/minio"
	"github.com/cloudstorex/backend/internal/shared/logger"
	"github.com/cloudstorex/backend/internal/shared/response"
	"github.com/cloudstorex/backend/internal/storage"
	"github.com/gin-gonic/gin"
	redis "github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type App struct {
	Config          *config.Config
	DB              *gorm.DB
	RedisClient     *redis.Client
	StorageRegistry provider.Registry
	ProviderService providerSvc.ProviderService
	StorageRouter   storage.Router
	StorageService  storage.Service
	PolicyService   policySvc.PolicyService
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

	// Initialize Provider Domain
	providerRepository := providerRepo.NewPostgresRepository(db)
	providerValidator := providerSvc.NewProviderValidator()
	providerService := providerSvc.NewProviderService(providerRepository, providerValidator)
	
	storageRegistry := provider.NewRegistry()

	// Sync MinIO from config to DB (if not exists)
	ctx := context.Background()
	defaultWorkspaceID := uuid.MustParse("00000000-0000-0000-0000-000000000000") // TODO: use real workspace ID logic later
	_, err = providerService.GetDefaultProvider(ctx, defaultWorkspaceID.String())
	if err != nil && errors.Is(err, storage.ErrProviderNotFound) {
		providerService.CreateProvider(ctx, dto.CreateProviderRequest{
			WorkspaceID:  defaultWorkspaceID,
			ProviderName: "minio",
			ProviderType: "MINIO",
			Endpoint:     cfg.MinioEndpoint,
			BucketPrefix: "cloudstorex-",
			IsDefault:    true,
		})
	}
	
	minioCfg := &minio.Config{
		Endpoint:      cfg.MinioEndpoint,
		AccessKey:     cfg.MinioAccessKey,
		SecretKey:     cfg.MinioSecretKey,
		DefaultBucket: cfg.MinioBucket,
		UseSSL:        cfg.MinioUseSSL,
	}
	
	// Register MinIO to runtime registry
	minioProv, err := minio.NewProvider(ctx, minioCfg, logger.Log)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize MinIO provider: %w", err)
	}
	if err := storageRegistry.Register("minio", minioProv); err != nil {
		return nil, fmt.Errorf("failed to register MinIO provider: %w", err)
	}

	// Setup AWS Provider if configured
	if cfg.AwsAccessKeyID != "" {
		awsCfg := &aws.Config{
			Region:          cfg.AwsRegion,
			Endpoint:        cfg.AwsEndpoint,
			AccessKeyID:     cfg.AwsAccessKeyID,
			SecretAccessKey: cfg.AwsSecretAccessKey,
			BucketPrefix:    cfg.AwsBucketPrefix,
		}
		awsProv, err := aws.NewProvider(ctx, awsCfg, logger.Log)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize AWS provider: %w", err)
		}
		if err := storageRegistry.Register("aws", awsProv); err != nil {
			return nil, fmt.Errorf("failed to register AWS provider: %w", err)
		}

		// Sync AWS to DB
		_, err = providerService.GetDefaultProvider(ctx, defaultWorkspaceID.String())
		// If AWS doesn't exist, we should probably fetch by name to avoid recreating.
		// Since we changed GetByName to require workspace_id, we'll just check if we have an AWS provider.
		awsProviders, _ := providerService.ListProviders(ctx, defaultWorkspaceID.String())
		hasAws := false
		for _, p := range awsProviders {
			if p.ProviderName == "aws" {
				hasAws = true
				break
			}
		}
		if !hasAws {
			providerService.CreateProvider(ctx, dto.CreateProviderRequest{
				WorkspaceID:  defaultWorkspaceID,
				ProviderName: "aws",
				ProviderType: "AWS_S3",
				Endpoint:     cfg.AwsEndpoint,
				Region:       cfg.AwsRegion,
				BucketPrefix: cfg.AwsBucketPrefix,
				IsDefault:    false, // Default stays minio initially
			})
		}
	}

	policyRepository := policyRepo.NewPostgresRepository(db)
	policyService := policySvc.NewPolicyService(policyRepository)
	policyRuleRegistry := rules.NewRegistry()
	policyRuleRegistry.Register(&rules.DefaultRule{})
	policyRuleRegistry.Register(&rules.RegionRule{})
	policyRuleRegistry.Register(&rules.ObjectSizeRule{})
	policyEngineEvaluator := engine.NewEvaluator(policyRuleRegistry)
	policyEventPub := policyEvents.NewLogPublisher(logger.Log)
	policyEngine := engine.NewPolicyEngine(policyRepository, policyEngineEvaluator, providerService, policyEventPub)

	storageRouter := storage.NewRouter(storageRegistry, policyEngine)

	metadataRepo := repository.NewPostgresMetadataRepository(db)
	metadataEvents := events.NewLogPublisher(logger.Log)
	metadataService := service.NewMetadataService(metadataRepo, metadataEvents)

	storageService := storage.NewService(storageRouter, metadataService)

	app := &App{
		Config:          cfg,
		DB:              db,
		RedisClient:     redisClient,
		StorageRegistry: storageRegistry,
		ProviderService: providerService,
		StorageRouter:   storageRouter,
		StorageService:  storageService,
		PolicyService:   policyService,
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

	// Provider routes
	providerHandler := provider.NewHandler(a.ProviderService)
	v1.GET("/providers", providerHandler.ListProviders)
	v1.POST("/providers", providerHandler.CreateProvider)
	v1.GET("/providers/:id", providerHandler.GetProvider)
	v1.PUT("/providers/:id", providerHandler.UpdateProvider)
	v1.DELETE("/providers/:id", providerHandler.DeleteProvider)
	v1.POST("/providers/:id/validate", providerHandler.ValidateProvider)
	v1.POST("/providers/:id/default", providerHandler.SetDefaultProvider)
	v1.POST("/providers/:id/enable", providerHandler.EnableProvider)
	v1.POST("/providers/:id/disable", providerHandler.DisableProvider)

	// Policy routes
	policyHandler := policyHandlerPkg.NewHandler(a.PolicyService)
	v1.GET("/policies", policyHandler.ListPolicies)
	v1.POST("/policies", policyHandler.CreatePolicy)
	v1.GET("/policies/:id", policyHandler.GetPolicy)
	v1.PUT("/policies/:id", policyHandler.UpdatePolicy)
	v1.DELETE("/policies/:id", policyHandler.DeletePolicy)
	v1.POST("/policies/:id/enable", policyHandler.EnablePolicy)
	v1.POST("/policies/:id/disable", policyHandler.DisablePolicy)
	v1.GET("/routing-decisions", policyHandler.ListRoutingDecisions)

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

		// Global Object APIs
		storageGroup.GET("/search", storageHandler.SearchObjects)
		storageGroup.GET("/objects/:id", storageHandler.GetObjectByID)
		storageGroup.GET("/objects/:id/metadata", storageHandler.GetObjectMetadata)
		storageGroup.POST("/objects/:id/tags", storageHandler.TagObject)
		storageGroup.DELETE("/objects/:id/tags/:key", storageHandler.UntagObject)
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
