package metrics_test

import (
	"testing"
	"time"

	"github.com/cloudstorex/backend/internal/observability/metrics"
	"github.com/stretchr/testify/assert"
)

func TestMetrics_RegistrationAndLabels(t *testing.T) {
	// Verify HTTP metrics with low-cardinality labels
	metrics.HTTPRequestsTotal.WithLabelValues("GET", "/api/v1/objects", "200").Inc()
	metrics.HTTPRequestDurationSeconds.WithLabelValues("GET", "/api/v1/objects", "200").Observe(0.125)
	metrics.HTTPActiveRequests.WithLabelValues("GET", "/api/v1/objects").Add(1)
	metrics.HTTPErrorRate.WithLabelValues("POST", "/api/v1/objects", "500").Inc()

	// Verify Storage metrics (Refinement 2: low-cardinality workspace, provider, bucket, status)
	metrics.StorageUploadsTotal.WithLabelValues("ws-1", "aws-s3", "test-bucket", "success").Inc()
	metrics.StorageUploadBytes.WithLabelValues("ws-1", "aws-s3", "test-bucket").Add(1048576)
	metrics.StorageDownloadsTotal.WithLabelValues("ws-1", "minio", "test-bucket", "success").Inc()
	metrics.StorageDownloadBytes.WithLabelValues("ws-1", "minio", "test-bucket").Add(2048)
	metrics.StorageDeleteOperations.WithLabelValues("ws-1", "aws-s3", "test-bucket", "success").Inc()

	// Verify Provider metrics
	metrics.ProviderLatencySeconds.WithLabelValues("aws-s3", "upload", "us-east-1", "success").Observe(0.045)
	metrics.ProviderErrorsTotal.WithLabelValues("minio", "download", "local").Inc()
	metrics.ProviderRequestsTotal.WithLabelValues("aws-s3", "upload", "us-east-1", "200").Inc()
	metrics.ProviderHealthStatus.WithLabelValues("aws-s3", "us-east-1").Set(1)

	// Verify Policy Engine metrics
	metrics.PolicyEvaluationsTotal.WithLabelValues("ws-1", "upload", "matched").Inc()
	metrics.PolicyRoutingDecisions.WithLabelValues("ws-1", "aws-s3", "upload", "SIZE_RULE").Inc()
	metrics.PolicyRuleMatches.WithLabelValues("ws-1", "SIZE_RULE").Inc()
	metrics.PolicyFallbackRoutes.WithLabelValues("ws-1", "upload").Inc()

	// Verify Metadata metrics
	metrics.MetadataQueriesTotal.WithLabelValues("ws-1", "search", "200").Inc()
	metrics.MetadataSearchLatencySeconds.WithLabelValues("ws-1", "200").Observe(0.012)
	metrics.MetadataObjectCatalogSize.WithLabelValues("ws-1", "test-bucket").Set(150)

	// Verify Background Jobs metrics
	metrics.JobsRunning.WithLabelValues("lifecycle-expire").Set(2)
	metrics.JobsCompletedTotal.WithLabelValues("lifecycle-expire", "success").Inc()
	metrics.JobsFailedTotal.WithLabelValues("lifecycle-expire").Inc()
	metrics.JobsRetryCountTotal.WithLabelValues("lifecycle-expire").Inc()
	metrics.JobsQueueDepth.WithLabelValues("default").Set(5)

	// Verify Infrastructure runtime collector
	stop := metrics.StartRuntimeMetricCollector(10 * time.Millisecond)
	time.Sleep(25 * time.Millisecond)
	stop()

	assert.NotNil(t, metrics.Goroutines)
	assert.NotNil(t, metrics.MemoryAllocBytes)
}
