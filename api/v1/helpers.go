package v1

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// represents the structure of the reponse Prometheus's API call
type PrometheusQueryResponse struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric map[string]string `json:"metric"`
			Value  [2]interface{}    `json:"value"`
		} `json:"result"`
	} `json:"data"`
}

func QueryPrometheus(queryURL string) (*PrometheusQueryResponse, error) {
	// queryURL := fmt.Sprintf("http://10.1.23.133:9090/api/v1/query?query=%s", query)

	resp, err := http.Get(queryURL)
	if err != nil {
		return nil, fmt.Errorf("failed to query")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("prom returned %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response %v", err)
	}

	var response PrometheusQueryResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse result %v", err)
	}
	return &response, nil

}

func PreparePromQuery(params map[string]string) string {
	query := "http://10.1.23.133:9090/api/v1/query?"
	for param, value := range params {
		query += fmt.Sprintf("&%s=%s", param, value)
	}
	return query
}
