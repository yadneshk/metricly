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
/api/v1/query?query=cpu_total&time=2024-11-21T09:18:00.001Z
/api/v1/query_range?query=up&start=2015-07-01T20:10:30.781Z&end=2015-07-01T20:11:00.781Z&step=15s

*/

package v1

import (
	"encoding/json"
	"net/http"
)

func PrometheusQueryHandler() http.HandlerFunc {
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
		response, _ := QueryPrometheus(PreparePromQuery(queryParams))
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
