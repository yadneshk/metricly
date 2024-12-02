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

		metricName := requestParams.Get("metric")
		operation := requestParams.Get("operation")
		window := requestParams.Get("window")

		if err := validateAggregateParams(metricName, operation, window); err != nil {
			http.Error(w, fmt.Sprintf("bad request: %s", err), http.StatusInternalServerError)
			return
		}

		baseQuery, _ := prometheus.NewQuery(conf, queryEndpoint)

		aggregateQuery := fmt.Sprintf("%s(%s[%s])", supportedOperations[operation], metricName, window)
		queryParams := map[string]string{
			"query": aggregateQuery,
		}

		promURL, err := baseQuery.BuildPrometheusURL(queryParams)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to build prometheus query %v", err), http.StatusBadRequest)
			return
		}

		var response PrometheusQueryResponse
		err = QueryPrometheus(promURL, &response)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to query Prometheus: %v %s", err, promURL), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to encode response: %s", err), http.StatusInternalServerError)
			return
		}
	}
}

// validates input params for aggregate query
func validateAggregateParams(metricName, operation, window string) error {
	if metricName == "" || operation == "" || window == "" {
		return fmt.Errorf("metric, operation and window, all required to aggregate metrics")
	}

	_, valid := supportedOperations[operation]
	if !valid {
		return fmt.Errorf("unsupported operation: %s", operation)
	}

	return nil
}
