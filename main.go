package main
import (
	//"fmt"
	"example.com/wb/db"
	"example.com/wb/cache"
	"example.com/wb/sub"
	"example.com/wb/server"

)

func main(){
	dbase := db.Connect_db()
	cach := cache.NewCache(dbase)

	go server.StartHttpServer(dbase, cach)
	go sub.Subscribe(dbase, cach)

	select {}
}