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

type SQSClient struct {
	client *sqs.Client
	log    *logging.Logger
}

// NewSQSClient initializes and returns an AWS SQS client.
func NewSQSClient(ctx context.Context, log *logging.Logger) (*SQSClient, error) {
	configAws, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.HandleError("E", "Failed to load AWS configuration", err)
		return nil, err
	}

	clienteSQS := sqs.NewFromConfig(configAws)
	return &SQSClient{
		client: clienteSQS,
		log:    log,
	}, nil
}

// SendPedido sends a Pedido object to the specified SQS queue.
func (s *SQSClient) SendPedido(ctx context.Context, queueURL string, pedido models.Pedido) error {
	pedidoJSON, err := json.Marshal(pedido)
	if err != nil {
		s.log.HandleError("E", "Failed to marshal pedido", err)
		return err
	}

	input := &sqs.SendMessageInput{
		MessageBody: aws.String(string(pedidoJSON)),
		QueueUrl:    aws.String(queueURL),
	}

	_, err = s.client.SendMessage(ctx, input)
	if err != nil {
		s.log.HandleError("E", "Failed to send pedido", err)
	}

	return err
}
