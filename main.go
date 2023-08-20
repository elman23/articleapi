package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/elman23/articleapi/pkg/auth"
	"github.com/elman23/articleapi/pkg/db"
	"github.com/elman23/articleapi/pkg/handlers"
	"github.com/gorilla/mux"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	// Printed to the console
	fmt.Println("Endpoint: [homePage].")
	// Home page message returned to the ResponseWriter
	fmt.Fprintf(w, "Welcome to the Article REST API!\nDeveloped with <3")
}

func welcome(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint: [welcome]")
	fmt.Fprintf(w, "Welcome to the home page!")
}

func handleRequests(DB *sql.DB) {
	// Database requests handler
	h := handlers.New(DB)

	// Create a new instance of the mux router
	myRouter := mux.NewRouter().StrictSlash(true)

	// Add routes and handle functions
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/articles", h.GetAllArticles).Methods(http.MethodGet)
	myRouter.HandleFunc("/articles/{id}", h.GetArticle).Methods(http.MethodGet)
	myRouter.HandleFunc("/articles", h.AddArticle).Methods(http.MethodPost)
	myRouter.HandleFunc("/articles/{id}", h.UpdateArticle).Methods(http.MethodPut)
	myRouter.HandleFunc("/articles/{id}", h.DeleteArticle).Methods(http.MethodDelete)

	// we will implement these handlers in the next sections
	myRouter.HandleFunc("/signin", auth.Signin)
	myRouter.Handle("/welcome", auth.IsAuthorized(welcome))
	myRouter.HandleFunc("/refresh", auth.Refresh)
	myRouter.HandleFunc("/logout", auth.Logout)

	// Log application startup
	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

func main() {
	// Retreive database connection information from enviroment variables
	url, port, user, password, dbname :=
		os.Getenv("POSTGRES_URL"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB")

	// Connect to the database
	DB := db.Connect(url, port, user, password, dbname)

	// Create the necessary tables in the database
	db.CreateTable(DB)

	// Handle requests
	handleRequests(DB)

	// Close database connection
	db.CloseConnection(DB)
}
