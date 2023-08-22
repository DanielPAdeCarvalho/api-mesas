package events

import (
	"context"
	"encoding/json"

	"mesas-api/pkg/logging"
	"mesas-api/pkg/models"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

// CreateClient initializes and returns an AWS SQS client.
func CreateClient(ctx context.Context, log *logging.Logger) (*sqs.Client, error) {
	// Load configuration from the environment (including IAM roles, if applicable)
	configAws, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.HandleError("E", "Failed to load AWS configuration", err)
		return nil, err
	}

	clienteSQS := sqs.NewFromConfig(configAws)
	return clienteSQS, nil
}

// EnviaPedido sends a Pedido object to the specified SQS queue.
func EnviaPedido(ctx context.Context, clienteSQS *sqs.Client, pedido models.Pedido, log *logging.Logger) error {
	pedidoJSON, err := json.Marshal(pedido)
	if err != nil {
		log.HandleError("E", "Failed to marshal pedido", err)
		return err
	}
	input := &sqs.SendMessageInput{
		MessageBody: aws.String(string(pedidoJSON)),
		QueueUrl:    aws.String(filas.FilaPedidos),
	}
	_, err = clienteSQS.SendMessage(ctx, input)
	if err != nil {
		log.HandleError("E", "Failed to send pedido", err)
	}
	return err
}
