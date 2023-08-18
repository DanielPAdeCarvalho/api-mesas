package driver

import (
	"context"
	"errors"
	"mesas-api/pkg/logging"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// dynamoDBClient is a concrete implementation of the DBClient interface.
type dynamoDBClient struct {
	client *dynamodb.Client
}

// ConfigAws initializes and returns a new DynamoDB client as a dynamoDBClient.
func ConfigAws(ctx context.Context) (*dynamoDBClient, error) {
	log := logging.NewLogger(logrus.New().Level) // Define logging level as needed

	// Initialize the viper with configuration
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".") // or path to your config file

	if err := v.ReadInConfig(); err != nil {
		log.HandleError("E", "failed to read the configuration file", err)
		return nil, err
	}

	credentialsFilePath := v.GetString("credentialsFilePath")
	configFilePath := v.GetString("configFilePath")

	// Load the AWS configuration
	configAws, err := config.LoadDefaultConfig(
		ctx,
		config.WithSharedCredentialsFiles([]string{credentialsFilePath}),
		config.WithSharedConfigFiles([]string{configFilePath}),
	)
	if err != nil {
		log.HandleError("E", "failed to load AWS configuration", err)
		return nil, errors.New("failed to load AWS configuration: " + err.Error())
	}

	// Create and return the DynamoDB client
	client := &dynamoDBClient{client: dynamodb.NewFromConfig(configAws)}
	return client, nil
}
