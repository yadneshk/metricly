package v1

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"

	"testing"
)

func TestPrometheusAggregateHandler(t *testing.T) {

	mockBody, _ := json.Marshal(mockPrometheusResponse)

	mockServer, _ := newMockServer(mockBody)
	defer mockServer.Close()

	conf.Prometheus.Address = strings.Split(mockServer.URL, ":")[1][2:]
	conf.Prometheus.Port = strings.Split(mockServer.URL, ":")[2]

	handler := PrometheusAggregateHandler(conf)

	req := httptest.NewRequest(http.MethodGet, "/aggregate?metric=metricly_cpu_total&operation=avg&window=1h", nil)
	w := httptest.NewRecorder()

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
		t.Errorf("expected value '100'; got %v", parsedResponse.Data.Result[0].Value[1])
	}

}
