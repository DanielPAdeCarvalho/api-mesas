package events

import (
	"context"
	"encoding/json"
	"fmt"
	"mesas-api/pkg/logging"
	"mesas-api/pkg/models"
	"mesas-api/pkg/utils"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type SQSClient struct {
	client *sqs.Client
	log    *logging.Logger
}

const queueURLSendPedido = "https://sqs.us-east-1.amazonaws.com/912225062963/PedidosCozinha"
const queueURLPedidoPronto = "https://sqs.us-east-1.amazonaws.com/912225062963/PedidoPronto"

// NewSQSClient initializes and returns an AWS SQS client.
func NewSQSClient(ctx context.Context, log *logging.Logger) (*SQSClient, error) {
	configAws := utils.ConfigAws(ctx, log)

	clienteSQS := sqs.NewFromConfig(configAws)
	return &SQSClient{
		client: clienteSQS,
		log:    log,
	}, nil
}

// SendPedido sends a Pedido object to the specified SQS queue.
func (s *SQSClient) SendPedido(ctx context.Context, item models.PedidoCozinha) error {
	itemJSON, err := json.Marshal(item)
	if err != nil {
		s.log.HandleError("E", "Failed to marshal pedido", err)
		return err
	}

	input := &sqs.SendMessageInput{
		MessageBody: aws.String(string(itemJSON)),
		QueueUrl:    aws.String(queueURLSendPedido),
	}

	_, err = s.client.SendMessage(ctx, input)
	if err != nil {
		s.log.HandleError("E", "Failed to send pedido", err)
	}
	return err
}

// SendPedido sends a Pedido object to the specified SQS queue.
func (s *SQSClient) PedidoPronto(ctx context.Context, nome string) error {
	input := &sqs.SendMessageInput{
		MessageBody: aws.String(nome),
		QueueUrl:    aws.String(queueURLPedidoPronto),
	}

	_, err := s.client.SendMessage(ctx, input)
	if err != nil {
		s.log.HandleError("E", "Failed to send pedido pronto", err)
	}
	fmt.Printf("Enviou o item para a lambda stock: %s\n", nome)
	return err
}
