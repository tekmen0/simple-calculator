package httpapi

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"mime"
	"net/http"
	"path"

	"github.com/tekmen0/simple-calculator/backend/internal/calculator"
)

const maxRequestBytes = 4 * 1024

type calculationResponse struct {
	Result float64 `json:"result"`
}

type healthResponse struct {
	Status string `json:"status"`
}

type errorResponse struct {
	Error errorDetail `json:"error"`
}

type errorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func NewHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/health", handleHealth)
	mux.HandleFunc("/api/calculations/{operation}", handleCalculation)
	mux.HandleFunc("/", notFound)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Keep the API's JSON 404 for paths that ServeMux would redirect.
		if path.Clean(r.URL.Path) != r.URL.Path {
			notFound(w, r)
			return
		}
		mux.ServeHTTP(w, r)
	})
	return http.MaxBytesHandler(logRequests(handler), maxRequestBytes)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		w.Header().Set("Allow", "GET, HEAD")
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Use GET or HEAD to check API health.")
		return
	}
	w.Header().Set("Cache-Control", "no-store")
	writeJSON(w, http.StatusOK, healthResponse{Status: "ok"})
}

func handleCalculation(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("operation")
	if !calculator.Supports(name) {
		notFound(w, r)
		return
	}
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Use POST to perform a calculation.")
		return
	}
	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || mediaType != "application/json" {
		writeError(w, http.StatusUnsupportedMediaType, "unsupported_media_type", "Use the application/json content type.")
		return
	}

	body, err := io.ReadAll(r.Body)
	if _, ok := errors.AsType[*http.MaxBytesError](err); ok {
		writeError(w, http.StatusRequestEntityTooLarge, "request_too_large", "Request body is too large.")
		return
	}
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request", "Could not read the request body.")
		return
	}
	input, ok := decodeOperands(body)
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid_request", "Send a JSON object with finite numbers for left and right.")
		return
	}

	result, err := calculator.Calculate(name, input.left, input.right)
	if err != nil {
		writeCalculationError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, calculationResponse{Result: result})
}

func notFound(w http.ResponseWriter, _ *http.Request) {
	writeError(w, http.StatusNotFound, "not_found", "Calculator endpoint not found.")
}

func writeCalculationError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, calculator.ErrDivisionByZero):
		writeError(w, http.StatusUnprocessableEntity, "division_by_zero", "Cannot divide by zero.")
	case errors.Is(err, calculator.ErrResultOutOfRange):
		writeError(w, http.StatusUnprocessableEntity, "result_out_of_range", "The result is outside the supported numeric range.")
	default:
		slog.Error("calculation failed", "err", err)
		writeError(w, http.StatusInternalServerError, "internal_error", "Could not perform the calculation.")
	}
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, errorResponse{Error: errorDetail{Code: code, Message: message}})
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		slog.Debug("could not write JSON response", "err", err)
	}
}
