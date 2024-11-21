package prometheus

import (
	"errors"
	"fmt"
	"net/url"
)

type QueryBuilder struct {
	BaseURL string
}

// only builds the prometheus query, doesn't prepend the prometheus server address
// promQL equivalent - cpu_total
func (qb *QueryBuilder) BuildQuery(metricName, time string) (map[string]string, error) {
	if metricName == "" {
		return nil, fmt.Errorf("metric cannot be empty")
	}

	params := map[string]string{
		"query": metricName,
	}
	// query := fmt.Sprintf("query=%s", metricName)
	if time != "" {
		params["time"] = time
	}
	return params, nil

}

// only builds the prometheus query, doesn't prepend the prometheus server address
// promQL equivalent - cpu_total
func (qb *QueryBuilder) BuildQueryRange(metricName, start, end, step string) (map[string]string, error) {
	if metricName == "" || start == "" || end == "" || step == "" {
		return nil, fmt.Errorf("metric, start, end, step cannot be empty")
	}

	return map[string]string{
		"query": metricName,
		"start": start,
		"end":   end,
		"step":  step,
	}, nil

}

// only builds the prometheus query, doesn't prepend the prometheus server address
// promQL equivalent - avg_over_time(cpu_total[1h])
func (qb *QueryBuilder) BuildAggregateQuery(metricName, operation, window string) (map[string]string, error) {
	if metricName == "" || operation == "" || window == "" {
		return nil, fmt.Errorf("metric, operation and window, all required to aggregate metrics")
	}

	supportedOperations := map[string]string{
		"avg": "avg_over_time",
		"max": "max_over_time",
		"min": "min_over_time",
	}

	opr, valid := supportedOperations[operation]
	if !valid {
		return nil, fmt.Errorf("unsupported operation: %s", operation)
	}

	return map[string]string{
		"query": fmt.Sprintf("%s(%s[%s])", opr, metricName, window),
	}, nil

}

func (qb *QueryBuilder) BuildPrometheusURL(query map[string]string, exprQuery string) (string, error) {
	if query == nil {
		return "", errors.New("query required")
	}

	queryParams := url.Values{}
	for key, value := range query {
		queryParams.Set(key, value)
	}

	return fmt.Sprintf("http://%s/api/v1/%s?%s", qb.BaseURL, exprQuery, queryParams.Encode()), nil

}
