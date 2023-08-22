package handlers

import (
	"mesas-api/pkg/driver"
	"mesas-api/pkg/logging"
	"mesas-api/pkg/models"
	"mesas-api/pkg/query"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	dynamoClient driver.DynamoDBClient
	log          *logging.Logger
}

func NewHandlers(dynamoClient driver.DynamoDBClient, log *logging.Logger) *Handlers {
	return &Handlers{dynamoClient: dynamoClient, log: log}
}

func (h *Handlers) ResponseOK(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, "Servidor up")
}

func (h *Handlers) GetMesa(c *gin.Context) {
	id := c.Param("id")
	mesa, err := query.SelectMesa(c.Request.Context(), id, h.dynamoClient.Client, h.log)
	if err != nil {
		h.log.HandleError("E", "Failed to get Mesa", err)
		c.IndentedJSON(http.StatusInternalServerError, "Failed to get Mesa")
		return
	}
	c.IndentedJSON(http.StatusOK, mesa)
}

func (h *Handlers) GetAllMesas(c *gin.Context) {
	mesas, err := query.SelectAllMesas(c.Request.Context(), h.dynamoClient.Client, h.log)
	if err != nil {
		h.log.HandleError("E", "Failed to get Mesas", err)
		c.IndentedJSON(http.StatusInternalServerError, "Failed to get Mesas")
		return
	}
	c.IndentedJSON(http.StatusOK, mesas)
}

func (h *Handlers) GetCardapio(c *gin.Context) {
	cardapio, err := query.SelectCardapio(c.Request.Context(), h.dynamoClient.Client, h.log)
	if err != nil {
		h.log.HandleError("E", "Failed to get Cardapio", err)
		c.IndentedJSON(http.StatusInternalServerError, "Failed to get Cardapio")
		return
	}
	c.IndentedJSON(http.StatusOK, cardapio)
}

// Associa um novo cliente a uma mesa no banco
func (h *Handlers) PutMesa(c *gin.Context) {
	var mesa models.Mesa
	err := c.BindJSON(&mesa)
	if err != nil {
		h.log.HandleError("E", "Failed to bind JSON", err)
		c.IndentedJSON(http.StatusInternalServerError, "Failed to bind JSON")
		return
	}
	query.UpdateMesa(c.Request.Context(), mesa, h.dynamoClient.Client, h.log)
	c.IndentedJSON(http.StatusOK, "Mesa atualizada")
}

// Cria um novo pedido para uma mesa
func (h *Handlers) PostPedido(c *gin.Context) {
	id := c.Param("id")
	var pedido models.Pedido
	err := c.BindJSON(&pedido)
	if err != nil {
		h.log.HandleError("E", "Failed to bind JSON", err)
		c.IndentedJSON(http.StatusInternalServerError, "Failed to bind JSON")
		return
	}

	mesa, err := query.SelectMesa(c.Request.Context(), id, h.dynamoClient.Client, h.log)
	if err != nil {
		h.log.HandleError("E", "Failed to select Mesa", err)
		c.IndentedJSON(http.StatusInternalServerError, "Failed to select Mesa")
		return
	}

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
	query.UpdateMesa(c.Request.Context(), mesa, h.dynamoClient.Client, h.log)

	//Envia o pedido para a cozinha fila PedidosCozinha
	if pedido.Cozinha == true {
		// Fazer a logica do envio para a tela de cozinha aqui
	}

	c.IndentedJSON(http.StatusOK, "Pedido adicionado")
}

// Remove um pedido de uma mesa
func (h *Handlers) DeletePedido(c *gin.Context) {
	id := c.Param("id")
	nomePedido := c.Param("nomePedido")

	mesa, err := query.SelectMesa(c.Request.Context(), id, h.dynamoClient.Client, h.log)
	if err != nil {
		h.log.HandleError("E", "Failed to select Mesa", err)
		c.IndentedJSON(http.StatusInternalServerError, "Failed to select Mesa")
		return
	}

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
	query.UpdateMesa(c.Request.Context(), mesa, h.dynamoClient.Client, h.log)
	c.IndentedJSON(http.StatusOK, "Pedido removido")
}

// Remove um cliente de uma mesa
func (h *Handlers) DeleteMesa(c *gin.Context) {
	id := c.Param("id")

	mesa, err := query.SelectMesa(c.Request.Context(), id, h.dynamoClient.Client, h.log)
	if err != nil {
		h.log.HandleError("E", "Failed to select Mesa", err)
		c.IndentedJSON(http.StatusInternalServerError, "Failed to select Mesa")
		return
	}
	// Remove the pedido from the mesa
	mesa.Cliente = ""
	mesa.Pedidos = nil
	query.UpdateMesa(c.Request.Context(), mesa, h.dynamoClient.Client, h.log)
	c.IndentedJSON(http.StatusOK, "Cliente removido")
}
