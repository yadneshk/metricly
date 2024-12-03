package v1

import (
	"metricly/config"
	"net/http"
	"net/http/httptest"
)

var (
	conf = &config.Config{
		Prometheus: struct {
			Address string "yaml:\"address\""
			Port    string "yaml:\"port\""
		}{
			Address: "localhost",
			Port:    "9090",
		},
	}

	mockPrometheusResponse = PrometheusResponse{
		Status: "success",
		Data: struct {
			ResultType string `json:"resultType"`
			Result     []struct {
				Metric map[string]string `json:"metric"`
				Value  [2]interface{}    `json:"value"`
			} `json:"result"`
		}{
			ResultType: "vector",
			Result: []struct {
				Metric map[string]string `json:"metric"`
				Value  [2]interface{}    `json:"value"`
			}{
				{
					Metric: map[string]string{"__name__": "test_metric"},
					Value:  [2]interface{}{"1689636523", "100"},
				},
			},
		},
	}
)

func newMockServer(response []byte) (*httptest.Server, error) {

	// emptyServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	w.WriteHeader(http.StatusNoContent) // 204 No Content
	// }))

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}))
	return mockServer, nil
}
