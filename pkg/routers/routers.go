package routers

import (
	"mesas-api/pkg/driver"
	"mesas-api/pkg/handlers"
	"mesas-api/pkg/logging"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/gin-gonic/gin"
)

func SetupRouter(dynamoClient driver.DynamoDBClient, clienteSQS *sqs.Client, log *logging.Logger) *gin.Engine {
	handlers := handlers.NewHandlers(dynamoClient, log)

	router := gin.Default()

	router.GET("/status", handlers.ResponseOK)
	router.GET("/mesa/:id", handlers.GetMesa)
	router.GET("/mesas", handlers.GetAllMesas)
	router.GET("/cardapio", handlers.GetCardapio)
	router.PUT("/mesa", handlers.PutMesa)
	router.POST("/pedido/:id", handlers.PostPedido)
	router.DELETE("/pedido/:id/:nomePedido", handlers.DeletePedido)
	router.DELETE("/mesa/:id", handlers.DeleteMesa)

	return router
}
