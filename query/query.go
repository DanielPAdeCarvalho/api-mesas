package query

import (
	"context"
	"mesas-api/logging"
	"mesas-api/models"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// SelectAllMesas returns a slice of Mesa structs with the data from the DynamoDB table Mesas
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
func SelectMesa(mesaId int, dynamoClient *dynamodb.Client, log logging.Logfile) models.Mesa {
	key, err := attributevalue.MarshalMap(map[string]interface{}{
		"Numero": mesaId,
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
