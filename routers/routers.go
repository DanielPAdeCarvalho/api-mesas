package routers

import (
	"mesas-api/logging"
	"mesas-api/query"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gin-gonic/gin"
)

func ResponseOK(c *gin.Context, log logging.Logfile) {
	c.IndentedJSON(http.StatusOK, "Servidor up")
}

func GetMesa(numero int, c *gin.Context, dynamoClient *dynamodb.Client, log logging.Logfile) {
	mesa := query.SelectMesa(numero, dynamoClient, log)
	c.IndentedJSON(http.StatusOK, mesa)
}

func GetAllMesas(c *gin.Context, dynamoClient *dynamodb.Client, log logging.Logfile) {
	mesas := query.SelectAllMesas(dynamoClient, log)
	c.IndentedJSON(http.StatusOK, mesas)
}
