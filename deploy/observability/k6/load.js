import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate, Trend } from 'k6/metrics';

export const errorRate = new Rate('errors');
export const uploadLatency = new Trend('upload_latency_ms');

export const options = {
  stages: [
    { duration: '30s', target: 10 }, // Ramp-up to 10 VUs
    { duration: '1m', target: 50 },  // Ramping to typical load of 50 VUs
    { duration: '3m', target: 50 },  // Stay at 50 VUs for 3 minutes
    { duration: '30s', target: 0 },  // Scale down
  ],
  thresholds: {
    errors: ['rate<0.02'], // Error rate under 2%
    http_req_duration: ['p(95)<1000', 'p(99)<2500'], // P95 < 1s, P99 < 2.5s
  },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';

export default function () {
  // Simulate API read traffic
  const res = http.get(`${BASE_URL}/api/v1/health`);
  check(res, {
    'status is 200': (r) => r.status === 200,
  }) || errorRate.add(1);

  sleep(0.5);
}
