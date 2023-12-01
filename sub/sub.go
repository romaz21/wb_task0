package sub

import (
	"encoding/json"
	"example.com/wb/db"
	"example.com/wb/cache"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/nats-io/stan.go"
	_ "github.com/lib/pq"
	"database/sql"
)

const (
	clusterID = "test-cluster"
	clientID  = "your-subscriber-id"
	subject   = "wb0"
	durableID = "your-durable-id"
)



func Subscribe(dbase *sql.DB, cach *cache.Cache) {
	// Подключение к кластеру NATS Streaming
	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL("nats://localhost:4222"))
	if err != nil {
		log.Fatal(err)
	}
	defer sc.Close()

	// Опции для подписки
	subOptions := []stan.SubscriptionOption{
		stan.SetManualAckMode(),
		stan.DurableName(durableID),
	}

	// Функция обработки полученных сообщений
	messageHandler := func(msg *stan.Msg) {
		var orderData db.Order

		// Распаковка JSON-данных из сообщения
		err := json.Unmarshal(msg.Data, &orderData)
		if err != nil {
			log.Println("Error decoding JSON:", err)
			return
		}
		// Ваша логика обработки сообщения здесь
		go db.Write_db(dbase, orderData)
		go cach.Set(orderData.OrderUID, string(msg.Data))
		//fmt.Println(orderData)
		// Подтверждение получения сообщения
		err = msg.Ack()
		if err != nil {
			log.Printf("Error acknowledging message: %v\n", err)
		}
	}

	// Подписка на канал с опциями
	sub, err := sc.Subscribe(subject, messageHandler, subOptions...)
	if err != nil {
		log.Fatal(err)
	}
	defer sub.Unsubscribe()

	// Ожидание сигналов завершения
	select {}
	//waitForSignal()
}

// waitForSignal ожидает сигналы завершения работы приложения (например, Ctrl+C)
func waitForSignal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-c
	fmt.Println("\nReceived signal. Shutting down...")
}
