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

// dynamoDBClient is a concrete implementation of the DBClient.
type DynamoDBClient struct {
	Client *dynamodb.Client
}

// ConfigAws initializes and returns a new DynamoDB client as a dynamoDBClient.
func ConfigAws(ctx context.Context) (*DynamoDBClient, error) {
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
	client := &DynamoDBClient{Client: dynamodb.NewFromConfig(configAws)}
	return client, nil
}

// DeleteItem implements query.UserRepository.
func (client *DynamoDBClient) DeleteItem(ctx context.Context, input *dynamodb.DeleteItemInput, opts ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error) {
	options := &dynamodb.Options{}
	for _, o := range opts {
		o(options)
	}

	return client.Client.DeleteItem(ctx, input, func(o *dynamodb.Options) {
		*o = *options
	})
}

// GetItem implements query.UserRepository.
func (client *DynamoDBClient) GetItem(ctx context.Context, input *dynamodb.GetItemInput, opts ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
	options := &dynamodb.Options{}
	for _, o := range opts {
		o(options)
	}

	return client.Client.GetItem(ctx, input, func(o *dynamodb.Options) {
		*o = *options
	})
}

// PutItem implements query.UserRepository.
func (client *DynamoDBClient) PutItem(ctx context.Context, input *dynamodb.PutItemInput, opts ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	options := &dynamodb.Options{}
	for _, o := range opts {
		o(options)
	}

	return client.Client.PutItem(ctx, input, func(o *dynamodb.Options) {
		*o = *options
	})
}

// Scan implements query.UserRepository.
func (client *DynamoDBClient) Scan(ctx context.Context, input *dynamodb.ScanInput, opts ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error) {
	options := &dynamodb.Options{}
	for _, o := range opts {
		o(options)
	}

	return client.Client.Scan(ctx, input, func(o *dynamodb.Options) {
		*o = *options
	})
}
