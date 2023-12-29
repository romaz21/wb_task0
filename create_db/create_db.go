package create_db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "1"
)

func Create_db() {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable", host, port, user, password)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully connected to the 'postgres' database!")

	_, err = db.Exec(`
		DROP DATABASE IF EXISTS test;
	`)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Database 'test' deleted successfully!")

	_, err = db.Exec(`
		CREATE DATABASE test
	`)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Database 'test' created successfully!")
	createDeliveriesTableQuery := `
		CREATE TABLE IF NOT EXISTS deliveries (
	    id SERIAL PRIMARY KEY,
	    name VARCHAR(255),
	    phone VARCHAR(20),
	    zip VARCHAR(10),
	    city VARCHAR(255),
	    address VARCHAR(255),
	    region VARCHAR(255),
	    email VARCHAR(255)
	);
	`
	createPaymentsTableQuery := `
		CREATE TABLE IF NOT EXISTS payments (
	    id SERIAL PRIMARY KEY,
	    transaction VARCHAR(255),
	    request_id VARCHAR(255),
	    currency VARCHAR(5),
	    provider VARCHAR(255),
	    amount INT,
	    payment_dt BIGINT,
	    bank VARCHAR(255),
	    delivery_cost INT,
	    goods_total INT,
	    custom_fee INT
	);
	`
	
	createItemsTableQuery := `
		CREATE TABLE IF NOT EXISTS items (
	    id SERIAL PRIMARY KEY,
	    chrt_id INT,
	    track_number VARCHAR(255),
	    price INT,
	    rid VARCHAR(255),
	    name VARCHAR(255),
	    sale INT,
	    size VARCHAR(10),
	    total_price INT,
	    nm_id INT,
	    brand VARCHAR(255),
	    status INT
	);
	`

	createOrdersTableQuery := `
		CREATE TABLE IF NOT EXISTS orders (
	    id SERIAL PRIMARY KEY,
	    order_uid VARCHAR(255),
	    track_number VARCHAR(255),
	    entry VARCHAR(255),
	    delivery_id INT REFERENCES deliveries(id),
	    payment_id INT REFERENCES payments(id),
	    locale VARCHAR(10),
	    internal_signature VARCHAR(255),
	    customer_id VARCHAR(255),
	    delivery_service VARCHAR(255),
	    shardkey VARCHAR(10),
	    sm_id INT,
	    date_created TIMESTAMP,
	    oof_shard VARCHAR(10)
	);
	`
	createOrdersItemsTableQuery := `
		CREATE TABLE IF NOT EXISTS orders_items (
	    order_id INT REFERENCES orders(id),
	    item_id INT REFERENCES items(id),
	    PRIMARY KEY (order_id, item_id)
	);
	`

	_, err = db.Exec(`
		DROP TABLE IF EXISTS orders_items;
	`)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Table 'orders_items' deleted successfully!")

	_, err = db.Exec(`
		DROP TABLE IF EXISTS orders;
	`)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Table 'orders' deleted successfully!")

	_, err = db.Exec(`
		DROP TABLE IF EXISTS items;
	`)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Table 'items' deleted successfully!")

	_, err = db.Exec(`
		DROP TABLE IF EXISTS deliveries;
	`)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Table 'deliveries' deleted successfully!")

	_, err = db.Exec(`
		DROP TABLE IF EXISTS payments;
	`)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Table 'payments' deleted successfully!")


	_, err = db.Exec(createDeliveriesTableQuery)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Table 'Deliveries' created successfully!")

	_, err = db.Exec(createPaymentsTableQuery)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Table 'Payments' created successfully!")

	_, err = db.Exec(createItemsTableQuery)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Table 'Items' created successfully!")

	_, err = db.Exec(createOrdersTableQuery)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Table 'Orders' created successfully!")

	_, err = db.Exec(createOrdersItemsTableQuery)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Table 'orders_items' created successfully!")

	
}
