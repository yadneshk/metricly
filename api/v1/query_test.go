package v1

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPrometheusQueryHandler(t *testing.T) {

	mockBody, _ := json.Marshal(mockPrometheusResponse)

	// Mock Prometheus server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(mockBody); err != nil {
			t.Error("failed to write mock response")
		}
	}))
	defer mockServer.Close()

	// Override the Prometheus address and port
	conf.Prometheus.Address = strings.Split(mockServer.URL, ":")[1][2:]
	conf.Prometheus.Port = strings.Split(mockServer.URL, ":")[2]

	// Instantiate handler
	handler := PrometheusQueryHandler(conf)

	// Create a mock HTTP request
	req := httptest.NewRequest(http.MethodGet, "/query?metric=mock_metric", nil)
	w := httptest.NewRecorder()

	// Call the handler
	handler.ServeHTTP(w, req)

	// Assertions
	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status %d; got %d", http.StatusOK, resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var parsedResponse PrometheusResponse
	err := json.Unmarshal(body, &parsedResponse)
	if err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	// Check the mock response data
	if parsedResponse.Status != "success" {
		t.Errorf("expected status 'success'; got %s", parsedResponse.Status)
	}
	if len(parsedResponse.Data.Result) != 1 {
		t.Errorf("expected 1 result; got %d", len(parsedResponse.Data.Result))
	}
	if parsedResponse.Data.Result[0].Metric["__name__"] != "test_metric" {
		t.Errorf("expected metric name 'mock_metric'; got %s", parsedResponse.Data.Result[0].Metric["__name__"])
	}
	if parsedResponse.Data.Result[0].Value[1] != "100" {
		t.Errorf("expected value '123.45'; got %v", parsedResponse.Data.Result[0].Value[1])
	}
}
