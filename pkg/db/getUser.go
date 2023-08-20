package db

import (
	"log"
	"os"

	"github.com/elman23/articleapi/pkg/handlers"
	"github.com/elman23/articleapi/pkg/models"
)

func GetUser(username string) (models.User, error) {
	log.Println("Query: [GetUser].")

	// Retreive database connection information from enviroment variables
	url, port, dbuser, password, dbname :=
		os.Getenv("POSTGRES_URL"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB")

	// Connect to the database
	DB := Connect(url, port, dbuser, password, dbname)
	h := handlers.New(DB)

	queryStmt := `SELECT * FROM users WHERE username = $1 ;`
	results, err := h.DB.Query(queryStmt, username)
	if err != nil {
		log.Println("Failed to execute query!", err)
		// w.WriteHeader(500)
		return models.User{}, err
	}

	var user models.User
	for results.Next() {
		err = results.Scan(&user.Id, &user.Username, &user.Password)
		if err != nil {
			log.Println("Failed to scan item!", err)
			// w.WriteHeader(500)
			return models.User{}, err
		}
	}

	CloseConnection(DB)
	log.Println("DB connection closed!")

	return user, nil
}
