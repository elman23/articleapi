package db

import (
	"log"

	"github.com/elman23/articleapi/pkg/models"
)

func (s *DbService) GetUser(username string) (models.User, error) {
	log.Println("Query: [GetUser].")

	queryStmt := `SELECT * FROM users WHERE username = $1 ;`
	results, err := s.DB.Query(queryStmt, username)
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

	return user, nil
}
