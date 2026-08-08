package metrics

import (
	"runtime"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Standard low-cardinality labels (Refinement 2)
// Do not include object names, user IDs, or dynamic filenames.

var (
	// ==========================================
	// 1. HTTP RED Metrics
	// ==========================================
	HTTPRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cloudstorex_http_requests_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"method", "route", "status"},
	)

	HTTPRequestDurationSeconds = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "cloudstorex_http_request_duration_seconds",
			Help:    "HTTP request latency distribution in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "route", "status"},
	)

	HTTPActiveRequests = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cloudstorex_http_active_requests",
			Help: "Number of currently active HTTP requests.",
		},
		[]string{"method", "route"},
	)

	HTTPErrorRate = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cloudstorex_http_errors_total",
			Help: "Total number of HTTP 4xx and 5xx errors.",
		},
		[]string{"method", "route", "status"},
	)

	// ==========================================
	// 2. Storage Operations Metrics
	// ==========================================
	StorageUploadsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cloudstorex_storage_uploads_total",
			Help: "Total number of object upload operations.",
		},
		[]string{"workspace", "provider", "bucket", "status"},
	)

	StorageDownloadsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cloudstorex_storage_downloads_total",
			Help: "Total number of object download operations.",
		},
		[]string{"workspace", "provider", "bucket", "status"},
	)

	StorageUploadBytes = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cloudstorex_storage_upload_bytes_total",
			Help: "Total bytes uploaded across storage providers.",
		},
		[]string{"workspace", "provider", "bucket"},
	)

	StorageDownloadBytes = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cloudstorex_storage_download_bytes_total",
			Help: "Total bytes downloaded across storage providers.",
		},
		[]string{"workspace", "provider", "bucket"},
	)

	StorageDeleteOperations = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cloudstorex_storage_deletes_total",
			Help: "Total delete operations executed.",
		},
		[]string{"workspace", "provider", "bucket", "status"},
	)

	// ==========================================
	// 3. Provider Metrics
	// ==========================================
	ProviderLatencySeconds = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "cloudstorex_provider_latency_seconds",
			Help:    "Storage provider operation latency in seconds.",
			Buckets: []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
		},
		[]string{"provider", "operation", "region", "status"},
	)

	ProviderErrorsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cloudstorex_provider_errors_total",
			Help: "Total provider operation errors.",
		},
		[]string{"provider", "operation", "region"},
	)

	ProviderRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cloudstorex_provider_requests_total",
			Help: "Total provider requests.",
		},
		[]string{"provider", "operation", "region", "status"},
	)

	ProviderHealthStatus = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cloudstorex_provider_health_status",
			Help: "Provider health status gauge (1=healthy, 0=unhealthy/degraded).",
		},
		[]string{"provider", "region"},
	)

	// ==========================================
	// 4. Policy Engine Metrics
	// ==========================================
	PolicyEvaluationsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cloudstorex_policy_evaluations_total",
			Help: "Total policy engine evaluations.",
		},
		[]string{"workspace", "operation", "status"},
	)

	PolicyRoutingDecisions = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cloudstorex_policy_routing_decisions_total",
			Help: "Total routing decisions made by policy engine.",
		},
		[]string{"workspace", "provider", "operation", "rule_type"},
	)

	PolicyRuleMatches = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cloudstorex_policy_rule_matches_total",
			Help: "Total policy rule matches.",
		},
		[]string{"workspace", "rule_type"},
	)

	PolicyFallbackRoutes = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cloudstorex_policy_fallback_routes_total",
			Help: "Total routing fallbacks triggered.",
		},
		[]string{"workspace", "operation"},
	)

	// ==========================================
	// 5. Metadata Layer Metrics
	// ==========================================
	MetadataQueriesTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cloudstorex_metadata_queries_total",
			Help: "Total metadata queries executed.",
		},
		[]string{"workspace", "operation", "status"},
	)

	MetadataSearchLatencySeconds = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "cloudstorex_metadata_search_latency_seconds",
			Help:    "Metadata search latency in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"workspace", "status"},
	)

	MetadataObjectCatalogSize = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cloudstorex_metadata_object_catalog_size",
			Help: "Current object catalog size in items.",
		},
		[]string{"workspace", "bucket"},
	)

	// ==========================================
	// 6. Background Jobs Metrics
	// ==========================================
	JobsRunning = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cloudstorex_jobs_running",
			Help: "Number of currently running background jobs.",
		},
		[]string{"job_type"},
	)

	JobsCompletedTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cloudstorex_jobs_completed_total",
			Help: "Total completed background jobs.",
		},
		[]string{"job_type", "status"},
	)

	JobsFailedTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cloudstorex_jobs_failed_total",
			Help: "Total failed background jobs.",
		},
		[]string{"job_type"},
	)

	JobsRetryCountTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cloudstorex_jobs_retry_count_total",
			Help: "Total job retry attempts.",
		},
		[]string{"job_type"},
	)

	JobsQueueDepth = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cloudstorex_jobs_queue_depth",
			Help: "Current background job queue depth.",
		},
		[]string{"queue"},
	)

	// ==========================================
	// 7. Infrastructure / Runtime Metrics (USE Method)
	// ==========================================
	PostgresConnections = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cloudstorex_postgres_connections",
			Help: "Active PostgreSQL database connections.",
		},
		[]string{"status"},
	)

	RedisConnections = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cloudstorex_redis_connections",
			Help: "Active Redis connections.",
		},
		[]string{"status"},
	)

	Goroutines = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "cloudstorex_goroutines",
			Help: "Current goroutine count.",
		},
	)

	MemoryAllocBytes = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "cloudstorex_memory_alloc_bytes",
			Help: "Current memory allocated bytes.",
		},
	)

	// ==========================================
	// 8. Reliability & Resilience Metrics
	// ==========================================
	ReliabilityReplicationBytesTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "cloudstorex_reliability_replication_bytes_total",
		Help: "Total bytes replicated across providers",
	}, []string{"workspace", "source_provider", "target_provider"})

	// Security Metrics (Epic 14)
	SecurityAuthenticationAttemptsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "cloudstorex_security_authentication_attempts_total",
		Help: "Total number of authentication attempts",
	}, []string{"provider", "status"})

	SecurityLoginFailuresTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "cloudstorex_security_login_failures_total",
		Help: "Total number of failed login attempts",
	}, []string{"provider", "reason"})

	SecurityMFAChallengesTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "cloudstorex_security_mfa_challenges_total",
		Help: "Total number of MFA challenges issued",
	}, []string{"status"})

	SecurityActiveSessionsTotal = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "cloudstorex_security_active_sessions_total",
		Help: "Total number of active user sessions",
	}, []string{"organization"})

	SecurityAPIKeyRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "cloudstorex_security_api_key_requests_total",
		Help: "Total number of requests authenticated via API Key",
	}, []string{"workspace"})

	SecurityAuthorizationDenialsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "cloudstorex_security_authorization_denials_total",
		Help: "Total number of denied authorization attempts",
	}, []string{"reason"})

	SecurityPolicyEvaluationsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "cloudstorex_security_policy_evaluations_total",
		Help: "Total number of security policy evaluations",
	}, []string{"policy_type", "result"})

	// AI Metrics (Epic 15)
	AIProviderRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "cloudstorex_ai_provider_requests_total",
		Help: "Total number of requests to the AI Provider",
	}, []string{"provider", "model", "status"})

	AIProviderLatency = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "cloudstorex_ai_provider_latency_seconds",
		Help: "Latency of AI Provider inferences",
		Buckets: []float64{0.1, 0.5, 1, 2, 5, 10, 30},
	}, []string{"provider", "model"})

	AIProviderCostTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "cloudstorex_ai_provider_cost_total",
		Help: "Estimated total cost of AI usage in USD",
	}, []string{"provider", "model"})

	AICacheHitsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "cloudstorex_ai_cache_hits_total",
		Help: "Total AI cache hits",
	}, []string{"type"})

	AIGuardrailBlocksTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "cloudstorex_ai_guardrail_blocks_total",
		Help: "Total AI interactions blocked by guardrails",
	}, []string{"reason"})

	ReliabilityReplicationErrorsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cloudstorex_reliability_replication_errors_total",
			Help: "Total replication errors.",
		},
		[]string{"workspace", "source_provider", "target_provider"},
	)

	ReliabilityReplicationLagSeconds = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "cloudstorex_reliability_replication_lag_seconds",
			Help:    "Cross-region replication lag in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"workspace", "source_provider", "target_provider"},
	)

	ReliabilityBackupBytesTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cloudstorex_reliability_backup_bytes_total",
			Help: "Total bytes backed up to cold storage.",
		},
		[]string{"workspace", "provider"},
	)

	ReliabilityBackupErrorsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cloudstorex_reliability_backup_errors_total",
			Help: "Total backup job errors.",
		},
		[]string{"workspace", "provider"},
	)

	ReliabilityRestoreBytesTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cloudstorex_reliability_restore_bytes_total",
			Help: "Total bytes restored from cold storage.",
		},
		[]string{"workspace", "provider"},
	)

	ReliabilityFailoversTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cloudstorex_reliability_failovers_total",
			Help: "Total automatic failovers triggered.",
		},
		[]string{"provider", "status"},
	)

	ReliabilityConsistencyChecksTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cloudstorex_reliability_consistency_checks_total",
			Help: "Total consistency checks performed.",
		},
		[]string{"workspace", "provider", "status"},
	)

	ReliabilityDiscrepanciesFoundTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cloudstorex_reliability_discrepancies_found_total",
			Help: "Total object inconsistencies found.",
		},
		[]string{"workspace", "provider", "type"},
	)

	ReliabilityRepairsCompletedTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cloudstorex_reliability_repairs_completed_total",
			Help: "Total self-healing repairs successfully completed.",
		},
		[]string{"workspace", "provider"},
	)

	ReliabilityChaosExperimentsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cloudstorex_reliability_chaos_experiments_total",
			Help: "Total chaos engineering experiments run.",
		},
		[]string{"workspace", "experiment_type", "status"},
	)

	ReliabilitySLAViolationsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cloudstorex_reliability_sla_violations_total",
			Help: "Total RTO or RPO SLA violations detected during drills.",
		},
		[]string{"workspace", "experiment_type", "sla_type"},
	)

)

// RecordRuntimeMetrics samples standard Go runtime statistics into Prometheus gauges.
func RecordRuntimeMetrics() {
	Goroutines.Set(float64(runtime.NumGoroutine()))
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	MemoryAllocBytes.Set(float64(m.Alloc))
}

// StartRuntimeMetricCollector starts a background goroutine to periodically update runtime gauges.
func StartRuntimeMetricCollector(interval time.Duration) func() {
	ticker := time.NewTicker(interval)
	done := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				RecordRuntimeMetrics()
			case <-done:
				ticker.Stop()
				return
			}
		}
	}()

	return func() {
		close(done)
	}
}
