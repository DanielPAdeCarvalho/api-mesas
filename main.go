package main

import (
	"log"
	"mesas-api/driver"
	"mesas-api/logging"
	"mesas-api/models"
	"mesas-api/routers"
	"net/http"
	"os"
	"strconv"

	"github.com/apex/gateway"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gin-gonic/gin"
)

var (
	dynamoClient *dynamodb.Client
	logs         logging.Logfile
)

func inLambda() bool {
	if lambdaTaskRoot := os.Getenv("LAMBDA_TASK_ROOT"); lambdaTaskRoot != "" {
		return true
	}
	return false
}

func setupRouter() *gin.Engine {
	apiRouter := gin.Default()

	apiRouter.GET("/", func(ctx *gin.Context) {
		logs.InfoLogger.Println("Servidor Ok")
		routers.ResponseOK(ctx, logs)
	})

	apiRouter.GET("/mesas", func(c *gin.Context) {
		routers.GetAllMesas(c, dynamoClient, logs)
	})

	apiRouter.GET("/mesa/:numero", func(c *gin.Context) {
		numeroStr := c.Param("numero")
		numero, err := strconv.Atoi(numeroStr)
		logging.Check(err, logs)
		routers.GetMesa(int(numero), c, dynamoClient, logs)
	})

	//Adicionar um cliente a uma mesa
	apiRouter.PUT("/mesa", func(c *gin.Context) {
		var mesa models.Mesa
		err := c.BindJSON(&mesa)
		logging.Check(err, logs)
		routers.PutMesa(mesa, c, dynamoClient, logs)
	})

	//Novo pedido para a mesa
	apiRouter.POST("/mesa/:numero", func(c *gin.Context) {
		numeroStr := c.Param("numero")
		numero, err := strconv.Atoi(numeroStr)
		logging.Check(err, logs)
		routers.PostPedido(numero, c, dynamoClient, logs)
	})

	//Remover um pedido da mesa
	apiRouter.DELETE("/mesa/:numero/:pedido", func(c *gin.Context) {
		numeroStr := c.Param("numero")
		numero, err := strconv.Atoi(numeroStr)
		logging.Check(err, logs)
		pedido := c.Param("pedido")
		routers.DeletePedido(numero, pedido, c, dynamoClient, logs)
	})

	//Remover um cliente da mesa
	apiRouter.DELETE("/mesa/:numero", func(c *gin.Context) {
		numeroStr := c.Param("numero")
		numero, err := strconv.Atoi(numeroStr)
		logging.Check(err, logs)
		routers.DeleteMesa(numero, c, dynamoClient, logs)
	})
	return apiRouter
}

// Para compilar o binario do sistema usamos:
//
//	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o assinatura-api .
//
// para criar o zip do projeto comando:
//
// zip lambda.zip assinatura-api
// .
func main() {
	InfoLogger := log.New(os.Stdout, " ", log.LstdFlags|log.Lshortfile)
	ErrorLogger := log.New(os.Stdout, " ", log.LstdFlags|log.Lshortfile)

	logs.InfoLogger = *InfoLogger
	logs.ErrorLogger = *ErrorLogger
	var err error
	// chamada de função para a criação da sessao de login com o banco
	dynamoClient, err = driver.ConfigAws()
	//chamada da função para revificar o erro retornado
	logging.Check(err, logs)

	if inLambda() {

		log.Fatal(gateway.ListenAndServe(":8080", setupRouter()))
	} else {

		log.Fatal(http.ListenAndServe(":8080", setupRouter()))
	}
}
