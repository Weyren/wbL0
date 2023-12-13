package nats

import (
	"WBL0/internal/common"
	"WBL0/internal/infrastructure/cache"
	"WBL0/internal/infrastructure/postgres"
	"context"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/nats-io/stan.go"
	"log"
)

func validateOrder(order common.Order) error {
	validate := validator.New()
	return validate.Struct(order)
}

type Nats struct {
	orderPostgres *postgres.OrderPostgres
	orderCache    *cache.OrderCache
}

func NewNats(orderPostgres *postgres.OrderPostgres, orderCache *cache.OrderCache) *Nats {
	return &Nats{orderPostgres: orderPostgres, orderCache: orderCache}
}

func (n Nats) ConnectNS(clusterID, clientID, natsURL string) (stan.Conn, error) {
	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL(natsURL))
	if err != nil {
		log.Fatalf("Error connecting to NATS Streaming: %v", err)
		return nil, err
	}

	return sc, nil
}
func (n Nats) Subscribe(sc stan.Conn, subject string) (stan.Subscription, error) {

	subscription, err := sc.Subscribe(subject, func(msg *stan.Msg) {
		var order common.Order

		err := json.Unmarshal(msg.Data, &order)
		if err != nil {
			log.Println("Error decoding JSON message, &v", err)
			return
		}
		if err := validateOrder(order); err != nil {
			log.Printf("Invalid order received: %v", err)
			return
		}

		err = n.orderPostgres.CreateOrder(context.TODO(), &order)
		if err != nil {
			log.Println(err)
			return
		}

		err = n.orderCache.PutOrder(order.OrderUID, order)
		if err != nil {
			return
		}

	})

	if err != nil {
		log.Fatalf("Error subscribing to NATS Streaming: %v", err)
	}
	return subscription, nil
}
