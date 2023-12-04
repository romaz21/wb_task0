package server

import (
	"database/sql"
	"example.com/wb/cache"
	"log"
	"html/template"
	"net/http"
)
type AppHandler struct {
    DB *sql.DB
    c *cache.Cache
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




