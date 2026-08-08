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

	"github.com/cloudstorex/backend/internal/analytics"
	"github.com/cloudstorex/backend/internal/audit"
	"github.com/cloudstorex/backend/internal/config"
	"github.com/cloudstorex/backend/internal/database"
	"github.com/cloudstorex/backend/internal/identity"

	"github.com/cloudstorex/backend/internal/ai"
	"github.com/cloudstorex/backend/internal/ai/prompts"
	"github.com/cloudstorex/backend/internal/ai/providers/anthropic"
	"github.com/cloudstorex/backend/internal/ai/providers/gemini"
	"github.com/cloudstorex/backend/internal/ai/providers/openai"

	"github.com/cloudstorex/backend/internal/security/abac"
	authEnginePkg "github.com/cloudstorex/backend/internal/security/engine"
	"github.com/cloudstorex/backend/internal/security/rbac"
	"github.com/cloudstorex/backend/internal/security/zerotrust"
	"github.com/cloudstorex/backend/internal/jobs"
	"github.com/cloudstorex/backend/internal/lifecycle"
	"github.com/cloudstorex/backend/internal/middleware"
	"github.com/cloudstorex/backend/internal/metadata/cache"
	"github.com/cloudstorex/backend/internal/metadata/events"
	metadataHandlerPkg "github.com/cloudstorex/backend/internal/metadata/handler"
	"github.com/cloudstorex/backend/internal/metadata/repository"
	"github.com/cloudstorex/backend/internal/metadata/service"
	"github.com/cloudstorex/backend/internal/policy/engine"
	policyEvents "github.com/cloudstorex/backend/internal/policy/events"
	policyHandlerPkg "github.com/cloudstorex/backend/internal/policy/handler"
	policyRepo "github.com/cloudstorex/backend/internal/policy/repository"
	policySvc "github.com/cloudstorex/backend/internal/policy/service"
	"github.com/cloudstorex/backend/internal/provider"
	"github.com/cloudstorex/backend/internal/provider/aws"
	"github.com/cloudstorex/backend/internal/provider/dto"
	providerRepo "github.com/cloudstorex/backend/internal/provider/repository"
	providerSvc "github.com/cloudstorex/backend/internal/provider/service"
	"github.com/cloudstorex/backend/internal/provider/minio"
	"github.com/cloudstorex/backend/internal/quota"
	"github.com/cloudstorex/backend/internal/shared/logger"
	"github.com/cloudstorex/backend/internal/observability/metrics"
	"github.com/cloudstorex/backend/internal/observability/tracing"
	"github.com/cloudstorex/backend/internal/shared/response"
	"github.com/cloudstorex/backend/internal/storage"
	"github.com/cloudstorex/backend/internal/workers"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
	StorageService   storage.Service
	MetadataService  service.MetadataService
	PolicyService    policySvc.PolicyService
	PolicyEngine     engine.PolicyEngine
	QuotaService     quota.Service
	LifecycleService lifecycle.Service
	AuditService     audit.Service
	AnalyticsService analytics.Service
	JobClient        jobs.Client
	JobDispatcher   jobs.Dispatcher
	AIManager       ai.AIManager
	PromptManager   prompts.Manager
	Router          *gin.Engine
}

