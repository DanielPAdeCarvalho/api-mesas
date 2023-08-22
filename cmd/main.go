package main

import (
	"context"
	"mesas-api/pkg/driver"
	"mesas-api/pkg/events"
	"mesas-api/pkg/logging"
	"mesas-api/pkg/routers"
	"net/http"
	"os"

	"github.com/apex/gateway"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var (
	dynamoClient *driver.DynamoDBClient
	Logs         *logging.Logger
	ClienteSQS   *sqs.Client
)

func main() {
	// Set log level as required
	logs := logging.NewLogger(logrus.InfoLevel)

	var err error
	dynamoClient, err = driver.ConfigAws(context.Background())
	logs.HandleError("E", "Failed to configure AWS", err)

	ClienteSQS, err = events.CreateClient(context.Background(), logs)

	router := setupRouter()

	if isRunningInLambda() {
		logs.HandleError("E", "Failed to start server", gateway.ListenAndServe(":8080", router))
	} else {
		logs.HandleError("E", "Failed to start server", http.ListenAndServe(":8080", router))
	}
}

// isRunningInLambda checks if the application is running within a Lambda environment.
func isRunningInLambda() bool {
	return os.Getenv("LAMBDA_TASK_ROOT") != ""
}

// setupRouter initializes the router and defines the routes for the application.
func setupRouter() *gin.Engine {
	router := routers.SetupRouter(dynamoClient, ClienteSQS, Logs)
	return router
}

// Para compilar o binario do sistema usamos:
//
//	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o mesas-api .
//
// para criar o zip do projeto comando:
//
// zip lambda.zip mesas-api
//
// main.go
