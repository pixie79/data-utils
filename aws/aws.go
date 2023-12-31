// Description: AWS utils
// Author: Pixie79
// ============================================================================
// package aws

package aws

import (
	"encoding/json"
	"fmt"

	"github.com/pixie79/data-utils/types"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/ssm"
	tuUtils "github.com/pixie79/tiny-utils/utils"
)

// FetchCredentials retrieves credentials from AWS Secrets Manager
func FetchCredentials(credentialsKey string) types.CredentialsType {
	credentialsString := GetSecretManagerValue(credentialsKey)
	credentials := types.CredentialsType{}
	err := json.Unmarshal([]byte(credentialsString), &credentials)
	tuUtils.MaybeDie(err, "could not explode credentials")
	tuUtils.Print("DEBUG", fmt.Sprintf("credentials retrieved: %s", credentialsKey))
	return credentials
}

// GetSecretManagerValue retrieves a secret from AWS Secrets Manager
func GetSecretManagerValue(passwordKey string) string {
	region := tuUtils.GetEnvDefault("AWS_REGION", "eu-west-1")
	tuUtils.Print("DEBUG", "creating AWS Session")
	sess, err := session.NewSession(&aws.Config{Region: aws.String(region)})
	tuUtils.MaybeDie(err, "could not connect to AWS")

	// Create Secrets Manager service
	smSvc := secretsmanager.New(sess, &aws.Config{
		Region:     aws.String(region),
		MaxRetries: aws.Int(3),
		Endpoint:   aws.String(fmt.Sprintf("https://secretsmanager.%s.amazonaws.com", region)),
	})

	// Get the secret
	resp, err := smSvc.GetSecretValue(&secretsmanager.GetSecretValueInput{
		SecretId: aws.String(passwordKey),
	})
	tuUtils.MaybeDie(err, fmt.Sprintf("failed to retrieve secret called %s ", passwordKey))

	tuUtils.Print("DEBUG", "Returning secret to calling function")
	result := *resp.SecretString
	return result
}

// GetSsmParam retrieves a parameter from AWS SSM Parameter Store
func GetSsmParam(parameterPath string) string {
	region := tuUtils.GetEnvDefault("AWS_REGION", "eu-west-1")
	tuUtils.Print("DEBUG", "creating AWS Session")
	sess, err := session.NewSession(&aws.Config{Region: aws.String(region)})
	tuUtils.MaybeDie(err, "could not connect to AWS")

	// Create SSM service
	ssmSvc := ssm.New(sess, &aws.Config{
		Region:     aws.String(region),
		MaxRetries: aws.Int(3),
		Endpoint:   aws.String(fmt.Sprintf("https://ssm.%s.amazonaws.com", region)),
	})

	// Get the parameter
	param, err := ssmSvc.GetParameter(&ssm.GetParameterInput{
		Name:           aws.String(parameterPath),
		WithDecryption: aws.Bool(true),
	})
	tuUtils.MaybeDie(err, fmt.Sprintf("failed to get Parameter from store: %+v", err))
	value := *param.Parameter.Value
	tuUtils.Print("DEBUG", fmt.Sprintf("parameter retrieved: %s", parameterPath))
	return value
}

// CreateCloudwatchMetric creates a metric in AWS Cloudwatch
func CreateCloudwatchMetric(metric []*cloudwatch.MetricDatum, namespace string) {
	region := tuUtils.GetEnvDefault("AWS_REGION", "eu-west-1")
	sess, err := session.NewSession(&aws.Config{Region: aws.String(region)})
	tuUtils.MaybeDie(err, "could not connect to AWS")

	// Create Cloudwatch service
	cwmSvc := cloudwatch.New(sess, &aws.Config{
		Region:              aws.String(region),
		MaxRetries:          aws.Int(3),
		Endpoint:            aws.String(fmt.Sprintf("https://events.%s.amazonaws.com", region)),
		STSRegionalEndpoint: endpoints.RegionalSTSEndpoint,
	})
	// Create the cloudwatch metric
	_, err = cwmSvc.PutMetricData(&cloudwatch.PutMetricDataInput{
		Namespace:  aws.String(namespace),
		MetricData: metric,
	})
	tuUtils.MaybeDie(err, fmt.Sprintf("failed to create cloudwatch metric: %+v", err))
	tuUtils.Print("DEBUG", fmt.Sprintf("cloudwatch metric created: %s", namespace))
}
