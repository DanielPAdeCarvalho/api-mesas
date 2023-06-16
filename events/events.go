package events

import (
	"context"
	"mesas-api/filas"
	"mesas-api/logging"
	"mesas-api/models"

	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func CreateClient(log logging.Logfile) *sqs.Client {
	configAws, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithSharedCredentialsFiles([]string{"driver/data/credentials.aws"}),
		config.WithSharedConfigFiles([]string{"driver/data/config.aws"}),
		config.WithSharedConfigProfile("sqs"),
	)
	logging.Check(err, log)

	clienteSQS := sqs.NewFromConfig(configAws)

	return clienteSQS
}

func EnviaPedido(clienteSQS *sqs.Client, log logging.Logfile, pedido models.Pedido) {
	pedidoJSON, err := json.Marshal(pedido)
	logging.Check(err, log)
	input := &sqs.SendMessageInput{
		MessageBody: aws.String(string(pedidoJSON)),
		QueueUrl:    aws.String(filas.FilaPedidos),
	}
	_, err = clienteSQS.SendMessage(context.Background(), input)
	logging.Check(err, log)
}
