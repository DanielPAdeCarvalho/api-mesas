package query

import (
	"context"
	"mesas-api/pkg/logging"
	"mesas-api/pkg/models"
	"sort"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type DynamoDBClient interface {
	Scan(ctx context.Context, params *dynamodb.ScanInput, optFns ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error)
	GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
}

const tableMesas = "Mesas"
const tableCardapio = "Cardapio"

// SelectAllMesas returns all occupied mesas
func SelectAllMesas(ctx context.Context, dynamoClient DynamoDBClient, log *logging.Logger) ([]models.Mesa, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(tableMesas),
	}
	result, err := dynamoClient.Scan(ctx, input)
	if err != nil {
		return nil, err
	}
	var mesas []models.Mesa
	for _, item := range result.Items {
		var mesa models.Mesa
		if err := attributevalue.UnmarshalMap(item, &mesa); err != nil {
			return nil, err
		}
		if mesa.Cliente != "" {
			mesas = append(mesas, mesa)
		}
	}
	return mesas, nil
}

// SelectMesa returns a Mesa struct with the data from the DynamoDB table Mesas
func SelectMesa(ctx context.Context, mesaId string, dynamoClient DynamoDBClient, log *logging.Logger) (models.Mesa, error) {
	key, err := attributevalue.MarshalMap(map[string]interface{}{
		"Id": mesaId,
	})
	if err != nil {
		return models.Mesa{}, err
	}

	input := &dynamodb.GetItemInput{
		TableName: aws.String(tableMesas),
		Key:       key,
	}
	result, err := dynamoClient.GetItem(ctx, input)
	if err != nil {
		return models.Mesa{}, err
	}

	var mesa models.Mesa
	if err := attributevalue.UnmarshalMap(result.Item, &mesa); err != nil {
		return models.Mesa{}, err
	}

	return mesa, nil
}

// SelectCardapio returns a slice of Cardapio structs with the data from the DynamoDB table Cardapio
func SelectCardapio(ctx context.Context, dynamoClient DynamoDBClient, log *logging.Logger) ([]models.Cardapio, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(tableCardapio),
	}
	result, err := dynamoClient.Scan(ctx, input)
	if err != nil {
		log.HandleError("E", "Failed to scan Cardapio", err)
		return nil, err
	}
	var cardapio []models.Cardapio
	for _, item := range result.Items {
		var itemCardapio models.Cardapio
		if err := attributevalue.UnmarshalMap(item, &itemCardapio); err != nil {
			log.HandleError("E", "Failed to unmarshal Cardapio item", err)
			return nil, err
		}
		cardapio = append(cardapio, itemCardapio)
	}
	sort.Slice(cardapio, func(i, j int) bool {
		return cardapio[i].Nome < cardapio[j].Nome
	})
	return cardapio, nil
}

// UpdateMesa updates the data from a Mesa struct in the DynamoDB table Mesas
func UpdateMesa(ctx context.Context, mesa models.Mesa, dynamoClient DynamoDBClient, log *logging.Logger) error {
	item, err := attributevalue.MarshalMap(mesa)
	if err != nil {
		log.HandleError("E", "Failed to marshal Mesa item", err)
		return err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(tableMesas),
		Item:      item,
	}
	if _, err := dynamoClient.PutItem(ctx, input); err != nil {
		log.HandleError("E", "Failed to put Mesa item", err)
		return err
	}
	return nil
}

// DeletePedido deletes a Pedido struct from the DynamoDB table Pedidos
// NOTE: Function implementation is missing
