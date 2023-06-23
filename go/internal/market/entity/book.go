package entity

import (
	"container/heap"
	"sync"
)

// EStrutura dados do book
type Book struct {
	Order         []*Order
	Transactions  []*Transaction
	OrdersChan    chan *Order     // Canal entrada ordens
	OrdersChanOut chan *Order     // Canal saida ordens
	Wg            *sync.WaitGroup //recurso Go await
}

// Função do book
func NewBook(orderChan chan *Order, orderChanOut chan *Order, wg *sync.WaitGroup) *Book {
	return &Book{
		Order:         []*Order{},
		Transactions:  []*Transaction{},
		OrdersChan:    orderChan,
		OrdersChanOut: orderChanOut,
		Wg:            wg,
	}
}

// Método trade, correlação de ordens de compra e venda
func (b *Book) Trade() {
	buyOrders := make(map[string]*OrderQueue)
	sellOrders := make(map[string]*OrderQueue)
	// buyOrders := NewOrderQueue()
	// sellOrders := NewOrderQueue()

	// heap.Init(buyOrders)
	// heap.Init(sellOrders)

	// loop roda infinitamente em uma tread separada verificando as orders do canal do Kafka
	for order := range b.OrdersChan {
		asset := order.Asset.ID

		if buyOrders[asset] == nil {
			buyOrders[asset] = NewOrderQueue()
			heap.Init(buyOrders[asset])
		}

		if sellOrders[asset] == nil {
			sellOrders[asset] = NewOrderQueue()
			heap.Init(sellOrders[asset])
		}

		if order.OrderType == "BUY" { // Ordens de compra
			buyOrders[asset].Push(order)
			if sellOrders[asset].Len() > 0 && sellOrders[asset].Orders[0].Price <= order.Price { // Existe order de venda com preço <= a ordem de compra
				sellOrder := sellOrders[asset].Pop().(*Order) // Remove da fila
				if sellOrder.PendingShares > 0 {              // Se cair em = 0 ela ja foi liquidada
					transaction := NewTransaction(sellOrder, order, order.Shares, sellOrder.Price) // Nova transação
					b.AddTransaction(transaction, b.Wg)
					sellOrder.Transactions = append(sellOrder.Transactions, transaction) // Adicionar ordem de venda
					order.Transactions = append(order.Transactions, transaction)         // Adicionar ordem de compra
					b.OrdersChanOut <- sellOrder                                         // Canais de saida pro kafka
					b.OrdersChanOut <- order                                             // Outra tread para isso
					if sellOrder.PendingShares > 0 {
						sellOrders[asset].Push(sellOrder) // Transação ainda não realizada faltou shares
					}
				}
			}
		} else if order.OrderType == "SELL" { // Ordens de venda
			sellOrders[asset].Push(order)
			if buyOrders[asset].Len() > 0 && buyOrders[asset].Orders[0].Price >= order.Price {
				buyOrder := buyOrders[asset].Pop().(*Order)
				if buyOrder.PendingShares > 0 {
					transaction := NewTransaction(order, buyOrder, order.Shares, buyOrder.Price)
					b.AddTransaction(transaction, b.Wg)
					buyOrder.Transactions = append(buyOrder.Transactions, transaction)
					order.Transactions = append(order.Transactions, transaction)
					b.OrdersChanOut <- buyOrder
					b.OrdersChanOut <- order
					if buyOrder.PendingShares > 0 {
						buyOrders[asset].Push(buyOrder)
					}
				}
			}
		}
	}
}

// Método que adiciona a transação na order (controle de pending shares)
func (b *Book) AddTransaction(transaction *Transaction, wg *sync.WaitGroup) {
	defer wg.Done() // Comando defer coloca para rodar por ultimo o wg.Done() Transação realizada

	sellingShares := transaction.SellingOrder.PendingShares
	buyingShares := transaction.BuyingOrder.PendingShares

	minShares := sellingShares
	if buyingShares < minShares {
		minShares = buyingShares
	}

	transaction.SellingOrder.Investor.UpdateAssetPosition(transaction.SellingOrder.Asset.ID, -minShares)
	transaction.AddSellOrderPendingShares(-minShares) // Subtrai de quem vende shares

	transaction.BuyingOrder.Investor.UpdateAssetPosition(transaction.BuyingOrder.Asset.ID, minShares)
	transaction.AddBuyOrderPendingShares(-minShares) // soma de quem compra shares

	transaction.CalculateTotal(transaction.Shares, transaction.BuyingOrder.Price)
	transaction.CloseBuyOrder()
	transaction.CloseSellOrder()
	b.Transactions = append(b.Transactions, transaction)
}
