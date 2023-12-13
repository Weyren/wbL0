package main

import (
	"WBL0/internal/config"
	"WBL0/internal/infrastructure/cache"
	"WBL0/internal/infrastructure/http"
	"WBL0/internal/infrastructure/nats"
	"WBL0/internal/infrastructure/postgres"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.GetConfig()
	log.Println("Get config")

	pool, err := postgres.ConnectDB(cfg.DB.User, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.Name)
	if err != nil {
		log.Fatal("Error connection to DB, %v", err)
	}
	defer pool.Close()
	log.Println("Create pool connection")

	orderpostgres := postgres.NewOrderPostgres(pool)
	ordercache := cache.NewOrderCache()

	err = ordercache.GetAllOrdersFromDB(orderpostgres)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Load cache from DB")

	nats := nats.NewNats(orderpostgres, ordercache)

	sc, err := nats.ConnectNS(
		cfg.Nats.Streaming.ClusterID,
		cfg.Nats.Streaming.ClientID,
		cfg.Nats.Streaming.NatsURL,
	)
	defer sc.Close()
	log.Println("Create nats-streaming connection")

	sub, err := nats.Subscribe(sc, cfg.Nats.Streaming.Subject)
	if err != nil {
		log.Fatal("failed subs, %v", err)
	}
	defer sub.Close()
	log.Println("Subscribed to nats-streaming")
	pub()
	handler := http.NewHandler(*ordercache, *orderpostgres)
	handler.RunServer()
	log.Println("Running server")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
}
