import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate } from 'k6/metrics';

export const errorRate = new Rate('errors');

export const options = {
  stages: [
    { duration: '10s', target: 5 },   // Normal baseline
    { duration: '10s', target: 150 }, // Immediate spike to 150 VUs
    { duration: '1m', target: 150 },  // Maintain spike for 1 minute
    { duration: '10s', target: 5 },   // Drop back down to baseline
    { duration: '30s', target: 0 },   // Exit
  ],
  thresholds: {
    errors: ['rate<0.05'],
    http_req_duration: ['p(95)<3000'],
  },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';

export default function () {
  const res = http.get(`${BASE_URL}/api/v1/health`);
  check(res, {
    'status is 200': (r) => r.status === 200,
  }) || errorRate.add(1);

  sleep(0.3);
}
