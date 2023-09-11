// Description: Prometheus utils to help parse prometheus metrics
// Author: Pixie79
// ============================================================================
// package prometheus

package prometheus

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/pixie79/data-utils/types"

	"github.com/pixie79/data-utils/aws"
	"github.com/pixie79/data-utils/utils"
)

var (
	metricsServer         string                                          // metricsServer is the hostname of the Prometheus metrics endpoint
	metricsSecure         string                                          // metricsSecure is a flag to indicate whether the Prometheus metrics endpoint is secure
	metricsProtocol       string                                          // metricsProtocol is the protocol to use to access the Prometheus metrics endpoint
	metricsUrlPath        string                                          // metricsUrlPath is the path to the Prometheus metrics endpoint
	MetricsUrl            string                                          // MetricsUrl is the URL to the Prometheus metrics endpoint
	MetricsCredentials    types.CredentialsType                           // MetricsCredentials is the credentials to access the Prometheus metrics endpoint
	metricsCredentialsKey string                                          // metricsCredentialsKey is the key to use to fetch the credentials from AWS Secrets Manager
	awsSecretsManager     string                                          // awsSecretsManager is a flag to indicate whether we're using AWS Secrets Manager
	reInitialSplit        = regexp.MustCompile(`(.+){(.+)}\s(\d+\.?\d*)`) // reInitialSplit is the regex to split the initial line of a Prometheus metric
	reTagSplit            = regexp.MustCompile(`(\w+)="(.+)"`)            // reTagSplit is the regex to split the tags of a Prometheus metric
	additionalTags        []types.TagsType
)

// init is called before main(), sets up the environment
func init() {
	metricsServer = utils.GetEnvDefault("METRICS_SERVER", "localhost:8080")
	metricsSecure = utils.GetEnvDefault("METRICS_SECURE", "true")
	metricsUrlPath = utils.GetEnvDefault("METRICS_URL_PATH", "/api/cloud/prometheus/public_metrics")
	if metricsSecure == "true" {
		metricsProtocol = "https://"
	} else {
		metricsProtocol = "http://"
	}

	MetricsUrl = fmt.Sprintf("%s%s%s", metricsProtocol, metricsServer, metricsUrlPath)

	awsSecretsManager = utils.GetEnvDefault("AWS_SECRETS_MANAGER", "false")
	// If we're using AWS Secrets Manager, fetch the credentials from there
	if awsSecretsManager == "true" {
		metricsCredentialsKey = utils.GetEnvDefault("METRICS_CREDENTIALS_KEY", "metrics-proxy")
		MetricsCredentials = aws.FetchCredentials(metricsCredentialsKey)
	} else {
		// Otherwise, use the credentials from the environment
		MetricsCredentials = types.CredentialsType{
			Username: utils.GetEnvDefault("METRICS_USERNAME", ""),
			Password: utils.GetEnvDefault("METRICS_PASSWORD", ""),
		}

	}

	additionalTags = []types.TagsType{
		{
			Name:  "environment",
			Value: utils.Environment,
		},
		{
			Name:  "metrics_server",
			Value: strings.Split(metricsServer, ":")[0],
		},
	}
}
