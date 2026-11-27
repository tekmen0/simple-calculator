package httpapi_test

import (
	"bytes"
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/tekmen0/simple-calculator/backend/internal/httpapi"
)

func TestRequestLog(t *testing.T) {
	var output bytes.Buffer
	previous := slog.Default()
	logOutput, logFlags := log.Writer(), log.Flags()
	slog.SetDefault(slog.New(slog.NewJSONHandler(&output, nil)))
	t.Cleanup(func() {
		slog.SetDefault(previous)
		log.SetOutput(logOutput)
		log.SetFlags(logFlags)
	})
	handler := httpapi.NewHandler()

	tests := []struct {
		method, target, body string
		route                string
		status               int
	}{
		{
			method: http.MethodGet, target: "/api/health?token=private-query",
			route: "/api/health", status: http.StatusOK,
		},
		{
			method: http.MethodPost, target: "/api/calculations/divide", body: `{"left":7,"right":0}`,
			route: "/api/calculations/{operation}", status: http.StatusUnprocessableEntity,
		},
		{
			method: http.MethodGet, target: "/api/private-path",
			route: "unmatched", status: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.method+" "+tt.route, func(t *testing.T) {
			output.Reset()
			r := httptest.NewRequest(tt.method, tt.target, strings.NewReader(tt.body))
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, r)

			var entry struct {
				Level    string `json:"level"`
				Message  string `json:"msg"`
				Method   string `json:"method"`
				Route    string `json:"route"`
				Status   int    `json:"status"`
				Duration *int64 `json:"duration"`
			}
			if err := json.Unmarshal(output.Bytes(), &entry); err != nil {
				t.Fatalf("request log = %q, want a single JSON record (decode error: %v)", output.String(), err)
			}
			if entry.Level != "INFO" || entry.Message != "request handled" || entry.Duration == nil || *entry.Duration < 0 {
				t.Errorf("request log = %s, want an INFO request record with a duration", output.String())
			}
			if entry.Method != tt.method || entry.Route != tt.route || entry.Status != tt.status || entry.Status != w.Code {
				t.Errorf("request log = %s, want method %s, route %s, and status %d (response status %d)",
					output.String(), tt.method, tt.route, tt.status, w.Code)
			}
			if strings.Contains(output.String(), "private-") {
				t.Error("request log contains private path or query data")
			}
		})
	}
}
