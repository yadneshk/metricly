package prometheus

import (
	"fmt"
	"metricly/config"
	"net/url"
)

type Query struct {
	Scheme            string
	PrometheusAddress string
	PrometheusPort    string
	Path              string
	Params            map[string]string
}

func NewQuery(config *config.Config, endpoint string) (*Query, error) {
	if config.Prometheus.Address == "" || config.Prometheus.Port == "" {
		return nil, fmt.Errorf("prometheus address and port must be configured")
	}
	return &Query{
		Scheme:            "http",
		PrometheusAddress: config.Prometheus.Address,
		PrometheusPort:    config.Prometheus.Port,
		Path:              fmt.Sprintf("/api/v1/%s", endpoint),
	}, nil
}

func (qb *Query) BuildPrometheusURL(queryParams map[string]string) string {

	queryURL := &url.URL{
		Scheme: qb.Scheme,
		Host:   fmt.Sprintf("%s:%s", qb.PrometheusAddress, qb.PrometheusPort),
		Path:   qb.Path,
	}
	baseURL := queryURL.Query()
	for qry, val := range queryParams {
		baseURL.Set(qry, val)
	}

	queryURL.RawQuery = baseURL.Encode()
	return queryURL.String()

}
