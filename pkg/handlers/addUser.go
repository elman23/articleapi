package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/elman23/articleapi/pkg/db"
	"github.com/elman23/articleapi/pkg/models"
)

func (h handler) AddUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Endpoint: [AddUser].")

	// Read to request body
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Fatalln(err)
		w.WriteHeader(500)
		return
	}
	var user models.User
	json.Unmarshal(body, &user)

	// article.Id = (uuid.New()).String()
	// queryStmt := `INSERT INTO articles (id,title,description,content) VALUES ($1, $2, $3, $4) RETURNING id;`
	// err = h.DB.QueryRow(queryStmt, &article.Id, &article.Title, &article.Desc, &article.Content).Scan(&article.Id)
	var singeton db.DbServiceSingleton
	dbService := singeton.GetService()
	userId, err := dbService.AddUser(user.Username, user.Password)
	if err != nil {
		log.Println("Failed to execute query!", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(userId)
}
