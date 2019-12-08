package lambdahandler

// Use this code snippet in your app.
// If you need more information about configurations or implementing the sample code, visit the AWS docs:
// https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/setting-up.html

import (
	"encoding/json"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/pkg/errors"
)

type DBConfig struct {
	Username            string `json:"username"`
	Password            string `json:"password"`
	Engine              string `json:"engine"`
	Host                string `json:"host"`
	Port                int    `json:"port"`
	DbClusterIdentifier string `json:"dbClusterIdentifier"`
	Name                string
}

func GetDBConfig() (DBConfig, error) {
	secretName := os.Getenv("DB_SECRET_NAME")
	if secretName == "" {
		return DBConfig{}, errors.New("configuration DB_SECRET_NAME is missing")
	}

	secretString, err := getSecret(secretName)
	if err != nil {
		return DBConfig{}, errors.New("failed to get secret by secrets manager")
	}

	c := DBConfig{}
	if err := json.Unmarshal([]byte(secretString), &c); err != nil {
		return DBConfig{}, errors.Wrap(err, "json decode error")
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		return DBConfig{}, errors.New("configuration DB_NAME is missing")
	}
	c.Name = dbName

	return c, nil
}

const defaultVersion = "AWSCURRENT"

func getSecret(secretName string) (string, error) {
	// Create a Secrets Manager client
	// FIXME session.New() is deprecated
	svc := secretsmanager.New(session.New())
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String(defaultVersion), // VersionStage defaults to AWSCURRENT if unspecified
	}

	// In this sample we only handle the specific exceptions for the 'GetSecretValue' API.
	// See https://docs.aws.amazon.com/secretsmanager/latest/apireference/API_GetSecretValue.html

	result, err := svc.GetSecretValue(input)
	if err != nil {
		return "", err
	}

	// Decrypts secret using the associated KMS CMK.
	// Depending on whether the secret is a string or binary, one of these fields will be populated.
	// TODO secretBinaryってどういうときにくるの
	if result.SecretString == nil {
		return "", errors.New("secret string is empty")
	}
	return *result.SecretString, nil
}
