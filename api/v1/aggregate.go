package v1

import (
	"encoding/json"
	"fmt"
	"metricly/config"
	"metricly/pkg/prometheus"
	"net/http"
)

func PrometheusAggregateHandler(conf *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		requestParams := r.URL.Query()

		metricName := requestParams.Get("metric")
		operation := requestParams.Get("operation")
		window := requestParams.Get("window")

		queryBuilder := prometheus.QueryBuilder{
			BaseURL: fmt.Sprintf("%s:%s", conf.Prometheus.Address, conf.Prometheus.Port),
		}

		query, err := queryBuilder.BuildAggregateQuery(metricName, operation, window)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to build aggregate query %v", err), http.StatusBadRequest)
			return
		}

		promURL, err := queryBuilder.BuildPrometheusURL(query, "query")
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
		json.NewEncoder(w).Encode(response)
	}
}
