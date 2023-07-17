package aws

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/ssm"
	data_utils "github.com/pixie79/data-utils"
	"github.com/pixie79/data-utils/utils"
)

func FetchCredentials(credentialsKey string) data_utils.CredentialsType {
	credentialsString := GetSecretManagerValue(credentialsKey)
	credentials := data_utils.CredentialsType{}
	err := json.Unmarshal([]byte(credentialsString), &credentials)
	utils.MaybeDie(err, "could not explode credentials")
	utils.Logger.Debug(fmt.Sprintf("credentials retrieved: %s", credentialsKey))
	return credentials
}

func GetSecretManagerValue(passwordKey string) string {
	region := utils.GetEnv("AWS_REGION", "eu-west-1")
	utils.Logger.Debug("creating AWS Session")
	sess, err := session.NewSession(&aws.Config{Region: aws.String(region)})
	utils.MaybeDie(err, "could not connect to AWS")

	smSvc := secretsmanager.New(sess, &aws.Config{
		Region:     aws.String(region),
		MaxRetries: aws.Int(3),
		Endpoint:   aws.String(fmt.Sprintf("https://secretsmanager.%s.amazonaws.com", region)),
	})

	resp, err := smSvc.GetSecretValue(&secretsmanager.GetSecretValueInput{
		SecretId: aws.String(passwordKey),
	})
	utils.MaybeDie(err, fmt.Sprintf("failed to retrieve secret called %s ", passwordKey))

	utils.Logger.Debug("Returning secret to calling function")
	result := *resp.SecretString
	return result
}

func GetSsmParam(parameterPath string) string {
	region := utils.GetEnv("AWS_REGION", "eu-west-1")
	utils.Logger.Debug("creating AWS Session")
	sess, err := session.NewSession(&aws.Config{Region: aws.String(region)})
	utils.MaybeDie(err, "could not connect to AWS")

	ssmSvc := ssm.New(sess, &aws.Config{
		Region:     aws.String(region),
		MaxRetries: aws.Int(3),
		Endpoint:   aws.String(fmt.Sprintf("https://ssm.%s.amazonaws.com", region)),
	})
	param, err := ssmSvc.GetParameter(&ssm.GetParameterInput{
		Name:           aws.String(parameterPath),
		WithDecryption: aws.Bool(true),
	})
	utils.MaybeDie(err, fmt.Sprintf("failed to get Parameter from store: %+v", err))
	value := *param.Parameter.Value
	utils.Logger.Debug(fmt.Sprintf("parameter retrieved: %s", parameterPath))
	return value
}

func CreateCloudwatchMetric(metric []*cloudwatch.MetricDatum, namespace string) {
	region := utils.GetEnv("AWS_REGION", "eu-west-1")
	sess, err := session.NewSession(&aws.Config{Region: aws.String(region)})
	utils.MaybeDie(err, "could not connect to AWS")

	cwmSvc := cloudwatch.New(sess, &aws.Config{
		Region:              aws.String(region),
		MaxRetries:          aws.Int(3),
		Endpoint:            aws.String(fmt.Sprintf("https://events.%s.amazonaws.com", region)),
		STSRegionalEndpoint: endpoints.RegionalSTSEndpoint,
	})

	_, err = cwmSvc.PutMetricData(&cloudwatch.PutMetricDataInput{
		Namespace:  aws.String(namespace),
		MetricData: metric,
	})
	utils.MaybeDie(err, fmt.Sprintf("failed to create cloudwatch metric: %+v", err))
	utils.Logger.Debug(fmt.Sprintf("cloudwatch metric created: %s", namespace))
}
