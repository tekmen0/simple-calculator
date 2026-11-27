package httpapi_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/tekmen0/simple-calculator/backend/internal/httpapi"
)

func TestHealth(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	w := httptest.NewRecorder()
	httpapi.NewHandler().ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("GET /api/health returned %d with %s, want 200", w.Code, w.Body)
	}
	if got := w.Header().Get("Content-Type"); got != "application/json" {
		t.Errorf("GET /api/health Content-Type = %q, want application/json", got)
	}
	if got := w.Header().Get("Cache-Control"); got != "no-store" {
		t.Errorf("GET /api/health Cache-Control = %q, want no-store", got)
	}
	var body struct {
		Status string `json:"status"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("GET /api/health response = %s, want JSON (decode error: %v)", w.Body, err)
	}
	if body.Status != "ok" {
		t.Errorf("GET /api/health status = %q, want ok", body.Status)
	}
}

func TestServeHTTPResults(t *testing.T) {
	const addBody = `{"left":7,"right":2}`
	handler := httpapi.NewHandler()
	tests := []struct {
		name        string
		operation   string
		contentType string
		body        string
		want        float64
	}{
		{name: "add", operation: "add", contentType: "application/json", body: addBody, want: 9},
		{name: "subtract", operation: "subtract", contentType: "application/json", body: addBody, want: 5},
		{name: "multiply", operation: "multiply", contentType: "application/json", body: addBody, want: 14},
		{name: "divide", operation: "divide", contentType: "application/json", body: addBody, want: 3.5},
		{
			name: "out of order fields and JSON media type parameters", operation: "add",
			contentType: "application/json; charset=utf-8", body: "{\"right\":-2,\"left\":2}\n ", want: 0,
		},
		{
			name: "body at size limit", operation: "add", contentType: "application/json",
			body: addBody + strings.Repeat(" ", 4096-len(addBody)), want: 9,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/api/calculations/"+tt.operation, strings.NewReader(tt.body))
			r.Header.Set("Content-Type", tt.contentType)
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, r)

			if w.Code != http.StatusOK {
				t.Fatalf("POST %s returned %d with %s, want 200", r.URL.Path, w.Code, w.Body)
			}
			if got := w.Header().Get("Content-Type"); got != "application/json" {
				t.Errorf("POST %s Content-Type = %q, want application/json", r.URL.Path, got)
			}
			var body struct {
				Result *float64 `json:"result"`
			}
			if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil || body.Result == nil {
				t.Fatalf("POST %s response = %s, want a JSON number in result (decode error: %v)", r.URL.Path, w.Body, err)
			}
			if got := *body.Result; got != tt.want {
				t.Errorf("POST %s result = %g, want %g", r.URL.Path, got, tt.want)
			}
		})
	}
}

func TestServeHTTPInvalidRequests(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{name: "empty body"},
		{name: "malformed JSON", body: `{"left":7,"right":`},
		{name: "unclosed object", body: `{"left":7,"right":2`},
		{name: "trailing comma", body: `{"left":7,"right":2,}`},
		{name: "array", body: `[{"left":7,"right":2}]`},
		{name: "null body", body: `null`},
		{name: "missing operand", body: `{"left":7}`},
		{name: "null operand", body: `{"left":null,"right":2}`},
		{name: "numeric string", body: `{"left":"7","right":2}`},
		{name: "boolean operand", body: `{"left":true,"right":2}`},
		{name: "out of range operand", body: `{"left":1e309,"right":2}`},
		{name: "unknown field", body: `{"left":7,"right":2,"extra":0}`},
		{name: "misnamed field", body: `{"LEFT":7,"right":2}`},
		{name: "duplicate field", body: `{"left":7,"right":2,"left":1}`},
		{name: "trailing JSON", body: `{"left":7,"right":2}null`},
		{name: "trailing junk", body: `{"left":7,"right":2}x`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/api/calculations/add", strings.NewReader(tt.body))
			r.Header.Set("Content-Type", "application/json")
			checkErrorResponse(t, r, http.StatusBadRequest, "invalid_request", "")
		})
	}
}

func TestServeHTTPUnreadableBody(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/api/calculations/add", iotest.ErrReader(io.ErrUnexpectedEOF))
	r.Header.Set("Content-Type", "application/json")
	checkErrorResponse(t, r, http.StatusBadRequest, "invalid_request", "")
}

func TestServeHTTPErrors(t *testing.T) {
	const path = "/api/calculations/divide"
	const validBody = `{"left":7,"right":2}`
	tests := []struct {
		name, method, path, contentType, body string
		status                                int
		code, allow                           string
	}{
		{
			name: "unknown path", method: http.MethodPost, path: "/api/unknown",
			status: http.StatusNotFound, code: "not_found",
		},
		{
			name: "unknown operation", method: http.MethodGet, path: "/api/calculations/power",
			status: http.StatusNotFound, code: "not_found",
		},
		{
			name: "extra path segment", method: http.MethodPost, path: path + "/extra",
			status: http.StatusNotFound, code: "not_found",
		},
		{
			name: "repeated slash", method: http.MethodPost, path: "/api//calculations/divide",
			status: http.StatusNotFound, code: "not_found",
		},
		{
			name: "wrong method", method: http.MethodGet, path: path,
			status: http.StatusMethodNotAllowed, code: "method_not_allowed", allow: http.MethodPost,
		},
		{
			name: "wrong health method", method: http.MethodPost, path: "/api/health",
			status: http.StatusMethodNotAllowed, code: "method_not_allowed", allow: "GET, HEAD",
		},
		{
			name: "missing content type", method: http.MethodPost, path: path, body: validBody,
			status: http.StatusUnsupportedMediaType, code: "unsupported_media_type",
		},
		{
			name: "wrong content type", method: http.MethodPost, path: path, body: validBody, contentType: "text/plain",
			status: http.StatusUnsupportedMediaType, code: "unsupported_media_type",
		},
		{
			name: "oversized body", method: http.MethodPost, path: path, contentType: "application/json",
			body: strings.Repeat("x", 4097), status: http.StatusRequestEntityTooLarge, code: "request_too_large",
		},
		{
			name: "division by zero", method: http.MethodPost, path: path, contentType: "application/json",
			body: `{"left":7,"right":0}`, status: http.StatusUnprocessableEntity, code: "division_by_zero",
		},
		{
			name: "division by negative zero", method: http.MethodPost, path: path, contentType: "application/json",
			body: `{"left":7,"right":-0}`, status: http.StatusUnprocessableEntity, code: "division_by_zero",
		},
		{
			name: "overflow", method: http.MethodPost, path: path, contentType: "application/json",
			body: `{"left":1e308,"right":0.1}`, status: http.StatusUnprocessableEntity, code: "result_out_of_range",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(tt.method, tt.path, strings.NewReader(tt.body))
			if tt.contentType != "" {
				r.Header.Set("Content-Type", tt.contentType)
			}
			checkErrorResponse(t, r, tt.status, tt.code, tt.allow)
		})
	}
}

func checkErrorResponse(t *testing.T, r *http.Request, status int, code, allow string) {
	t.Helper()
	w := httptest.NewRecorder()
	httpapi.NewHandler().ServeHTTP(w, r)

	if w.Code != status {
		t.Errorf("%s %s returned %d with %s, want %d", r.Method, r.URL.Path, w.Code, w.Body, status)
	}
	if got := w.Header().Get("Content-Type"); got != "application/json" {
		t.Errorf("%s %s Content-Type = %q, want application/json", r.Method, r.URL.Path, got)
	}
	if got := w.Header().Get("Allow"); got != allow {
		t.Errorf("%s %s Allow = %q, want %q", r.Method, r.URL.Path, got, allow)
	}
	var body struct {
		Error struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("%s %s response = %s, want JSON error (decode error: %v)", r.Method, r.URL.Path, w.Body, err)
	}
	if got := body.Error.Code; got != code {
		t.Errorf("%s %s error code = %q, want %q", r.Method, r.URL.Path, got, code)
	}
	if strings.TrimSpace(body.Error.Message) == "" {
		t.Errorf("%s %s returned an empty error message, want a human-readable message", r.Method, r.URL.Path)
	}
}
