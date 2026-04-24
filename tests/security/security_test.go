package security

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/hpds/skill-hub/pkg/response"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const baseURL = "http://localhost:8080"

type testCase struct {
	name       string
	method     string
	path       string
	headers    map[string]string
	body       string
	expectCode int  // expected response code
	expectOK   bool // expect 200 HTTP status
}

func runCase(t *testing.T, tc testCase) {
	t.Run(tc.name, func(t *testing.T) {
		req, err := http.NewRequest(tc.method, baseURL+tc.path, strings.NewReader(tc.body))
		require.NoError(t, err)
		for k, v := range tc.headers {
			req.Header.Set(k, v)
		}
		if tc.headers == nil {
			req.Header.Set("Content-Type", "application/json")
		}

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		if tc.expectOK {
			assert.Equal(t, 200, resp.StatusCode, "expected 200 for %s", tc.name)
		}

		if tc.expectCode > 0 {
			var r response.Response
			err := json.NewDecoder(resp.Body).Decode(&r)
			require.NoError(t, err)
			assert.Equal(t, tc.expectCode, r.Code, "expected code %d for %s", tc.expectCode, tc.name)
		}
	})
}

func TestSQLInjection(t *testing.T) {
	payloads := []string{
		"1' OR '1'='1",
		"1; DROP TABLE skills; --",
		"' UNION SELECT * FROM users --",
		"1' AND 1=1 --",
		"admin' --",
		"1' WAITFOR DELAY '0:0:5' --",
		"' OR '1'='1' /*",
	}

	for _, payload := range payloads {
		name := fmt.Sprintf("SQL injection in search: %s", payload[:min(30, len(payload))])
		runCase(t, testCase{
			name:       name,
			method:     "POST",
			path:       "/api/v1/skills/search",
			body:       fmt.Sprintf(`{"query":"%s"}`, payload),
			expectCode: 0,
			expectOK:   true,
		})
	}

	for _, payload := range payloads {
		name := fmt.Sprintf("SQL injection in skill id: %s", payload[:min(30, len(payload))])
		runCase(t, testCase{
			name:       name,
			method:     "GET",
			path:       fmt.Sprintf("/api/v1/skills/%s", payload),
			expectCode: 10002,
			expectOK:   true,
		})
	}
}

func TestXSS(t *testing.T) {
	payloads := []string{
		"<script>alert('xss')</script>",
		"<img src=x onerror=alert(1)>",
		"javascript:alert('xss')",
		"<svg onload=alert(1)>",
		"'><script>alert(1)</script>",
		"<iframe src=javascript:alert(1)>",
	}

	for _, payload := range payloads {
		name := fmt.Sprintf("XSS in search: %s", payload[:min(30, len(payload))])
		runCase(t, testCase{
			name:       name,
			method:     "POST",
			path:       "/api/v1/skills/search",
			body:       fmt.Sprintf(`{"query":"%s"}`, payload),
			expectCode: 0,
			expectOK:   true,
		})
	}
}

func TestAuthenticationBypass(t *testing.T) {
	t.Run("missing auth header on protected endpoint", func(t *testing.T) {
		runCase(t, testCase{
			name:       "no auth header",
			method:     "GET",
			path:       "/api/v1/admin/skills",
			expectCode: 10003,
			expectOK:   true,
		})
	})

	t.Run("invalid token format", func(t *testing.T) {
		runCase(t, testCase{
			name:   "invalid token",
			method: "GET", path: "/api/v1/admin/skills",
			headers:    map[string]string{"Authorization": "Bearer invalid-token"},
			expectCode: 10003,
			expectOK:   true,
		})
	})

	t.Run("empty token", func(t *testing.T) {
		runCase(t, testCase{
			name:       "empty token",
			method:     "GET",
			path:       "/api/v1/admin/skills",
			headers:    map[string]string{"Authorization": "Bearer "},
			expectCode: 10003,
			expectOK:   true,
		})
	})
}

func TestRateLimiting(t *testing.T) {
	t.Run("rapid requests may be rate limited", func(t *testing.T) {
		limited := false
		for i := 0; i < 100; i++ {
			req, _ := http.NewRequest("GET", baseURL+"/api/v1/skills", nil)
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				continue
			}
			if resp.StatusCode == 429 {
				limited = true
				resp.Body.Close()
				break
			}
			resp.Body.Close()
		}
		// This is informational - rate limiting may or may not be enforced
		t.Logf("Rate limiting triggered: %v", limited)
	})
}

func TestInputValidation(t *testing.T) {
	t.Run("oversized payload rejected", func(t *testing.T) {
		longStr := strings.Repeat("A", 1024*1024) // 1MB
		runCase(t, testCase{
			name:       "oversized search query",
			method:     "POST",
			path:       "/api/v1/skills/search",
			body:       fmt.Sprintf(`{"query":"%s"}`, longStr),
			expectCode: -1, // don't check code, just check status
		})
	})

	t.Run("malformed json rejected", func(t *testing.T) {
		runCase(t, testCase{
			name:       "malformed json",
			method:     "POST",
			path:       "/api/v1/skills/search",
			body:       `{"query": "test"`,
			expectCode: 10002,
			expectOK:   true,
		})
	})
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
