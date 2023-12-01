package server

import (
	//"fmt"
	"database/sql"
	"example.com/wb/cache"
	//"example.com/wb/db"
	"log"
	"html/template"
	"net/http"
)
type AppHandler struct {
    DB *sql.DB
    c *cache.Cache
}

type Result struct {
	Column1 string
	Column2 string
	Column3 string
}

type PageData struct {
	Order string
}

func (a *AppHandler) homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		id := r.FormValue("inputText")
		order, _ := a.c.Get(id)
		data := PageData{Order: order,}
		tmpl, err := template.ParseFiles("static/index.html")
		if err != nil {
			log.Fatal(err)
		}
		tmpl.Execute(w, data)
	} else {
		tmpl, err := template.ParseFiles("static/index.html")
		if err != nil {
			log.Fatal(err)
		}
		tmpl.Execute(w, nil)
	}
}

func StartHttpServer(db *sql.DB, c *cache.Cache) {
	appHandler := &AppHandler{DB: db, c: c}
	http.HandleFunc("/", appHandler.homeHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}




