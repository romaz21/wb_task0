package sub

import (
	"encoding/json"
	"example.com/wb/db"
	"example.com/wb/cache"
	"log"
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
	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL("nats://localhost:4222"))
	if err != nil {
		log.Fatal(err)
	}
	defer sc.Close()

	subOptions := []stan.SubscriptionOption{
		stan.SetManualAckMode(),
		stan.DurableName(durableID),
	}

	messageHandler := func(msg *stan.Msg) {
		var orderData db.Order

		err := json.Unmarshal(msg.Data, &orderData)
		if err != nil {
			log.Println("Error decoding JSON:", err)
			return
		}
		go db.Write_db(dbase, orderData)
		go cach.Set(orderData.OrderUID, string(msg.Data))
		err = msg.Ack()
		if err != nil {
			log.Printf("Error acknowledging message: %v\n", err)
		}
	}

	sub, err := sc.Subscribe(subject, messageHandler, subOptions...)
	if err != nil {
		log.Fatal(err)
	}
	defer sub.Unsubscribe()

	select {}
}

