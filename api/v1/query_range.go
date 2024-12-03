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
	"time"
)

var (
	queryRangeEndpoint = "query_range"
)

func prometheusQueryRangeHandler(conf *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		requestParams := r.URL.Query()

		supportedParams := map[string]bool{
			"metric": true,
			"start":  true,
			"end":    true,
			"last":   true,
			"step":   true,
		}

		queryParams, err := processQueryParams(requestParams, supportedParams, []string{})
		if err != nil {
			sendErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		queryParams, err = processRangeParams(queryParams)
		if err != nil {
			sendErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		baseQuery, _ := prometheus.NewQuery(conf, queryRangeEndpoint)

		promURL := baseQuery.BuildPrometheusURL(queryParams)

		var response PrometheusQueryRangeResponse
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

func processRangeParams(requestParams map[string]string) (map[string]string, error) {

	if requestParams["last"] != "" {
		if requestParams["start"] != "" || requestParams["end"] != "" {
			return nil, fmt.Errorf("start, end and last cannot be used together, either use start and end or just last to get a range of datapoints")
		}

		duration, err := time.ParseDuration(requestParams["last"])
		if err != nil {
			return nil, fmt.Errorf("failed to parse time %s to epoch time", requestParams["last"])
		}
		end := time.Now().Unix()
		start := time.Now().Add(-duration).Unix()

		delete(requestParams, "last")
		requestParams["start"] = fmt.Sprint(start)
		requestParams["end"] = fmt.Sprint(end)
	} else if requestParams["start"] == "" || requestParams["end"] == "" {
		return nil, fmt.Errorf("start and end both required to get a range of datapoints")
	}

	if _, stepExists := requestParams["step"]; !stepExists {
		requestParams["step"] = "15s"
	}

	return requestParams, nil

}
