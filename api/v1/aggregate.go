package v1

import (
	"encoding/json"
	"fmt"
	"metricly/config"
	"metricly/pkg/prometheus"
	"net/http"
)

var (
	supportedOperations = map[string]string{
		"avg":   "avg_over_time",
		"max":   "max_over_time",
		"min":   "min_over_time",
		"count": "count_over_time",
	}
)

func PrometheusAggregateHandler(conf *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		requestParams := r.URL.Query()

		supportedParams := map[string]bool{
			"metric":    true,
			"operation": true,
			"window":    true,
		}
		requiredParams := []string{"metric", "operation", "window"}

		queryParams, err := processQueryParams(requestParams, supportedParams, requiredParams)
		if err != nil {
			sendErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// needed to replace "avg" with "avg_over_time" to make it compatible with Prom
		queryParams, err = processAggregateParams(queryParams)
		if err != nil {
			sendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("bad request: %s", err))
			return
		}

		baseQuery, err := prometheus.NewQuery(conf, queryEndpoint)
		if err != nil {
			sendErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		promURL := baseQuery.BuildPrometheusURL(queryParams)

		var response PrometheusResponse
		err = QueryPrometheus(promURL, &response)
		if err != nil {
			sendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to query Prometheus: %v %s", err, promURL))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			sendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to encode response: %s", err))
			return
		}
	}
}

// validates input params for aggregate query
func processAggregateParams(queryParams map[string]string) (map[string]string, error) {

	_, valid := supportedOperations[queryParams["operation"]]
	if !valid {
		return nil, fmt.Errorf("unsupported operation: %s", queryParams["operation"])
	}

	queryParams["operation"] = supportedOperations[queryParams["operation"]]

	aggregateQuery := fmt.Sprintf("%s(%s[%s])", queryParams["operation"], queryParams["query"], queryParams["window"])

	return map[string]string{
		"query": aggregateQuery,
	}, nil
}
