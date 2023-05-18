package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/elman23/articleapi/pkg/db"
	"github.com/elman23/articleapi/pkg/handlers"
	"github.com/gorilla/mux"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the Article REST API!")
	fmt.Println("Article REST API.")
}

func handleRequests(DB *sql.DB) {
	// DB handler
	h := handlers.New(DB)
	// Create a new instance of the mux router.
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/articles", h.GetAllArticles).Methods(http.MethodGet)
	myRouter.HandleFunc("/articles/{id}", h.GetArticle).Methods(http.MethodGet)
	myRouter.HandleFunc("/articles", h.AddArticle).Methods(http.MethodPost)
	myRouter.HandleFunc("/articles/{id}", h.UpdateArticle).Methods(http.MethodPut)
	myRouter.HandleFunc("/articles/{id}", h.DeleteArticle).Methods(http.MethodDelete)
	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

func main() {
	DB := db.Connect()
	handleRequests(DB)
	db.CloseConnection(DB)
}
