package prometheus

import (
	"fmt"
	"github.com/pixie79/data-utils"
	"github.com/pixie79/data-utils/aws"
	"github.com/pixie79/data-utils/utils"
	"regexp"
)

var (
	metricsServer         string
	metricsSecure         string
	metricsProtocol       string
	metricsUrlPath        string
	MetricsUrl            string
	MetricsCredentials    data_utils.CredentialsType
	metricsCredentialsKey string
	awsSecretsManager     string
	reInitialSplit        = regexp.MustCompile(`(.+){(.+)}\s(\d+\.?\d*)`)
	reTagSplit            = regexp.MustCompile(`(\w+)="(.+)"`)
)

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

	awsSecretsManager = utils.GetEnv("AWS_SECRETS_MANAGER", "true")
	if awsSecretsManager == "true" {
		metricsCredentialsKey = utils.GetEnv("METRICS_CREDENTIALS_KEY", "metrics-proxy")
		MetricsCredentials = aws.FetchCredentials(metricsCredentialsKey)
	} else {
		MetricsCredentials = data_utils.CredentialsType{
			Username: utils.GetEnv("METRICS_USERNAME", ""),
			Password: utils.GetEnv("METRICS_PASSWORD", ""),
		}

	}
}
