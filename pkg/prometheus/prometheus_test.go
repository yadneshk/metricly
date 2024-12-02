package prometheus

import (
	"metricly/config"
	"net/url"
	"testing"
)

func TestNewQuery(t *testing.T) {
	tests := []struct {
		name      string
		config    *config.Config
		endpoint  string
		expectErr bool
	}{
		{
			name: "Valid Config",
			config: &config.Config{
				Prometheus: struct {
					Address string "yaml:\"address\""
					Port    string "yaml:\"port\""
				}{
					"localhost", "1111",
				},
			},
			endpoint:  "query",
			expectErr: false,
		},
		{
			name: "Missing Prometheus Address",
			config: &config.Config{
				Prometheus: struct {
					Address string "yaml:\"address\""
					Port    string "yaml:\"port\""
				}{
					Address: "",
					Port:    "9090",
				},
			},
			endpoint:  "query",
			expectErr: true,
		},
		{
			name: "Missing Prometheus Port",
			config: &config.Config{
				Prometheus: struct {
					Address string "yaml:\"address\""
					Port    string "yaml:\"port\""
				}{
					Address: "localhost",
					Port:    "",
				},
			},
			endpoint:  "query",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q, err := NewQuery(tt.config, tt.endpoint)
			if tt.expectErr {
				if err == nil {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("did not expect error but got %v", err)
				}
				if q.Path != "/api/v1/"+tt.endpoint {
					t.Errorf("expected path to be /api/v1/%s, got %s", tt.endpoint, q.Path)
				}
			}
		})
	}
}

func TestBuildPrometheusURL(t *testing.T) {
	tests := []struct {
		name        string
		queryParams map[string]string
		expectedURL string
		expectErr   bool
	}{
		{
			name: "Valid Query Params",
			queryParams: map[string]string{
				"query": "cpu_usage",
				"step":  "10",
			},
			expectedURL: "http://localhost:9090/api/v1/query?query=cpu_usage&step=10",
			expectErr:   false,
		},
		{
			name:        "Empty Query Params",
			queryParams: map[string]string{},
			expectedURL: "http://localhost:9090/api/v1/query",
			expectErr:   false,
		},
	}

	qb := Query{
		Scheme:            "http",
		PrometheusAddress: "localhost",
		PrometheusPort:    "9090",
		Path:              "/api/v1/query",
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotURL, err := qb.BuildPrometheusURL(tt.queryParams)
			if tt.expectErr {
				if err == nil {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("did not expect error but got %v", err)
				}
				parsedGot, _ := url.Parse(gotURL)
				parsedExpected, _ := url.Parse(tt.expectedURL)
				if parsedGot.Scheme != parsedExpected.Scheme ||
					parsedGot.Host != parsedExpected.Host ||
					parsedGot.Path != parsedExpected.Path ||
					parsedGot.Query().Encode() != parsedExpected.Query().Encode() {
					t.Errorf("expected URL %s, got %s", tt.expectedURL, gotURL)
				}
			}
		})
	}
}
