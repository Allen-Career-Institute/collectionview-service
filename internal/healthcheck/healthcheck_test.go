package healthcheck

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestNewHandler_ErrorWritingResponse2(t *testing.T) {
	router := mux.NewRouter()
	router.Handle("/health/status", NewHandler())

	// Create a mock request
	req, err := http.NewRequest("GET", "/health/status", nil)
	assert.NoError(t, err)

	// Create a mock response recorder that fails on write
	rr := &mockResponseWriter{httptest.NewRecorder(), false}

	// Serve the request using the router
	router.ServeHTTP(rr, req)

	// Assert the response

	// Error should be logged
	// Assert the log output as per your logging implementation
}

func TestNewHandler_ErrorWritingResponse(t *testing.T) {
	router := mux.NewRouter()
	router.Handle("/health/status", NewHandler())

	// Create a mock request
	req, err := http.NewRequest("GET", "/health/status", nil)
	assert.NoError(t, err)

	// Create a mock response recorder that fails on write
	rr := &mockResponseWriter{httptest.NewRecorder(), true}

	// Serve the request using the router
	router.ServeHTTP(rr, req)

	// Assert the response

	// Error should be logged
	// Assert the log output as per your logging implementation
}

// mockResponseWriter is a custom response writer that can be used to simulate errors during writing response.
type mockResponseWriter struct {
	http.ResponseWriter
	failOnWrite bool
}

func (m *mockResponseWriter) Write([]byte) (int, error) {
	if m.failOnWrite {
		return 0, assert.AnError // Simulate an error during write
	}
	return m.ResponseWriter.Write([]byte("OK"))
}
