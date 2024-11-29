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

		requestParams := r.URL.Query()

		metricName := requestParams.Get("metric")
		time := requestParams.Get("timestamp")

		queryBuilder := prometheus.QueryBuilder{
			BaseURL: fmt.Sprintf("%s:%s", conf.Prometheus.Address, conf.Prometheus.Port),
		}

		query, err := queryBuilder.BuildQuery(metricName, time)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to build aggregate query %v", err), http.StatusBadRequest)
			return
		}

		promURL, err := queryBuilder.BuildPrometheusURL(query, "query")
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to build Prometheus query %v", err), http.StatusBadRequest)
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
