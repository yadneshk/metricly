package v1

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
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

// log API requests
func logAPIRequests(r *http.Request, duration int64, statusCode int) {
	slog.Info(
		fmt.Sprintf(
			"%s %s %s status: %d len: %d time: %dms", r.RemoteAddr, r.Method, r.URL, statusCode, r.ContentLength, duration,
		),
	)
}
