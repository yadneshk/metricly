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

		queryBuilder := prometheus.QueryBuilder{
			BaseURL: fmt.Sprintf("%s:%s", conf.Prometheus.Address, conf.Prometheus.Port),
		}

		query, err := queryBuilder.BuildQueryRange(metricName, start, end, step)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to build aggregate query %v", err), http.StatusBadRequest)
			return
		}

		promURL, err := queryBuilder.BuildPrometheusURL(query, "query_range")
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
