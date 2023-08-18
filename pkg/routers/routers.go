package routers

import (
	"mesas-api/pkg/logging"
	"mesas-api/pkg/models"
	"mesas-api/pkg/query"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/gin-gonic/gin"
)

func ResponseOK(c *gin.Context, log logging.Logfile) {
	c.IndentedJSON(http.StatusOK, "Servidor up")
}

func GetMesa(id string, c *gin.Context, dynamoClient *dynamodb.Client, log logging.Logfile) {
	mesa := query.SelectMesa(id, dynamoClient, log)
	c.IndentedJSON(http.StatusOK, mesa)
}

func GetAllMesas(c *gin.Context, dynamoClient *dynamodb.Client, log logging.Logfile) {
	mesas := query.SelectAllMesas(dynamoClient, log)
	c.IndentedJSON(http.StatusOK, mesas)
}

func GetCardapio(c *gin.Context, dynamoClient *dynamodb.Client, log logging.Logfile) {
	cardapio := query.SelectCardapio(dynamoClient, log)
	c.IndentedJSON(http.StatusOK, cardapio)
}

// Associa um novo cliente a uma mesa no banco
func PutMesa(mesa models.Mesa, c *gin.Context, dynamoClient *dynamodb.Client, log logging.Logfile) {
	query.UpdateMesa(mesa, dynamoClient, log)
	c.IndentedJSON(http.StatusOK, "Mesa atualizada")
}

// Cria um novo pedido para uma mesa
func PostPedido(id string, c *gin.Context, dynamoClient *dynamodb.Client, clienteSQS *sqs.Client, log logging.Logfile) {
	var pedido models.Pedido
	err := c.BindJSON(&pedido)
	logging.Check(err, log)

	mesa := query.SelectMesa(id, dynamoClient, log)

	// Get the mesa from the database and create the Pedidos map if it's nil
	if mesa.Pedidos == nil {
		mesa.Pedidos = make(map[string]models.Pedido)
		pedido.Quantidade = 1
		mesa.Pedidos[pedido.Nome] = pedido
	} else {
		// Check if the pedido already exists
		pedidoTmp, ok := mesa.Pedidos[pedido.Nome]
		if ok {
			// If it does, increment the quantity
			pedidoTmp.Quantidade++
			mesa.Pedidos[pedido.Nome] = pedidoTmp
		} else {
			// If it doesn't, create a new pedido
			pedido.Quantidade = 1
			mesa.Pedidos[pedido.Nome] = pedido
		}
	}
	// Add the pedido to the mesa
	query.UpdateMesa(mesa, dynamoClient, log)

	//Envia o pedido para a cozinha fila PedidosCozinha
	if pedido.Cozinha == true {

	}

	c.IndentedJSON(http.StatusOK, "Pedido adicionado")
}

// Remove um pedido de uma mesa
func DeletePedido(id string, nomePedido string, c *gin.Context, dynamoClient *dynamodb.Client, log logging.Logfile) {
	// Get the mesa from the database
	mesa := query.SelectMesa(id, dynamoClient, log)
	// Remove the pedido from the mesa
	if pedido, ok := mesa.Pedidos[nomePedido]; ok {
		if pedido.Quantidade > 1 {
			pedido.Quantidade--
			mesa.Pedidos[pedido.Nome] = pedido
		} else {
			delete(mesa.Pedidos, nomePedido)
			if len(mesa.Pedidos) == 0 {
				mesa.Pedidos = nil
			}
		}
	}
	query.UpdateMesa(mesa, dynamoClient, log)
	c.IndentedJSON(http.StatusOK, "Pedido removido")
}

// Remove um cliente de uma mesa
func DeleteMesa(id string, c *gin.Context, dynamoClient *dynamodb.Client, log logging.Logfile) {
	// Get the mesa from the database
	mesa := query.SelectMesa(id, dynamoClient, log)
	// Remove the pedido from the mesa
	mesa.Cliente = ""
	mesa.Pedidos = nil
	query.UpdateMesa(mesa, dynamoClient, log)
	c.IndentedJSON(http.StatusOK, "Cliente removido")
}
