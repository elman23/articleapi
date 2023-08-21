package db

import (
	"log"

	"github.com/elman23/articleapi/pkg/hashing"
	"github.com/google/uuid"
)

func (s *DbService) AddUser(username string, password string) (string, error) {
	log.Println("Query: [AddUser].")

	id := (uuid.New()).String()
	queryStmt := `INSERT INTO users (id,username,password) VALUES ($1, $2, $3) RETURNING id;`
	hashedPassword, err := hashing.HashPassword(password)
	err = s.DB.QueryRow(queryStmt, &id, &username, &hashedPassword).Scan(&id)
	if err != nil {
		log.Println("Failed to scan item!", err)
		return "", err
	}

	return id, nil
}
