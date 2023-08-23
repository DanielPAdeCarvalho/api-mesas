package driver

import (
	"context"
	"mesas-api/pkg/logging"
	"mesas-api/pkg/utils"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// dynamoDBClient is a concrete implementation of the DBClient.
type DynamoDBClient struct {
	Client *dynamodb.Client
}

// NewDynamoClient initializes and returns a new DynamoDB client as a dynamoDBClient.
func NewDynamoClient(ctx context.Context, log *logging.Logger) (*DynamoDBClient, error) {
	configAWS := utils.ConfigAws(ctx, log)
	// Create and return the DynamoDB client
	client := &DynamoDBClient{Client: dynamodb.NewFromConfig(configAWS)}
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
