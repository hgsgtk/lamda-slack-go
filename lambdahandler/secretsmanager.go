package lambdahandler

// Use this code snippet in your app.
// If you need more information about configurations or implementing the sample code, visit the AWS docs:
// https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/setting-up.html

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/pkg/errors"
	"os"
)

type DBConfig struct {
	Username            string `json:"username"`
	Password            string `json:"password"`
	Engine              string `json:"engine"`
	Host                string `json:"host"`
	Port                int    `json:"port"`
	DbClusterIdentifier string `json:"dbClusterIdentifier"`
}

func GetSecret() (DBConfig, error) {
	secretName := os.Getenv("DB_SECRET_NAME")

	// Create a Secrets Manager client
	// FIXME session.New() is deprecated
	svc := secretsmanager.New(session.New())
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}

	// In this sample we only handle the specific exceptions for the 'GetSecretValue' API.
	// See https://docs.aws.amazon.com/secretsmanager/latest/apireference/API_GetSecretValue.html

	result, err := svc.GetSecretValue(input)
	if err != nil {
		return DBConfig{}, err
	}

	// Decrypts secret using the associated KMS CMK.
	// Depending on whether the secret is a string or binary, one of these fields will be populated.
	// TODO secretBinaryってどういうときにくるの
	var secretString string
	if result.SecretString == nil {
		return DBConfig{}, errors.New("secret string is empty")
	}
	secretString = *result.SecretString

	c := DBConfig{}
	if err := json.Unmarshal([]byte(secretString), &c); err != nil {
		return DBConfig{}, errors.Wrap(err, "json decode error")
	}

	return c, nil
}
