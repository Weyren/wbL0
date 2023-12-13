package cache

import (
	"WBL0/internal/common"
	"WBL0/internal/infrastructure/postgres"
	"context"
	"fmt"
	"sync"
)

type OrderCache struct {
	Mutex sync.Mutex
	Data  map[string]common.Order
}

func NewOrderCache() *OrderCache {
	return &OrderCache{
		Mutex: sync.Mutex{},
		Data:  make(map[string]common.Order, 0),
	}
}

func (oc *OrderCache) PutOrder(uid string, order common.Order) error {
	oc.Mutex.Lock()
	defer oc.Mutex.Unlock()
	oc.Data[uid] = order
	return nil

}

func (oc *OrderCache) GetOrder(uid string) (common.Order, error) {
	oc.Mutex.Lock()
	defer oc.Mutex.Unlock()

	if order, exists := oc.Data[uid]; exists {
		return order, nil
	}
	err := fmt.Errorf("cache error")
	return common.Order{}, err
}

func (oc *OrderCache) GetAllOrdersFromDB(op *postgres.OrderPostgres) error {

	ordersFromDB, err := op.GetAllOrders(context.Background())
	if err != nil {
		return err
	}

	for _, order := range ordersFromDB {
		err := oc.PutOrder(order.OrderUID, order)
		if err != nil {
			return fmt.Errorf("failed to put to cache")
		}
	}
	return nil
}
