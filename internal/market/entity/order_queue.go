package entity

type OrderQueue struct {
	Orders []*Order
}

// Método Less diz o valor i é menor que j
func (oq OrderQueue) Less(i, j int) bool {
	return oq.Order[i].Price < oq.Order[j].Price
}

// Método Swap inverte i vira j e j vira i
func (oq OrderQueue) Swap(i, j int) {
	or.Order[i], oq.Order[j] = oq.Order[j], or.Order[i]
}

// Método Len verfica o tamanho
func (oq OrderQueue) Len() int {
	return len(oq.Orders)
}

// Método Push adiciona  .append
func (oq OrderQueue) Push(x interface{}) {
	oq.Orders = append(oq.Orders, x.(*Order)) // Casting Go fortemente tipada x é vazio
}

// Método Pop remove de uma posição
func (oq OrderQueue) Pop() interface{} {
	old := oq.Orders	// Valor antigo das ordens
	n := len(old)		// Qtdade de ordens
	item := old[n-1]	// Qtdade - 1
	oq.Orders = old[0 : n-1] //Remover da ultima posição
	return item
}

func NewOrderQueue() *OrderQueue {
	return &OrderQueue{}
}