package v1

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockPrometheusResponse struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric map[string]string `json:"metric"`
			Values [][]interface{}   `json:"values"`
		} `json:"result"`
	} `json:"data"`
}

func TestQueryPrometheus(t *testing.T) {
	// Mock response data
	mockResponse := MockPrometheusResponse{
		Status: "success",
		Data: struct {
			ResultType string `json:"resultType"`
			Result     []struct {
				Metric map[string]string `json:"metric"`
				Values [][]interface{}   `json:"values"`
			} `json:"result"`
		}{
			ResultType: "vector",
			Result: []struct {
				Metric map[string]string `json:"metric"`
				Values [][]interface{}   `json:"values"`
			}{
				{
					Metric: map[string]string{"__name__": "test_metric"},
					Values: [][]interface{}{
						{1.0, "10"},
						{2.0, "20"},
					},
				},
			},
		},
	}

	mockBody, _ := json.Marshal(mockResponse)

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET request, got %s", r.Method)
		}
		if r.URL.String() != "/api/v1/query" {
			t.Errorf("unexpected URL: got %s", r.URL.String())
		}
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(mockBody); err != nil {
			t.Error("failed to write response")
		}
	}))
	defer mockServer.Close()

	// Test the QueryPrometheus function
	var result MockPrometheusResponse
	err := QueryPrometheus(mockServer.URL+"/api/v1/query", &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Validate the response
	if result.Status != "success" {
		t.Errorf("unexpected status: got %s, want success", result.Status)
	}
	if len(result.Data.Result) != 1 {
		t.Errorf("unexpected number of results: got %d, want 1", len(result.Data.Result))
	}
	if result.Data.Result[0].Metric["__name__"] != "test_metric" {
		t.Errorf("unexpected metric name: got %s, want test_metric", result.Data.Result[0].Metric["__name__"])
	}
}

func TestQueryPrometheusErrorResponse(t *testing.T) {
	// Create a mock server with an error response
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := w.Write([]byte(`{"status":"error","errorType":"internal","error":"something went wrong"}`)); err != nil {
			t.Error("failed to write mock response")
		}
	}))
	defer mockServer.Close()

	var result MockPrometheusResponse
	err := QueryPrometheus(mockServer.URL, &result)

	// Validate the error
	if err == nil || err.Error() != "prom returned 500 Internal Server Error" {
		t.Fatalf("unexpected error: got %v, want prom returned 500 Internal Server Error", err)
	}
}
