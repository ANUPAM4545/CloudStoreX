import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate } from 'k6/metrics';

export const errorRate = new Rate('errors');

export const options = {
  stages: [
    { duration: '1m', target: 50 },  // Ramp to normal load
    { duration: '2m', target: 150 }, // Push to heavy load
    { duration: '3m', target: 200 }, // Peak saturation limit test
    { duration: '1m', target: 0 },   // Cool down
  ],
  thresholds: {
    errors: ['rate<0.05'], // Allow up to 5% errors during peak stress
    http_req_duration: ['p(95)<2500'],
  },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';

export default function () {
  const res = http.get(`${BASE_URL}/api/v1/health`);
  check(res, {
    'status is 200': (r) => r.status === 200,
  }) || errorRate.add(1);

  sleep(0.2);
}
