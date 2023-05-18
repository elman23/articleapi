package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/elman23/articleapi/pkg/handlers"
	"github.com/gorilla/mux"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the Article REST API!")
	fmt.Println("Article REST API.")
}

func handleRequests() {
	// Create a new instance of the mux router.
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/articles", handlers.GetAllArticles).Methods(http.MethodGet)
	myRouter.HandleFunc("/articles/{id}", handlers.GetArticle).Methods(http.MethodGet)
	myRouter.HandleFunc("/articles", handlers.AddArticle).Methods(http.MethodPost)
	myRouter.HandleFunc("/articles/{id}", handlers.UpdateArticle).Methods(http.MethodPut)
	myRouter.HandleFunc("/articles/{id}", handlers.DeleteArticle).Methods(http.MethodDelete)
	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

func main() {
	handleRequests()
}
