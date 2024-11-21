/*

Sample response from Prometheus
http://localhost:9090/api/v1/query?query=cpu_steal
{
  "status": "success",
  "data": {
    "resultType": "vector",
    "result": [
      {
        "metric": {
          "__name__": "cpu_steal",
          "hostname": "mynode",
          "instance": "127.0.0.1:8080",
          "job": "metricly"
        },
        "value": [1732165754.764, "0"
        ]
      }
    ]
  }
}

Implements endpoints:
/api/v1/query?query=cpu_total&timestamp=2024-11-21T09:18:00.001Z


*/

package v1

import (
	"encoding/json"
	"metricly/config"
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

func PrometheusQueryHandler(conf *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		input_query := r.URL.Query().Get("query")

		if input_query == "" {
			http.Error(w, "Empty 'query' param", http.StatusBadRequest)
			return
		}

		queryParams := make(map[string]string)
		queryParams["query"] = input_query
		timestamp := r.URL.Query().Get("timestamp")
		if timestamp != "" {
			queryParams["time"] = timestamp
		}

		promQuery := PreparePromQuery(conf, "query", queryParams)
		var response PrometheusQueryResponse
		_ = QueryPrometheus(promQuery, &response)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
