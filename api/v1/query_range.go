/*
Implements endpoint
/api/v1/query_range?query=cpu_total&start=2024-11-21T09:18:00.001Z&end=2024-11-21T10:18:00.001Z&step=15s
*/
package v1

import (
	"encoding/json"
	"fmt"
	"metricly/config"
	"metricly/pkg/prometheus"
	"net/http"
)

var (
	queryRangeEndpoint = "query_range"
)

// represents the structure of the reponse Prometheus's API call
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

func prometheusQueryRangeHandler(conf *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		requestParams := r.URL.Query()

		metricName := requestParams.Get("metric")
		start := requestParams.Get("start")
		end := requestParams.Get("end")
		step := requestParams.Get("step")

		if err := validateQueryParams(metricName, start, end, step); err != nil {
			http.Error(w, fmt.Sprintf("bad request: %s", err), http.StatusBadRequest)
			return
		}

		baseQuery, _ := prometheus.NewQuery(conf, queryRangeEndpoint)

		queryParams := map[string]string{
			"query": metricName,
			"start": start,
			"end":   end,
			"step":  step,
		}

		promURL, err := baseQuery.BuildPrometheusURL(queryParams)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to build prometheus query %v", err), http.StatusBadRequest)
			return
		}

		var response PrometheusQueryRangeResponse
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
func validateQueryParams(metricName, start, end, step string) error {
	if metricName == "" || start == "" || end == "" || step == "" {
		return fmt.Errorf("metric, start, end and step, all required to get range metrics")
	}

	return nil
}
