import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate } from 'k6/metrics';

export const errorRate = new Rate('errors');

export const options = {
  stages: [
    { duration: '2m', target: 30 },  // Ramp to moderate sustained load
    { duration: '2h', target: 30 },  // Run for 2 hours to detect memory leaks and goroutine exhaustion
    { duration: '2m', target: 0 },   // Graceful shutdown
  ],
  thresholds: {
    errors: ['rate<0.01'], // Require < 1% error rate over endurance run
    http_req_duration: ['p(95)<750'], // Latency should not degrade over time
  },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';

export default function () {
  const res = http.get(`${BASE_URL}/api/v1/health`);
  check(res, {
    'status is 200': (r) => r.status === 200,
  }) || errorRate.add(1);

  sleep(1);
}
