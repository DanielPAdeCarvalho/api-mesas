package query

import (
	"context"
	"mesas-api/logging"
	"mesas-api/models"
	"sort"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// SelectAllMesas retorna todas as mesas que estao ocupadas
func SelectAllMesas(dynamoClient *dynamodb.Client, log logging.Logfile) []models.Mesa {
	input := &dynamodb.ScanInput{
		TableName: aws.String("Mesas"),
	}
	result, err := dynamoClient.Scan(context.Background(), input)
	logging.Check(err, log)
	var mesas []models.Mesa
	for _, item := range result.Items {
		var mesa models.Mesa
		err = attributevalue.UnmarshalMap(item, &mesa)
		logging.Check(err, log)

		// Add the Mesa struct to the slice if the Cliente field is not empty
		if mesa.Cliente != "" {
			mesas = append(mesas, mesa)
		}
	}
	return mesas
}

// SelectMesa returns a Mesa struct with the data from the DynamoDB table Mesas
func SelectMesa(mesaId string, dynamoClient *dynamodb.Client, log logging.Logfile) models.Mesa {
	key, err := attributevalue.MarshalMap(map[string]interface{}{
		"Id": mesaId,
	})
	logging.Check(err, log)

	input := &dynamodb.GetItemInput{
		TableName: aws.String("Mesas"),
		Key:       key,
	}
	result, err := dynamoClient.GetItem(context.Background(), input)
	logging.Check(err, log)

	var mesa models.Mesa
	err = attributevalue.UnmarshalMap(result.Item, &mesa)
	logging.Check(err, log)

	return mesa
}

// SelectCardapio returns a slice of Cardapio structs with the data from the DynamoDB table Cardapio
func SelectCardapio(dynamoClient *dynamodb.Client, log logging.Logfile) []models.Cardapio {
	input := &dynamodb.ScanInput{
		TableName: aws.String("Cardapio"),
	}
	result, err := dynamoClient.Scan(context.Background(), input)
	logging.Check(err, log)
	var cardapio []models.Cardapio
	for _, item := range result.Items {
		var itemCardapio models.Cardapio
		err = attributevalue.UnmarshalMap(item, &itemCardapio)
		logging.Check(err, log)
		cardapio = append(cardapio, itemCardapio)
	}
	sort.Slice(cardapio, func(i, j int) bool {
		return cardapio[i].Nome < cardapio[j].Nome
	})
	return cardapio
}

// UpdateMesa updates the data from a Mesa struct in the DynamoDB table Mesas
func UpdateMesa(mesa models.Mesa, dynamoClient *dynamodb.Client, log logging.Logfile) {
	item, err := attributevalue.MarshalMap(mesa)
	logging.Check(err, log)

	input := &dynamodb.PutItemInput{
		TableName: aws.String("Mesas"),
		Item:      item,
	}
	_, err = dynamoClient.PutItem(context.Background(), input)
	logging.Check(err, log)
}

// DeletePedido deletes a Pedido struct from the DynamoDB table Pedidos
