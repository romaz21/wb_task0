package pub

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"github.com/nats-io/stan.go"
	"example.com/wb/db"
)

const (
	clusterID = "test-cluster"
	clientID  = "your-publisher-id"
	subject   = "wb0"
)

func Pub() {
	fileContent, err := ioutil.ReadFile("model.json")
	if err != nil {
		log.Fatal(err)
	}

	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL("nats://10.55.171.64:4222"))
	if err != nil {
		log.Fatal(err)
	}
	defer sc.Close()

	for i := 0; i < 1; i++ {
		err = sc.Publish(subject, fileContent)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Published message to channel %s, number: %d \n", subject, i)
	}
	

	var order db.Order

	err = json.Unmarshal(fileContent, &order)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", order)
}