func NewApp() (*App, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	logger.InitLogger(cfg.AppEnv)

	_, err = tracing.InitTracer("cloudstorex-backend", cfg.AppEnv)
	if err != nil {
		logger.Log.Warn("Failed to initialize OpenTelemetry tracer", slog.String("error", err.Error()))
	}
	_ = metrics.StartRuntimeMetricCollector(15 * time.Second)

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
	// Default workspace ID used for system bootstrap
	defaultWorkspaceID := uuid.MustParse("00000000-0000-0000-0000-000000000000")
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
		_, _ = providerService.GetDefaultProvider(ctx, defaultWorkspaceID.String())
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
	
	policyEngineEvaluator := engine.NewEvaluator()
	policyEventPub := policyEvents.NewLogPublisher(logger.Log, policyRepository)
	
	policyProviderValidator := storage.NewProviderValidator(storageRegistry)
	policyEngine := engine.NewPolicyEngine(policyRepository, policyEngineEvaluator, providerService, policyProviderValidator, policyEventPub)

	storageRouter := storage.NewRouter(storageRegistry, policyEngine)

	jobRegistry := jobs.NewRegistry()
	jobClient := jobs.NewClient(db, redisClient)
	jobDispatcher := jobs.NewDispatcher(db, redisClient, jobRegistry, 5)

	metadataCache := cache.NewRedisMetadataCache(redisClient)
	metadataRepo := repository.NewPostgresMetadataRepository(db)
	metadataEvents := events.NewJobEventPublisher(jobClient, logger.Log)
	metadataService := service.NewMetadataService(metadataRepo, metadataEvents, metadataCache)

	quotaRepo := quota.NewPostgresRepository(db)
	quotaService := quota.NewService(quotaRepo, logger.Log)

	storageService := storage.NewService(storageRouter, metadataService, quotaService)

	lifecycleRepo := lifecycle.NewPostgresRepository(db)
	lifecycleService := lifecycle.NewService(lifecycleRepo, logger.Log)

	auditRepo := audit.NewPostgresRepository(db)
	auditService := audit.NewService(auditRepo, logger.Log)

	analyticsRepo := analytics.NewPostgresRepository(db)
	analyticsService := analytics.NewService(analyticsRepo, logger.Log)

	// Initialize AI Provider Abstraction
	aiManager := ai.NewManager(ai.ProviderType(cfg.AIDefaultProvider), logger.Log)
	
	if cfg.OpenAIApiKey != "" {
		aiManager.RegisterProvider(openai.NewProvider(cfg.OpenAIApiKey))
	}
	if cfg.AnthropicApiKey != "" {
		aiManager.RegisterProvider(anthropic.NewProvider(cfg.AnthropicApiKey))
	}
	if cfg.GeminiApiKey != "" {
		geminiProv, err := gemini.NewProvider(context.Background(), cfg.GeminiApiKey)
		if err == nil {
			aiManager.RegisterProvider(geminiProv)
		} else {
			logger.Log.Warn("Failed to initialize Gemini provider", slog.String("error", err.Error()))
		}
	}

	promptManager, err := prompts.NewManager()
	if err != nil {
		logger.Log.Warn("Failed to initialize prompt manager", slog.String("error", err.Error()))
	}

	// Register Jobs
	jobRegistry.Register("ProcessEvent", workers.NewEventIntegrationWorker(quotaService, logger.Log))
	jobRegistry.Register("SoftDeleteCleanup", workers.NewSoftDeleteCleanupWorker(storageRouter, metadataService, logger.Log))
	jobRegistry.Register("VersionCleanup", workers.NewVersionCleanupWorker(logger.Log))
	jobRegistry.Register("Lifecycle", workers.NewLifecycleWorker(lifecycleRepo, metadataService, storageService, logger.Log))
	jobRegistry.Register("AuditLog", workers.NewAuditWorker(auditService, logger.Log))
	jobRegistry.Register("AnalyticsAggregation", workers.NewAnalyticsWorker(analyticsService, quotaService, logger.Log))

	app := &App{
		Config:          cfg,
		DB:              db,
		RedisClient:     redisClient,
		StorageRegistry: storageRegistry,
		ProviderService: providerService,
		StorageRouter:    storageRouter,
		StorageService:   storageService,
		MetadataService:  metadataService,
		PolicyService:    policyService,
		PolicyEngine:     policyEngine,
		QuotaService:     quotaService,
		LifecycleService: lifecycleService,
		AuditService:     auditService,
		AnalyticsService: analyticsService,
		JobClient:        jobClient,
		JobDispatcher:   jobDispatcher,
		AIManager:       aiManager,
		PromptManager:   promptManager,
		Router:          gin.New(), // Create without default middlewares
	}

	app.setupMiddlewares()
	app.setupRoutes()

	return app, nil
}

func (a *App) setupMiddlewares() {
	a.Router.Use(middleware.Recovery())
	a.Router.Use(middleware.OpenTelemetryMiddleware("cloudstorex-backend"))
	a.Router.Use(middleware.PrometheusMiddleware())
	a.Router.Use(middleware.RequestLogger())
	a.Router.Use(middleware.SecurityHeaders())
	a.Router.Use(middleware.CORS(a.Config))
}

