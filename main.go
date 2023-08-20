package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/elman23/articleapi/pkg/auth"
	"github.com/elman23/articleapi/pkg/db"
	"github.com/elman23/articleapi/pkg/handlers"
	"github.com/gorilla/mux"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	// Printed to the console
	log.Println("Endpoint: [homePage].")
	// Home page message returned to the ResponseWriter
	fmt.Fprintf(w, "Welcome to the Article REST API!\nDeveloped with <3")
}

func welcome(w http.ResponseWriter, r *http.Request) {
	log.Println("Endpoint: [welcome]")
	fmt.Fprintf(w, "Welcome to the home page!")
}

func handleRequests(DB *sql.DB) {

	// Create a new instance of the mux router
	myRouter := mux.NewRouter().StrictSlash(true)

	// we will implement these handlers in the next sections
	myRouter.HandleFunc("/signin", auth.Signin).Methods(http.MethodPost)
	myRouter.Handle("/welcome", auth.IsAuthorized(welcome)).Methods(http.MethodGet)
	myRouter.Handle("/refresh", auth.IsAuthorized(auth.Refresh)).Methods(http.MethodGet)
	myRouter.HandleFunc("/logout", auth.Logout).Methods(http.MethodGet)

	// Database requests handler
	h := handlers.New(DB)
	// Add routes and handle functions
	myRouter.HandleFunc("/", homePage)
	myRouter.Handle("/articles", auth.IsAuthorized(h.GetAllArticles)).Methods(http.MethodGet)
	myRouter.Handle("/articles/{id}", auth.IsAuthorized(h.GetArticle)).Methods(http.MethodGet)
	myRouter.Handle("/articles", auth.IsAuthorized(h.AddArticle)).Methods(http.MethodPost)
	myRouter.Handle("/articles/{id}", auth.IsAuthorized(h.UpdateArticle)).Methods(http.MethodPut)
	myRouter.Handle("/articles/{id}", auth.IsAuthorized(h.DeleteArticle)).Methods(http.MethodDelete)

	// Log application startup
	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

func main() {
	var singleton db.DbServiceSingleton
	dbService := singleton.GetService()

	// Handle requests
	handleRequests(dbService.DB)

	// Close database connection
	singleton.CloseConnection()
}
