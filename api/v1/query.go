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
/api/v1/query?query=cpu_total
/api/v1/query?query=cpu_total&timestamp=2024-11-21T09:18:00.001Z


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
	queryEndpoint = "query"
)

func PrometheusQueryHandler(conf *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		requestParams := r.URL.Query()

		supportedParams := map[string]bool{
			"metric": true,
			"time":   true,
		}
		requiredParams := []string{"metric"}

		queryParams, err := processQueryParams(requestParams, supportedParams, requiredParams)
		if err != nil {
			sendErrorResponse(w, http.StatusBadRequest, err.Error())
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
