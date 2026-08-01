import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate, Trend } from 'k6/metrics';

export const errorRate = new Rate('errors');
export const apiLatency = new Trend('api_latency_ms');

export const options = {
  vus: 3,
  duration: '30s',
  thresholds: {
    errors: ['rate<0.01'], // < 1% error rate
    http_req_duration: ['p(95)<500'], // 95% of requests must complete under 500ms
  },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';

export default function () {
  // 1. Health check
  const healthRes = http.get(`${BASE_URL}/api/v1/health`);
  check(healthRes, {
    'health check is 200': (r) => r.status === 200,
    'status is up': (r) => r.json('status') === 'up',
  }) || errorRate.add(1);
  apiLatency.add(healthRes.timings.duration);

  // 2. Metrics check
  const metricsRes = http.get(`${BASE_URL}/metrics`);
  check(metricsRes, {
    'metrics status is 200': (r) => r.status === 200,
    'metrics contains http requests total': (r) => r.body.includes('cloudstorex_http_requests_total'),
  }) || errorRate.add(1);

  sleep(1);
}
