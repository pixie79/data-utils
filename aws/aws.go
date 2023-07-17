package aws

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/pixie79/dataUtils/utils"
)

func FetchCredentials(credentialsKey string) CredentialsType {
	credentialsString := GetSecretManagerValue(credentialsKey)
	credentials := CredentialsType{}
	err := json.Unmarshal([]byte(credentialsString), &credentials)
	utils.MaybeDie(err, "could not explode credentials")
	utils.Logger.Debug("credentials retrieved")
	return credentials
}

func GetSecretManagerValue(passwordKey string) string {
	region := utils.GetEnv("AWS_REGION", "eu-west-1")
	utils.Logger.Debug("Creating AWS Session")
	sess, err := session.NewSession(&aws.Config{Region: aws.String(region)})
	utils.MaybeDie(err, "Could not connect to AWS")

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
