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
	"github.com/sirupsen/logrus"
)

func main() {
	// Set log level from environment or config
	logLevel := logrus.InfoLevel
	logs := logging.NewLogger(logLevel)

	dynamoClient, err := driver.NewDynamoClient(context.Background(), logs)
	if err != nil {
		logs.HandleError("F", "Failed to configure AWS", err)
		return
	}

	ClienteSQS, err := events.NewSQSClient(context.Background(), logs)
	if err != nil {
		logs.HandleError("F", "Failed to configure SQS", err)
		return
	}

	router := routers.SetupRouter(dynamoClient, ClienteSQS, logs)

	serverPort := ":8080" // Default port

	if isRunningInLambda() {
		logs.HandleError("E", "Failed to start server", gateway.ListenAndServe(serverPort, router))
	} else {
		logs.HandleError("E", "Failed to start server", http.ListenAndServe(serverPort, router))
	}
}

// isRunningInLambda checks if the application is running within a Lambda environment.
func isRunningInLambda() bool {
	return os.Getenv("LAMBDA_TASK_ROOT") != ""
}

// Para compilar o binario do sistema usamos:
//
//	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o mesas-api .
//
// para criar o zip do projeto comando:
//
// zip lambda.zip mesas-api
//
