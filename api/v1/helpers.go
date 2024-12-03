package v1

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
)

// represents the structure of the reponse Prometheus's API call
// used for /query & /aggregate apis
type PrometheusResponse struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric map[string]string `json:"metric"`
			Value  [2]interface{}    `json:"value"`
		} `json:"result"`
	} `json:"data"`
}

// represents the structure of the reponse Prometheus's API call
// used for /query_range api
type PrometheusQueryRangeResponse struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric map[string]string `json:"metric"`
			Values [][]interface{}   `json:"values"`
		} `json:"result"`
	} `json:"data"`
}

// represents error reponse
type ErrorResponse struct {
	Status string `json:"status"`
	Data   struct {
		StatusCode int    `json:"status_code"`
		Message    string `json:"message"`
	} `json:"data"`
}

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

func sendErrorResponse(w http.ResponseWriter, statusCode int, message string) {

	error := ErrorResponse{
		Status: "failed",
		Data: struct {
			StatusCode int    "json:\"status_code\""
			Message    string "json:\"message\""
		}{
			StatusCode: statusCode,
			Message:    message,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(error); err != nil {
		slog.Error("failed to parse error into json")
	}
}

func processQueryParams(reqParam url.Values, supportedParams map[string]bool, requiredParams []string) (map[string]string, error) {

	result := map[string]string{}
	for param, value := range reqParam {
		if !supportedParams[param] {
			return nil, fmt.Errorf("unknown parameter in request: %s", param)
		}
		if len(value) == 0 {
			return nil, fmt.Errorf("parameter %s found empty", param)
		}
		for _, v := range value {
			if v == "" {
				return nil, fmt.Errorf("value for parameter %s found empty", param)
			}
			result[param] = v
		}
	}

	for _, req := range requiredParams {
		if _, exists := result[req]; !exists {
			return nil, fmt.Errorf("parameter %s required", req)
		}
	}

	tmp := result["metric"]
	delete(result, "metric")
	result["query"] = tmp
	return result, nil
}
