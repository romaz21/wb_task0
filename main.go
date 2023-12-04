package main
import (
	"example.com/wb/create_db"
	"example.com/wb/db"
	"example.com/wb/cache"
	"example.com/wb/sub"
	"example.com/wb/server"

)

func main(){
	if 1==1{
		create_db.Create_db()
	}
	dbase := db.Connect_db()
	cach := cache.NewCache(dbase)
	go server.StartHttpServer(dbase, cach)
	go sub.Subscribe(dbase, cach)

	select {}
}