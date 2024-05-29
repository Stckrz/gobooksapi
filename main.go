package main

import (
	"database/sql"
	"fmt"
	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "postgres"
)

var db *sql.DB

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	var err error
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("unable to connect to db: %v\n", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Ok"))
	})
	r.Mount("/books", BookRoutes())

	fmt.Println("Starting server")
	http.ListenAndServe(":3000", r)
}

func BookRoutes() chi.Router {
	r := chi.NewRouter()
	BookHandler := BookHandler{}
	r.Get("/", BookHandler.listBooks)
	r.Get("/{id}", BookHandler.getBook)
	r.Patch("/{id}", BookHandler.editBook)
	r.Post("/", BookHandler.postBook)
	return r
}
