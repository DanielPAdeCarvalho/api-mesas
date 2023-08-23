package models

type Mesa struct {
	Id      string            `json:"id"`
	Cliente string            `json:"cliente"`
	Pedidos map[string]Pedido `json:"pedidos"`
}

type Pedido struct {
	Nome       string  `json:"nome"`
	Preco      float64 `json:"preco"`
	Quantidade int     `json:"quantidade"`
	Cozinha    bool    `json:"cozinha"`
}

type Cardapio struct {
	Nome    string  `json:"nome"`
	Preco   float64 `json:"preco"`
	Cozinha bool    `json:"cozinha"`
}

type PedidoCozinha struct {
	Item string `json:"item"`
	Nome string `json:"nome"`
}
