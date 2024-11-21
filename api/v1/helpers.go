package v1

import (
	"encoding/json"
	"fmt"
	"io"
	"metricly/config"
	"net/http"
)

// make QueryPrometheus a generic function so that it can serve both reuests, query & query_range
func QueryPrometheus[T any](queryURL string, target *T) error {

	resp, err := http.Get(queryURL)
	if err != nil {
		return fmt.Errorf("failed to query: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("prom returned %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response %v", err)
	}

	// var response PrometheusQueryResponse
	err = json.Unmarshal(body, target)
	if err != nil {
		return fmt.Errorf("failed to parse result %v", err)
	}
	return nil

}

func PreparePromQuery(conf *config.Config, promExpr string, params map[string]string) string {

	// query := "http://10.1.23.133:9090/api/v1/query?"
	query := fmt.Sprintf("http://%s:%s/api/v1/%s?", conf.Prometheus.Address, conf.Prometheus.Port, promExpr)
	for param, value := range params {
		query += fmt.Sprintf("&%s=%s", param, value)
	}
	return query
}
