package db

import (
	"database/sql"
	"log"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "1"
)



func Connect_db() *sql.DB{
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable", host, port, user, password)
	
	// Создание подключения к базе данных
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("Unable to connect to the database:", err)
	}
	//defer db.Close()

	// Проверка подключения
	err = db.Ping()
	if err != nil {
		fmt.Println("Error pinging the database:", err)
	}
	return db
}

func Write_db(db *sql.DB, order Order) {
	// Замените параметры подключения на ваши
	

	// Пример вызова функции записи в базу данных
	

	err := insert(db, order)
	if err != nil {
		fmt.Println("Error inserting order:", err)
		return
	}

	fmt.Println("Order inserted successfully.")
}

func insertDelivery(db *sql.DB, delivery Delivery) (int, error) {
	var deliveryID int
	err := db.QueryRow("INSERT INTO deliveries (name, phone, zip, city, address, region, email) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id",
		delivery.Name, delivery.Phone, delivery.Zip, delivery.City, delivery.Address, delivery.Region, delivery.Email).Scan(&deliveryID)
	return deliveryID, err
}

func insertPayment(db *sql.DB, payment Payment) (int, error) {
	var paymentID int
	err := db.QueryRow("INSERT INTO payments (transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id",
		payment.Transaction, payment.RequestID, payment.Currency, payment.Provider, payment.Amount, payment.PaymentDt, payment.Bank, payment.DeliveryCost, payment.GoodsTotal, payment.CustomFee).Scan(&paymentID)
	return paymentID, err
}

func insertItem(db *sql.DB, item Item) (int, error) {
	var itemID int
	err := db.QueryRow("INSERT INTO items (chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id",
		item.ChrtID, item.TrackNumber, item.Price, item.RID, item.Name, item.Sale, item.Size, item.TotalPrice, item.NmID, item.Brand, item.Status).Scan(&itemID)
	return itemID, err
}

func insertOrder(db *sql.DB, order Order, deliveryID int, paymentID int) (int, error) {
	var orderID int
	err := db.QueryRow("INSERT INTO orders (order_uid, track_number, entry, delivery_id, payment_id, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) RETURNING id",
		order.OrderUID, order.TrackNumber, order.Entry, deliveryID, paymentID, order.Locale, order.InternalSignature, order.CustomerID, order.DeliveryService, order.ShardKey, order.SmID, order.DateCreated, order.OofShard).Scan(&orderID)
	return orderID, err
}

func insertOrderItem(db *sql.DB, orderID int, itemID int) (error) {
	_, err := db.Exec("INSERT INTO orders_items (order_id, item_id) VALUES ($1, $2)",
		orderID, itemID)
	return err
}

func insert(db *sql.DB, order Order) error {
	deliveryID, err := insertDelivery(db, order.Delivery)
	if err != nil {
		log.Fatal(err)
	}

	paymentID, err := insertPayment(db, order.Payment)
	if err != nil {
		log.Fatal(err)
	}

	orderID, err := insertOrder(db, order, deliveryID, paymentID)
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range order.Items {
		itemID, err := insertItem(db, item)
		if err != nil {
			log.Fatal(err)	
		}
		err = insertOrderItem(db, orderID, itemID)
		if err != nil {
			log.Fatal(err)	
		}
	}

	

	fmt.Printf("Inserted data: Delivery ID=%d, Payment ID=%d, Order ID=%d\n", deliveryID, paymentID, orderID)
	return nil
}

func QueryId(db *sql.DB, id string) (*sql.Rows, error) {
	// Выполнение запроса к базе данных
	rows, err := db.Query("SELECT order_uid, track_number, payment_id FROM orders WHERE id = $1", id)
	if err != nil {
		return nil, err
	}

	return rows, nil
}
