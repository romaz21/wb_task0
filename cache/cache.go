package cache

import (
	"fmt"
	"sync"
	"example.com/wb/db"
	"encoding/json"
	"database/sql"
)

type Cache struct {
	mu    sync.Mutex
	data  map[string]string
	db *sql.DB
	size int64
}

func NewCache(db *sql.DB) *Cache {
	cache := &Cache{}
	cache.Init(db)
	return cache

}

func (c *Cache) Init(db *sql.DB) {
	c.size = 100
	c.db = db
	c.data = make(map[string]string, c.size)
	res, err := c.getDataFromDb()
	if err != nil {
		fmt.Println("Error getting data from the database:", err)
	}
	for _, k := range res {
		jsonData, err := json.Marshal(k)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		//fmt.Println(string(jsonData))
		//fmt.Println(string(k.OrderUID))
		c.Set(k.OrderUID, string(jsonData))
	}
}

func (c *Cache) getCacheFromDb() {
	fmt.Printf("GETTING CACHE FROM DB")

}

func (c *Cache) getDataFromDb() ([]db.Order, error) {
	var results []db.Order
	rows, err := c.db.Query("SELECT order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard FROM orders")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var order db.Order
		if err := rows.Scan(&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale, &order.InternalSignature, &order.CustomerID, &order.DeliveryService, &order.ShardKey, &order.SmID, &order.DateCreated, &order.OofShard); err != nil {
			return nil, err
		}
		//fmt.Printf(string(result) + "\n")
		results = append(results, order)
	}

	return results, nil

}

func (c *Cache) Get(key string) (string, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	value, ok := c.data[key]
	return value, ok
}

func (c *Cache) Set(key string, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = value
}


