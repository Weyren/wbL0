package main

import (
	"github.com/nats-io/stan.go"
	"log"
	"os"
)

func pub() {
	clusterID := "test-cluster"
	clientID := "publisher"
	subject := "orders"
	natsURL := "nats://localhost:4222"

	// Создание подключения
	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL(natsURL))
	if err != nil {
		log.Fatalf("Ошибка подключения к NATS Streaming: %v", err)
	}
	defer sc.Close()

	// Сообщение для отправки

	for i := 0; i < 2; i++ {
		message, err := os.ReadFile("/Users/pavelsvinkin/GolandProjects/wbL0/internal/common/model.json")
		err = sc.Publish(subject, message)
		if err != nil {
			log.Fatalf("Ошибка отправки сообщения: %v", err)
		}

		log.Printf("Сообщение отправлено успешно: %s", message)
	}
}
