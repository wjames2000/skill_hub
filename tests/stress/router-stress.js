// k6 stress test for Skill Hub Router API
// Usage: k6 run router-stress.js
// Requires: k6 installed (https://k6.io/docs/getting-started/installation/)

import http from 'k6/http';
import { check, sleep, group } from 'k6';
import { Rate, Trend, Counter } from 'k6/metrics';

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';

const matchErrorRate = new Rate('match_errors');
const matchDuration = new Trend('match_duration_ms');
const executeErrorRate = new Rate('execute_errors');
const executeDuration = new Trend('execute_duration_ms');
const searchErrorRate = new Rate('search_errors');
const searchDuration = new Trend('search_duration_ms');

const queries = [
  'analyze excel data trends',
  'create a presentation',
  'generate a chart',
  'convert pdf to excel',
  'summarize document',
  'check grammar',
  'translate to chinese',
  'send email to team',
  'schedule meeting',
  'analyze customer feedback',
  'merge pdf files',
  'resize images',
  'optimize sql query',
  'backup database',
  'monitor website uptime',
];

export const options = {
  stages: [
    { duration: '1m', target: 10 },   // Ramp up to 10 users
    { duration: '2m', target: 50 },   // Ramp up to 50 users
    { duration: '3m', target: 100 },  // Ramp up to 100 users
    { duration: '2m', target: 100 },  // Stay at 100 users
    { duration: '1m', target: 0 },    // Ramp down
  ],
  thresholds: {
    http_req_duration: ['p(95)<5000'],   // 95% of requests < 5s
    match_errors: ['rate<0.10'],          // < 10% error rate
    execute_errors: ['rate<0.15'],        // < 15% error rate
    search_errors: ['rate<0.10'],         // < 10% error rate
  },
};

function randomQuery() {
  return queries[Math.floor(Math.random() * queries.length)];
}

export default function () {
  group('Router API Stress Test', function () {
    // POST /api/v1/router/match
    group('Match endpoint', function () {
      const payload = JSON.stringify({
        query: randomQuery(),
        top_k: 5,
      });
      const params = {
        headers: { 'Content-Type': 'application/json' },
        tags: { name: 'router_match' },
      };
      const start = Date.now();
      const res = http.post(`${BASE_URL}/api/v1/router/match`, payload, params);
      const duration = Date.now() - start;
      matchDuration.add(duration);
      matchErrorRate.add(res.status !== 200);

      check(res, {
        'match status is 200': (r) => r.status === 200,
        'match response has success code': (r) => {
          try {
            return JSON.parse(r.body).code === 0;
          } catch { return false; }
        },
      });
    });

    sleep(0.5);

    // POST /api/v1/router/execute
    group('Execute endpoint', function () {
      const payload = JSON.stringify({
        query: randomQuery(),
        skill_id: Math.floor(Math.random() * 10) + 1,
      });
      const params = {
        headers: { 'Content-Type': 'application/json' },
        tags: { name: 'router_execute' },
      };
      const start = Date.now();
      const res = http.post(`${BASE_URL}/api/v1/router/execute`, payload, params);
      const duration = Date.now() - start;
      executeDuration.add(duration);
      executeErrorRate.add(res.status !== 200);

      check(res, {
        'execute status is 200': (r) => r.status === 200,
      });
    });

    sleep(0.5);

    // POST /api/v1/skills/search
    group('Search endpoint', function () {
      const payload = JSON.stringify({
        query: randomQuery(),
      });
      const params = {
        headers: { 'Content-Type': 'application/json' },
        tags: { name: 'skills_search' },
      };
      const start = Date.now();
      const res = http.post(`${BASE_URL}/api/v1/skills/search`, payload, params);
      const duration = Date.now() - start;
      searchDuration.add(duration);
      searchErrorRate.add(res.status !== 200);

      check(res, {
        'search status is 200': (r) => r.status === 200,
        'search response has success code': (r) => {
          try {
            return JSON.parse(r.body).code === 0;
          } catch { return false; }
        },
      });
    });

    sleep(1);
  });
}
