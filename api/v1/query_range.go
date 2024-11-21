/*
Implements endpoint
/api/v1/query_range?query=cpu_total&start=2024-11-21T09:18:00.001Z&end=2024-11-21T10:18:00.001Z&step=15s
*/
package v1

import (
	"encoding/json"
	"metricly/config"
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
		input_query := r.URL.Query().Get("query")
		if input_query == "" {
			http.Error(w, "Empty 'query' param", http.StatusBadRequest)
			return
		}
		queryParams := make(map[string]string)
		queryParams["query"] = input_query

		start := r.URL.Query().Get("start")
		end := r.URL.Query().Get("end")
		step := r.URL.Query().Get("step")
		if start == "" || end == "" || step == "" {
			http.Error(w, "Missing one or more params - start,end,step", http.StatusBadRequest)
			return
		}
		queryParams["start"] = start
		queryParams["end"] = end
		queryParams["step"] = step

		promQuery := PreparePromQuery(conf, "query_range", queryParams)
		var response PrometheusQueryRangeResponse
		_ = QueryPrometheus(promQuery, &response)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