func (a *App) setupRoutes() {
	livenessHandler := func(c *gin.Context) {
		response.Success(c, http.StatusOK, gin.H{
			"status": "up",
		})
	}

	readinessHandler := func(c *gin.Context) {
		if a.DB == nil {
			response.Error(c, http.StatusServiceUnavailable, "SERVICE_UNAVAILABLE", "database unreachable")
			return
		}
		sqlDB, err := a.DB.DB()
		if err != nil || sqlDB.Ping() != nil {
			response.Error(c, http.StatusServiceUnavailable, "SERVICE_UNAVAILABLE", "database unreachable")
			return
		}
		if a.RedisClient == nil || a.RedisClient.Ping(c.Request.Context()).Err() != nil {
			response.Error(c, http.StatusServiceUnavailable, "SERVICE_UNAVAILABLE", "redis unreachable")
			return
		}
		if a.StorageRegistry == nil || len(a.StorageRegistry.List()) == 0 {
			response.Error(c, http.StatusServiceUnavailable, "SERVICE_UNAVAILABLE", "provider registry uninitialized or empty")
			return
		}
		if a.PolicyService == nil {
			response.Error(c, http.StatusServiceUnavailable, "SERVICE_UNAVAILABLE", "policy engine uninitialized")
			return
		}
		if a.JobDispatcher == nil || a.JobClient == nil {
			response.Error(c, http.StatusServiceUnavailable, "SERVICE_UNAVAILABLE", "background worker queue disconnected")
			return
		}

		response.Success(c, http.StatusOK, gin.H{
			"status":        "ready",
			"postgres":      "up",
			"redis":         "up",
			"providers":     len(a.StorageRegistry.List()),
			"policy_engine": "up",
			"workers":       "up",
		})
	}

	a.Router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	enablePprof := false
	pprofToken := ""
	if a.Config != nil {
		enablePprof = a.Config.EnablePprof
		pprofToken = a.Config.PprofToken
	}
	RegisterPprof(a.Router, enablePprof, pprofToken)

	a.Router.GET("/healthz", livenessHandler)
	a.Router.GET("/livez", livenessHandler)
	a.Router.GET("/readyz", readinessHandler)

	v1 := a.Router.Group("/api/v1")

	v1.GET("/live", livenessHandler)
	v1.GET("/ready", readinessHandler)

	// Health check
	v1.GET("/health", func(c *gin.Context) {
		dbStatus := "down"
		if a.DB != nil {
			if sqlDB, err := a.DB.DB(); err == nil && sqlDB.Ping() == nil {
				dbStatus = "up"
			}
		}
		
		redisStatus := "down"
		if a.RedisClient != nil && a.RedisClient.Ping(c.Request.Context()).Err() == nil {
			redisStatus = "up"
		}

		minioBucket := "cloudstorex-default"
		if a.Config != nil && a.Config.MinioBucket != "" {
			minioBucket = a.Config.MinioBucket
		}

		minioStatus := "down"
		if a.StorageService != nil {
			if exists, err := a.StorageService.ObjectExists(c.Request.Context(), minioBucket, "health-check-dummy-key"); err == nil || errors.Is(err, storage.ErrObjectNotFound) || !exists {
				minioStatus = "up"
			}
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
	policyHandler := policyHandlerPkg.NewHandler(a.PolicyService, a.PolicyEngine)
	v1.GET("/policies", policyHandler.ListPolicies)
	v1.POST("/policies", policyHandler.CreatePolicy)
	v1.POST("/policies/evaluate", policyHandler.Evaluate)
	v1.GET("/policies/:id", policyHandler.GetPolicy)
	v1.PUT("/policies/:id", policyHandler.UpdatePolicy)
	v1.DELETE("/policies/:id", policyHandler.DeletePolicy)
	v1.POST("/policies/:id/enable", policyHandler.EnablePolicy)
	v1.POST("/policies/:id/disable", policyHandler.DisablePolicy)
	v1.GET("/routing-decisions", policyHandler.ListRoutingDecisions)

	// Identity dependencies
	jwtSecret := "default_jwt_secret"
	maxUploadSizeMB := int64(100)
	if a.Config != nil {
		if a.Config.JWTSecret != "" {
			jwtSecret = a.Config.JWTSecret
		}
		if a.Config.MaxUploadSizeMB > 0 {
			maxUploadSizeMB = a.Config.MaxUploadSizeMB
		}
	}
	tokenService := identity.NewTokenService(jwtSecret)
	
	authEngine := authEnginePkg.NewAuthorizationEngine(rbac.NewEvaluator(), abac.NewEvaluator(), zerotrust.NewEvaluator())

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
		protected.Use(middleware.AuthMiddleware(tokenService, authEngine))
		protected.GET("", func(c *gin.Context) {
			userID, _ := c.Get("user_id")
			response.Success(c, http.StatusOK, gin.H{"user_id": userID})
		})
	}

	// Storage routes (protected by auth middleware)
	storageHandler := storage.NewHandler(a.StorageService, maxUploadSizeMB)
	storageGroup := v1.Group("/storage")
	storageGroup.Use(middleware.AuthMiddleware(tokenService, authEngine))
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
		storageGroup.POST("/buckets/:bucket/presigned-upload/*key", storageHandler.GeneratePresignedUploadURL)
		storageGroup.GET("/buckets/:bucket/presigned-download/*key", storageHandler.GeneratePresignedDownloadURL)

		// Global Object APIs
		storageGroup.GET("/search", storageHandler.SearchObjects)
		storageGroup.GET("/objects/:id", storageHandler.GetObjectByID)
		storageGroup.GET("/objects/:id/metadata", storageHandler.GetObjectMetadata)
		storageGroup.POST("/objects/:id/tags", storageHandler.TagObject)
		storageGroup.DELETE("/objects/:id/tags/:key", storageHandler.UntagObject)
		storageGroup.GET("/objects/:id/versions", storageHandler.ListObjectVersions)
		storageGroup.POST("/objects/:id/restore", storageHandler.RestoreObject)
		storageGroup.POST("/objects/:id/legal-hold", storageHandler.SetLegalHold)
		storageGroup.POST("/objects/:id/retention", storageHandler.SetRetention)
	}

	// Metadata APIs
	metadataHandler := metadataHandlerPkg.NewMetadataHandler(a.MetadataService)
	metadataHandler.RegisterRoutes(v1)

	// Quota APIs
	quotaHandler := quota.NewHandler(a.QuotaService)
	quotaGroup := v1.Group("/quotas")
	quotaGroup.Use(middleware.AuthMiddleware(tokenService, authEngine))
	{
		quotaGroup.POST("/workspaces/:workspace_id", quotaHandler.SetWorkspaceQuota)
		quotaGroup.GET("/workspaces/:workspace_id", quotaHandler.GetWorkspaceQuota)
	}

	// Lifecycle APIs
	lifecycleHandler := lifecycle.NewHandler(a.LifecycleService)
	lifecycleGroup := v1.Group("/lifecycle")
	lifecycleGroup.Use(middleware.AuthMiddleware(tokenService, authEngine))
	{
		lifecycleGroup.POST("/buckets/:bucket_id/rules", lifecycleHandler.CreateRule)
		lifecycleGroup.GET("/buckets/:bucket_id/rules", lifecycleHandler.ListRules)
		lifecycleGroup.DELETE("/rules/:id", lifecycleHandler.DeleteRule)
	}

	// Audit APIs
	auditHandler := audit.NewHandler(a.AuditService)
	auditGroup := v1.Group("/audit-logs")
	auditGroup.Use(middleware.AuthMiddleware(tokenService, authEngine))
	{
		auditGroup.GET("/workspaces/:workspace_id", auditHandler.ListLogs)
	}

	// Analytics APIs
	analyticsHandler := analytics.NewHandler(a.AnalyticsService)
	analyticsGroup := v1.Group("/analytics")
	analyticsGroup.Use(middleware.AuthMiddleware(tokenService, authEngine))
	{
		analyticsGroup.GET("/workspaces/:workspace_id/history", analyticsHandler.GetAnalyticsHistory)
		analyticsGroup.GET("/workspaces/:workspace_id/latest", analyticsHandler.GetLatestMetrics)
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

	// Start Job Dispatcher
	ctx, cancel := context.WithCancel(context.Background())
	a.JobDispatcher.Start(ctx)

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("Shutting down server...")

	cancel() // Cancel context to stop fetching jobs
	a.JobDispatcher.Stop()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("Server forced to shutdown", slog.String("error", err.Error()))
		return err
	}

	slog.Info("Server exiting")
	return nil
}
