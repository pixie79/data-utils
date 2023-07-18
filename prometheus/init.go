// Description: Prometheus utils to help parse prometheus metrics
// Author: Pixie79
// ============================================================================
// package prometheus

package prometheus

import (
	"fmt"
	"regexp"

	data_utils "github.com/pixie79/data-utils"
	"github.com/pixie79/data-utils/aws"
	"github.com/pixie79/data-utils/utils"
)

var (
	metricsServer         string                                          // metricsServer is the hostname of the Prometheus metrics endpoint
	metricsSecure         string                                          // metricsSecure is a flag to indicate whether the Prometheus metrics endpoint is secure
	metricsProtocol       string                                          // metricsProtocol is the protocol to use to access the Prometheus metrics endpoint
	metricsUrlPath        string                                          // metricsUrlPath is the path to the Prometheus metrics endpoint
	MetricsUrl            string                                          // MetricsUrl is the URL to the Prometheus metrics endpoint
	MetricsCredentials    data_utils.CredentialsType                      // MetricsCredentials is the credentials to access the Prometheus metrics endpoint
	metricsCredentialsKey string                                          // metricsCredentialsKey is the key to use to fetch the credentials from AWS Secrets Manager
	awsSecretsManager     string                                          // awsSecretsManager is a flag to indicate whether we're using AWS Secrets Manager
	reInitialSplit        = regexp.MustCompile(`(.+){(.+)}\s(\d+\.?\d*)`) // reInitialSplit is the regex to split the initial line of a Prometheus metric
	reTagSplit            = regexp.MustCompile(`(\w+)="(.+)"`)            // reTagSplit is the regex to split the tags of a Prometheus metric
)

// init is called before main(), sets up the environment
func init() {
	metricsServer = utils.GetEnv("METRICS_SERVER", "localhost:8080")
	metricsSecure = utils.GetEnv("METRICS_SECURE", "true")
	metricsUrlPath = utils.GetEnv("METRICS_URL_PATH", "/api/cloud/prometheus/public_metrics")
	if metricsSecure == "true" {
		metricsProtocol = "https://"
	} else {
		metricsProtocol = "http://"
	}

	MetricsUrl = fmt.Sprintf("%s%s%s", metricsProtocol, metricsServer, metricsUrlPath)

	awsSecretsManager = utils.GetEnv("AWS_SECRETS_MANAGER", "false")
	// If we're using AWS Secrets Manager, fetch the credentials from there
	if awsSecretsManager == "true" {
		metricsCredentialsKey = utils.GetEnv("METRICS_CREDENTIALS_KEY", "metrics-proxy")
		MetricsCredentials = aws.FetchCredentials(metricsCredentialsKey)
	} else {
		// Otherwise, use the credentials from the environment
		MetricsCredentials = data_utils.CredentialsType{
			Username: utils.GetEnv("METRICS_USERNAME", ""),
			Password: utils.GetEnv("METRICS_PASSWORD", ""),
		}

	}
}
